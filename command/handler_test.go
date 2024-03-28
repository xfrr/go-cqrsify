package command_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/xfrr/go-cqrsify/command"
)

type mockBus struct {
	lock sync.RWMutex

	dispatchCalls int
	dispatchFn    func(ctx context.Context, subject string, cmd command.Command[any]) error

	subscribeCalls int
	subscribeFn    func(ctx context.Context, subject string) (<-chan command.Context[any], error)
}

func (m *mockBus) Dispatch(ctx context.Context, subject string, cmd command.Command[any]) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.dispatchCalls++
	if m.dispatchFn != nil {
		return m.dispatchFn(ctx, subject, cmd)
	}

	return nil
}

func (m *mockBus) Subscribe(ctx context.Context, subject string) (<-chan command.Context[any], error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.subscribeCalls++
	if m.subscribeFn != nil {
		return m.subscribeFn(ctx, subject)
	}

	return nil, nil
}

func TestHandler(t *testing.T) {
	var (
		mockSubscriber = &mockBus{}
	)

	t.Run("New", func(t *testing.T) {
		type args struct {
			name string
			fn   command.HandlerFunc[any]
		}

		cases := []struct {
			name string
			args args
		}{
			{
				name: "should create a new handler",
				args: args{
					name: "name",
					fn: func(ctx command.Context[any]) error {
						return nil
					},
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				h := command.NewHandler[any](mockSubscriber)
				if h == nil {
					t.Error("expected handler to not be nil")
				}
			})
		}
	})
}

func TestHandle(t *testing.T) {
	var (
		mockSubject = "subject"
	)

	t.Run("should return an error when handler is nil", func(t *testing.T) {
		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return make(<-chan command.Context[any]), nil
			},
		}

		_, err := command.Handle[MockCommandPayload](context.Background(), mockSubscriber, "subject", nil)
		if err == nil || !errors.Is(err, command.ErrNilHandler) {
			t.Fatalf("expected error to be %v, got %v", command.ErrNilHandler, err)
		}
	})

	t.Run("should return an error when subscribe fails", func(t *testing.T) {
		mockErr := errors.New("something went wrong")
		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return nil, mockErr
			},
		}

		_, err := command.Handle[MockCommandPayload](context.Background(), mockSubscriber, "subject", func(ctx command.Context[MockCommandPayload]) error {
			return nil
		})
		if err, ok := err.(command.ErrSubscribeFailed); !ok {
			t.Fatalf("expected error to be %v, got %v", command.ErrSubscribeFailed{}, err)
		}

		expected := command.ErrSubscribeFailed{}.Wrap(mockErr)
		if err.Error() != expected.Error() {
			t.Fatalf("expected error to be %v, got %v", expected.Error(), err.Error())
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped == nil || unwrapped.Error() != mockErr.Error() {
			t.Fatalf("expected error to be %v, got %v", mockErr, unwrapped)
		}
	})

	t.Run("should return an error when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return make(<-chan command.Context[any]), nil
			},
		}

		errs, err := command.Handle[MockCommandPayload](ctx, mockSubscriber, "subject", func(ctx command.Context[MockCommandPayload]) error {
			return nil
		})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		select {
		case err, ok := <-errs:
			if !ok {
				t.Fatal("expected errors to be open")
			}

			if !errors.Is(err, context.Canceled) {
				t.Fatalf("expected error to be %v, got %v", context.Canceled, err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("expected context to be canceled")
		}
	})

	t.Run("should return an error when casting context fails", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ch := make(chan command.Context[any])

		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return ch, nil
			},
		}

		errs, err := command.Handle[MockCommandPayload](ctx, mockSubscriber, "subject",
			func(ctx command.Context[MockCommandPayload]) error {
				return nil
			})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		// dispatch invalid command context
		cctx := command.WithContext(ctx, command.New("id", "invalid").Any())
		ch <- cctx

		defer cancel()
		// wait for context to be handled
		select {
		case err, ok := <-errs:
			if !ok {
				t.Fatal("expected errors to be open")
			}

			if errors.Unwrap(err) != command.ErrCastContext {
				t.Fatalf("expected error to be %v, got %v", command.ErrCastContext, err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("expected context to be canceled")
		}
	})

	t.Run("should return an error when handling command fails", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ch := make(chan command.Context[any])
		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return ch, nil
			},
		}

		errs, err := command.Handle[MockCommandPayload](ctx, mockSubscriber, mockSubject,
			func(ctx command.Context[MockCommandPayload]) error {
				return errors.New("handler failed")
			})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		cmd := command.New("id", MockCommandPayload{
			Greeting: "hello",
		})

		// dispatch command context
		ch <- command.WithContext(ctx, cmd.Any())

		defer cancel()
		// wait for command to be handled
		select {
		case err, ok := <-errs:
			if !ok {
				t.Fatal("expected errors to be open")
			}

			if err == nil || err.Error() != "handler failed" {
				t.Fatalf("expected error to be %v, got %v", "handler failed", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("expected command to be handled")
		}
	})

	t.Run("should stop handling commands when context channel is closed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ch := make(chan command.Context[any])
		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return ch, nil
			},
		}

		errs, err := command.Handle[MockCommandPayload](ctx, mockSubscriber, mockSubject,
			func(ctx command.Context[MockCommandPayload]) error {
				return nil
			})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		// close command context channel
		close(ch)

		defer cancel()
		// wait for command to be handled
		select {
		case err, ok := <-errs:
			if ok {
				t.Fatalf("expected errors to be closed, got %v", err)
			}

			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("expected command to be handled")
		}
	})

	t.Run("should handle a command without errors", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ch := make(chan command.Context[any])
		mockSubscriber := &mockBus{
			subscribeFn: func(ctx context.Context, subject string) (<-chan command.Context[any], error) {
				return ch, nil
			},
		}

		handled := make(chan struct{})
		errs, err := command.Handle[MockCommandPayload](
			ctx, mockSubscriber, mockSubject,
			func(ctx command.Context[MockCommandPayload]) error {
				close(handled)
				return nil
			})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		cmd := command.New("id", MockCommandPayload{
			Greeting: "hello",
		})

		// dispatch command context
		ch <- command.WithContext(ctx, cmd.Any())

		defer cancel()
		// wait for command to be handled
		select {
		case <-handled:
		case err, ok := <-errs:
			if !ok {
				t.Fatal("expected errors to be open")
			}
			if err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("expected command to be handled")
		}

	})
}
