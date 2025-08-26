// Package policy provides DDD Policy pattern implementation for encapsulating business rules
package policy

import (
	"context"
)

// Policy defines the interface for all business rule policies
type Policy[T any] interface {
	Evaluate(ctx context.Context, subject T) Result
	Name() string
}

// BasePolicy provides common functionality for concrete policies
type BasePolicy struct {
	name string
}

// NewBasePolicy creates a new base policy
func NewBasePolicy(name string) BasePolicy {
	return BasePolicy{name: name}
}

// Name returns the policy name
func (bp BasePolicy) Name() string {
	return bp.name
}
