package cqrs_test

import (
	"context"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

type cmdStringer struct{}

func (s cmdStringer) String() string {
	return "cmdStringer"
}

type cmdGoStringer struct{}

func (s cmdGoStringer) GoString() string {
	return "cmdGoStringer"
}

type cmd struct{}

func (s cmd) CommandName() string {
	return "CommandName"
}

func BenchmarkCommandDispatch(b *testing.B) {
	b.Run("CommandDispatchString", func(b *testing.B) {
		cmd := "cmd"
		benchmarkCommandDispatch(b, cmd)
	})

	b.Run("CommandDispatchInt", func(b *testing.B) {
		cmd := 1
		benchmarkCommandDispatch(b, cmd)
	})

	b.Run("CommandDispatchStruct", func(b *testing.B) {
		cmd := struct{}{}
		benchmarkCommandDispatch(b, cmd)
	})

	b.Run("CommandDispatchStringer", func(b *testing.B) {
		cmd := cmdStringer{}
		benchmarkCommandDispatch(b, cmd)
	})

	b.Run("CommandDispatchGoStringer", func(b *testing.B) {
		cmd := cmdGoStringer{}
		benchmarkCommandDispatch(b, cmd)
	})

	b.Run("CommandDispatchCommand", func(b *testing.B) {
		cmd := cmd{}
		benchmarkCommandDispatch(b, cmd)
	})
}

func benchmarkCommandDispatch(b *testing.B, cmd interface{}) {
	b.Helper()

	ctx := context.Background()

	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, nil
	}

	bus := &mockBus{
		dispatch: func(ctx context.Context, cmdname string, cmd interface{}, opts ...cqrs.DispatchOption) (response interface{}, err error) {
			return handler(ctx, cmd)
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res, err := cqrs.Dispatch[any](ctx, bus, cmd)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}

		if res != nil {
			b.Fatalf("expected response to be nil, got %v", res)
		}
	}
}
