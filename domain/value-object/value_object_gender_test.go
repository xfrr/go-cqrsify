package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestGenderSuite(t *testing.T) {
	suite.Run(t, new(GenderTestSuite))
}

// GenderTestSuite groups phone number-related tests
type GenderTestSuite struct {
	suite.Suite
}

func (suite *GenderTestSuite) TestValidGender() {
	genderVO, err := valueobject.NewGender(valueobject.MaleGenderType)

	suite.Require().NoError(err)
	suite.Equal(valueobject.MaleGenderType, genderVO.Value())
	suite.Equal(valueobject.MaleGenderType.String(), genderVO.String())
}

func (suite *GenderTestSuite) TestInvalidGender() {
	_, err := valueobject.NewGender("invalid")
	suite.Require().Error(err)

	expectedErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{
			{
				Field:   "value",
				Message: "invalid gender type",
			},
		},
	}

	suite.ErrorAs(err, &expectedErr)
}

func (suite *GenderTestSuite) TestGenderEquality() {
	gender1, _ := valueobject.NewGender(valueobject.MaleGenderType)
	gender2, _ := valueobject.NewGender(valueobject.MaleGenderType)
	gender3, _ := valueobject.NewGender(valueobject.FemaleGenderType)

	suite.True(gender1.Equals(gender2))
	suite.False(gender1.Equals(gender3))
	suite.False(gender1.Equals(nil))
}

func (suite *GenderTestSuite) TestParseGender() {
	genderVO, err := valueobject.ParseGender(string(valueobject.FemaleGenderType))
	suite.Require().NoError(err)
	suite.Equal(valueobject.FemaleGenderType, genderVO.Value())

	_, err = valueobject.ParseGender("unknown")
	suite.Require().Error(err)
}
