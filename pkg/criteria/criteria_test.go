package criteria_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xfrr/go-cqrsify/pkg/criteria"
)

// Test data structures
type Person struct {
	Name string
	Age  int
	City string
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

// Run the test suite
func TestCriteriaTestSuite(t *testing.T) {
	suite.Run(t, new(CriteriaTestSuite))
}

// CriteriaTestSuite provides a test suite for all criteria functionality
type CriteriaTestSuite struct {
	suite.Suite
	people   []Person
	products []Product
}

func (suite *CriteriaTestSuite) SetupTest() {
	suite.people = []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Charlie", Age: 35, City: "New York"},
		{Name: "David", Age: 25, City: "Paris"},
		{Name: "Eve", Age: 40, City: "London"},
	}

	suite.products = []Product{
		{ID: 1, Name: "Laptop", Price: 999.99},
		{ID: 2, Name: "Mouse", Price: 25.50},
		{ID: 3, Name: "Keyboard", Price: 75.00},
		{ID: 4, Name: "Monitor", Price: 299.99},
	}
}

// PredicateCriteria Tests
func (suite *CriteriaTestSuite) TestPredicateCriteria_WithValidPredicate() {
	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 30
	})

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Charlie", Age: 35, City: "New York"},
		{Name: "Eve", Age: 40, City: "London"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestPredicateCriteria_WithNoMatches() {
	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age > 50
	})

	result := criteria.MeetsCriteria(suite.people)

	assert.Empty(suite.T(), result)
}

func (suite *CriteriaTestSuite) TestPredicateCriteria_WithAllMatches() {
	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return len(p.Name) > 0
	})

	result := criteria.MeetsCriteria(suite.people)

	assert.Equal(suite.T(), suite.people, result)
}

func (suite *CriteriaTestSuite) TestPredicateCriteria_WithEmptySlice() {
	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 30
	})

	result := criteria.MeetsCriteria([]Person{})

	assert.Empty(suite.T(), result)
}

// FieldCriteria Tests
func (suite *CriteriaTestSuite) TestFieldCriteria_EqualComparison() {
	criteria := criteria.NewFieldCriteria(
		func(p Person) int { return p.Age },
		25,
		criteria.ComparisonFunctions[int]{}.Equal(),
	)

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "David", Age: 25, City: "Paris"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestFieldCriteria_StringComparison() {
	criteria := criteria.NewFieldCriteria(
		func(p Person) string { return p.City },
		"London",
		criteria.ComparisonFunctions[string]{}.Equal(),
	)

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Eve", Age: 40, City: "London"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestFieldCriteria_OrderedComparisons() {
	testCases := []struct {
		name      string
		compareFn func(int, int) bool
		value     int
		expected  []Person
	}{
		{
			name:      "GreaterThan",
			compareFn: criteria.OrderedComparisonFunctions[int]{}.GreaterThan(),
			value:     30,
			expected: []Person{
				{Name: "Charlie", Age: 35, City: "New York"},
				{Name: "Eve", Age: 40, City: "London"},
			},
		},
		{
			name:      "LessThan",
			compareFn: criteria.OrderedComparisonFunctions[int]{}.LessThan(),
			value:     30,
			expected: []Person{
				{Name: "Alice", Age: 25, City: "New York"},
				{Name: "David", Age: 25, City: "Paris"},
			},
		},
		{
			name:      "GreaterThanOrEqual",
			compareFn: criteria.OrderedComparisonFunctions[int]{}.GreaterThanOrEqual(),
			value:     30,
			expected: []Person{
				{Name: "Bob", Age: 30, City: "London"},
				{Name: "Charlie", Age: 35, City: "New York"},
				{Name: "Eve", Age: 40, City: "London"},
			},
		},
		{
			name:      "LessThanOrEqual",
			compareFn: criteria.OrderedComparisonFunctions[int]{}.LessThanOrEqual(),
			value:     30,
			expected: []Person{
				{Name: "Alice", Age: 25, City: "New York"},
				{Name: "Bob", Age: 30, City: "London"},
				{Name: "David", Age: 25, City: "Paris"},
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			criteria := criteria.NewFieldCriteria(
				func(p Person) int { return p.Age },
				tc.value,
				tc.compareFn,
			)

			result := criteria.MeetsCriteria(suite.people)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func (suite *CriteriaTestSuite) TestFieldCriteria_StringSpecificComparisons() {
	testCases := []struct {
		name      string
		compareFn func(string, string) bool
		value     string
		expected  []Person
	}{
		{
			name:      "Contains",
			compareFn: criteria.StringComparisonFunctions{}.Contains(),
			value:     "New",
			expected: []Person{
				{Name: "Alice", Age: 25, City: "New York"},
				{Name: "Charlie", Age: 35, City: "New York"},
			},
		},
		{
			name:      "HasPrefix",
			compareFn: criteria.StringComparisonFunctions{}.HasPrefix(),
			value:     "L",
			expected: []Person{
				{Name: "Bob", Age: 30, City: "London"},
				{Name: "Eve", Age: 40, City: "London"},
			},
		},
		{
			name:      "HasSuffix",
			compareFn: criteria.StringComparisonFunctions{}.HasSuffix(),
			value:     "on",
			expected: []Person{
				{Name: "Bob", Age: 30, City: "London"},
				{Name: "Eve", Age: 40, City: "London"},
			},
		},
		{
			name:      "EqualFold",
			compareFn: criteria.StringComparisonFunctions{}.EqualFold(),
			value:     "LONDON",
			expected: []Person{
				{Name: "Bob", Age: 30, City: "London"},
				{Name: "Eve", Age: 40, City: "London"},
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			criteria := criteria.NewFieldCriteria(
				func(p Person) string { return p.City },
				tc.value,
				tc.compareFn,
			)

			result := criteria.MeetsCriteria(suite.people)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// AndCriteria Tests
func (suite *CriteriaTestSuite) TestAndCriteria_BothCriteriaMet() {
	ageCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 25
	})
	cityCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.City == "New York"
	})

	andCriteria := criteria.NewAndCriteria(ageCriteria, cityCriteria)
	result := andCriteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Charlie", Age: 35, City: "New York"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestAndCriteria_NoCriteriaMatch() {
	ageCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age > 50
	})
	cityCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.City == "New York"
	})

	andCriteria := criteria.NewAndCriteria(ageCriteria, cityCriteria)
	result := andCriteria.MeetsCriteria(suite.people)

	assert.Empty(suite.T(), result)
}

// OrCriteria Tests
func (suite *CriteriaTestSuite) TestOrCriteria_EitherCriteriaMet() {
	ageCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 40
	})
	cityCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.City == "Paris"
	})

	orCriteria := criteria.NewOrCriteria(ageCriteria, cityCriteria)
	result := orCriteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "David", Age: 25, City: "Paris"},
		{Name: "Eve", Age: 40, City: "London"},
	}
	assert.ElementsMatch(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestOrCriteria_NoDuplicates() {
	ageCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age == 25
	})
	cityCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.City == "New York"
	})

	orCriteria := criteria.NewOrCriteria(ageCriteria, cityCriteria)
	result := orCriteria.MeetsCriteria(suite.people)

	// Alice should appear only once even though she matches both criteria
	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "David", Age: 25, City: "Paris"},
		{Name: "Charlie", Age: 35, City: "New York"},
	}
	assert.ElementsMatch(suite.T(), expected, result)
	assert.Len(suite.T(), result, 3)
}

// NotCriteria Tests
func (suite *CriteriaTestSuite) TestNotCriteria_InvertsPredicate() {
	ageCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 35
	})

	notCriteria := criteria.NewNotCriteria(ageCriteria)
	result := notCriteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "David", Age: 25, City: "Paris"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestNotCriteria_WithAllMatching() {
	allMatchCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return len(p.Name) > 0
	})

	notCriteria := criteria.NewNotCriteria(allMatchCriteria)
	result := notCriteria.MeetsCriteria(suite.people)

	assert.Empty(suite.T(), result)
}

func (suite *CriteriaTestSuite) TestNotCriteria_WithNoMatches() {
	noMatchCriteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age > 50
	})

	notCriteria := criteria.NewNotCriteria(noMatchCriteria)
	result := notCriteria.MeetsCriteria(suite.people)

	assert.Equal(suite.T(), suite.people, result)
}

// CriteriaBuilder Tests
func (suite *CriteriaTestSuite) TestCriteriaBuilder_FluentInterface() {
	criteria := criteria.NewCriteriaBuilder[Person]().
		WithCriteria(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.Age >= 25
		})).
		And(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.City == "New York"
		})).
		Build()

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Charlie", Age: 35, City: "New York"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestCriteriaBuilder_ComplexChaining() {
	criteria := criteria.NewCriteriaBuilder[Person]().
		WithCriteria(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.Age <= 30
		})).
		Or(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.City == "New York"
		})).
		Build()

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Charlie", Age: 35, City: "New York"},
		{Name: "David", Age: 25, City: "Paris"},
	}
	assert.ElementsMatch(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestCriteriaBuilder_WithNot() {
	criteria := criteria.NewCriteriaBuilder[Person]().
		WithCriteria(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.Age >= 35
		})).
		Not().
		Build()

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "David", Age: 25, City: "Paris"},
	}
	assert.Equal(suite.T(), expected, result)
}

func (suite *CriteriaTestSuite) TestCriteriaBuilder_EmptyBuilder() {
	builder := criteria.NewCriteriaBuilder[Person]()
	criteria := builder.Build()

	assert.Nil(suite.T(), criteria)
}

func (suite *CriteriaTestSuite) TestCriteriaBuilder_AddingToEmptyBuilder() {
	criteria := criteria.NewCriteriaBuilder[Person]().
		And(criteria.NewPredicateCriteria(func(p Person) bool {
			return p.Age >= 30
		})).
		Build()

	result := criteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Charlie", Age: 35, City: "New York"},
		{Name: "Eve", Age: 40, City: "London"},
	}
	assert.Equal(suite.T(), expected, result)
}

// Integration tests
func (suite *CriteriaTestSuite) TestIntegration_ComplexCriteriaChain() {
	// Find people who are either:
	// - From London and age >= 30, OR
	// - From New York and age <= 30
	londonAdults := criteria.NewAndCriteria(
		criteria.NewFieldCriteria(
			func(p Person) string { return p.City },
			"London",
			criteria.ComparisonFunctions[string]{}.Equal(),
		),
		criteria.NewFieldCriteria(
			func(p Person) int { return p.Age },
			30,
			criteria.OrderedComparisonFunctions[int]{}.GreaterThanOrEqual(),
		),
	)

	newYorkYoung := criteria.NewAndCriteria(
		criteria.NewFieldCriteria(
			func(p Person) string { return p.City },
			"New York",
			criteria.ComparisonFunctions[string]{}.Equal(),
		),
		criteria.NewFieldCriteria(
			func(p Person) int { return p.Age },
			30,
			criteria.OrderedComparisonFunctions[int]{}.LessThanOrEqual(),
		),
	)

	finalCriteria := criteria.NewOrCriteria(londonAdults, newYorkYoung)
	result := finalCriteria.MeetsCriteria(suite.people)

	expected := []Person{
		{Name: "Alice", Age: 25, City: "New York"},
		{Name: "Bob", Age: 30, City: "London"},
		{Name: "Eve", Age: 40, City: "London"},
	}
	assert.ElementsMatch(suite.T(), expected, result)
}

// Edge cases and error scenarios
func (suite *CriteriaTestSuite) TestEdgeCases_EmptyInput() {
	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 25
	})

	result := criteria.MeetsCriteria([]Person{})
	assert.Empty(suite.T(), result)
}

// Additional individual tests for specific functionality
func TestCriteria_ComparisonFunctions_NotEqual(t *testing.T) {
	comp := criteria.ComparisonFunctions[int]{}
	notEqualFn := comp.NotEqual()

	assert.True(t, notEqualFn(5, 10))
	assert.False(t, notEqualFn(5, 5))
}

func TestCriteria_OrderedCriteria_ComparisonFunctions__NotEqual(t *testing.T) {
	comp := criteria.ComparisonFunctions[int]{}
	notEqualFn := comp.NotEqual()

	assert.True(t, notEqualFn(5, 10))
	assert.False(t, notEqualFn(5, 5))
}

func TestStringCriteria_ComparisonFunctions_AllMethods(t *testing.T) {
	comp := criteria.StringComparisonFunctions{}

	tests := []struct {
		name     string
		fn       func(string, string) bool
		field    string
		value    string
		expected bool
	}{
		{"Contains", comp.Contains(), "hello world", "world", true},
		{"Contains_NotFound", comp.Contains(), "hello world", "xyz", false},
		{"HasPrefix", comp.HasPrefix(), "hello world", "hello", true},
		{"HasPrefix_NotFound", comp.HasPrefix(), "hello world", "world", false},
		{"HasSuffix", comp.HasSuffix(), "hello world", "world", true},
		{"HasSuffix_NotFound", comp.HasSuffix(), "hello world", "hello", false},
		{"EqualFold", comp.EqualFold(), "Hello", "HELLO", true},
		{"EqualFold_NotEqual", comp.EqualFold(), "Hello", "World", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.field, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark test
func (suite *CriteriaTestSuite) TestPerformance_LargeDataset() {
	// Create a larger dataset for performance testing
	largePeople := make([]Person, 1000)
	for i := range largePeople {
		largePeople[i] = Person{
			Name: "Person" + string(rune(i)),
			Age:  20 + (i % 60),
			City: []string{"New York", "London", "Paris", "Tokyo"}[i%4],
		}
	}

	criteria := criteria.NewPredicateCriteria(func(p Person) bool {
		return p.Age >= 30 && p.City == "London"
	})

	result := criteria.MeetsCriteria(largePeople)

	// Should find people with age >= 30 from London
	assert.True(suite.T(), len(result) > 0)
	for _, person := range result {
		assert.True(suite.T(), person.Age >= 30)
		assert.Equal(suite.T(), "London", person.City)
	}
}
