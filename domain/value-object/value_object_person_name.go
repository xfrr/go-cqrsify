package valueobject

import (
	"fmt"
	"strings"
)

var _ ValueObject = (*PersonName)(nil)

// PersonName value object
type PersonName struct {
	BaseValueObject
	firstName string
	lastName  string
}

// NewPersonName creates a new PersonName value object
func NewPersonName(firstName, lastName string) (PersonName, error) {
	pn := PersonName{
		firstName: strings.TrimSpace(firstName),
		lastName:  strings.TrimSpace(lastName),
	}
	if err := pn.Validate(); err != nil {
		return PersonName{}, err
	}
	return pn, nil
}

func (pn PersonName) FirstName() string {
	return pn.firstName
}

func (pn PersonName) LastName() string {
	return pn.lastName
}

func (pn PersonName) FullName() string {
	return fmt.Sprintf("%s %s", pn.firstName, pn.lastName)
}

func (pn PersonName) String() string {
	return pn.FullName()
}

func (pn PersonName) Validate() error {
	var errs []ValidationError

	if pn.firstName == "" {
		errs = append(errs, ValidationError{Field: "firstName", Message: "cannot be empty"})
	}

	if pn.lastName == "" {
		errs = append(errs, ValidationError{Field: "lastName", Message: "cannot be empty"})
	}

	if len(errs) > 0 {
		return MultiValidationError{Errors: errs}
	}
	return nil
}

func (pn PersonName) Equals(other ValueObject) bool {
	if otherName, ok := other.(PersonName); ok {
		return pn.firstName == otherName.firstName &&
			pn.lastName == otherName.lastName
	}
	return false
}
