package cqrs

import "context"

// DispatchOption is a function that modifies the context and the request before dispatching it.
type DispatchOption func(ctx context.Context, req interface{}) context.Context

// Dispatch generates a unique identifier for the request and dispatches it to the bus.
//
// The identifier is the unique name of the request that the handler will be registered for.
//
// - If the payload is a cqrs.Command, the identifier will be obtained from the CommandName() method.
//
// - If the payload is a struct, the identifier will be the name of the struct.
//
// - If the payload is a fmt.Stringer, the identifier will be the string representation of the request.
//
// - If the payload is a fmt.GoStringer, the identifier will be the Go string representation of the request.
//
// - If the payload is another type, the identifier will be the type name of the request (fmt.Sprintf("%T", request)).
//
// The request is dispatched to the bus and the response is returned.
func Dispatch[Request any](
	ctx context.Context,
	bus Bus,
	req Request,
	opts ...DispatchOption,
) (interface{}, error) {
	if bus == nil {
		return nil, ErrNilBus
	}

	id := getIdentifier(req)
	if id == "" {
		return nil, ErrInvalidRequest
	}

	return bus.Dispatch(ctx, id, req, opts...)
}
