package policy

// Result represents the outcome of a policy evaluation
type Result struct {
	Allowed bool
	Reason  string
	Code    string
}

// NewResult creates a new policy result
func NewResult(allowed bool, reason, code string) Result {
	return Result{
		Allowed: allowed,
		Reason:  reason,
		Code:    code,
	}
}

// Allow creates a successful policy result
func Allow(reason string) Result {
	return NewResult(true, reason, "ALLOWED")
}

// Deny creates a failed policy result with reason
func Deny(reason, code string) Result {
	return NewResult(false, reason, code)
}
