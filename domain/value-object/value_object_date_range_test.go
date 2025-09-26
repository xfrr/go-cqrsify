package valueobject_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestDateRangeSuite(t *testing.T) {
	suite.Run(t, new(DateRangeTestSuite))
}

// DateRangeTestSuite groups date range-related tests
type DateRangeTestSuite struct {
	suite.Suite
	startDate time.Time
	endDate   time.Time
}

func (suite *DateRangeTestSuite) SetupTest() {
	suite.startDate = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.endDate = time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
}

func (suite *DateRangeTestSuite) TestValidDateRange() {
	dateRange, err := valueobject.NewDateRange(suite.startDate, suite.endDate)

	suite.Require().NoError(err)
	suite.True(dateRange.StartDate().Equal(suite.startDate))
	suite.True(dateRange.EndDate().Equal(suite.endDate))
}

func (suite *DateRangeTestSuite) TestInvalidDateRange() {
	_, err := valueobject.NewDateRange(suite.endDate, suite.startDate) // end before start

	suite.Require().Error(err)
	suite.Contains(err.Error(), "start date must be before")
}

func (suite *DateRangeTestSuite) TestContainsDate() {
	dateRange, _ := valueobject.NewDateRange(suite.startDate, suite.endDate)
	midDate := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	outsideDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	suite.True(dateRange.Contains(midDate))
	suite.True(dateRange.Contains(suite.startDate)) // boundary inclusive
	suite.True(dateRange.Contains(suite.endDate))   // boundary inclusive
	suite.False(dateRange.Contains(outsideDate))
}

func (suite *DateRangeTestSuite) TestOverlaps() {
	range1, _ := valueobject.NewDateRange(
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
	)
	range2, _ := valueobject.NewDateRange(
		time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
	)
	range3, _ := valueobject.NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
	)

	suite.True(range1.Overlaps(range2))
	suite.False(range1.Overlaps(range3))
}

func (suite *DateRangeTestSuite) TestDateRangeEquality() {
	range1, _ := valueobject.NewDateRange(suite.startDate, suite.endDate)
	range2, _ := valueobject.NewDateRange(suite.startDate, suite.endDate)
	range3, _ := valueobject.NewDateRange(suite.startDate, suite.startDate)

	suite.True(range1.Equals(range2))
	suite.False(range1.Equals(range3))
	suite.False(range1.Equals(nil))
}

func (suite *DateRangeTestSuite) TestDuration() {
	dateRange, _ := valueobject.NewDateRange(suite.startDate, suite.endDate)
	expectedDuration := suite.endDate.Sub(suite.startDate)

	suite.Equal(expectedDuration, dateRange.Duration())
}

func (suite *DateRangeTestSuite) TestDateRangeString() {
	dateRange, _ := valueobject.NewDateRange(suite.startDate, suite.endDate)
	expected := "2023-01-01 to 2023-12-31"

	suite.Equal(expected, dateRange.String())
}
