package saga_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/xfrr/go-cqrsify/pkg/retry"
	"github.com/xfrr/go-cqrsify/saga"
)

type mockLocker struct {
	mu        sync.Mutex
	lockedKey string
	tryOK     bool
	tryErr    error
	unlockErr error
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

type mockStore struct {
	mu    sync.Mutex
	data  map[string]*saga.Instance
	failC bool
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
	s.data[inst.ID] = inst
	return nil
}

type mockUUID struct {
	mu  sync.Mutex
	cur int
}

func (g *mockUUID) New() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cur++
	return "saga-" + strconv.Itoa(g.cur)
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
	started, completed         int
	compensating, compFinished int
	stepStart, stepOK, stepKO  int
	compOK, compKO             int
}

func (h *hookRecorder) hooks() saga.Hooks {
	return saga.Hooks{
		OnSagaStarted: func(_ context.Context, _ *saga.Instance) { h.inc(&h.started) },
		OnSagaCompleted: func(_ context.Context, _ *saga.Instance) {
			h.inc(&h.completed)
		},
		OnSagaCompensating: func(_ context.Context, _ *saga.Instance, _ int) {
			h.inc(&h.compensating)
		},
		OnSagaCompensatingFinished: func(_ context.Context, _ *saga.Instance) {
			h.inc(&h.compFinished)
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
	s.Require().NoError(err)

	inst, err := s.store.Load(s.T().Context(), id)
	s.Require().NoError(err)

	// Final status should be one of compensation terminal states; with all compensations OK -> success
	s.Equal(saga.StatusCompleted, inst.Status)
	s.Equal(1, s.hrec.stepKO) // one failure
	s.Equal(1, s.hrec.compOK) // one successful compensation
	s.Equal(0, s.hrec.compKO) // no failed compensations
	s.Equal(1, callsComp)     // only first step compensated
	s.Equal(1, inst.Current)  // failed at step index 1
	s.Equal(saga.StatusCompensateSuccess, inst.Steps[0].Status)
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
	s.Equal(saga.StatusCompleted, inst.Status)
	s.Equal(1, s.hrec.compensating)
	s.Equal(1, s.hrec.compFinished)
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
