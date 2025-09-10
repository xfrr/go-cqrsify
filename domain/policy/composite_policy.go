package domainpolicy

import (
	"context"
	"fmt"
)

// CompositePolicy allows combining multiple policies with logical operators
type CompositePolicy[T any] struct {
	name     string
	policies []Policy[T]
	operator LogicalOperator
}

// LogicalOperator defines how multiple policies are combined
type LogicalOperator string

const (
	AND LogicalOperator = "AND"
	OR  LogicalOperator = "OR"
)

// NewCompositePolicy creates a new composite policy
func NewCompositePolicy[T any](name string, operator LogicalOperator, policies ...Policy[T]) *CompositePolicy[T] {
	return &CompositePolicy[T]{
		name:     name,
		policies: policies,
		operator: operator,
	}
}

// Evaluate evaluates all contained policies based on the logical operator
func (cp *CompositePolicy[T]) Evaluate(ctx context.Context, subject T) Result {
	if len(cp.policies) == 0 {
		return Allow("No policies to evaluate")
	}

	switch cp.operator {
	case AND:
		return cp.evaluateAND(ctx, subject)
	case OR:
		return cp.evaluateOR(ctx, subject)
	default:
		return Deny("Invalid logical operator", "INVALID_OPERATOR")
	}
}

// Name returns the policy name
func (cp *CompositePolicy[T]) Name() string {
	return cp.name
}

// evaluateAND returns true only if all policies allow
func (cp *CompositePolicy[T]) evaluateAND(ctx context.Context, subject T) Result {
	for _, policy := range cp.policies {
		result := policy.Evaluate(ctx, subject)
		if !result.Allowed {
			return Deny(
				fmt.Sprintf("Policy '%s' denied: %s", policy.Name(), result.Reason),
				result.Code,
			)
		}
	}
	return Allow("All policies passed")
}

// evaluateOR returns true if at least one policy allows
func (cp *CompositePolicy[T]) evaluateOR(ctx context.Context, subject T) Result {
	var lastResult Result
	for _, policy := range cp.policies {
		result := policy.Evaluate(ctx, subject)
		if result.Allowed {
			return Allow(fmt.Sprintf("Policy '%s' allowed: %s", policy.Name(), result.Reason))
		}
		lastResult = result
	}
	return Deny("No policies allowed the action", lastResult.Code)
}

// Add adds a new policy to the composite
func (cp *CompositePolicy[T]) Add(policy Policy[T]) {
	cp.policies = append(cp.policies, policy)
}
