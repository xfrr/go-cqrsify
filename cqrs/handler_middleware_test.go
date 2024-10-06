package cqrs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

func TestRecoverMiddleware(t *testing.T) {
	var recovered interface{}
	hook := func(r interface{}) {
		recovered = r
	}

	middleware := cqrs.RecoverMiddleware(hook)
	handler := middleware(func(ctx context.Context, cmd interface{}) (interface{}, error) {
		panic("test panic")
	})

	_, err := handler(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if recovered != "test panic" {
		t.Fatalf("expected recovered to be 'test panic', got %v", recovered)
	}
}

func TestChainMiddleware(t *testing.T) {
	middleware1 := func(next func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
		return func(ctx context.Context, cmd interface{}) (interface{}, error) {
			return next(ctx, "middleware1")
		}
	}

	middleware2 := func(next func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
		return func(ctx context.Context, cmd interface{}) (interface{}, error) {
			if cmd != "middleware1" {
				return nil, errors.New("unexpected request")
			}
			return next(ctx, "middleware2")
		}
	}

	handler := func(ctx context.Context, cmd interface{}) (interface{}, error) {
		if cmd != "middleware2" {
			return nil, errors.New("unexpected request")
		}
		return "success", nil
	}

	chain := cqrs.ChainMiddleware(middleware1, middleware2)
	chainedHandler := chain(handler)

	result, err := chainedHandler(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "success" {
		t.Fatalf("expected result to be 'success', got %v", result)
	}
}
