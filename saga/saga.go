package saga

import (
	"context"
	"errors"
	"time"

	"github.com/xfrr/go-cqrsify/pkg/retry"
)

// Status of a saga or step.
type Status string

const (
	StatusPending           Status = "PENDING"
	StatusRunning           Status = "RUNNING"
	StatusFailed            Status = "FAILED"
	StatusCompensating      Status = "COMPENSATING"
	StatusCompensateSuccess Status = "COMPENSATE_SUCCESS"
	StatusCompensateFailed  Status = "COMPENSATE_FAILED"
	StatusCompleted         Status = "COMPLETED"
	StatusCancelled         Status = "CANCELLED"
)

var (
	ErrAlreadyTerminal = errors.New("saga is already in a terminal state")
	ErrConflict        = errors.New("concurrent modification detected")
	ErrLocked          = errors.New("resource is locked by another worker")
	ErrNotFound        = errors.New("not found")
)

type StepAction func(ctx context.Context, ex *Execution) error
type StepCompensation func(ctx context.Context, ex *Execution) error

type Step struct {
	Name                     string
	Action                   StepAction
	Compensate               StepCompensation        // optional
	Timeout                  time.Duration           // optional per-step timeout
	IdempotencyFn            func(*Execution) string // optional idempotency key
	RetryOptions             retry.Options           // optional; uses coordinator defaults if zero-valued
	CompensationRetryOptions retry.Options           // optional; uses coordinator defaults if zero-valued
}

type Definition struct {
	Name  string
	Steps []Step
}

type StepState struct {
	Index      int
	Name       string
	Status     Status
	Attempt    int
	StartedAt  time.Time
	FinishedAt time.Time
	ErrorMsg   string
	// TODO: use generics instead of map[string]any
	Data map[string]any // step-scoped stored values (serialized by Store)
}

type Instance struct {
	ID        string
	Name      string
	Input     map[string]any
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	Current   int // index of next step to run, 0..len(Steps)
	Revision  int // optimistic concurrency control
	Steps     []StepState
	Metadata  map[string]string
}

func (si *Instance) IncrementRevision() { si.Revision++ }

func (si *Instance) Terminal() bool {
	switch si.Status {
	case StatusCompleted, StatusFailed, StatusCancelled:
		return true
	case StatusCompensating, StatusPending, StatusRunning, StatusCompensateSuccess, StatusCompensateFailed:
		return false
	default:
		return false
	}
}
