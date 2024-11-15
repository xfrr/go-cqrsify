package cqrs

import "context"

// Bus exposes the API for registering and dispatching requests to their respective handlers.
type Bus interface {
	Dispatcher
	Registerer
}

type Dispatcher interface {
	Dispatch(
		ctx context.Context,
		name string,
		payload interface{},
		opts ...DispatchOption,
	) (response interface{}, err error)
}

type Registerer interface {
	RegisterHandler(
		ctx context.Context,
		name string,
		handler HandlerFuncAny,
	) error
}
