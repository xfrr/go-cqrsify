package valueobject

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var _ ValueObject = (*Money)(nil)

// Money value object
type Money struct {
	BaseValueObject
	amountCents int64
	currencyISO string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) (Money, error) {
	const conversionFactor = 100

	m := Money{
		amountCents: int64(math.Round(amount * conversionFactor)),
		currencyISO: strings.ToUpper(strings.TrimSpace(currency)),
	}

	if err := m.Validate(); err != nil {
		return Money{}, err
	}
	return m, nil
}

func (m Money) Amount() float64 {
	const conversionFactor = 100
	return math.Round(float64(m.amountCents)) / conversionFactor
}

func (m Money) AmountInCents() int64 {
	return m.amountCents
}

func (m Money) Currency() string {
	return m.currencyISO
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount(), m.currencyISO)
}

func (m Money) Validate() error {
	const currencyLength = 3

	var errs []ValidationError
	if m.amountCents < 0 {
		errs = append(errs, ValidationError{Field: "amount", Message: "cannot be negative"})
	}

	if m.currencyISO == "" {
		errs = append(errs, ValidationError{Field: "currency", Message: "cannot be empty"})
	} else if len(m.currencyISO) != currencyLength {
		errs = append(errs, ValidationError{Field: "currency", Message: fmt.Sprintf("must be %d characters (ISO 4217)", currencyLength)})
	}

	return ValidationErrors(errs)
}

func (m Money) Equals(other ValueObject) bool {
	if otherMoney, ok := other.(Money); ok {
		return m.amountCents == otherMoney.amountCents &&
			m.currencyISO == otherMoney.currencyISO
	}
	return false
}

// Add adds two Money values (same currency only)
func (m Money) Add(other Money) (Money, error) {
	if m.currencyISO != other.currencyISO {
		return Money{}, errors.New("cannot add money with different currencies")
	}

	return Money{
		amountCents: m.amountCents + other.amountCents,
		currencyISO: m.currencyISO,
	}, nil
}

// Subtract subtracts two Money values (same currency only)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currencyISO != other.currencyISO {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}

	if m.amountCents < other.amountCents {
		return Money{}, errors.New("resulting amount cannot be negative")
	}

	return Money{
		amountCents: m.amountCents - other.amountCents,
		currencyISO: m.currencyISO,
	}, nil
}
