package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test@example.com", email.Value())
}

func (suite *EmailTestSuite) TestEmailNormalization() {
	email, err := valueobject.NewEmail("  TEST@EXAMPLE.COM  ")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test@example.com", email.Value())
}

func (suite *EmailTestSuite) TestInvalidEmail() {
	_, err := valueobject.NewEmail("invalid-email")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid email format")
}

func (suite *EmailTestSuite) TestEmptyEmail() {
	_, err := valueobject.NewEmail("")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be empty")
}

func (suite *EmailTestSuite) TestEmailEquality() {
	email1, _ := valueobject.NewEmail("test@example.com")
	email2, _ := valueobject.NewEmail("test@example.com")
	email3, _ := valueobject.NewEmail("other@example.com")

	assert.True(suite.T(), email1.Equals(email2))
	assert.False(suite.T(), email1.Equals(email3))
	assert.False(suite.T(), email1.Equals(nil))
}

func (suite *EmailTestSuite) TestMustNewEmail() {
	email := valueobject.MustNewEmail("test@example.com")
	assert.Equal(suite.T(), "test@example.com", email.Value())

	assert.Panics(suite.T(), func() {
		valueobject.MustNewEmail("invalid-email")
	})
}
