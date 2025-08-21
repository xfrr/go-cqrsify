package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestPersonNameSuite(t *testing.T) {
	suite.Run(t, new(PersonNameTestSuite))
}

// PersonNameTestSuite groups person name-related tests
type PersonNameTestSuite struct {
	suite.Suite
}

func (suite *PersonNameTestSuite) TestValidPersonName() {
	name, err := valueobject.NewPersonName("John", "Doe")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John", name.FirstName())
	assert.Equal(suite.T(), "Doe", name.LastName())
	assert.Equal(suite.T(), "John Doe", name.FullName())
	assert.Equal(suite.T(), "John Doe", name.String())
}

func (suite *PersonNameTestSuite) TestEmptyFirstName() {
	_, err := valueobject.NewPersonName("", "Doe")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "firstName")
}

func (suite *PersonNameTestSuite) TestEmptyLastName() {
	_, err := valueobject.NewPersonName("John", "")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "lastName")
}

func (suite *PersonNameTestSuite) TestPersonNameEquality() {
	name1, _ := valueobject.NewPersonName("John", "Doe")
	name2, _ := valueobject.NewPersonName("John", "Doe")
	name3, _ := valueobject.NewPersonName("Jane", "Doe")

	assert.True(suite.T(), name1.Equals(name2))
	assert.False(suite.T(), name1.Equals(name3))
	assert.False(suite.T(), name1.Equals(nil))
}

func (suite *PersonNameTestSuite) TestNameTrimming() {
	name, err := valueobject.NewPersonName("  John  ", "  Doe  ")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John", name.FirstName())
	assert.Equal(suite.T(), "Doe", name.LastName())
}
