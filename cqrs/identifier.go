package cqrs

import "fmt"

// getIdentifier returns the unique identifier of the given interface based on its type.
func getIdentifier(req interface{}) string {
	if req == nil {
		return ""
	}

	switch v := req.(type) {
	case Command:
		return v.CommandName()
	case Query:
		return v.QueryName()
	case fmt.Stringer:
		return v.String()
	case fmt.GoStringer:
		return v.GoString()
	default:
		return fmt.Sprintf("%T", req)
	}
}
