package valueobject

import (
	"regexp"
	"strings"
)

// Email value object
type Email struct {
	BaseValueObject
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object
func NewEmail(email string) (*Email, error) {
	e := &Email{value: strings.TrimSpace(strings.ToLower(email))}
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return e, nil
}

// MustNewEmail creates a new Email value object or panics
func MustNewEmail(email string) *Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

func (e *Email) Value() string {
	return e.value
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) Validate() error {
	if e.value == "" {
		return ValidationError{Field: "email", Message: "cannot be empty"}
	}
	if !emailRegex.MatchString(e.value) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}
	return nil
}

func (e *Email) Equals(other ValueObject) bool {
	if otherEmail, ok := other.(*Email); ok {
		return e.value == otherEmail.value
	}
	return false
}
