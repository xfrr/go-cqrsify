package cqrs

import "fmt"

// Query is an interface that represents a query request.
type Query interface {
	QueryName() string
}

// getQueryIdentifier returns the unique identifier of the given interface.
func getQueryIdentifier(anyQuery interface{}) string {
	switch queryType := anyQuery.(type) {
	case nil:
		return ""
	case Query:
		return queryType.QueryName()
	case fmt.Stringer:
		return queryType.String()
	case fmt.GoStringer:
		return queryType.GoString()
	default:
		return fmt.Sprintf("%T", queryType)
	}
}
