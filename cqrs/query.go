package cqrs

// Query is an interface that represents a query request.
type Query interface {
	QueryName() string
}
