package command

import "context"

// Bus represents a command bus that can dispatch commands and subscribe to them.
// The behaviour of the bus is defined by the implementation of the Dispatcher and Subscriber interfaces.
type Bus interface {
	Dispatcher
	Subscriber
}

// A Dispatcher dispatches commands to subscribed handlers based on the topic.
type Dispatcher interface {
	// Dispatch dispatches the provided command to the subscribers.
	// The behavior of this method depends on the implementation.
	Dispatch(ctx context.Context, topic string, cmd Command[any]) error
}

// A Subscriber subscribes to commands with a given topic.
type Subscriber interface {
	// Subscribe subscribes to the command with the provided topic.
	// The returned channels are closed when the context is canceled.
	// The behavior of this method depends on the implementation.
	Subscribe(ctx context.Context, topic string) (<-chan Context[any], error)
}
