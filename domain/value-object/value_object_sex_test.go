package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestSexSuite(t *testing.T) {
	suite.Run(t, new(SexTestSuite))
}

// SexTestSuite groups phone number-related tests
type SexTestSuite struct {
	suite.Suite
}

func (suite *SexTestSuite) TestValidSex() {
	sexVO, err := valueobject.NewSex(valueobject.MaleSexType)

	suite.Require().NoError(err)
	suite.Equal(valueobject.MaleSexType, sexVO.Value())
	suite.Equal(valueobject.MaleSexType.String(), sexVO.String())
}

func (suite *SexTestSuite) TestInvalidSex() {
	_, err := valueobject.NewSex("invalid")
	suite.Require().Error(err)

	expectedErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{
			{
				Field:   "value",
				Message: "invalid sex type",
			},
		},
	}

	suite.ErrorAs(err, &expectedErr)
}

func (suite *SexTestSuite) TestSexEquality() {
	sex1, _ := valueobject.NewSex(valueobject.MaleSexType)
	sex2, _ := valueobject.NewSex(valueobject.MaleSexType)
	sex3, _ := valueobject.NewSex(valueobject.FemaleSexType)

	suite.True(sex1.Equals(sex2))
	suite.False(sex1.Equals(sex3))
	suite.False(sex1.Equals(nil))
}

func (suite *SexTestSuite) TestSexIsValid() {
	sex1, _ := valueobject.NewSex(valueobject.MaleSexType)
	sex2, _ := valueobject.NewSex("invalid")

	suite.True(sex1.IsValid())
	suite.False(sex2.IsValid())
}

func (suite *SexTestSuite) TestSexIsZero() {
	sex1, _ := valueobject.NewSex(valueobject.MaleSexType)
	sex2 := valueobject.Sex{}

	suite.False(sex1.IsZero())
	suite.True(sex2.IsZero())
}

func (suite *SexTestSuite) TestSexIs() {
	sex1, _ := valueobject.NewSex(valueobject.MaleSexType)
	sex2, _ := valueobject.NewSex(valueobject.FemaleSexType)

	suite.True(sex1.Is(valueobject.MaleSexType))
	suite.False(sex1.Is(valueobject.FemaleSexType))
	suite.True(sex2.Is(valueobject.FemaleSexType))
	suite.False(sex2.Is(valueobject.MaleSexType))
}

func (suite *SexTestSuite) TestParseSex() {
	testCases := []struct {
		input       string
		expectedSex valueobject.SexType
		expectedErr error
	}{
		{"male", valueobject.MaleSexType, nil},
		{"female", valueobject.FemaleSexType, nil},
		{"intersex", valueobject.IntersexSexType, nil},
		{"hermaphroditism", valueobject.HermaphroditismSexType, nil},
		{"sequential_hermaphroditism", valueobject.SequentialHermaphroditismSexType, nil},
		{"parthenogenesis", valueobject.ParthenogenesisSexType, nil},
		{"asexual", valueobject.AsexualSexType, nil},
		{"complex", valueobject.ComplexSexType, nil},
		{"unknown", "", valueobject.ValidationError{
			Field:   "value",
			Message: "invalid sex type",
		}},
	}

	for _, tc := range testCases {
		suite.Run(tc.input, func() {
			sexVO, err := valueobject.ParseSex(tc.input)
			if tc.expectedErr != nil {
				suite.Require().Error(err)
				suite.Equal(tc.expectedErr, err)
			} else {
				suite.Require().NoError(err)
				suite.Equal(tc.expectedSex, sexVO.Value())
			}
		})
	}
}
