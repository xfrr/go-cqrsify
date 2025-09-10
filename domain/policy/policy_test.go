package domainpolicy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	policy "github.com/xfrr/go-cqrsify/domain/policy"
)

type User struct {
	Age      int
	IsActive bool
}

type agePolicy struct {
	policy.BasePolicy
	MinAge int
}

func (p *agePolicy) Evaluate(ctx context.Context, user User) policy.Result {
	if user.Age < p.MinAge {
		return policy.Deny("User does not meet age requirement", "INSUFFICIENT_AGE")
	}
	return policy.Allow("User meets age requirement")
}

func newAgePolicy(minAge int) *agePolicy {
	return &agePolicy{
		BasePolicy: policy.NewBasePolicy("age-policy"),
		MinAge:     minAge,
	}
}

type isActiveUserPolicy struct {
	policy.BasePolicy
}

func (p *isActiveUserPolicy) Evaluate(ctx context.Context, user User) policy.Result {
	if !user.IsActive {
		return policy.Deny("User is not active", "INACTIVE_USER")
	}
	return policy.Allow("User is active")
}

func newActiveUserPolicy() *isActiveUserPolicy {
	return &isActiveUserPolicy{
		BasePolicy: policy.NewBasePolicy("active-user-policy"),
	}
}

func TestResult(t *testing.T) {
	t.Run("NewResult creates result correctly", func(t *testing.T) {
		result := policy.NewResult(true, "test reason", "TEST_CODE")
		assert.True(t, result.Allowed)
		assert.Equal(t, "test reason", result.Reason)
		assert.Equal(t, "TEST_CODE", result.Code)
	})

	t.Run("Allow creates successful result", func(t *testing.T) {
		result := policy.Allow("success")
		assert.True(t, result.Allowed)
		assert.Equal(t, "success", result.Reason)
		assert.Equal(t, "ALLOWED", result.Code)
	})

	t.Run("Deny creates failed result", func(t *testing.T) {
		result := policy.Deny("failed", "FAIL_CODE")
		assert.False(t, result.Allowed)
		assert.Equal(t, "failed", result.Reason)
		assert.Equal(t, "FAIL_CODE", result.Code)
	})
}

func TestCompositePolicy(t *testing.T) {
	agePolicy := newAgePolicy(18)
	activePolicy := newActiveUserPolicy()

	t.Run("AND operator - all policies pass", func(t *testing.T) {
		composite := policy.NewCompositePolicy("admin-access", policy.AND, agePolicy, activePolicy)
		user := User{Age: 25, IsActive: true}

		result := composite.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Equal(t, "All policies passed", result.Reason)
		assert.Equal(t, "admin-access", composite.Name())
	})

	t.Run("AND operator - one policy fails", func(t *testing.T) {
		composite := policy.NewCompositePolicy("admin-access", policy.AND, agePolicy, activePolicy)
		user := User{Age: 25, IsActive: false}

		result := composite.Evaluate(context.Background(), user)

		assert.False(t, result.Allowed)
		assert.Contains(t, result.Reason, "Policy 'active-user-policy' denied")
	})

	t.Run("OR operator - at least one policy passes", func(t *testing.T) {
		composite := policy.NewCompositePolicy("any-access", policy.OR, agePolicy, activePolicy)
		user := User{Age: 25, IsActive: false} // Age passes, active fails

		result := composite.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Contains(t, result.Reason, "Policy 'age-policy' allowed")
	})

	t.Run("OR operator - no policies pass", func(t *testing.T) {
		composite := policy.NewCompositePolicy("any-access", policy.OR, agePolicy, activePolicy)
		user := User{Age: 16, IsActive: false} // Both fail

		result := composite.Evaluate(context.Background(), user)

		assert.False(t, result.Allowed)
		assert.Equal(t, "No policies allowed the action", result.Reason)
	})

	t.Run("Empty composite policy allows by default", func(t *testing.T) {
		composite := policy.NewCompositePolicy[User]("empty", policy.AND)
		user := User{}

		result := composite.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Equal(t, "No policies to evaluate", result.Reason)
	})

	t.Run("Invalid operator returns error", func(t *testing.T) {
		composite := policy.NewCompositePolicy("invalid", "INVALID", agePolicy)
		user := User{}

		result := composite.Evaluate(context.Background(), user)

		assert.False(t, result.Allowed)
		assert.Equal(t, "INVALID_OPERATOR", result.Code)
	})

	t.Run("Add method adds policy to composite", func(t *testing.T) {
		composite := policy.NewCompositePolicy[User]("dynamic", policy.AND)
		composite.Add(agePolicy)
		composite.Add(activePolicy)

		user := User{Age: 25, IsActive: true}
		result := composite.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
	})
}

func TestPolicyEngine(t *testing.T) {
	engine := policy.NewPolicyEngine[User]()
	agePolicy := newAgePolicy(18)
	activePolicy := newActiveUserPolicy()

	t.Run("Register and evaluate specific policy", func(t *testing.T) {
		engine.Register(agePolicy)

		user := User{Age: 25}
		result, err := engine.Evaluate(context.Background(), "age-policy", user)

		require.NoError(t, err)
		assert.True(t, result.Allowed)
	})

	t.Run("Evaluate non-existent policy returns error", func(t *testing.T) {
		user := User{}
		result, err := engine.Evaluate(context.Background(), "non-existent", user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "policy 'non-existent' not found")
		assert.Equal(t, policy.Result{}, result)
	})

	t.Run("EvaluateAll with multiple policies", func(t *testing.T) {
		engine.Register(agePolicy)
		engine.Register(activePolicy)

		user := User{Age: 25, IsActive: true}
		result := engine.EvaluateAll(context.Background(), user)

		assert.True(t, result.Allowed)
	})

	t.Run("EvaluateAll with failing policy", func(t *testing.T) {
		user := User{Age: 16, IsActive: false}
		result := engine.EvaluateAll(context.Background(), user)

		assert.False(t, result.Allowed)
	})

	t.Run("EvaluateAll with no policies", func(t *testing.T) {
		emptyEngine := policy.NewPolicyEngine[User]()
		user := User{}
		result := emptyEngine.EvaluateAll(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Equal(t, "No policies registered", result.Reason)
	})

	t.Run("GetPolicyNames returns all registered policy names", func(t *testing.T) {
		names := engine.GetPolicyNames()

		assert.Len(t, names, 2)
		assert.Contains(t, names, "age-policy")
		assert.Contains(t, names, "active-user-policy")
	})
}

func TestConditionalPolicy(t *testing.T) {
	basePolicy := newAgePolicy(18)

	t.Run("Evaluates wrapped policy when condition is true", func(t *testing.T) {
		condition := func(ctx context.Context, user User) bool {
			return user.IsActive
		}

		conditional := policy.NewConditionalPolicy("vip-age-check", condition, basePolicy)
		user := User{Age: 25, IsActive: true}

		result := conditional.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Contains(t, result.Reason, "User meets age requirement")
	})

	t.Run("Skips wrapped policy when condition is false", func(t *testing.T) {
		condition := func(ctx context.Context, user User) bool {
			return user.IsActive
		}

		conditional := policy.NewConditionalPolicy("vip-age-check", condition, basePolicy)
		user := User{Age: 16, IsActive: false} // Would fail age check, but condition is false

		result := conditional.Evaluate(context.Background(), user)

		assert.True(t, result.Allowed)
		assert.Equal(t, "Condition not met, policy skipped", result.Reason)
	})

	t.Run("Name returns correct policy name", func(t *testing.T) {
		condition := func(ctx context.Context, user User) bool { return true }
		conditional := policy.NewConditionalPolicy("test-conditional", condition, basePolicy)

		assert.Equal(t, "test-conditional", conditional.Name())
	})
}

// MockPolicy for testing purposes
type MockPolicy struct {
	name      string
	result    policy.Result
	callCount int
}

func NewMockPolicy(name string, result policy.Result) *MockPolicy {
	return &MockPolicy{
		name:   name,
		result: result,
	}
}

func (mp *MockPolicy) Name() string {
	return mp.name
}

func (mp *MockPolicy) Evaluate(ctx context.Context, subject User) policy.Result {
	mp.callCount++
	return mp.result
}

func (mp *MockPolicy) CallCount() int {
	return mp.callCount
}

func TestMockPolicy(t *testing.T) {
	t.Run("Mock policy behavior", func(t *testing.T) {
		mockResult := policy.Allow("mock success")
		mock := NewMockPolicy("mock-policy", mockResult)

		result := mock.Evaluate(context.Background(), User{})

		assert.Equal(t, "mock-policy", mock.Name())
		assert.Equal(t, mockResult, result)
		assert.Equal(t, 1, mock.CallCount())
	})
}

func TestPolicyEngineIntegration(t *testing.T) {
	t.Run("Complex policy evaluation scenario", func(t *testing.T) {
		engine := policy.NewPolicyEngine[User]()

		// Register multiple policies
		engine.Register(newAgePolicy(21))
		engine.Register(newActiveUserPolicy())

		// Test user that should pass all policies
		validUser := User{
			Age:      25,
			IsActive: true,
		}

		result := engine.EvaluateAll(context.Background(), validUser)
		assert.True(t, result.Allowed)

		// Test user that should fail age policy
		youngUser := User{
			Age:      20,
			IsActive: true,
		}

		result = engine.EvaluateAll(context.Background(), youngUser)
		assert.False(t, result.Allowed)
		assert.Contains(t, result.Reason, "age-policy")
	})
}

// Benchmark tests
func BenchmarkPolicyEvaluation(b *testing.B) {
	policy := newAgePolicy(18)
	user := User{Age: 25}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		policy.Evaluate(ctx, user)
	}
}

func BenchmarkCompositePolicyEvaluation(b *testing.B) {
	composite := policy.NewCompositePolicy(
		"benchmark",
		policy.AND,
		newAgePolicy(18),
		newActiveUserPolicy(),
	)
	user := User{Age: 25, IsActive: true}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composite.Evaluate(ctx, user)
	}
}

func BenchmarkPolicyEngine(b *testing.B) {
	engine := policy.NewPolicyEngine[User]()
	engine.Register(newAgePolicy(18))
	engine.Register(newActiveUserPolicy())

	user := User{Age: 25, IsActive: true}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.EvaluateAll(ctx, user)
	}
}
