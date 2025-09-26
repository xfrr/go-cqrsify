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
