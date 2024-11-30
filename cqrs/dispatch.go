package cqrs

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNilDispatcher is returned when dispatching a request with a nil dispatcher.
	ErrNilDispatcher = errors.New("the provided dispatcher is nil")
)

type EmptyRequestResponse struct{}

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
func Dispatch[Response, Request any](
	ctx context.Context,
	dispatcher Dispatcher,
	req Request,
	opts ...DispatchOption,
) (Response, error) {
	var res Response

	if dispatcher == nil {
		return res, ErrNilDispatcher
	}

	id := getIdentifier(req)
	if id == "" {
		return res, ErrBadRequest
	}

	rawResponse, err := dispatcher.Dispatch(ctx, id, req, opts...)
	if err != nil {
		return res, fmt.Errorf("dispatching request: %w", err)
	}

	if rawResponse == nil {
		return res, nil
	}

	response, ok := rawResponse.(Response)
	if !ok {
		return res, fmt.Errorf("invalid response type: %T, expected: %T", rawResponse, res)
	}

	return response, nil
}
