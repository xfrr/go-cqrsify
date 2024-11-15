package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

type MockCommandPayload struct {
	Greeting string
}

func TestHandle(t *testing.T) {
	t.Run("should return an error when bus is nil", func(t *testing.T) {
		err := cqrs.Handle(context.Background(), nil, func(_ context.Context, _ MockCommandPayload) (interface{}, error) {
			return nil, nil
		})
		if !errors.Is(err, cqrs.ErrNilBus) {
			t.Fatalf("expected error to be %v, got %v", cqrs.ErrNilBus, err)
		}
	})

	t.Run("should return an error when handler is nil", func(t *testing.T) {
		mockBus := &mockBus{
			register: func(_ context.Context, _ string, _ cqrs.HandlerFuncAny) error {
				return nil
			},
		}

		err := cqrs.Handle[MockCommandPayload, any](context.Background(), mockBus, nil)
		if !errors.Is(err, cqrs.ErrNilHandler) {
			t.Fatalf("expected error to be %v, got %v", cqrs.ErrNilHandler, err)
		}
	})

	t.Run("should return an error when command is not valid", func(t *testing.T) {
		mockBus := &mockBus{
			register: func(_ context.Context, _ string, handler cqrs.HandlerFuncAny) error {
				_, err := handler(context.Background(), nil)
				return err
			},
		}

		err := cqrs.Handle(
			context.Background(),
			mockBus,
			func(_ context.Context, _ MockCommandPayload) (interface{}, error) {
				return nil, nil
			})
		if !errors.Is(err, cqrs.ErrInvalidRequest) {
			t.Fatalf("expected error to be %v, got %v", cqrs.ErrInvalidRequest, err)
		}

		if len(mockBus.registerCalls) != 1 {
			t.Fatalf("expected registerCalls to be %d, got %d", 1, len(mockBus.registerCalls))
		}
	})

	t.Run("should return an error when registering a handler fails", func(t *testing.T) {
		mockErr := errors.New("something went wrong")
		mockBus := &mockBus{
			register: func(_ context.Context, _ string, _ cqrs.HandlerFuncAny) error {
				return mockErr
			},
		}

		err := cqrs.Handle(
			context.Background(),
			mockBus,
			func(_ context.Context, _ MockCommandPayload) (interface{}, error) {
				return nil, nil
			})
		if !errors.Is(err, mockErr) {
			t.Fatalf("expected error to be %v, got %v", mockErr, err)
		}

		if len(mockBus.registerCalls) != 1 {
			t.Fatalf("expected registerCalls to be %d, got %d", 1, len(mockBus.registerCalls))
		}
	})

	t.Run("should handle a command without errors", func(t *testing.T) {
		mockBus := &mockBus{
			register: func(_ context.Context, _ string, handler cqrs.HandlerFuncAny) error {
				_, err := handler(context.Background(), MockCommandPayload{})
				return err
			},
		}

		err := cqrs.Handle(
			context.Background(),
			mockBus,
			func(_ context.Context, _ MockCommandPayload) (interface{}, error) {
				return nil, nil
			})
		if err != nil {
			t.Fatalf("expected error to be nil, got %v", err)
		}

		if len(mockBus.registerCalls) != 1 {
			t.Fatalf("expected registerCalls to be %d, got %d", 1, len(mockBus.registerCalls))
		}
	})
}
