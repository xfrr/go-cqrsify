package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestValidationSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

// ValidationTestSuite groups validation-related tests
type ValidationTestSuite struct {
	suite.Suite
}

func (suite *ValidationTestSuite) TestMultipleValidationErrors() {
	_, err := valueobject.NewAddress("", "", "", "", "")

	require.Error(suite.T(), err)

	multiErr, ok := err.(valueobject.MultiValidationError)
	require.True(suite.T(), ok, "expected MultiValidationError")
	assert.Len(suite.T(), multiErr.Errors, 5)

	// Check that all required fields are mentioned in the error
	errorStr := err.Error()
	assert.Contains(suite.T(), errorStr, "street")
	assert.Contains(suite.T(), errorStr, "city")
	assert.Contains(suite.T(), errorStr, "state")
	assert.Contains(suite.T(), errorStr, "zipCode")
	assert.Contains(suite.T(), errorStr, "country")
}

func (suite *ValidationTestSuite) TestValidationErrorString() {
	err := valueobject.ValidationError{Field: "testField", Message: "test message"}
	expected := "validation error on field 'testField': test message"

	assert.Equal(suite.T(), expected, err.Error())
}
