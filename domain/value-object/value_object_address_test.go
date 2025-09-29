package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

// AddressTestSuite groups address-related tests
type AddressTestSuite struct {
	suite.Suite
}

func (suite *AddressTestSuite) TestValidAddress() {
	addr, err := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")

	suite.Require().NoError(err)
	suite.Equal("123 Main St", addr.Street())
	suite.Equal("Anytown", addr.City())
	suite.Equal("CA", addr.State())
	suite.Equal("12345", addr.ZipCode())
	suite.Equal("USA", addr.Country())
}

func (suite *AddressTestSuite) TestAddressString() {
	addr, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	expected := "123 Main St, Anytown, CA 12345, USA"

	suite.Equal(expected, addr.String())
}

func (suite *AddressTestSuite) TestMissingRequiredFields() {
	testCases := []struct {
		name        string
		street      string
		city        string
		state       string
		zipCode     string
		country     string
		expectedErr string
	}{
		{"missing street", "", "Anytown", "CA", "12345", "USA", "street"},
		{"missing city", "123 Main St", "", "CA", "12345", "USA", "city"},
		{"missing state", "123 Main St", "Anytown", "", "12345", "USA", "state"},
		{"missing zipCode", "123 Main St", "Anytown", "CA", "", "USA", "zipCode"},
		{"missing country", "123 Main St", "Anytown", "CA", "12345", "", "country"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := valueobject.NewAddress(tc.street, tc.city, tc.state, tc.zipCode, tc.country)
			suite.Require().Error(err)
			suite.Contains(err.Error(), tc.expectedErr)
		})
	}
}

func (suite *AddressTestSuite) TestAddressEquality() {
	addr1, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	addr2, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	addr3, _ := valueobject.NewAddress("456 Oak Ave", "Anytown", "CA", "12345", "USA")

	suite.True(addr1.Equals(addr2))
	suite.False(addr1.Equals(addr3))
	suite.False(addr1.Equals(nil))
}

func (suite *AddressTestSuite) TestAddressTrimming() {
	addr, err := valueobject.NewAddress("  123 Main St  ", "  Anytown  ", "  CA  ", "  12345  ", "  USA  ")

	suite.Require().NoError(err)
	suite.Equal("123 Main St", addr.Street())
	suite.Equal("Anytown", addr.City())
}

func (suite *AddressTestSuite) TestIsEmptyStreetError() {
	err := valueobject.ErrEmptyStreet
	suite.Require().True(valueobject.IsEmptyStreetError(err), "Expected IsEmptyStreetError to return true for ErrEmptyStreet")

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrEmptyStreet, valueobject.ErrEmptyCity},
	}
	suite.Require().True(valueobject.IsEmptyStreetError(multiErr), "Expected IsEmptyStreetError to return true for MultiValidationError containing ErrEmptyStreet")

	otherErr := valueobject.ErrEmptyCity
	suite.Require().False(valueobject.IsEmptyStreetError(otherErr), "Expected IsEmptyStreetError to return false for ErrEmptyCity")
}

func (suite *AddressTestSuite) TestIsEmptyCityError() {
	err := valueobject.ErrEmptyCity
	suite.Require().True(valueobject.IsEmptyCityError(err), "Expected IsEmptyCityError to return true for ErrEmptyCity")

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrEmptyStreet, valueobject.ErrEmptyCity},
	}
	suite.Require().True(valueobject.IsEmptyCityError(multiErr), "Expected IsEmptyCityError to return true for MultiValidationError containing ErrEmptyCity")

	otherErr := valueobject.ErrEmptyStreet
	suite.Require().False(valueobject.IsEmptyCityError(otherErr), "Expected IsEmptyCityError to return false for ErrEmptyStreet")
}

func (suite *AddressTestSuite) TestIsEmptyStateError() {
	err := valueobject.ErrEmptyState
	suite.Require().True(valueobject.IsEmptyStateError(err), "Expected IsEmptyStateError to return true for ErrEmptyState")

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrEmptyState, valueobject.ErrEmptyCity},
	}
	suite.Require().True(valueobject.IsEmptyStateError(multiErr), "Expected IsEmptyStateError to return true for MultiValidationError containing ErrEmptyState")

	otherErr := valueobject.ErrEmptyCity
	suite.Require().False(valueobject.IsEmptyStateError(otherErr), "Expected IsEmptyStateError to return false for ErrEmptyCity")
}

func (suite *AddressTestSuite) TestIsEmptyZipCodeError() {
	err := valueobject.ErrEmptyZipCode
	suite.Require().True(valueobject.IsEmptyZipCodeError(err), "Expected IsEmptyZipCodeError to return true for ErrEmptyZipCode")

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrEmptyZipCode, valueobject.ErrEmptyCity},
	}
	suite.Require().True(valueobject.IsEmptyZipCodeError(multiErr), "Expected IsEmptyZipCodeError to return true for MultiValidationError containing ErrEmptyZipCode")

	otherErr := valueobject.ErrEmptyCity
	suite.Require().False(valueobject.IsEmptyZipCodeError(otherErr), "Expected IsEmptyZipCodeError to return false for ErrEmptyCity")
}

func (suite *AddressTestSuite) TestIsEmptyCountryError() {
	err := valueobject.ErrEmptyCountry
	suite.Require().True(valueobject.IsEmptyCountryError(err), "Expected IsEmptyCountryError to return true for ErrEmptyCountry")

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrEmptyCountry, valueobject.ErrEmptyCity},
	}
	suite.Require().True(valueobject.IsEmptyCountryError(multiErr), "Expected IsEmptyCountryError to return true for MultiValidationError containing ErrEmptyCountry")

	otherErr := valueobject.ErrEmptyCity
	suite.Require().False(valueobject.IsEmptyCountryError(otherErr), "Expected IsEmptyCountryError to return false for ErrEmptyCity")
}
