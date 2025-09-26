package valueobject

import (
	"slices"
)

var _ ValueObject = (*Sex)(nil)

// SexType defines the type for biological sex.
type SexType string

func (s SexType) String() string {
	return string(s)
}

const (
	// MaleSexType represents organisms that have XY chromosomes, testes, and sperm production.
	MaleSexType SexType = "male"
	// FemaleSexType represents organisms that have XX chromosomes, ovaries, and egg production.
	FemaleSexType SexType = "female"
	// IntersexSexType represents a range of natural variations in chromosomes, gonads, hormones, or genitalia.
	// E.g. Klinefelter syndrome (XXY), Turner syndrome (XO), Androgen Insensitivity Syndrome.
	IntersexSexType SexType = "intersex"
	// HermaphroditismSexType represents organisms that have both male and female reproductive organs.
	// E.g. earthworms, snails, many plants.
	HermaphroditismSexType SexType = "hermaphroditism"
	// SequentialHermaphroditismSexType represents organisms that change sex during life.
	// E.g. Clownfish (male → female), wrasses, parrotfish.
	SequentialHermaphroditismSexType SexType = "sequential_hermaphroditism"
	// ParthenogenesisSexType represents asexual reproduction where an organism
	// can reproduce without fertilization.
	// E.g. some reptiles, amphibians, and fish.
	ParthenogenesisSexType SexType = "parthenogenesis"
	// AsexualSexType represents organisms that do not have a biological sex.
	// E.g. bacteria, archaea, and many protists.
	AsexualSexType SexType = "asexual"
	// ComplexSexType represents organisms with complex or not well-defined sexual characteristics.
	// E.g. some animals (like bees, ants, wasps) have haplodiploid sex determination,
	// where sex depends on the number of chromosome sets.
	ComplexSexType SexType = "complex"
)

var AllSexTypes = []SexType{
	MaleSexType,
	FemaleSexType,
	IntersexSexType,
	HermaphroditismSexType,
	SequentialHermaphroditismSexType,
	ParthenogenesisSexType,
	AsexualSexType,
	ComplexSexType,
}

// Sex represents a value object for biological sex,
// by which most organisms are classified on the basis
// of their reproductive organs and functions.
type Sex struct {
	BaseValueObject
	value SexType
}

func (s Sex) Value() SexType {
	return s.value
}

func (s Sex) String() string {
	return string(s.value)
}

func (s Sex) Equals(other ValueObject) bool {
	if otherSex, ok := other.(Sex); ok {
		return s.value == otherSex.value
	}
	return false
}

func (s Sex) Validate() error {
	var errs []ValidationError
	if !slices.Contains(AllSexTypes, s.value) {
		errs = append(errs, ValidationError{
			Field:   "value",
			Message: "invalid sex type",
		})
	}
	return ValidationErrors(errs)
}

// NewSex creates a new Sex value object.
func NewSex(value SexType) (Sex, error) {
	s := Sex{
		BaseValueObject: BaseValueObject{},
		value:           value,
	}
	if err := s.Validate(); err != nil {
		return Sex{}, err
	}
	return s, nil
}
