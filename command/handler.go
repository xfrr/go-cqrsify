package command

import (
	"context"
	"fmt"
)

// HandlerFunc is a function to handle commands.
type HandlerFunc[Payload any] func(Context[Payload]) error

// Handler wraps a Bus to provide a convenient way to subscribe to and handle commands.
type Handler[Payload any] struct {
	bus Subscriber
}

// NewHandler wraps the provided Bus in a *Handler.
func NewHandler[Payload any](s Subscriber) *Handler[Payload] {
	return &Handler[Payload]{
		bus: s,
	}
}

// Handle subscribes to the provided subject and handles the commands asynchronously with the provided handler.
func (h *Handler[Payload]) Handle(ctx context.Context, subject string, handler HandlerFunc[Payload]) (<-chan error, error) {
	if handler == nil {
		return nil, ErrNilHandler
	}

	contextCh, err := h.bus.Subscribe(ctx, subject)
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

// Handle is a shortcut for creating a new handler and subscribe it to the provided bus with given subject.
func Handle[Payload any](ctx context.Context, bus Bus, subject string, handler HandlerFunc[Payload]) (<-chan error, error) {
	return NewHandler[Payload](bus).Handle(ctx, subject, handler)
}
