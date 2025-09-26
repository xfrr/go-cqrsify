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
