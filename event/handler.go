package event

import (
	"context"
	"fmt"
)

// HandlerFunc is a function to handle events.
type HandlerFunc[Payload any] func(Context[Payload]) error

// Handler wraps a Bus to provide a convenient way to subscribe to and handle events.
type Handler[Payload any] struct {
	bus Subscriber
}

// NewHandler wraps the provided Bus in a *Handler.
func NewHandler[Payload any](s Subscriber) *Handler[Payload] {
	return &Handler[Payload]{
		bus: s,
	}
}

// Handle subscribes to the provided reason and handles the events asynchronously with the provided handler.
func (h *Handler[Payload]) Handle(ctx context.Context, reason string, handler HandlerFunc[Payload]) (<-chan error, error) {
	if handler == nil {
		return nil, ErrNilHandler
	}

	contextCh, err := h.bus.Subscribe(ctx, reason)
	if err != nil {
		return nil, ErrSubscribeFailed{}.Wrap(err)
	}

	errs := make(chan error)
	go h.handle(ctx, handler, contextCh, errs)

	return errs, nil
}

func (h *Handler[Payload]) handle(ctx context.Context, handlefn HandlerFunc[Payload], contextCh <-chan anyContext, errs chan<- error) {
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

		casted, ok := CastContext[Payload](cctx)
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

// Subscribe is a shortcut for creating a new handler and subscribe it to the provided bus with given reason.
func Subscribe[Payload any](ctx context.Context, bus Bus, reason string, handler HandlerFunc[Payload]) (<-chan error, error) {
	return NewHandler[Payload](bus).Handle(ctx, reason, handler)
}
