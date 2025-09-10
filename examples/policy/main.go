package main

import (
	"context"
	"fmt"

	policy "github.com/xfrr/go-cqrsify/domain/policy"
)

var _ policy.Policy[User] = (*AgeRestrictionPolicy)(nil)

type AgeRestrictionPolicy struct {
	policy.BasePolicy
	MinAge int
	MaxAge int
}

func newUserAgePolicy(minAge int, maxAge int) AgeRestrictionPolicy {
	return AgeRestrictionPolicy{
		BasePolicy: policy.NewBasePolicy("user-age-policy"),
		MinAge:     minAge,
		MaxAge:     maxAge,
	}
}

func (p AgeRestrictionPolicy) Evaluate(ctx context.Context, user User) policy.Result {
	if user.Age >= p.MinAge && user.Age <= p.MaxAge {
		return policy.Allow(fmt.Sprintf("User age %d meets age requirements of %d - %d", user.Age, p.MinAge, p.MaxAge))
	}
	return policy.Deny(
		fmt.Sprintf("User age %d does not meet age requirements of %d - %d", user.Age, p.MinAge, p.MaxAge),
		"AGE_RESTRICTED",
	)
}

type User struct {
	Age int
}

func main() {
	// Create a new policy engine (optional)
	policyEngine := policy.NewPolicyEngine[User]()

	// Create a new user age policy
	agePolicy := newUserAgePolicy(18, 65)

	// 	Register the policy with the engine
	policyEngine.Register(agePolicy)

	// Evaluate the policy using the policy engine
	result := policyEngine.EvaluateAll(context.Background(), User{Age: 20})
	fmt.Println("Policy Result:")
	fmt.Println("- Allowed:", result.Allowed)
	fmt.Println("- Code:", result.Code)
	fmt.Println("- Reason:", result.Reason)
	fmt.Println()

	// Now let's evaluate a user with age 17
	result = policyEngine.EvaluateAll(context.Background(), User{Age: 17})
	fmt.Println("Policy Result:")
	fmt.Println("- Allowed:", result.Allowed)
	fmt.Println("- Code:", result.Code)
	fmt.Println("- Reason:", result.Reason)

	// You can add more policies and evaluate them similarly
	// Hope this helps!
}
