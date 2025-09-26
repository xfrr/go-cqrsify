package valueobject_test

import (
	"errors"
	"testing"

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

	suite.Require().Error(err)

	var multiErr valueobject.MultiValidationError
	ok := errors.As(err, &multiErr)
	suite.Require().True(ok, "expected MultiValidationError")
	suite.Require().Len(multiErr.Errors, 5)

	// Check that all required fields are mentioned in the error
	errorStr := err.Error()
	suite.Contains(errorStr, "street")
	suite.Contains(errorStr, "city")
	suite.Contains(errorStr, "state")
	suite.Contains(errorStr, "zipCode")
	suite.Contains(errorStr, "country")
}

func (suite *ValidationTestSuite) TestValidationErrorString() {
	err := valueobject.ValidationError{Field: "testField", Message: "test message"}
	expected := "validation error on field 'testField': test message"

	suite.Equal(expected, err.Error())
}
