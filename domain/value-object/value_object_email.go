package valueobject

import (
	"regexp"
	"strings"
)

var _ ValueObject = (*Email)(nil)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email value object
type Email struct {
	BaseValueObject
	value string
}

// NewEmail creates a new Email value object
func NewEmail(email string) (Email, error) {
	e := Email{value: strings.TrimSpace(strings.ToLower(email))}
	if err := e.Validate(); err != nil {
		return Email{}, err
	}
	return e, nil
}

// MustNewEmail creates a new Email value object or panics
func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}

func (e Email) Validate() error {
	var errs []ValidationError
	switch {
	case e.value == "":
		errs = append(errs, ValidationError{Field: "value", Message: "email cannot be empty"})
	case !emailRegex.MatchString(e.value):
		errs = append(errs, ValidationError{Field: "value", Message: "invalid email format"})
	}
	return ValidationErrors(errs)
}

func (e Email) Equals(other ValueObject) bool {
	if otherEmail, ok := other.(Email); ok {
		return e.value == otherEmail.value
	}
	return false
}
