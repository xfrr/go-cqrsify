package aggregate

// ID represents the unique identifier of an aggregate.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}
