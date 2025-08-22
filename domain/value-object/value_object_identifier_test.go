package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

type IdentifierTestSuite struct {
	suite.Suite
}

func TestIdentifierSuite(t *testing.T) {
	suite.Run(t, new(IdentifierTestSuite))
}

// TestNewIdentifier tests the constructor function
func (suite *IdentifierTestSuite) TestNewIdentifier() {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{"string identifier", "test-id", "test-id"},
		{"int identifier", 123, 123},
		{"uint identifier", uint(456), uint(456)},
		{"float64 identifier", 3.14, 3.14},
		{"int64 identifier", int64(789), int64(789)},
		{"uint32 identifier", uint32(101112), uint32(101112)},
		{"float32 identifier", float32(2.71), float32(2.71)},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			switch v := tt.value.(type) {
			case string:
				id := valueobject.NewIdentifier(v)
				// Since value field is not exported, we test through String() method
				assert.Equal(suite.T(), tt.expected, id.Value())
			case int:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			case uint:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			case float64:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			case int64:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			case uint32:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			case float32:
				id := valueobject.NewIdentifier(v)
				assert.Equal(suite.T(), tt.expected, id.Value())
			}
		})
	}
}

// TestString tests the String method for various types
func (suite *IdentifierTestSuite) TestString() {
	tests := []struct {
		name     string
		id       interface{}
		expected string
	}{
		// String types
		{"string value", valueobject.NewIdentifier("hello"), "hello"},
		{"empty string", valueobject.NewIdentifier(""), ""},

		// Integer types
		{"int value", valueobject.NewIdentifier(123), "123"},
		{"negative int", valueobject.NewIdentifier(-456), "-456"},
		{"int8 value", valueobject.NewIdentifier(int8(127)), "127"},
		{"int16 value", valueobject.NewIdentifier(int16(32767)), "32767"},
		{"int32 value", valueobject.NewIdentifier(int32(2147483647)), "2147483647"},
		{"int64 value", valueobject.NewIdentifier(int64(9223372036854775807)), "9223372036854775807"},

		// Unsigned integer types
		{"uint value", valueobject.NewIdentifier(uint(123)), "123"},
		{"uint8 value", valueobject.NewIdentifier(uint8(255)), "255"},
		{"uint16 value", valueobject.NewIdentifier(uint16(65535)), "65535"},
		{"uint32 value", valueobject.NewIdentifier(uint32(4294967295)), "4294967295"},
		{"uint64 value", valueobject.NewIdentifier(uint64(18446744073709551615)), "18446744073709551615"},

		// Float types
		{"float32 value", valueobject.NewIdentifier(float32(3.14)), "3.140000104904175"},
		{"float64 value", valueobject.NewIdentifier(float64(2.718281828)), "2.718281828"},
		{"float with no decimals", valueobject.NewIdentifier(float64(42)), "42"},
		{"negative float", valueobject.NewIdentifier(float64(-3.14)), "-3.14"},

		// Boolean (falls to default case)
		{"bool value", valueobject.NewIdentifier(true), "true"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, ok := tt.id.(valueobject.ValueObject)
			assert.True(suite.T(), ok)
			assert.Equal(suite.T(), tt.expected, result.String())
		})
	}
}

// TestEquals tests the equality comparison
func (suite *IdentifierTestSuite) TestEquals() {
	// String identifiers
	suite.Run("equal string identifiers", func() {
		id1 := valueobject.NewIdentifier("test-id")
		id2 := valueobject.NewIdentifier("test-id")
		assert.True(suite.T(), id1.Equals(id2))
	})

	suite.Run("different string identifiers", func() {
		id1 := valueobject.NewIdentifier("test-id-1")
		id2 := valueobject.NewIdentifier("test-id-2")
		assert.False(suite.T(), id1.Equals(id2))
	})

	// Integer identifiers
	suite.Run("equal int identifiers", func() {
		id1 := valueobject.NewIdentifier(123)
		id2 := valueobject.NewIdentifier(123)
		assert.True(suite.T(), id1.Equals(id2))
	})

	suite.Run("different int identifiers", func() {
		id1 := valueobject.NewIdentifier(123)
		id2 := valueobject.NewIdentifier(456)
		assert.False(suite.T(), id1.Equals(id2))
	})

	// Float identifiers
	suite.Run("equal float identifiers", func() {
		id1 := valueobject.NewIdentifier(3.14)
		id2 := valueobject.NewIdentifier(3.14)
		assert.True(suite.T(), id1.Equals(id2))
	})

	suite.Run("different float identifiers", func() {
		id1 := valueobject.NewIdentifier(3.14)
		id2 := valueobject.NewIdentifier(2.71)
		assert.False(suite.T(), id1.Equals(id2))
	})

	// Boolean identifiers
	suite.Run("equal bool identifiers", func() {
		id1 := valueobject.NewIdentifier(true)
		id2 := valueobject.NewIdentifier(true)
		assert.True(suite.T(), id1.Equals(id2))
	})

	suite.Run("different bool identifiers", func() {
		id1 := valueobject.NewIdentifier(true)
		id2 := valueobject.NewIdentifier(false)
		assert.False(suite.T(), id1.Equals(id2))
	})
}

// TestValidate tests the validation logic
func (suite *IdentifierTestSuite) TestValidate() {
	// Valid cases
	validCases := []struct {
		name string
		id   interface{}
	}{
		{"valid string", valueobject.NewIdentifier("valid-id")},
		{"valid positive int", valueobject.NewIdentifier(123)},
		{"valid negative int", valueobject.NewIdentifier(-123)},
		{"valid positive uint", valueobject.NewIdentifier(uint(456))},
		{"valid positive float", valueobject.NewIdentifier(3.14)},
		{"valid negative float", valueobject.NewIdentifier(-3.14)},
		{"valid int64", valueobject.NewIdentifier(int64(789))},
		{"valid uint32", valueobject.NewIdentifier(uint32(101112))},
		{"valid float32", valueobject.NewIdentifier(float32(2.71))},
	}

	for _, tt := range validCases {
		suite.Run(tt.name, func() {
			err := tt.id.(interface{ Validate() error }).Validate()
			assert.NoError(suite.T(), err)
		})
	}

	// Invalid cases
	invalidCases := []struct {
		name          string
		id            interface{}
		expectedError string
	}{
		{"empty string", valueobject.NewIdentifier(""), `invalid identifier: ""`},
		{"zero int", valueobject.NewIdentifier(0), "invalid identifier: 0"},
		{"zero uint", valueobject.NewIdentifier(uint(0)), "invalid identifier: 0"},
		{"zero float", valueobject.NewIdentifier(0.0), "invalid identifier: 0.000000"},
		{"zero int64", valueobject.NewIdentifier(int64(0)), "invalid identifier: 0"},
		{"zero uint32", valueobject.NewIdentifier(uint32(0)), "invalid identifier: 0"},
		{"zero float32", valueobject.NewIdentifier(float32(0)), "invalid identifier: 0.000000"},
		{"unsupported type", valueobject.NewIdentifier(true), "invalid identifier: true"},
	}

	for _, tt := range invalidCases {
		suite.Run(tt.name, func() {
			err := tt.id.(interface{ Validate() error }).Validate()
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), tt.expectedError, err.Error())
		})
	}
}

// Additional unit tests using basic testing functions for specific edge cases
func TestIdentifierStringEdgeCases(t *testing.T) {
	t.Run("very large numbers", func(t *testing.T) {
		// Test edge values for different integer types
		maxInt64 := valueobject.NewIdentifier(int64(9223372036854775807))
		assert.Equal(t, "9223372036854775807", maxInt64.String())

		maxUint64 := valueobject.NewIdentifier(uint64(18446744073709551615))
		assert.Equal(t, "18446744073709551615", maxUint64.String())
	})

	t.Run("scientific notation floats", func(t *testing.T) {
		// Very small float
		smallFloat := valueobject.NewIdentifier(0.000000001)
		result := smallFloat.String()
		assert.Contains(t, result, "0.000000001")

		// Very large float
		largeFloat := valueobject.NewIdentifier(1e20)
		result2 := largeFloat.String()
		assert.Equal(t, "100000000000000000000", result2)
	})
}

func TestIdentifierValidationEdgeCases(t *testing.T) {
	t.Run("complex type", func(t *testing.T) {
		type customType struct {
			value string
		}
		customId := valueobject.NewIdentifier(customType{value: "test"})
		err := customId.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid identifier")
	})
}

// Benchmark tests
func BenchmarkNewIdentifier(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueobject.NewIdentifier("test-id")
		}
	})

	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			valueobject.NewIdentifier(123)
		}
	})
}

func BenchmarkString(b *testing.B) {
	id := valueobject.NewIdentifier("test-id")

	for b.Loop() {
		_ = id.String()
	}
}

func BenchmarkEquals(b *testing.B) {
	id1 := valueobject.NewIdentifier("test-id")
	id2 := valueobject.NewIdentifier("test-id")

	for b.Loop() {
		id1.Equals(id2)
	}
}

func BenchmarkValidate(b *testing.B) {
	id := valueobject.NewIdentifier("test-id")

	for b.Loop() {
		id.Validate()
	}
}
