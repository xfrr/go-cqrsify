package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xfrr/go-cqrsify/pkg/retry"
)

var defaultRetrier = retry.New(retry.Options{
	MaxAttempts: 5,
	Hooks: retry.Hooks{
		OnAttempt: func(i int) { /* metrics.Inc("attempt") */ },
		OnRetry: func(i int, err error, d time.Duration) {
			fmt.Printf("attempt %d failed: %v, retrying in %s\n", i+1, err, d)
		},
		OnGiveUp: func(i int, finalErr error, cause error) {
			fmt.Printf("attempt %d gave up: %v, cause: %v\n", i+1, finalErr, cause)
		},
	},
	Classifier: retry.RetryOn{Predicate: func(err error) bool {
		// retry only on transient network errors, 5xx, or context.DeadlineExceeded from server, etc.
		return isTransient(err)
	}},
})

func main() {
	ctx := context.Background()
	// Example 1: Default exponential + full jitter, 5 attempts, 10s max.
	// Simulate an operation that may fail and needs to be retried.
	// In this example, myOp returns an error that is considered retryable.
	// You can modify the logic in myOp and isTransient to test different scenarios.
	doDefaultRetries(ctx)

	// Example 2: Result-aware retry with backoff
	doResultAwareRetries(ctx)

	// Example 3: Batch retry with shared budget
	doBatchRetries(ctx)

	// Example 4: Retries with Idempotency Key (implemented in memory repository in this example)
	doIdempotentRetries(ctx)
}

func doDefaultRetries(ctx context.Context) {
	fmt.Println("Default Retry with Backoff")
	fmt.Println("-------------------------------------------------------")

	err := defaultRetrier.Do(ctx, func(ctx context.Context) error {
		// Simulate an operation that may fail and needs to be retried.
		// In this example, myOp returns an error that is considered retryable.
		return myOp(true, true)
	})
	if err != nil {
		// all attempts failed, handle error
		fmt.Println("> Operation failed after retries:", err)
		return
	}
	fmt.Println("> Operation succeeded")
}

func doResultAwareRetries(ctx context.Context) {
	fmt.Println()
	fmt.Println("Retry with Backoff with Result-aware strategy")
	fmt.Println("-------------------------------------------------------")

	result, err := retry.DoResult(ctx, defaultRetrier, func(ctx context.Context) (string, error) {
		err := myOp(true, true)
		if err != nil {
			return "", err
		}
		return "operation successful", nil
	})
	if err != nil {
		// all attempts failed, handle error
		fmt.Println("> Operation failed after retries:", err)
		return
	}
	fmt.Println("> Operation succeeded:", result)
}

func doBatchRetries(ctx context.Context) {
	fmt.Println()
	fmt.Println("Batch Retry with Backoff")
	fmt.Println("-------------------------------------------------------")

	tasks := []func(context.Context) error{
		func(ctx context.Context) error { return myOp(false, false) }, // should succeed
		func(ctx context.Context) error { return myOp(true, true) },   // should retry and eventually succeed
		func(ctx context.Context) error { return myOp(true, false) },  // should fail immediately
	}

	for i, task := range tasks {
		fmt.Printf("Starting task %d\n", i+1)
		err := defaultRetrier.Do(ctx, task)
		if err != nil {
			fmt.Printf("> Task %d failed after retries: %v\n", i+1, err)
		} else {
			fmt.Printf("> Task %d succeeded\n", i+1)
		}
	}
}

func doIdempotentRetries(ctx context.Context) {
	fmt.Println()
	fmt.Println("Idempotent Retry with Backoff")
	fmt.Println("-------------------------------------------------------")

	// In-memory store for idempotency keys
	// In production environments, you should use a persistent store like Redis or a database.
	dedupStore := retry.NewInMemoryDedupeStore()

	// Unique idempotency key for the operation
	dedupKey := "unique-operation-id-12345"

	dedupTTL := 5 * time.Minute // Time to keep the idempotency key

	// Wrap the operation with idempotency
	// If the operation has already been completed successfully, it will be skipped.
	// If it fails, it will be retried according to the retry policy.
	err := defaultRetrier.Do(ctx, retry.WithIdempotency(dedupStore, dedupKey, dedupTTL, func(ctx context.Context) error {
		// Simulate an operation that may fail and needs to be retried.
		return myOp(true, true)
	}))
	if err != nil {
		// all attempts failed, handle error
		fmt.Println("> Operation failed after retries:", err)
		return
	}
	fmt.Println("> Operation succeeded (idempotent)")
}

func myOp(throwError, retryable bool) error {
	fmt.Println("Doing some operation...")
	// simulate an operation that may fail
	if throwError && !retryable {
		return errors.New("permanent failure")
	} else if throwError && retryable {
		return context.DeadlineExceeded
	}
	// simulate work
	time.Sleep(100 * time.Millisecond)
	return nil
}

func isTransient(err error) bool {
	// example: treat all errors as transient except context.DeadlineExceeded
	return errors.Is(err, context.DeadlineExceeded)
}
