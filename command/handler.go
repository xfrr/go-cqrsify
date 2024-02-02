package command

import (
	"context"
	"fmt"
)

// HandlerFunc is a function to handle commands.
type HandlerFunc[M any] func(Context[M]) error

// Handler wraps a Bus to provide a convenient way to subscribe to and handle commands.
type Handler[M any] struct {
	bus Subscriber
}

// NewHandler wraps the provided Bus in a *Handler.
func NewHandler[M any](s Subscriber) *Handler[M] {
	return &Handler[M]{
		bus: s,
	}
}

// Handle subscribes to the provided topic and handles the commands asynchronously with the provided handler.
func (h *Handler[M]) Handle(ctx context.Context, topic string, handler HandlerFunc[M]) (<-chan error, error) {
	if handler == nil {
		return nil, ErrNilHandler
	}

	contextCh, err := h.bus.Subscribe(ctx, topic)
	if err != nil {
		return nil, ErrSubscribeFailed{}.Wrap(err)
	}

	errs := make(chan error)
	go h.handle(ctx, handler, contextCh, errs)

	return errs, nil
}

func (h *Handler[P]) handle(ctx context.Context, handlefn HandlerFunc[P], contextCh <-chan anyContext, errs chan<- error) {
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

		casted, ok := CastContext[P](cctx)
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

// Handle is a shortcut for creating a new handler and subscribe it to the provided bus with given topic.
func Handle[M any](ctx context.Context, bus Bus, topic string, handler HandlerFunc[M]) (<-chan error, error) {
	return NewHandler[M](bus).Handle(ctx, topic, handler)
}
