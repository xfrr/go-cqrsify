package policy

import (
	"context"
	"fmt"
)

// PolicyEngine manages and evaluates policies
type PolicyEngine[T any] struct {
	policies map[string]Policy[T]
}

// NewPolicyEngine creates a new policy engine
func NewPolicyEngine[T any]() *PolicyEngine[T] {
	return &PolicyEngine[T]{
		policies: make(map[string]Policy[T]),
	}
}

// Register registers a new policy
func (pe *PolicyEngine[T]) Register(policy Policy[T]) {
	pe.policies[policy.Name()] = policy
}

// Evaluate evaluates a specific policy by name
func (pe *PolicyEngine[T]) Evaluate(ctx context.Context, policyName string, subject T) (Result, error) {
	policy, exists := pe.policies[policyName]
	if !exists {
		return Result{}, fmt.Errorf("policy '%s' not found", policyName)
	}
	return policy.Evaluate(ctx, subject), nil
}

// EvaluateAll evaluates all registered policies with AND logic
func (pe *PolicyEngine[T]) EvaluateAll(ctx context.Context, subject T) Result {
	if len(pe.policies) == 0 {
		return Allow("No policies registered")
	}

	var policies []Policy[T]
	for _, policy := range pe.policies {
		policies = append(policies, policy)
	}

	composite := NewCompositePolicy("all-policies", AND, policies...)
	return composite.Evaluate(ctx, subject)
}

// GetPolicyNames returns all registered policy names
func (pe *PolicyEngine[T]) GetPolicyNames() []string {
	var names []string
	for name := range pe.policies {
		names = append(names, name)
	}
	return names
}
