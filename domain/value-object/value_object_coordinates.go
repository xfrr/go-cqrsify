package valueobject

import "fmt"

var _ ValueObject = (*Coordinates)(nil)

// Address value object
type Coordinates struct {
	BaseValueObject

	latitude  float64
	longitude float64
}

// NewCoordinates creates a new Coordinates value object
func NewCoordinates(latitude, longitude float64) (*Coordinates, error) {
	coords := &Coordinates{
		latitude:  latitude,
		longitude: longitude,
	}
	if err := coords.Validate(); err != nil {
		return nil, err
	}
	return coords, nil
}

func (c *Coordinates) Latitude() float64  { return c.latitude }
func (c *Coordinates) Longitude() float64 { return c.longitude }
func (c *Coordinates) Equals(vo ValueObject) bool {
	if other, ok := vo.(*Coordinates); ok {
		return c.latitude == other.latitude && c.longitude == other.longitude
	}
	return false
}

func (c *Coordinates) String() string {
	return fmt.Sprintf("Latitude: %.6f, Longitude: %.6f", c.latitude, c.longitude)
}

func (c *Coordinates) Validate() error {
	var errs []ValidationError

	if c.latitude < -90 || c.latitude > 90 {
		errs = append(errs, ValidationError{Field: "latitude", Message: "must be between -90 and 90"})
	}
	if c.longitude < -180 || c.longitude > 180 {
		errs = append(errs, ValidationError{Field: "longitude", Message: "must be between -180 and 180"})
	}

	return ValidationErrors(errs)
}
