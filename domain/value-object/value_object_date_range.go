package valueobject

import (
	"fmt"
	"time"
)

var _ ValueObject = (*DateRange)(nil)

// DateRange value object
type DateRange struct {
	BaseValueObject
	startDate time.Time
	endDate   time.Time
}

// NewDateRange creates a new DateRange value object
func NewDateRange(startDate, endDate time.Time) (*DateRange, error) {
	// Truncate to day precision
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	dr := &DateRange{
		startDate: start,
		endDate:   end,
	}
	if err := dr.Validate(); err != nil {
		return nil, err
	}
	return dr, nil
}

func (dr *DateRange) StartDate() time.Time {
	return dr.startDate
}

func (dr *DateRange) EndDate() time.Time {
	return dr.endDate
}

func (dr *DateRange) Duration() time.Duration {
	return dr.endDate.Sub(dr.startDate)
}

func (dr *DateRange) String() string {
	return fmt.Sprintf("%s to %s",
		dr.startDate.Format("2006-01-02"),
		dr.endDate.Format("2006-01-02"))
}

func (dr *DateRange) Validate() error {
	if dr.startDate.After(dr.endDate) {
		return ValidationError{Field: "dateRange", Message: "start date must be before or equal to end date"}
	}
	return nil
}

func (dr *DateRange) Equals(other ValueObject) bool {
	if otherRange, ok := other.(*DateRange); ok {
		return dr.startDate.Equal(otherRange.startDate) && dr.endDate.Equal(otherRange.endDate)
	}
	return false
}

// Contains checks if a date is within the range (inclusive)
func (dr *DateRange) Contains(date time.Time) bool {
	checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return (checkDate.Equal(dr.startDate) || checkDate.After(dr.startDate)) &&
		(checkDate.Equal(dr.endDate) || checkDate.Before(dr.endDate))
}

// Overlaps checks if this date range overlaps with another
func (dr *DateRange) Overlaps(other *DateRange) bool {
	return dr.startDate.Before(other.endDate) && dr.endDate.After(other.startDate)
}
