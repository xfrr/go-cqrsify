# Retry Backoff Package

This package provides a simple and flexible way to implement retry logic with customizable backoff strategies in Go applications. It is designed to help developers handle transient errors and improve the resilience of their applications.

## Features

- **Flexible Backoff Strategies**
  - Constant, Exponential, and Decorrelated backoff
  - **Jitter policies**: None, Full, Equal, Decorrelated
  - **HTTP 429/503 Retry-After support** via `StrategyWithHint` (honors server hints while capping with your policy)

- **Result-Aware Retries**
  - `DoResult[T]` returns typed results (`func(ctx) (T, error)`) with retry orchestration

- **Pluggable Stop Policies**
  - `Stopper` interface to short-circuit retries (e.g., circuit breakers, maintenance windows)
  - **TokenBucketStopper** to halt retries under cluster backpressure or latency SLO breaches

- **Idempotency Helpers**
  - `WithIdempotency` and `WithIdempotencyResult` wrappers ensure at-least-once safety with dedupe tokens
  - In-memory store provided; easily swap for Redis, DynamoDB, SQL, etc.

- **Batch Orchestration**
  - `DoN` helper to retry a batch of items with shared budget/deadline
  - Configurable concurrency, fail-fast mode, and error aggregation

- **Observability Hooks**
  - Lifecycle callbacks: `OnAttempt`, `OnRetry`, `OnGiveUp`
  - Integrate seamlessly with Prometheus, OpenTelemetry, or logging

- **Context-Aware & Testable**
  - Full `context.Context` support for cancellation and deadlines
  - `Sleeper` and RNG sources are injectable for deterministic testing

- **Production-Grade Design**
  - SOLID & Clean Architecture principles
  - Immutable, concurrency-safe `Retrier`
  - Extensible interfaces for strategies, jitter, classifiers, and stoppers
