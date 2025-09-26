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

func (suite *GenderTestSuite) TestGenderIsValid() {
	gender1, _ := valueobject.NewGender(valueobject.MaleGenderType)
	gender2, _ := valueobject.NewGender("invalid")

	suite.True(gender1.IsValid())
	suite.False(gender2.IsValid())
}

func (suite *GenderTestSuite) TestGenderIsEmpty() {
	gender1, _ := valueobject.NewGender(valueobject.MaleGenderType)
	gender2, _ := valueobject.NewGender("")

	suite.False(gender1.IsEmpty())
	suite.True(gender2.IsEmpty())
}

func (suite *GenderTestSuite) TestGenderIs() {
	gender1, _ := valueobject.NewGender(valueobject.MaleGenderType)
	gender2, _ := valueobject.NewGender(valueobject.FemaleGenderType)

	suite.True(gender1.Is(valueobject.MaleGenderType))
	suite.False(gender1.Is(valueobject.FemaleGenderType))
	suite.True(gender2.Is(valueobject.FemaleGenderType))
	suite.False(gender2.Is(valueobject.MaleGenderType))
}

func (suite *GenderTestSuite) TestParseGender() {
	testCases := []struct {
		input    string
		expected valueobject.GenderType
		isValid  bool
	}{
		{"male", valueobject.MaleGenderType, true},
		{"female", valueobject.FemaleGenderType, true},
		{"other", valueobject.OtherGenderType, true},
		{"unknown", "", false},
	}

	for _, tc := range testCases {
		suite.Run(tc.input, func() {
			genderVO, err := valueobject.ParseGender(tc.input)
			if tc.isValid {
				suite.Require().NoError(err)
				suite.Equal(tc.expected, genderVO.Value())
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
