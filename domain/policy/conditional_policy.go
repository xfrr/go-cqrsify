package policy

import "context"

// ConditionalPolicy allows dynamic policy evaluation based on conditions
type ConditionalPolicy[T any] struct {
	BasePolicy
	condition func(ctx context.Context, subject T) bool
	policy    Policy[T]
}

// NewConditionalPolicy creates a conditional policy
func NewConditionalPolicy[T any](name string, condition func(ctx context.Context, subject T) bool, policy Policy[T]) *ConditionalPolicy[T] {
	return &ConditionalPolicy[T]{
		BasePolicy: NewBasePolicy(name),
		condition:  condition,
		policy:     policy,
	}
}

// Evaluate evaluates the wrapped policy only if condition is met
func (cp *ConditionalPolicy[T]) Evaluate(ctx context.Context, subject T) Result {
	if !cp.condition(ctx, subject) {
		return Allow("Condition not met, policy skipped")
	}
	return cp.policy.Evaluate(ctx, subject)
}
