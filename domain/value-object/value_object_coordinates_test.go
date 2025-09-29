package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	valueobject "github.com/xfrr/go-cqrsify/domain/value-object"
)

func TestCoordinatesSuite(t *testing.T) {
	suite.Run(t, new(CoordinatesTestSuite))
}

// CoordinatesTestSuite groups coordinates-related tests
type CoordinatesTestSuite struct {
	suite.Suite
}

func (suite *CoordinatesTestSuite) TestValidCoordinates() {
	coords, err := valueobject.NewCoordinates(37.7749, -122.4194)

	suite.Require().NoError(err)
	suite.InEpsilon(37.7749, coords.Latitude(), 0.0001)
	suite.InEpsilon(-122.4194, coords.Longitude(), 0.0001)
}

func (suite *CoordinatesTestSuite) TestCoordinatesString() {
	coords, _ := valueobject.NewCoordinates(37.7749, -122.4194)
	expected := "Latitude: 37.774900, Longitude: -122.419400"

	suite.Equal(expected, coords.String())
}

func (suite *CoordinatesTestSuite) TestInvalidLatitude() {
	_, err := valueobject.NewCoordinates(100.0, -122.4194) // Invalid latitude
	suite.Require().Error(err)

	expected := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{
			{
				Field:   "latitude",
				Message: "must be between -90 and 90",
			},
		},
	}
	suite.ErrorAs(err, &expected)
}

func (suite *CoordinatesTestSuite) TestInvalidLongitude() {
	_, err := valueobject.NewCoordinates(37.7749, -200.0) // Invalid longitude
	suite.Require().Error(err)

	expected := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{
			{
				Field:   "longitude",
				Message: "must be between -180 and 180",
			},
		},
	}
	suite.ErrorAs(err, &expected)
}

func (suite *CoordinatesTestSuite) TestCoordinatesEquality() {
	coords1, _ := valueobject.NewCoordinates(37.7749, -122.4194)
	coords2, _ := valueobject.NewCoordinates(37.7749, -122.4194)
	coords3, _ := valueobject.NewCoordinates(34.0522, -118.2437)

	suite.True(coords1.Equals(coords2))
	suite.False(coords1.Equals(coords3))
	suite.False(coords1.Equals(nil))
}

func (suite *CoordinatesTestSuite) TestIsZero() {
	coords1, _ := valueobject.NewCoordinates(0, 0)
	coords2, _ := valueobject.NewCoordinates(37.7749, -122.4194)

	suite.True(coords1.IsZero())
	suite.False(coords2.IsZero())
}

func (suite *CoordinatesTestSuite) TestIsInvalidLatitudeError() {
	err := valueobject.ErrInvalidLatitude
	suite.Require().True(valueobject.IsInvalidLatitudeError(err))

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrInvalidLongitude, valueobject.ErrInvalidLatitude},
	}
	suite.Require().True(valueobject.IsInvalidLatitudeError(multiErr))

	otherErr := valueobject.ErrInvalidLongitude
	suite.Require().False(valueobject.IsInvalidLatitudeError(otherErr))
}

func (suite *CoordinatesTestSuite) TestIsInvalidLongitudeError() {
	err := valueobject.ErrInvalidLongitude
	suite.Require().True(valueobject.IsInvalidLongitudeError(err))

	multiErr := valueobject.MultiValidationError{
		Errors: []valueobject.ValidationError{valueobject.ErrInvalidLongitude, valueobject.ErrInvalidLatitude},
	}
	suite.Require().True(valueobject.IsInvalidLongitudeError(multiErr))

	otherErr := valueobject.ErrInvalidLatitude
	suite.Require().False(valueobject.IsInvalidLongitudeError(otherErr))
}
