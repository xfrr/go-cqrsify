package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "+1", phone.CountryCode())
	assert.Equal(suite.T(), "1234567890", phone.Number())
	assert.Equal(suite.T(), "+11234567890", phone.FullNumber())
}

func (suite *PhoneNumberTestSuite) TestPhoneNumberWithSpaces() {
	phone, err := valueobject.NewPhoneNumber("+1", "123 456 7890")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "1234567890", phone.Number())
}

func (suite *PhoneNumberTestSuite) TestInvalidPhoneNumber() {
	_, err := valueobject.NewPhoneNumber("+1", "invalid")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid phone number format")
}

func (suite *PhoneNumberTestSuite) TestEmptyCountryCode() {
	_, err := valueobject.NewPhoneNumber("", "1234567890")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "countryCode")
}

func (suite *PhoneNumberTestSuite) TestEmptyNumber() {
	_, err := valueobject.NewPhoneNumber("+1", "")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "number")
}

func (suite *PhoneNumberTestSuite) TestPhoneNumberEquality() {
	phone1, _ := valueobject.NewPhoneNumber("+1", "1234567890")
	phone2, _ := valueobject.NewPhoneNumber("+1", "1234567890")
	phone3, _ := valueobject.NewPhoneNumber("+1", "0987654321")

	assert.True(suite.T(), phone1.Equals(phone2))
	assert.False(suite.T(), phone1.Equals(phone3))
	assert.False(suite.T(), phone1.Equals(nil))
}
