package cqrs_test

import (
	"context"
	"testing"

	"github.com/xfrr/go-cqrsify/cqrs"
)

func TestWithHeader(t *testing.T) {
	ctx := context.Background()
	key := "testKey"
	value := "testValue"

	option := cqrs.WithHeader(key, value)
	newCtx := option(ctx, nil)

	header, ok := cqrs.HeaderFromContext(newCtx)
	if !ok {
		t.Fatal("expected header to be present in context")
	}

	if header[key] != value {
		t.Fatalf("expected header[%s] to be %v, got %v", key, value, header[key])
	}
}

func TestHeaderFromContext(t *testing.T) {
	ctx := context.Background()
	key := "testKey"
	value := "testValue"

	ctx = cqrs.WithHeader(key, value)(ctx, nil)
	retrievedHeader, ok := cqrs.HeaderFromContext(ctx)
	if !ok {
		t.Fatal("expected header to be present in context")
	}

	if retrievedHeader[key] != value {
		t.Fatalf("expected header[%s] to be %v, got %v", key, value, retrievedHeader[key])
	}
}

func TestHeaderFromContext_NoHeader(t *testing.T) {
	ctx := context.Background()

	_, ok := cqrs.HeaderFromContext(ctx)
	if ok {
		t.Fatal("expected no header to be present in context")
	}
}
