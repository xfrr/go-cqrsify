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
		cmdname string,
		cmd interface{},
		opts ...DispatchOption,
	) (response interface{}, err error)
}

type Registerer interface {
	RegisterHandler(
		ctx context.Context,
		cmdname string,
		handler HandlerFuncAny,
	) error
}
