package saga

import (
	"context"
)

type contextKeySagaID struct{}
type contextKeyStepIndex struct{}
type contextKeyStepName struct{}
type contextKeyStepAttempt struct{}

type Execution struct {
	SagaID    string
	Def       *Definition
	Instance  *Instance
	StepIndex int
	StepData  map[string]any
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
	ctx = context.WithValue(ctx, contextKeySagaID{}, e.SagaID)
	ctx = context.WithValue(ctx, contextKeyStepIndex{}, e.StepIndex)
	ctx = context.WithValue(ctx, contextKeyStepName{}, e.Def.Steps[e.StepIndex].Name)
	ctx = context.WithValue(ctx, contextKeyStepAttempt{}, e.Instance.Steps[e.StepIndex].Attempt)
	return ctx
}

func SagaIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeySagaID{}).(string)
	return v, ok
}

func StepIndexFromContext(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(contextKeyStepIndex{}).(int)
	return v, ok
}

func StepNameFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyStepName{}).(string)
	return v, ok
}

func StepAttemptFromContext(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(contextKeyStepAttempt{}).(int)
	return v, ok
}
