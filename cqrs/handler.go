package cqrs

import (
	"context"
)

// HandlerFunc is a function to handle requests.
type HandlerFunc[Request any] func(ctx context.Context, req Request) (interface{}, error)

// HandlerFuncAny is a handler function that handles a request.
type HandlerFuncAny func(ctx context.Context, req interface{}) (interface{}, error)

// Handle registers the provided handler function to the bus.
//
// The handler function will be called when a request is dispatched to the bus.
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
func Handle[Request any](
	ctx context.Context,
	bus Bus,
	handler HandlerFunc[Request],
) error {
	if bus == nil {
		return ErrNilBus
	}

	if handler == nil {
		return ErrNilHandler
	}

	cmdname := getIdentifier(*new(Request))
	return bus.RegisterHandler(ctx, cmdname, wrapHandler(handler))
}

func wrapHandler[Request any](handler HandlerFunc[Request]) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(Request)
		if !ok {
			return nil, ErrInvalidRequest
		}

		return handler(ctx, cmd)
	}
}