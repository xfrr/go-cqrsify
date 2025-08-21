package valueobject

import (
	"fmt"
	"regexp"
	"strings"
)

// PhoneNumber value object
type PhoneNumber struct {
	BaseValueObject
	countryCode string
	number      string
}

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// NewPhoneNumber creates a new PhoneNumber value object
func NewPhoneNumber(countryCode, number string) (*PhoneNumber, error) {
	pn := &PhoneNumber{
		countryCode: strings.TrimSpace(countryCode),
		number:      strings.ReplaceAll(strings.TrimSpace(number), " ", ""),
	}
	if err := pn.Validate(); err != nil {
		return nil, err
	}
	return pn, nil
}

func (pn *PhoneNumber) CountryCode() string {
	return pn.countryCode
}

func (pn *PhoneNumber) Number() string {
	return pn.number
}

func (pn *PhoneNumber) FullNumber() string {
	return fmt.Sprintf("%s%s", pn.countryCode, pn.number)
}

func (pn *PhoneNumber) String() string {
	return pn.FullNumber()
}

func (pn *PhoneNumber) Validate() error {
	var errs []ValidationError

	if pn.countryCode == "" {
		errs = append(errs, ValidationError{Field: "countryCode", Message: "cannot be empty"})
	}

	if pn.number == "" {
		errs = append(errs, ValidationError{Field: "number", Message: "cannot be empty"})
	}

	fullNumber := pn.FullNumber()
	if !phoneRegex.MatchString(fullNumber) {
		errs = append(errs, ValidationError{Field: "phoneNumber", Message: "invalid phone number format"})
	}

	if len(errs) > 0 {
		return MultiValidationError{Errors: errs}
	}
	return nil
}

func (pn *PhoneNumber) Equals(other ValueObject) bool {
	if otherPhone, ok := other.(*PhoneNumber); ok {
		return pn.countryCode == otherPhone.countryCode && pn.number == otherPhone.number
	}
	return false
}
