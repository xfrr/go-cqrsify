package valueobject

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Money value object
type Money struct {
	BaseValueObject
	amount   int64 // Store as cents to avoid floating point issues
	currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) (*Money, error) {
	m := &Money{
		amount:   int64(math.Round(amount * 100)), // Convert to cents with proper rounding
		currency: strings.ToUpper(strings.TrimSpace(currency)),
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// NewMoneyFromCents creates Money from cents (avoiding float conversion)
func NewMoneyFromCents(cents int64, currency string) (*Money, error) {
	m := &Money{
		amount:   cents,
		currency: strings.ToUpper(strings.TrimSpace(currency)),
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Money) Amount() float64 {
	return math.Round(float64(m.amount)) / 100
}

func (m *Money) AmountInCents() int64 {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}

func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount(), m.currency)
}

func (m *Money) Validate() error {
	var errs []ValidationError

	if m.amount < 0 {
		errs = append(errs, ValidationError{Field: "amount", Message: "cannot be negative"})
	}

	if m.currency == "" {
		errs = append(errs, ValidationError{Field: "currency", Message: "cannot be empty"})
	} else if len(m.currency) != 3 {
		errs = append(errs, ValidationError{Field: "currency", Message: "must be 3 characters (ISO 4217)"})
	}

	if len(errs) > 0 {
		return MultiValidationError{Errors: errs}
	}
	return nil
}

func (m *Money) Equals(other ValueObject) bool {
	if otherMoney, ok := other.(*Money); ok {
		return m.amount == otherMoney.amount && m.currency == otherMoney.currency
	}
	return false
}

// Add adds two Money values (same currency only)
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, errors.New("cannot add money with different currencies")
	}
	return NewMoneyFromCents(m.amount+other.amount, m.currency)
}

// Subtract subtracts two Money values (same currency only)
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, errors.New("cannot subtract money with different currencies")
	}
	return NewMoneyFromCents(m.amount-other.amount, m.currency)
}
