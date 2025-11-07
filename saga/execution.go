package saga

import (
	"context"
)

type Execution struct {
	SagaID    string
	Def       *Definition
	Instance  *Instance
	StepIndex int
	StepData  map[string]any
	Store     Store
}

// Accessors for shared data.
func (e *Execution) Get(key string) (any, bool) {
	if e.StepData == nil {
		return nil, false
	}
	v, ok := e.StepData[key]
	return v, ok
}
func (e *Execution) Set(key string, v any) {
	if e.StepData == nil {
		e.StepData = map[string]any{}
	}
	e.StepData[key] = v
}

// Helper to compute idempotency token if defined.
func (e *Execution) IdempotencyKey() string {
	step := e.Def.Steps[e.StepIndex]
	if step.IdempotencyFn == nil {
		return ""
	}
	return step.IdempotencyFn(e)
}

// Provide a step-scoped context if needed (e.g., propagating saga IDs).
func (e *Execution) WithStepContext(ctx context.Context) context.Context {
	return ctx
}
