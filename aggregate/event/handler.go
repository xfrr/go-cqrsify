package event

import (
	"context"
	"fmt"
)

// HandlerFunc is a function to handle events.
type HandlerFunc[ID comparable, Payload any] func(Context[ID, Payload]) error

// Handler wraps a Bus to provide a convenient way to subscribe to and handle events.
type Handler[ID comparable, Payload any] struct {
	bus Subscriber
}

// NewHandler wraps the provided Bus in a *Handler.
func NewHandler[ID comparable, Payload any](s Subscriber) *Handler[ID, Payload] {
	return &Handler[ID, Payload]{bus: s}
}

// Handle subscribes to the provided name and handles the events asynchronously with the provided handler.
func (h *Handler[ID, Payload]) Handle(ctx context.Context, name string, handler HandlerFunc[ID, Payload]) (<-chan error, error) {
	if handler == nil {
		return nil, ErrNilHandler
	}

	contextCh, err := h.bus.Subscribe(ctx, name)
	if err != nil {
		return nil, ErrSubscribeFailed{}.Wrap(err)
	}

	errs := make(chan error)
	go h.handle(ctx, handler, contextCh, errs)

	return errs, nil
}

func (h *Handler[ID, Payload]) handle(ctx context.Context, handlefn HandlerFunc[ID, Payload], contextCh <-chan Context[any, any], errs chan<- error) {
	defer close(errs)

	for {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return
		default:
		}

		cctx, ok := <-contextCh
		if !ok {
			break
		}

		casted, ok := CastContext[ID, Payload](cctx)
		if !ok {
			errs <- fmt.Errorf("%w [from=%T, to=%T]", ErrCastContext, cctx, casted)
			continue
		}

		err := handlefn(casted)
		if err != nil {
			errs <- err
			continue
		}
	}
}

// Subscribe is a shortcut for creating a new handler and subscribe it to the provided bus with given name.
func Subscribe[ID comparable, Payload any](ctx context.Context, bus Bus, name string, handler HandlerFunc[ID, Payload]) (<-chan error, error) {
	return NewHandler[ID, Payload](bus).Handle(ctx, name, handler)
}
