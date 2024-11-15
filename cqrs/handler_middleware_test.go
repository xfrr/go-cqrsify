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
	handler := middleware(func(_ context.Context, _ interface{}) (interface{}, error) {
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
	middleware1 := func(next cqrs.HandlerFuncAny) cqrs.HandlerFuncAny {
		return func(ctx context.Context, _ interface{}) (interface{}, error) {
			return next(ctx, "middleware1")
		}
	}

	middleware2 := func(next cqrs.HandlerFuncAny) cqrs.HandlerFuncAny {
		return func(ctx context.Context, payload interface{}) (interface{}, error) {
			if payload != "middleware1" {
				return nil, errors.New("unexpected request")
			}
			return next(ctx, "middleware2")
		}
	}

	handler := func(_ context.Context, payload interface{}) (interface{}, error) {
		if payload != "middleware2" {
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
