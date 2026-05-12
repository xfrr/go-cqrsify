package saga_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/xfrr/go-cqrsify/pkg/retry"
	"github.com/xfrr/go-cqrsify/saga"
)

type mockLocker struct {
	mu          sync.Mutex
	lockedKey   string
	tryOK       bool
	tryErr      error
	unlockErr   error
	failRenewAt int
	renewErr    error
	renewCalls  int
}

func (l *mockLocker) TryLock(_ context.Context, key string, _ time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tryErr != nil {
		return false, l.tryErr
	}
	if !l.tryOK {
		return false, nil
	}
	l.lockedKey = key
	return true, nil
}

func (l *mockLocker) Unlock(_ context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if key != l.lockedKey {
		return nil
	}
	return l.unlockErr
}

func (l *mockLocker) Refresh(_ context.Context, key string, _ time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if key != l.lockedKey {
		return false, nil
	}
	return true, nil
}

func (l *mockLocker) Renew(_ context.Context, key string, _ time.Duration) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if key != l.lockedKey {
		return false, nil
	}
	l.renewCalls++
	if l.failRenewAt > 0 && l.renewCalls >= l.failRenewAt {
		if l.renewErr != nil {
			return false, l.renewErr
		}
		return false, nil
	}
	return true, nil
}

type mockStore struct {
	mu         sync.Mutex
	data       map[string]*saga.Instance
	failC      bool
	failSaveAt int
	saveCalls  int
}

func newMemStore() *mockStore {
	return &mockStore{data: map[string]*saga.Instance{}}
}

func (s *mockStore) Create(_ context.Context, inst *saga.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failC {
		return errors.New("create error")
	}
	s.data[inst.ID] = inst
	return nil
}

func (s *mockStore) Load(_ context.Context, id string) (*saga.Instance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	inst, ok := s.data[id]
	if !ok {
		return nil, saga.ErrNotFound
	}
	return inst, nil
}

func (s *mockStore) Save(_ context.Context, inst *saga.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.saveCalls++
	if s.failSaveAt > 0 && s.saveCalls == s.failSaveAt {
		return errors.New("save error")
	}
	s.data[inst.ID] = inst
	return nil
}

type mockUUID struct {
	mu  sync.Mutex
	cur int
}

func (g *mockUUID) New() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cur++
	return "saga-" + strconv.Itoa(g.cur), nil
}

type mockClock struct {
	mu   sync.Mutex
	t    time.Time
	step time.Duration
	// if fixedAfterFirst is set, all calls after the first return t+fixedAfterFirst
	fixedAfterFirst *time.Duration
	calls           int
}

func newScriptClock(start time.Time, step time.Duration) *mockClock {
	return &mockClock{t: start, step: step}
}

func (c *mockClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls++
	if c.fixedAfterFirst != nil && c.calls > 1 {
		return c.t.Add(*c.fixedAfterFirst)
	}
	if c.calls == 1 {
		return c.t
	}
	c.t = c.t.Add(c.step)
	return c.t
}

type hookRecorder struct {
	mu                         sync.Mutex
	started, completed, failed int
	compensating, compFinished int
	stepStart, stepOK, stepKO  int
	compOK, compKO             int
	compFinishedStatuses       []saga.Status
}

func (h *hookRecorder) hooks() saga.Hooks {
	return saga.Hooks{
		OnSagaStarted: func(_ context.Context, _ *saga.Instance) { h.inc(&h.started) },
		OnSagaCompleted: func(_ context.Context, _ *saga.Instance) {
			h.inc(&h.completed)
		},
		OnSagaFailed: func(_ context.Context, _ *saga.Instance, _ error) {
			h.inc(&h.failed)
		},
		OnSagaCompensating: func(_ context.Context, _ *saga.Instance, _ int) {
			h.inc(&h.compensating)
		},
		OnSagaCompensatingFinished: func(_ context.Context, inst *saga.Instance) {
			h.inc(&h.compFinished)
			h.addCompFinishedStatus(inst.Status)
		},
		OnStepStart: func(_ context.Context, _ *saga.Instance, _ saga.StepState) {
			h.inc(&h.stepStart)
		},
		OnStepSuccess: func(_ context.Context, _ *saga.Instance, _ saga.StepState) {
			h.inc(&h.stepOK)
		},
		OnStepFailure: func(_ context.Context, _ *saga.Instance, _ saga.StepState, _ error) {
			h.inc(&h.stepKO)
		},
		OnStepCompensationOK: func(_ context.Context, _ *saga.Instance, _ saga.StepState) {
			h.inc(&h.compOK)
		},
		OnStepCompensationKO: func(_ context.Context, _ *saga.Instance, _ saga.StepState, _ error) {
			h.inc(&h.compKO)
		},
	}
}

func (h *hookRecorder) inc(p *int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	*p++
}

func (h *hookRecorder) addCompFinishedStatus(status saga.Status) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.compFinishedStatuses = append(h.compFinishedStatuses, status)
}

type CoordinatorSuite struct {
	suite.Suite

	store  *mockStore
	locker *mockLocker
	uuid   *mockUUID
	clock  *mockClock
	hrec   *hookRecorder

	def saga.Definition
	cfg saga.CoordinatorConfig
}

func (s *CoordinatorSuite) SetupTest() {
	s.store = newMemStore()
	s.locker = &mockLocker{tryOK: true}
	s.uuid = &mockUUID{}
	s.clock = newScriptClock(time.Date(2025, 11, 6, 10, 0, 0, 0, time.UTC), time.Second)
	s.hrec = &hookRecorder{}

	// definition with two simple steps that pass
	s.def = saga.Definition{
		Name: "test-saga",
		Steps: []saga.Step{
			{
				Name: "step-1",
				Action: func(_ context.Context, ex *saga.Execution) error {
					ex.StepData["a"] = 1
					return nil
				},
				Compensate: func(_ context.Context, ex *saga.Execution) error {
					ex.StepData["a-comp"] = true
					return nil
				},
			},
			{
				Name: "step-2",
				Action: func(_ context.Context, ex *saga.Execution) error {
					ex.StepData["b"] = 2
					return nil
				},
				Compensate: func(_ context.Context, ex *saga.Execution) error {
					ex.StepData["b-comp"] = true
					return nil
				},
			},
		},
	}

	// single attempt, no sleeps, to keep tests snappy and deterministic
	fastRetry := func(_ saga.Step) *retry.Retrier {
		opts := retry.DefaultOptions()
		opts.MaxAttempts = 1
		return retry.New(opts)
	}

	s.cfg = saga.CoordinatorConfig{
		LockTTL:      5 * time.Second,
		Hooks:        s.hrec.hooks(),
		TimeProvider: s.clock,
		UUIDProvider: s.uuid,
		MaxCompTime:  0,
		RetryFactory: fastRetry,
	}
}

func (s *CoordinatorSuite) newCoordinator() *saga.Coordinator {
	return saga.NewCoordinator(&s.def, s.store, s.locker, s.cfg)
}

func (s *CoordinatorSuite) TestStart_InitializesInstanceAndCallsHook() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), map[string]any{"k": "v"}, map[string]string{"m": "n"})
	s.Require().NoError(err)
	s.Require().NotEmpty(id)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)

	s.Len(inst.Steps, 2)
	s.Equal("test-saga", inst.Name)
	s.Equal(saga.StatusPending, inst.Status)
	s.Equal(0, inst.Current)
	s.Equal("step-1", inst.Steps[0].Name)
	s.Equal("step-2", inst.Steps[1].Name)
	s.Equal(saga.StatusPending, inst.Steps[0].Status)
	s.Equal(saga.StatusPending, inst.Steps[1].Status)
	s.Equal(map[string]any{"k": "v"}, inst.Input)
	s.Equal(map[string]string{"m": "n"}, inst.Metadata)
	s.Equal(1, s.hrec.started)
}

func (s *CoordinatorSuite) TestRun_AllStepsSucceed_CompletesSaga() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)

	s.Equal(saga.StatusCompleted, inst.Status)
	s.Equal(2, inst.Current)
	s.Equal(saga.StatusCompleted, inst.Steps[0].Status)
	s.Equal(saga.StatusCompleted, inst.Steps[1].Status)
	s.Equal(2, s.hrec.stepStart) // both steps started once
	s.Equal(2, s.hrec.stepOK)
	s.Equal(0, s.hrec.stepKO)
	s.Equal(1, s.hrec.completed)
}

func (s *CoordinatorSuite) TestRun_ActionFails_TriggersCompensationAndMarksStatuses() {
	// Replace step-2 with failing action
	failErr := errors.New("boom")
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		return failErr
	}
	callsComp := 0
	s.def.Steps[0].Compensate = func(_ context.Context, _ *saga.Execution) error {
		callsComp++
		return nil
	}
	// Force single attempt and no delay
	s.cfg.RetryFactory = func(_ saga.Step) *retry.Retrier {
		opts := retry.DefaultOptions()
		opts.MaxAttempts = 1
		return retry.New(opts)
	}
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "boom")

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)

	// Saga-level terminal status after compensation is COMPENSATE_SUCCESS when all compensations succeed.
	s.Equal(saga.StatusCompensateSuccess, inst.Status)
	s.Equal(1, s.hrec.stepKO) // one failure
	s.Equal(1, s.hrec.failed)
	s.Equal(1, s.hrec.compOK) // one successful compensation
	s.Equal(0, s.hrec.compKO) // no failed compensations
	s.Equal(1, callsComp)     // only first step compensated
	s.Equal(1, inst.Current)  // failed at step index 1
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
	s.Equal(saga.StatusFailed, inst.Steps[1].Status)
	s.Equal("boom", inst.Steps[1].ErrorMsg)
	s.False(inst.Steps[1].FinishedAt.IsZero())
	s.Equal("step_action_failed", inst.FailureReason)
	s.Require().Len(s.hrec.compFinishedStatuses, 1)
	s.Equal(saga.StatusCompensateSuccess, s.hrec.compFinishedStatuses[0])
}

func (s *CoordinatorSuite) TestRun_ActionPanic_TriggersCompensationAndReturnsError() {
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		panic("kaboom")
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "panicked")

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal(saga.StatusCompensateSuccess, inst.Status)
	s.Equal(saga.StatusFailed, inst.Steps[1].Status)
	s.Contains(inst.Steps[1].ErrorMsg, "panicked")
	s.False(inst.Steps[1].FinishedAt.IsZero())
	s.Equal("step_action_panicked", inst.FailureReason)
	s.Equal(1, s.hrec.failed)
	s.Equal(1, s.hrec.compOK)
}

func (s *CoordinatorSuite) TestRun_CompensationPanic_MarksFailed() {
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		return errors.New("step boom")
	}
	s.def.Steps[0].Compensate = func(_ context.Context, _ *saga.Execution) error {
		panic("compensate panic")
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "compensation")

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal(saga.StatusCompensateFailed, inst.Status)
	s.Equal("compensation_failed", inst.FailureReason)
	s.Equal(saga.StatusCompensateFailed, inst.Steps[0].Status)
	s.GreaterOrEqual(s.hrec.compKO, 1)
	s.Require().Len(s.hrec.compFinishedStatuses, 1)
	s.Equal(saga.StatusCompensateFailed, s.hrec.compFinishedStatuses[0])
}

func (s *CoordinatorSuite) TestRun_ResumesCompensationWhenStatusCompensating() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	inst.Status = saga.StatusCompensating
	inst.Current = 2
	inst.Steps[0].Status = saga.StatusCompleted
	inst.Steps[1].Status = saga.StatusCompleted
	err = s.store.Save(s.T().Context(), inst)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err = s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCompensateSuccess, inst.Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[1].Status)
	s.Equal(2, s.hrec.compOK)
	s.Equal(0, s.hrec.compKO)
	s.Require().Len(s.hrec.compFinishedStatuses, 1)
	s.Equal(saga.StatusCompensateSuccess, s.hrec.compFinishedStatuses[0])
}

func (s *CoordinatorSuite) TestRun_ResumesCompensationWhenStatusCancelledAndIncomplete() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	inst.Status = saga.StatusCancelled
	inst.Current = 2
	inst.Steps[0].Status = saga.StatusCompleted
	inst.Steps[1].Status = saga.StatusCompensateSuccess
	err = s.store.Save(s.T().Context(), inst)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err = s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCancelled, inst.Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[1].Status)
	s.Equal(1, s.hrec.compOK)
	s.Equal(0, s.hrec.compKO)
	s.Require().Len(s.hrec.compFinishedStatuses, 1)
	s.Equal(saga.StatusCancelled, s.hrec.compFinishedStatuses[0])
}

func (s *CoordinatorSuite) TestRun_CancelledAndFullyCompensated_NoOp() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	inst.Status = saga.StatusCancelled
	inst.Current = 2
	inst.Steps[0].Status = saga.StatusCompensateSuccess
	inst.Steps[1].Status = saga.StatusCompensateSuccess
	err = s.store.Save(s.T().Context(), inst)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err = s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCancelled, inst.Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[1].Status)
	s.Equal(0, s.hrec.compensating)
	s.Equal(0, s.hrec.compFinished)
}

func (s *CoordinatorSuite) TestRun_CancelledWithCompletedStepWithoutCompensator_NoReentry() {
	s.def.Steps[0].Compensate = nil
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	inst.Status = saga.StatusCancelled
	inst.Current = 1
	inst.Steps[0].Status = saga.StatusCompleted
	err = s.store.Save(s.T().Context(), inst)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err = s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCancelled, inst.Status)
	s.Equal(saga.StatusCompleted, inst.Steps[0].Status)
	s.Equal(0, s.hrec.compensating)
	s.Equal(0, s.hrec.compFinished)
	s.Equal(0, s.hrec.compOK)
	s.Equal(0, s.hrec.compKO)
}

func (s *CoordinatorSuite) TestRun_UsesCoordinatorCompensationRetryFactory() {
	factoryCalls := 0
	s.cfg.CompensationRetryFactory = func(_ saga.Step) *retry.Retrier {
		factoryCalls++
		opts := retry.DefaultOptions()
		opts.MaxAttempts = 1
		opts.Strategy = retry.ConstantStrategy{Delay: 0}
		return retry.New(opts)
	}

	failErr := errors.New("boom")
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		return failErr
	}
	s.def.Steps[0].Compensate = func(_ context.Context, _ *saga.Execution) error {
		return nil
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "boom")
	s.Equal(1, factoryCalls)
}

func (s *CoordinatorSuite) TestRun_CustomCompensationRetryFactory_RetriesCompensation() {
	failErr := errors.New("boom")
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		return failErr
	}

	compAttempts := 0
	s.def.Steps[0].Compensate = func(_ context.Context, _ *saga.Execution) error {
		compAttempts++
		if compAttempts == 1 {
			return errors.New("transient compensation failure")
		}
		return nil
	}

	s.cfg.CompensationRetryFactory = func(_ saga.Step) *retry.Retrier {
		opts := retry.DefaultOptions()
		opts.MaxAttempts = 2
		opts.Strategy = retry.ConstantStrategy{Delay: 0}
		return retry.New(opts)
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "boom")
	s.Equal(2, compAttempts)

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal(saga.StatusCompensateSuccess, inst.Status)
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
}

func (s *CoordinatorSuite) TestRun_StepCompensationRetryOptions_OverrideCoordinatorDefaults() {
	failErr := errors.New("boom")
	s.def.Steps[1].Action = func(_ context.Context, _ *saga.Execution) error {
		return failErr
	}

	compAttempts := 0
	s.def.Steps[0].Compensate = func(_ context.Context, _ *saga.Execution) error {
		compAttempts++
		if compAttempts == 1 {
			return errors.New("transient compensation failure")
		}
		return nil
	}
	s.def.Steps[0].CompensationRetryOptions = retry.Options{
		MaxAttempts: 1,
		Strategy:    retry.ConstantStrategy{Delay: 0},
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "compensation failed")
	s.Equal(1, compAttempts)

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal(saga.StatusCompensateFailed, inst.Status)
	s.Equal(saga.StatusCompensateFailed, inst.Steps[0].Status)
}

func (s *CoordinatorSuite) TestRun_StepTimeout_ReturnsError() {
	s.def.Steps = []saga.Step{
		{
			Name:    "slow-step",
			Timeout: 10 * time.Millisecond,
			Action: func(ctx context.Context, _ *saga.Execution) error {
				<-time.After(30 * time.Millisecond)
				return ctx.Err()
			},
		},
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, context.DeadlineExceeded.Error())

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal("step_action_failed", inst.FailureReason)
}

func (s *CoordinatorSuite) TestRun_SaveFailureOnStepError_IsReturned() {
	s.def.Steps[0].Action = func(_ context.Context, _ *saga.Execution) error {
		return errors.New("boom")
	}
	s.store.failSaveAt = 3 // running-status save, step-start save, failure-save

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "save failure state failed")
}

func (s *CoordinatorSuite) TestRun_LeaseLost_ReturnsErrLockLost() {
	s.cfg.LockTTL = 15 * time.Millisecond
	s.locker.failRenewAt = 1
	s.def.Steps[0].Action = func(_ context.Context, _ *saga.Execution) error {
		<-time.After(40 * time.Millisecond)
		return nil
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, saga.ErrLockLost)

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal("lock_lost", inst.FailureReason)
}

func (s *CoordinatorSuite) TestRun_TracksIdempotencyAndStepContext() {
	s.def.Steps = []saga.Step{
		{
			Name: "ctx-step",
			IdempotencyFn: func(ex *saga.Execution) string {
				return "idemp-" + ex.SagaID
			},
			Action: func(ctx context.Context, ex *saga.Execution) error {
				sagaID, ok := saga.SagaIDFromContext(ctx)
				if !ok || sagaID != ex.SagaID {
					return errors.New("missing saga id in context")
				}
				stepName, ok := saga.StepNameFromContext(ctx)
				if !ok || stepName != "ctx-step" {
					return errors.New("missing step name in context")
				}
				attempt, ok := saga.StepAttemptFromContext(ctx)
				if !ok || attempt != 1 {
					return errors.New("missing step attempt in context")
				}
				return nil
			},
		},
	}

	c := s.newCoordinator()
	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal("idemp-"+id, inst.Steps[0].LastIdempotencyKey)
}

func (s *CoordinatorSuite) TestCancel_MarksCancelledAndCompensates() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	// Progress saga to have one completed step
	err = c.Run(s.T().Context(), id)
	s.Require().NoError(err)

	// Reset store to a state where first step is completed and second pending
	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	inst.Status = saga.StatusRunning
	inst.Current = 1
	inst.Steps[0].Status = saga.StatusCompleted
	inst.Steps[1].Status = saga.StatusPending
	err = s.store.Save(s.T().Context(), inst)
	s.Require().NoError(err)

	// Cancel should compensate step 0
	err = c.Cancel(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err = s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCancelled, inst.Status)
	s.Equal(1, s.hrec.compensating)
	s.Equal(1, s.hrec.compFinished)
	s.Require().Len(s.hrec.compFinishedStatuses, 1)
	s.Equal(saga.StatusCancelled, s.hrec.compFinishedStatuses[0])
}

func (s *CoordinatorSuite) TestCancel_AcquiresSagaLockKey() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	err = c.Cancel(s.T().Context(), id)
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)
	s.Equal(saga.StatusCancelled, inst.Status)

	s.locker.mu.Lock()
	lockedKey := s.locker.lockedKey
	s.locker.mu.Unlock()
	s.Equal(fmt.Sprintf("saga:%s", id), lockedKey)
}

func (s *CoordinatorSuite) TestCancel_RespectsLock_ReturnsErrLocked() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	s.locker.tryOK = false
	err = c.Cancel(s.T().Context(), id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, saga.ErrLocked)

	inst, loadErr := s.store.Load(s.T().Context(), id)
	s.Require().NoError(loadErr)
	s.Equal(saga.StatusPending, inst.Status)
}

func (s *CoordinatorSuite) TestRun_RespectsLock_ReturnsErrLocked() {
	c := s.newCoordinator()

	id, err := c.Start(s.T().Context(), nil, nil)
	s.Require().NoError(err)

	s.locker.tryOK = false // cannot acquire
	err = c.Run(s.T().Context(), id)
	s.Require().Error(err)
	s.ErrorIs(err, saga.ErrLocked)
}

func (s *CoordinatorSuite) TestStart_WithNilDefinition_ReturnsError() {
	c := saga.NewCoordinator(nil, s.store, s.locker, s.cfg)
	_, err := c.Start(s.T().Context(), nil, nil)
	s.Require().Error(err)
}

func TestCoordinatorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CoordinatorSuite))
}
