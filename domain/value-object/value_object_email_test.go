package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

// EmailTestSuite groups email-related tests
type EmailTestSuite struct {
	suite.Suite
}

func (suite *EmailTestSuite) TestValidEmail() {
	email, err := valueobject.NewEmail("test@example.com")

	suite.Require().NoError(err)
	suite.Equal("test@example.com", email.Value())
}

func (suite *EmailTestSuite) TestEmailNormalization() {
	email, err := valueobject.NewEmail("  TEST@EXAMPLE.COM  ")

	suite.Require().NoError(err)
	suite.Equal("test@example.com", email.Value())
}

func (suite *EmailTestSuite) TestInvalidEmail() {
	_, err := valueobject.NewEmail("invalid-email")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "invalid email format")
}

func (suite *EmailTestSuite) TestEmptyEmail() {
	_, err := valueobject.NewEmail("")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "cannot be empty")
}

func (suite *EmailTestSuite) TestEmailEquality() {
	email1, _ := valueobject.NewEmail("test@example.com")
	email2, _ := valueobject.NewEmail("test@example.com")
	email3, _ := valueobject.NewEmail("other@example.com")

	suite.True(email1.Equals(email2))
	suite.False(email1.Equals(email3))
	suite.False(email1.Equals(nil))
}

func (suite *EmailTestSuite) TestMustNewEmail() {
	email := valueobject.MustNewEmail("test@example.com")
	suite.Equal("test@example.com", email.Value())

	suite.Panics(func() {
		valueobject.MustNewEmail("invalid-email")
	})
}
