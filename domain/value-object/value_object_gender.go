package valueobject

import "slices"

// GenderType defines the type for human gender.
type GenderType string

func (s GenderType) String() string {
	return string(s)
}

const (
	// MaleGenderType represents individuals who identify as male.
	MaleGenderType GenderType = "male"
	// FemaleGenderType represents individuals who identify as female.
	FemaleGenderType GenderType = "female"
	// OtherGenderType represents individuals who identify outside the traditional gender binary.
	OtherGenderType GenderType = "other"
)

var AllGenderTypes = []GenderType{
	MaleGenderType,
	FemaleGenderType,
	OtherGenderType,
}

type Gender struct {
	BaseValueObject
	value GenderType
}

func (g Gender) Equals(other ValueObject) bool {
	if otherGender, ok := other.(Gender); ok {
		return g.value == otherGender.value
	}
	return false
}

func (g Gender) Value() GenderType {
	return g.value
}

func (g Gender) String() string {
	return string(g.value)
}

func (g Gender) Validate() error {
	var errs []ValidationError
	if !slices.Contains(AllGenderTypes, g.value) {
		errs = append(errs, ValidationError{
			Field:   "value",
			Message: "invalid gender type",
		})
	}
	return ValidationErrors(errs)
}

// NewGender creates a new Gender value object.
func NewGender(value GenderType) (Gender, error) {
	g := Gender{
		BaseValueObject: BaseValueObject{},
		value:           value,
	}
	if err := g.Validate(); err != nil {
		return Gender{}, err
	}
	return g, nil
}

// ParseGender parses a string into a Gender value object.
func ParseGender(value string) (Gender, error) {
	switch GenderType(value) {
	case MaleGenderType:
		return NewGender(MaleGenderType)
	case FemaleGenderType:
		return NewGender(FemaleGenderType)
	case OtherGenderType:
		return NewGender(OtherGenderType)
	default:
		return Gender{}, ValidationError{
			Field:   "value",
			Message: "invalid gender type",
		}
	}
}
