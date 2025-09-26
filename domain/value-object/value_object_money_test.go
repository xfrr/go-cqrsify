package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestMoneySuite(t *testing.T) {
	suite.Run(t, new(MoneyTestSuite))
}

// MoneyTestSuite groups money-related tests
type MoneyTestSuite struct {
	suite.Suite
}

func (suite *MoneyTestSuite) TestValidMoney() {
	money, err := valueobject.NewMoney(19.99, "USD")

	suite.Require().NoError(err)
	suite.InDelta(19.99, money.Amount(), 0.001)
	suite.Equal("USD", money.Currency())
	suite.Equal(int64(1999), money.AmountInCents())
}

func (suite *MoneyTestSuite) TestNegativeAmount() {
	_, err := valueobject.NewMoney(-10.0, "USD")

	suite.Require().Error(err)
	suite.Contains(err.Error(), "cannot be negative")
}

func (suite *MoneyTestSuite) TestInvalidCurrency() {
	tests := []struct {
		name     string
		currency string
	}{
		{"empty currency", ""},
		{"invalid length", "INVALID"},
		{"too short", "US"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			_, err := valueobject.NewMoney(10.0, tt.currency)
			suite.Require().Error(err)
		})
	}
}

func (suite *MoneyTestSuite) TestAddMoney() {
	money1, _ := valueobject.NewMoney(10.0, "USD")
	money2, _ := valueobject.NewMoney(5.0, "USD")

	result, err := money1.Add(money2)

	suite.Require().NoError(err)
	suite.InDelta(15.0, result.Amount(), 0.001)
}

func (suite *MoneyTestSuite) TestAddDifferentCurrencies() {
	money1, _ := valueobject.NewMoney(10.0, "USD")
	money2, _ := valueobject.NewMoney(5.0, "EUR")

	_, err := money1.Add(money2)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "different currencies")
}

func (suite *MoneyTestSuite) TestSubtractMoney() {
	money1, _ := valueobject.NewMoney(15.0, "USD")
	money2, _ := valueobject.NewMoney(5.0, "USD")

	result, err := money1.Subtract(money2)

	suite.Require().NoError(err)
	suite.InDelta(10.0, result.Amount(), 0.001)
}

func (suite *MoneyTestSuite) TestSubtractResultingInNegative() {
	money1, _ := valueobject.NewMoney(5.0, "USD")
	money2, _ := valueobject.NewMoney(10.0, "USD")

	_, err := money1.Subtract(money2)
	suite.Require().Error(err)
}

func (suite *MoneyTestSuite) TestMoneyEquality() {
	money1, _ := valueobject.NewMoney(10.0, "USD")
	money2, _ := valueobject.NewMoney(10.0, "USD")
	money3, _ := valueobject.NewMoney(10.0, "EUR")
	money4, _ := valueobject.NewMoney(15.0, "USD")

	suite.True(money1.Equals(money2))
	suite.False(money1.Equals(money3))
	suite.False(money1.Equals(money4))
	suite.False(money1.Equals(nil))
}

func (suite *MoneyTestSuite) TestMoneyString() {
	money, _ := valueobject.NewMoney(19.99, "USD")

	suite.Equal("19.99 USD", money.String())
}
