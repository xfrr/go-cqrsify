package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "123 Main St", addr.Street())
	assert.Equal(suite.T(), "Anytown", addr.City())
	assert.Equal(suite.T(), "CA", addr.State())
	assert.Equal(suite.T(), "12345", addr.ZipCode())
	assert.Equal(suite.T(), "USA", addr.Country())
}

func (suite *AddressTestSuite) TestAddressString() {
	addr, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	expected := "123 Main St, Anytown, CA 12345, USA"

	assert.Equal(suite.T(), expected, addr.String())
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
			assert.Error(suite.T(), err)
			assert.Contains(suite.T(), err.Error(), tc.expectedErr)
		})
	}
}

func (suite *AddressTestSuite) TestAddressEquality() {
	addr1, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	addr2, _ := valueobject.NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	addr3, _ := valueobject.NewAddress("456 Oak Ave", "Anytown", "CA", "12345", "USA")

	assert.True(suite.T(), addr1.Equals(addr2))
	assert.False(suite.T(), addr1.Equals(addr3))
	assert.False(suite.T(), addr1.Equals(nil))
}

func (suite *AddressTestSuite) TestAddressTrimming() {
	addr, err := valueobject.NewAddress("  123 Main St  ", "  Anytown  ", "  CA  ", "  12345  ", "  USA  ")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "123 Main St", addr.Street())
	assert.Equal(suite.T(), "Anytown", addr.City())
}
