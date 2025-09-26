package valueobject_test

import (
	"testing"

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

	suite.Require().NoError(err)
	suite.Equal("John", name.FirstName())
	suite.Equal("Doe", name.LastName())
	suite.Equal("John Doe", name.FullName())
	suite.Equal("John Doe", name.String())
}

func (suite *PersonNameTestSuite) TestEmptyFirstName() {
	_, err := valueobject.NewPersonName("", "Doe")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "firstName")
}

func (suite *PersonNameTestSuite) TestEmptyLastName() {
	_, err := valueobject.NewPersonName("John", "")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "lastName")
}

func (suite *PersonNameTestSuite) TestPersonNameEquality() {
	name1, _ := valueobject.NewPersonName("John", "Doe")
	name2, _ := valueobject.NewPersonName("John", "Doe")
	name3, _ := valueobject.NewPersonName("Jane", "Doe")

	suite.True(name1.Equals(name2))
	suite.False(name1.Equals(name3))
	suite.False(name1.Equals(nil))
}

func (suite *PersonNameTestSuite) TestNameTrimming() {
	name, err := valueobject.NewPersonName("  John  ", "  Doe  ")

	suite.Require().NoError(err)
	suite.Equal("John", name.FirstName())
	suite.Equal("Doe", name.LastName())
}
