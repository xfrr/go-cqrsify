package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestPhoneNumberSuite(t *testing.T) {
	suite.Run(t, new(PhoneNumberTestSuite))
}

// PhoneNumberTestSuite groups phone number-related tests
type PhoneNumberTestSuite struct {
	suite.Suite
}

func (suite *PhoneNumberTestSuite) TestValidPhoneNumber() {
	phone, err := valueobject.NewPhoneNumber("+1", "1234567890")

	suite.Require().NoError(err)
	suite.Equal("+1", phone.CountryCode())
	suite.Equal("1234567890", phone.Number())
	suite.Equal("+11234567890", phone.FullNumber())
}

func (suite *PhoneNumberTestSuite) TestPhoneNumberWithSpaces() {
	phone, err := valueobject.NewPhoneNumber("+1", "123 456 7890")

	suite.Require().NoError(err)
	suite.Equal("1234567890", phone.Number())
}

func (suite *PhoneNumberTestSuite) TestInvalidPhoneNumber() {
	_, err := valueobject.NewPhoneNumber("+1", "invalid")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "invalid phone number format")
}

func (suite *PhoneNumberTestSuite) TestEmptyCountryCode() {
	_, err := valueobject.NewPhoneNumber("", "1234567890")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "countryCode")
}

func (suite *PhoneNumberTestSuite) TestEmptyNumber() {
	_, err := valueobject.NewPhoneNumber("+1", "")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "number")
}

func (suite *PhoneNumberTestSuite) TestPhoneNumberEquality() {
	phone1, _ := valueobject.NewPhoneNumber("+1", "1234567890")
	phone2, _ := valueobject.NewPhoneNumber("+1", "1234567890")
	phone3, _ := valueobject.NewPhoneNumber("+1", "0987654321")

	suite.True(phone1.Equals(phone2))
	suite.False(phone1.Equals(phone3))
	suite.False(phone1.Equals(nil))
}
