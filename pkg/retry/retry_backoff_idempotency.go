package retry

import (
	"context"
	"fmt"
	"time"
)

type ctxKey string

const idemKey ctxKey = "retry.idempotency.token"

// GetIdempotencyToken returns the token stored in context (if any).
func GetIdempotencyToken(ctx context.Context) (string, bool) {
	if v := ctx.Value(idemKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s, true
		}
	}
	return "", false
}

// WithIdempotencyToken injects a token into ctx (used by helpers below).
func WithIdempotencyToken(ctx context.Context, token string) context.Context {
	if token == "" {
		return ctx
	}
	return context.WithValue(ctx, idemKey, token)
}

// DedupeStore coordinates idempotent execution across workers/instances.
// Typical implementations: Redis (SET NX PX), DynamoDB, SQL row with unique key, etc.
type DedupeStore interface {
	// Begin reserves the token with TTL. Returns proceed=true if the caller should execute the work.
	Begin(ctx context.Context, token string, ttl time.Duration) (proceed bool, err error)
	// Commit marks successful completion (optional persistence/cleanup).
	Commit(ctx context.Context, token string) error
	// Rollback releases reservation after a failed attempt (or leave for TTL expiry).
	Rollback(ctx context.Context, token string) error
}

// Helper: wrap a side-effecting function to enforce idempotency with a token.
// If Begin returns proceed=false, the function is skipped and nil is returned (treat as already-done).
func WithIdempotency(store DedupeStore, token string, ttl time.Duration, fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if token == "" {
			// No token -> proceed normally
			return fn(ctx)
		}
		ctx = WithIdempotencyToken(ctx, token)
		ok, err := store.Begin(ctx, token, ttl)
		if err != nil {
			return fmt.Errorf("idempotency begin: %w", err)
		}
		if !ok {
			// Another worker already did (or is doing) the work; treat as no-op for idempotency
			return nil
		}
		defer func() {
			// Best-effort cleanup on panic: attempt rollback; re-panic.
			if rec := recover(); rec != nil {
				_ = store.Rollback(ctx, token)
				panic(rec)
			}
		}()

		if err = fn(ctx); err != nil {
			_ = store.Rollback(ctx, token)
			return err
		}
		return store.Commit(ctx, token)
	}
}

// Result variant: returns (T, error). Returns zero T if proceed=false (already-done).
func WithIdempotencyResult[T any](store DedupeStore, token string, ttl time.Duration, fn func(ctx context.Context) (T, error)) func(ctx context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		var zero T
		if token == "" {
			return fn(ctx)
		}

		ctx = WithIdempotencyToken(ctx, token)
		ok, err := store.Begin(ctx, token, ttl)
		if err != nil {
			return zero, fmt.Errorf("idempotency begin: %w", err)
		}
		if !ok {
			return zero, nil
		}
		defer func() {
			if rec := recover(); rec != nil {
				_ = store.Rollback(ctx, token)
				panic(rec)
			}
		}()
		v, e := fn(ctx)
		if e != nil {
			_ = store.Rollback(ctx, token)
			return zero, e
		}

		if err = store.Commit(ctx, token); err != nil {
			return zero, fmt.Errorf("idempotency commit: %w", err)
		}
		return v, nil
	}
}
