package valueobject

import (
	"fmt"
	"strings"
)

// Address value object
type Address struct {
	BaseValueObject
	street  string
	city    string
	state   string
	zipCode string
	country string
}

// NewAddress creates a new Address value object
func NewAddress(street, city, state, zipCode, country string) (*Address, error) {
	addr := &Address{
		street:  strings.TrimSpace(street),
		city:    strings.TrimSpace(city),
		state:   strings.TrimSpace(state),
		zipCode: strings.TrimSpace(zipCode),
		country: strings.TrimSpace(country),
	}
	if err := addr.Validate(); err != nil {
		return nil, err
	}
	return addr, nil
}

func (a *Address) Street() string  { return a.street }
func (a *Address) City() string    { return a.city }
func (a *Address) State() string   { return a.state }
func (a *Address) ZipCode() string { return a.zipCode }
func (a *Address) Country() string { return a.country }

func (a *Address) String() string {
	return fmt.Sprintf("%s, %s, %s %s, %s", a.street, a.city, a.state, a.zipCode, a.country)
}

func (a *Address) Validate() error {
	var errs []ValidationError

	if a.street == "" {
		errs = append(errs, ValidationError{Field: "street", Message: "cannot be empty"})
	}
	if a.city == "" {
		errs = append(errs, ValidationError{Field: "city", Message: "cannot be empty"})
	}
	if a.state == "" {
		errs = append(errs, ValidationError{Field: "state", Message: "cannot be empty"})
	}
	if a.zipCode == "" {
		errs = append(errs, ValidationError{Field: "zipCode", Message: "cannot be empty"})
	}
	if a.country == "" {
		errs = append(errs, ValidationError{Field: "country", Message: "cannot be empty"})
	}

	if len(errs) > 0 {
		return MultiValidationError{Errors: errs}
	}
	return nil
}

func (a *Address) Equals(other ValueObject) bool {
	if otherAddr, ok := other.(*Address); ok {
		return a.street == otherAddr.street &&
			a.city == otherAddr.city &&
			a.state == otherAddr.state &&
			a.zipCode == otherAddr.zipCode &&
			a.country == otherAddr.country
	}
	return false
}
