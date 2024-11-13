package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

func TestNewBus(t *testing.T) {
	b := cqrs.NewInMemoryBus()
	if b == nil {
		t.Fatal("expected non-nil bus")
	}
}

func TestBus_RegisterHandler(t *testing.T) {
	ctx := context.Background()

	b := cqrs.NewInMemoryBus()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	err := b.RegisterHandler(ctx, "testHandler", handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = b.RegisterHandler(ctx, "testHandler", handler)
	if !errors.Is(err, cqrs.ErrHandlerAlreadyRegistered) {
		t.Fatalf("expected ErrHandlerAlreadyRegistered, got %v", err)
	}
}

func TestBus_UnregisterHandler(t *testing.T) {
	ctx := context.Background()
	b := cqrs.NewInMemoryBus()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	err := b.RegisterHandler(ctx, "testHandler", handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	b.UnregisterHandler(ctx, "testHandler")
	if b.Exists("testHandler") {
		t.Fatal("expected handler to be unregistered")
	}
}

func TestBus_Dispatch(t *testing.T) {
	ctx := context.Background()

	b := cqrs.NewInMemoryBus()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		header, ok := cqrs.HeaderFromContext(ctx)
		if !ok {
			t.Fatal("expected header to be found")
		}

		if header["key"] != "value" {
			t.Fatalf("expected header key to be value, got %v", header["key"])
		}

		return "response", nil
	}

	err := b.RegisterHandler(ctx, "testHandler", handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var res interface{}

	res, err = b.Dispatch(ctx, "testHandler", "request", cqrs.WithHeader("key", "value"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res != "response" {
		t.Fatalf("expected response to be response, got %v", res)
	}

	res, err = b.Dispatch(ctx, "nonExistentHandler", "request", cqrs.WithHeader("key", "value"))
	if !errors.Is(err, cqrs.ErrHandlerNotFound) {
		t.Fatalf("expected ErrHandlerNotFound, got %v", err)
	}
	if res != nil {
		t.Fatalf("expected response to be nil, got %v", res)
	}
}

func TestBus_Close(t *testing.T) {
	ctx := context.Background()

	b := cqrs.NewInMemoryBus()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	err := b.RegisterHandler(ctx, "testHandler", handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	b.Close()
	if b.Exists("testHandler") {
		t.Fatal("expected all handlers to be unregistered")
	}
}

func TestBus_Use(t *testing.T) {
	ctx := context.Background()

	middleware := func(next cqrs.HandlerFuncAny) cqrs.HandlerFuncAny {
		return func(ctx context.Context, cmd interface{}) (interface{}, error) {
			return next(ctx, cmd)
		}
	}

	b := cqrs.NewInMemoryBus()
	b.Use(middleware)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	err := b.RegisterHandler(ctx, "testHandler", handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	res, err := b.Dispatch(ctx, "testHandler", "request")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res != "response" {
		t.Fatalf("expected response to be response, got %v", res)
	}
}
