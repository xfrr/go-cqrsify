package saga

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xfrr/go-cqrsify/pkg/lock"
	"github.com/xfrr/go-cqrsify/pkg/multierror"
	"github.com/xfrr/go-cqrsify/pkg/retry"
)

var (
	ErrNilDefinition = errors.New("saga definition is nil")
	ErrEmptySagaID   = errors.New("saga ID is empty")
	ErrLockLost      = errors.New("saga lock lease lost")
)

// RetryFactory builds a retrier for a given step.
type RetryFactory func(step Step) *retry.Retrier

type CoordinatorConfig struct {
	// LockTTL is the duration for which the saga lock is held.
	// If zero or negative, the lock is held without expiration (if supported by the locker).
	LockTTL time.Duration
	// Hooks are the saga lifecycle hooks.
	Hooks Hooks
	// TimeProvider is the provider for the current time.
	TimeProvider TimeProvider
	// UUIDProvider is the provider for generating UUIDs.
	UUIDProvider UUIDProvider
	// MaxCompTime is the maximum duration for compensation.
	MaxCompTime time.Duration
	// RetryFactory builds a retrier for saga step executions.
	// If nil, a default retrier factory is used.
	RetryFactory RetryFactory
	// CompensationRetryFactory builds a retrier for saga step compensations.
	// If nil, a default retrier factory is used.
	CompensationRetryFactory RetryFactory
}

type Coordinator struct {
	def    *Definition
	store  Store
	locker lock.Locker
	cfg    CoordinatorConfig
}

type compensationTrigger string

const (
	compensationTriggerFailure compensationTrigger = "failure"
	compensationTriggerCancel  compensationTrigger = "cancel"
)

// NewCoordinator creates a new saga Coordinator.
func NewCoordinator(def *Definition, store Store, locker lock.Locker, cfg CoordinatorConfig) *Coordinator {
	if cfg.RetryFactory == nil {
		cfg.RetryFactory = stepActionRetryFactory()
	}
	if cfg.CompensationRetryFactory == nil {
		cfg.CompensationRetryFactory = stepCompensationRetryFactory()
	}
	if cfg.UUIDProvider == nil {
		cfg.UUIDProvider = DefaultUUIDProvider
	}
	if cfg.TimeProvider == nil {
		cfg.TimeProvider = DefaultTimeProvider
	}

	return &Coordinator{def: def, store: store, locker: locker, cfg: cfg}
}

// Start initiates a new saga instance with the given input and metadata.
// It returns the ID of the newly created saga instance.
func (c *Coordinator) Start(ctx context.Context, input map[string]any, metadata map[string]string) (string, error) {
	if c.def == nil {
		return "", ErrNilDefinition
	}

	inst, err := c.newInstance(input, metadata)
	if err != nil {
		return "", err
	}

	if createErr := c.store.Create(ctx, inst); createErr != nil {
		return "", createErr
	}

	if c.cfg.Hooks.OnSagaStarted != nil {
		c.cfg.Hooks.OnSagaStarted(ctx, inst)
	}
	return inst.ID, nil
}

// Run executes the saga with the given ID.
// It acquires a lock to ensure exclusive execution.
func (c *Coordinator) Run(ctx context.Context, sagaID string) error {
	if c.def == nil {
		return ErrNilDefinition
	}
	if sagaID == "" {
		return ErrEmptySagaID
	}

	lockCtx, cleanup, keepaliveLost, err := c.acquireLockWithKeepalive(ctx, sagaID)
	if err != nil {
		return err
	}
	defer cleanup()

	inst, err := c.store.Load(ctx, sagaID)
	if err != nil {
		return err
	}
	if shouldResume, trigger := c.resumeCompensationTrigger(inst); shouldResume {
		if compensationErr := c.compensate(lockCtx, inst, trigger); compensationErr != nil {
			return compensationErr
		}
		if trigger == compensationTriggerCancel {
			return nil
		}
	}
	if inst.Terminal() {
		return nil
	}
	if runningStatusError := c.ensureRunningStatus(ctx, inst); runningStatusError != nil {
		return runningStatusError
	}

	if stepExecutionErr := c.runSteps(lockCtx, inst, keepaliveLost); stepExecutionErr != nil {
		if c.cfg.Hooks.OnSagaFailed != nil {
			c.cfg.Hooks.OnSagaFailed(ctx, inst, stepExecutionErr)
		}
		// mark failed and compensate
		if setErr := c.setSagaStatus(ctx, inst, StatusFailed); setErr != nil {
			return setErr
		}
		if compensationErr := c.compensate(lockCtx, inst, compensationTriggerFailure); compensationErr != nil {
			return fmt.Errorf("step execution failed: %w; compensation failed: %v", stepExecutionErr, compensationErr)
		}
		return stepExecutionErr
	}

	if setStatusError := c.setSagaStatus(ctx, inst, StatusCompleted); setStatusError != nil {
		return setStatusError
	}
	if c.cfg.Hooks.OnSagaCompleted != nil {
		c.cfg.Hooks.OnSagaCompleted(ctx, inst)
	}
	return nil
}

// Cancel aborts the saga with the given ID and triggers compensation.
// It returns an immediate error if the saga is already in a terminal state.
func (c *Coordinator) Cancel(ctx context.Context, sagaID string) error {
	if c.def == nil {
		return ErrNilDefinition
	}
	if sagaID == "" {
		return ErrEmptySagaID
	}

	lockCtx, cleanup, _, err := c.acquireLockWithKeepalive(ctx, sagaID)
	if err != nil {
		return err
	}
	defer cleanup()

	inst, err := c.store.Load(lockCtx, sagaID)
	if err != nil {
		return err
	}
	if inst.Terminal() {
		return ErrAlreadyTerminal
	}
	if sagaStatusErr := c.setSagaStatus(lockCtx, inst, StatusCancelled); sagaStatusErr != nil {
		return sagaStatusErr
	}
	return c.compensate(lockCtx, inst, compensationTriggerCancel)
}

func (c *Coordinator) newInstance(input map[string]any, metadata map[string]string) (*Instance, error) {
	now := c.now()
	id, err := c.cfg.UUIDProvider.New()
	if err != nil {
		return nil, fmt.Errorf("failed to generate saga ID: %w", err)
	}

	inst := &Instance{
		ID:        id,
		Name:      c.def.Name,
		Input:     input,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
		Current:   0,
		Revision:  0,
		Steps:     make([]StepState, len(c.def.Steps)),
		Metadata:  metadata,
	}
	for i, st := range c.def.Steps {
		inst.Steps[i] = StepState{
			Index:  i,
			Name:   st.Name,
			Status: StatusPending,
			Data:   map[string]any{},
		}
	}
	return inst, nil
}

func (c *Coordinator) ensureRunningStatus(ctx context.Context, inst *Instance) error {
	if inst.Status != StatusPending {
		return nil
	}
	return c.setSagaStatus(ctx, inst, StatusRunning)
}

// lock acquisition with keepalive support.
func (c *Coordinator) acquireLockWithKeepalive(
	ctx context.Context,
	sagaID string,
) (lockCtx context.Context, cleanup func(), keepaliveLost <-chan error, err error) {
	lockKey := fmt.Sprintf("saga:%s", sagaID)

	// 	try to acquire lock
	isLocked, err := c.locker.TryLock(ctx, lockKey, c.cfg.LockTTL)
	if err != nil {
		return nil, func() {}, nil, err
	}
	if !isLocked {
		return nil, func() {}, nil, ErrLocked
	}

	// start keepalive if supported
	var wg sync.WaitGroup
	lockCtx, cancel := context.WithCancel(ctx)
	stopKeepalive := make(chan struct{})
	keepAliveLost := make(chan error, 1)
	if renewer, isRenewer := c.locker.(lock.Renewer); isRenewer && c.cfg.LockTTL > 0 {
		wg.Go(func() { c.keepalive(lockCtx, renewer, lockKey, stopKeepalive, keepAliveLost) })
	}

	cleanup = func() {
		close(stopKeepalive)
		wg.Wait()
		cancel()
		_ = c.locker.Unlock(ctx, lockKey)
	}
	return lockCtx, cleanup, keepAliveLost, nil
}

// keepalive periodically refreshes the lock lease. If it fails or the lease is lost,
// it reports once through keepaliveLost and exits.
func (c *Coordinator) keepalive(
	ctx context.Context,
	renewer lock.Renewer,
	key string,
	stop <-chan struct{},
	keepaliveLost chan<- error,
) {
	interval := c.cfg.LockTTL / 3
	if interval <= 0 {
		interval = time.Second
	}
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ctx.Done():
			return
		case <-t.C:
			ok, err := renewer.Renew(ctx, key, c.cfg.LockTTL)
			if err != nil || !ok {
				select {
				case keepaliveLost <- ifErr(err, ErrLockLost):
				case <-ctx.Done():
				}
				return
			}
		}
	}
}

func (c *Coordinator) runSteps(ctx context.Context, inst *Instance, keepaliveLost <-chan error) error {
	for inst.Current < len(c.def.Steps) {
		if err := c.checkKeepalive(keepaliveLost); err != nil {
			inst.FailureReason = "lock_lost"
			return err
		}

		step := c.def.Steps[inst.Current]
		ss := &inst.Steps[inst.Current]

		// skip if already persisted as completed.
		if ss.Status == StatusCompleted {
			inst.Current++
			continue
		}

		if err := c.runSingleStep(ctx, inst, &step, ss); err != nil {
			return err
		}
	}
	return nil
}

func (c *Coordinator) runSingleStep(ctx context.Context, inst *Instance, step *Step, ss *StepState) error {
	r := c.cfg.RetryFactory(*step)

	return r.Do(ctx, func(runCtx context.Context) error {
		execCtx, cancel := c.withStepTimeout(runCtx, step.Timeout)
		if cancel != nil {
			defer cancel()
		}

		ex := c.buildExecution(inst, ss)
		if err := c.markStepAttemptStart(execCtx, inst, *step, ss, ex.IdempotencyKey()); err != nil {
			return err
		}

		if err := c.executeAction(execCtx, inst, ss, step, ex); err != nil {
			return err
		}
		return c.markStepSuccess(execCtx, inst, ss, ex.StepData)
	})
}

func (c *Coordinator) checkKeepalive(ka <-chan error) error {
	select {
	case err := <-ka:
		if err == nil {
			return ErrLockLost
		}
		return err
	default:
		return nil
	}
}

func (c *Coordinator) withStepTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, timeout)
}

func (c *Coordinator) buildExecution(inst *Instance, ss *StepState) *Execution {
	return &Execution{
		SagaID:    inst.ID,
		Def:       c.def,
		Instance:  inst,
		StepIndex: inst.Current,
		StepData:  ss.Data,
	}
}

func (c *Coordinator) resumeCompensationTrigger(inst *Instance) (bool, compensationTrigger) {
	if inst == nil {
		return false, ""
	}

	if !c.hasIncompleteCompensation(inst) {
		return false, ""
	}

	switch inst.Status {
	case StatusCompensating:
		return true, compensationTriggerFailure
	case StatusCancelled:
		return true, compensationTriggerCancel
	default:
		return false, ""
	}
}

func (c *Coordinator) hasIncompleteCompensation(inst *Instance) bool {
	for i := c.compensationStartIndex(inst); i >= 0; i-- {
		if i >= len(inst.Steps) {
			continue
		}
		ss := inst.Steps[i]
		if ss.Status == StatusCompleted {
			return true
		}
	}
	return false
}

func (c *Coordinator) compensationStartIndex(inst *Instance) int {
	if inst == nil {
		return -1
	}
	if len(inst.Steps) == 0 {
		return -1
	}
	if inst.Current <= 0 {
		return 0
	}
	startIdx := inst.Current - 1
	if startIdx >= len(inst.Steps) {
		return len(inst.Steps) - 1
	}
	return startIdx
}

func (c *Coordinator) compensate(ctx context.Context, inst *Instance, trigger compensationTrigger) error {
	startIdx := c.compensationStartIndex(inst)
	if startIdx < 0 {
		c.finishCompensationStatus(inst, trigger, false, false, 0, 0)
		return c.store.Save(ctx, inst)
	}

	if c.cfg.Hooks.OnSagaCompensating != nil {
		c.cfg.Hooks.OnSagaCompensating(ctx, inst, startIdx)
	}

	if err := c.beginCompensation(ctx, inst); err != nil {
		return err
	}

	merr, attempted, succeeded, deadlineExceeded := c.compensateSteps(ctx, inst, startIdx)
	c.finishCompensationStatus(inst, trigger, merr.HasErrors(), deadlineExceeded, attempted, succeeded)

	if err := c.store.Save(ctx, inst); err != nil {
		return err
	}
	if c.cfg.Hooks.OnSagaCompensatingFinished != nil {
		c.cfg.Hooks.OnSagaCompensatingFinished(ctx, inst)
	}
	return merr.ErrorOrNil()
}

func (c *Coordinator) beginCompensation(ctx context.Context, inst *Instance) error {
	inst.Status = StatusCompensating
	inst.UpdatedAt = c.now()
	return c.store.Save(ctx, inst)
}

func (c *Coordinator) compensateSteps(
	ctx context.Context,
	inst *Instance,
	startIdx int,
) (*multierror.MultiError, int, int, bool) {
	merr := multierror.New()

	var deadline time.Time
	if c.cfg.MaxCompTime > 0 {
		deadline = c.now().Add(c.cfg.MaxCompTime)
	}

	attempted, succeeded := 0, 0
	deadlineExceeded := false

	for i := startIdx; i >= 0; i-- {
		if !deadline.IsZero() && c.now().After(deadline) {
			deadlineExceeded = true
			break
		}

		step := c.def.Steps[i]
		ss := &inst.Steps[i]
		if ss.Status != StatusCompleted || step.Compensate == nil {
			continue
		}

		attempted++
		if err := c.compensateStep(ctx, inst, i, &step, ss); err != nil {
			merr.Append(err)
			continue
		}
		succeeded++
	}

	return merr, attempted, succeeded, deadlineExceeded
}

func (c *Coordinator) compensateStep(
	ctx context.Context,
	inst *Instance,
	index int,
	step *Step,
	ss *StepState,
) error {
	ss.Status = StatusCompensating
	ss.StartedAt = c.now()
	inst.UpdatedAt = ss.StartedAt
	if err := c.store.Save(ctx, inst); err != nil {
		return err
	}

	ex := &Execution{
		SagaID:    inst.ID,
		Def:       c.def,
		Instance:  inst,
		StepIndex: index,
		StepData:  ss.Data,
	}

	cr := c.cfg.CompensationRetryFactory(*step)
	compensationErr := cr.Do(ctx, func(runCtx context.Context) error {
		return c.executeCompensation(runCtx, inst, ss, step, ex)
	})

	if compensationErr != nil {
		ss.Status = StatusCompensateFailed
		ss.ErrorMsg = "compensation: " + compensationErr.Error()
		inst.FailureReason = "compensation_failed"
		ss.FinishedAt = c.now()
		inst.UpdatedAt = ss.FinishedAt
		if err := c.store.Save(ctx, inst); err != nil {
			return err
		}
		return compensationErr
	}

	ss.Status = StatusCompensateSuccess
	ss.ErrorMsg = ""
	ss.FinishedAt = c.now()
	inst.UpdatedAt = ss.FinishedAt
	if c.cfg.Hooks.OnStepCompensationOK != nil {
		c.cfg.Hooks.OnStepCompensationOK(ctx, inst, *ss)
	}
	return c.store.Save(ctx, inst)
}

func (c *Coordinator) finishCompensationStatus(
	inst *Instance,
	trigger compensationTrigger,
	hasErrors bool,
	deadlineExceeded bool,
	attempted,
	succeeded int,
) {
	inst.UpdatedAt = c.now()
	switch {
	case hasErrors:
		inst.Status = StatusFailed
		if inst.FailureReason == "" {
			inst.FailureReason = "compensation_failed"
		}
	case deadlineExceeded && attempted != succeeded:
		inst.Status = StatusFailed
		inst.FailureReason = "compensation_deadline_exceeded"
	default:
		if trigger == compensationTriggerCancel {
			inst.Status = StatusCancelled
			return
		}
		inst.Status = StatusCompleted
	}
}

func (c *Coordinator) executeAction(ctx context.Context, inst *Instance, ss *StepState, step *Step, ex *Execution) (err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("step %q panicked: %v", step.Name, r)
			c.onStepFailure(ctx, inst, ss, panicErr)
			if saveErr := c.store.Save(ctx, inst); saveErr != nil {
				err = fmt.Errorf("step panic and save failure state: %w: %v", panicErr, saveErr)
				return
			}
			err = panicErr
		}
	}()

	if actionErr := step.Action(ex.WithStepContext(ctx), ex); actionErr != nil {
		c.onStepFailure(ctx, inst, ss, actionErr)
		inst.FailureReason = "step_action_failed"
		if saveErr := c.store.Save(ctx, inst); saveErr != nil {
			return fmt.Errorf("step action failed: %w; save failure state failed: %v", actionErr, saveErr)
		}
		return actionErr
	}

	return nil
}

func (c *Coordinator) executeCompensation(ctx context.Context, inst *Instance, ss *StepState, step *Step, ex *Execution) (err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("compensation of step %q panicked: %v", step.Name, r)
			c.onCompensationFailure(ctx, inst, ss, panicErr)
			if saveErr := c.store.Save(ctx, inst); saveErr != nil {
				err = fmt.Errorf("compensation panic and save failure state: %w: %v", panicErr, saveErr)
				return
			}
			err = panicErr
		}
	}()

	if compensationErr := step.Compensate(ex.WithStepContext(ctx), ex); compensationErr != nil {
		c.onCompensationFailure(ctx, inst, ss, compensationErr)
		inst.FailureReason = "compensation_failed"
		if saveErr := c.store.Save(ctx, inst); saveErr != nil {
			return fmt.Errorf("compensation failed: %w; save failure state failed: %v", compensationErr, saveErr)
		}
		return compensationErr
	}

	return nil
}

func (c *Coordinator) setSagaStatus(ctx context.Context, inst *Instance, s Status) error {
	inst.Status = s
	inst.UpdatedAt = c.now()
	return c.store.Save(ctx, inst)
}

func (c *Coordinator) markStepAttemptStart(ctx context.Context, inst *Instance, step Step, ss *StepState, idempotencyKey string) error {
	ss.Status = StatusRunning
	ss.Attempt++
	if step.IdempotencyFn != nil {
		ss.LastIdempotencyKey = idempotencyKey
	}
	ss.StartedAt = c.now()
	inst.UpdatedAt = ss.StartedAt

	if c.cfg.Hooks.OnStepStart != nil {
		c.cfg.Hooks.OnStepStart(ctx, inst, *ss)
	}
	return c.store.Save(ctx, inst)
}

func (c *Coordinator) markStepSuccess(ctx context.Context, inst *Instance, ss *StepState, newData map[string]any) error {
	ss.Status = StatusCompleted
	ss.FinishedAt = c.now()
	ss.ErrorMsg = ""
	if newData != nil {
		ss.Data = newData
	}
	inst.Current++
	inst.UpdatedAt = ss.FinishedAt

	if c.cfg.Hooks.OnStepSuccess != nil {
		c.cfg.Hooks.OnStepSuccess(ctx, inst, *ss)
	}
	return c.store.Save(ctx, inst)
}

func (c *Coordinator) onStepFailure(ctx context.Context, inst *Instance, ss *StepState, err error) {
	ss.ErrorMsg = err.Error()
	if c.cfg.Hooks.OnStepFailure != nil {
		c.cfg.Hooks.OnStepFailure(ctx, inst, *ss, err)
	}
}

func (c *Coordinator) onCompensationFailure(ctx context.Context, inst *Instance, ss *StepState, err error) {
	ss.ErrorMsg = "compensation: " + err.Error()
	if c.cfg.Hooks.OnStepCompensationKO != nil {
		c.cfg.Hooks.OnStepCompensationKO(ctx, inst, *ss, err)
	}
}

func (c *Coordinator) now() time.Time {
	if c.cfg.TimeProvider != nil {
		return c.cfg.TimeProvider.Now()
	}
	return time.Now()
}

func ifErr(err error, fallback error) error {
	if err != nil {
		return err
	}
	return fallback
}
