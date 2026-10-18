# Project Overview

go-cqrsify is a Go library for building Domain-Driven Design, CQRS, event-sourced,
and event-driven systems. The primary language is Go, with Markdown docs and a small
Makefile-based developer workflow. The repository is a workspace: the root module
contains the core library, while sibling workspace modules provide examples and
optional adapters for HTTP, NATS, and PostgreSQL-backed unit-of-work integration.

## Repository Structure

- `domain/`: Core domain model abstractions. Defines aggregates, events, repositories,
  event stores, search criteria, and policy helpers. `domain/inmemory/` contains
  in-memory repository implementations used as reference adapters and test targets.
- `messaging/`: CQRS and event-driven messaging contracts plus in-memory buses for
  commands, queries, and events. `messaging/http/` adds HTTP transport adapters.
  `messaging/nats/` adds NATS JetStream and Pub/Sub adapters.
- `uow/`: Generic unit-of-work, transactions, and savepoint abstractions.
  `uow/postgres/` contains the PostgreSQL transaction manager adapter.
- `errors/`: Error classification helpers for permanent, temporary, and retryable
  failure semantics.
- `pkg/`: Shared support packages such as criteria, locking, retry, API helpers,
  and multi-error aggregation.
- `examples/`: Standalone example applications showing library usage patterns.
- `docs/`: Documentation site assets and usage guides.
- Root files: `go.mod` is the main module, `go.work` wires the workspace,
  `Makefile` contains developer commands, and `README.md` is the entry doc.

## Build & Development Commands

Run commands from the workspace root unless a task explicitly targets a nested module.

```bash
go test ./...
make test
make test-cover
make test-cover-out
make cover-html
make bench
make lint
go generate ./...
go run ./examples/<example>
```

Notes:

- Go 1.27 is declared in `go.mod` and `go.work`.
- The project is a library; there is no single root binary build target.
- `go generate ./...` is required after interface changes that affect generated mocks.

## Code Style & Conventions

- Prefer small interfaces and composition over large concrete types.
- Core APIs are generic. Keep explicit type parameters on aggregates, repositories,
  stores, and message handlers.
- Package layout follows `interface first, implementation second`. Public interfaces
  are declared near the top of each package.
- Default implementations use `Base*` types such as `BaseAggregate` and `BaseEvent`.
  Prefer embedding these instead of re-implementing common behavior.
- Functional options are used for extensibility, for example `WithEventTimestamp(...)`
  and in-memory bus configuration functions.
- Message handlers and subscriptions are keyed by string names. Event and command
  names must stay stable and human-readable.
- Context is threaded through repositories, buses, and unit-of-work boundaries.
  Preserve cancellation and deadline propagation.
- Tests use `testing`, `require`, and `assert` from `github.com/stretchr/testify`.
  Many tests call `t.Parallel()`.
- Keep comments short and behavioral. Match existing naming and file organization.

## Architecture Notes

- `domain` owns the business model: aggregates expose IDs, versions, and optional
  uncommitted events. Event-sourced aggregates restore state by replaying history.
- `messaging` transports commands, queries, and events between producers and
  handlers. In-memory buses are the default process-local implementation.
- `uow` binds repositories to a transaction boundary. Nested calls rely on
  savepoints only when the backend implements `Savepointer` and savepoints are
  enabled in config.
- `errors` and `pkg/retry` support retry-aware flows without coupling retry logic
  to a specific transport or storage backend.
- Adapters live at the edges. HTTP, NATS, and PostgreSQL packages integrate the core
  abstractions without changing domain contracts.

High-level flow:

1. Application code creates or loads an aggregate through a repository.
2. Aggregate methods record domain events and update in-memory state.
3. Repositories persist state or event history, then commit aggregate events.
4. Message buses publish commands, queries, or events to subscribed handlers.
5. Optional unit-of-work boundaries coordinate repository access in a transaction.

## Testing Strategy

- Primary test scope is package-level unit testing with Go's `testing` package.
- Use `make test` or `go test ./...` for the full workspace.
- Use `make test-cover` and `make test-cover-out` for coverage metrics.
- Use `make cover-html` to generate and open HTML coverage output.
- Use `make bench` for benchmark runs.
- When changing public interfaces in `messaging`, regenerate mocks with
  `go generate ./...` before running tests.
- Keep new tests adjacent to the implementation in `*_test.go` files.

## Security & Compliance

- Do not hardcode credentials, tokens, or DSNs in source, examples, or docs.
- No environment-based secret loading was detected in the core packages. If new
  secret handling is introduced, keep it outside the domain layer and document it.
- Preserve error classification semantics in `errors/`; retry behavior can depend on
  whether an error is permanent, temporary, or retryable.
- Preserve aggregate and event integrity: events should keep a valid aggregate
  reference, a stable name, and consistent version ordering.
- Do not bypass context propagation in handlers, repositories, or transactions.
- Validate external payloads at adapter boundaries, not inside pure domain types.

## Agent Guardrails

- Never edit `LICENSE`.
- Never hand-edit files under `messaging/mock/`; regenerate them with
  `go generate ./...` after interface changes.
- Never fix library behavior only in `examples/`; correct the owning package first,
  then update examples if necessary.
- Never change `go.work` or nested module manifests unless the task explicitly
  requires workspace or dependency boundary changes.
- Never introduce transport-specific concerns into `domain/`.

## Extensibility Hooks

- Add new aggregates by embedding `*domain.BaseAggregate[ID]` and registering event
  handlers with `HandleEvent`.
- Add new events by implementing `domain.Event` or using `BaseEvent` plus event
  options such as `WithEventTimestamp`.
- Add repositories by implementing the smallest needed interface:
  `Repository`, `VersionedRepository`, or `EventSourcedRepository`.
- Add transports by implementing the messaging interfaces in `messaging/` instead of
  coupling application code to a concrete bus.
- Add transaction backends by implementing `uow.TransactionManager`, `uow.Tx`, and
  optionally `uow.Savepointer`.
- Extend search and filtering through `pkg/criteria` and the domain search helpers.

## Further Reading

- [README](./README.md)
- [Docs](./docs)
- [Usage Guide](./docs/usage.md)
- [Aggregate Example](./examples/aggregate)
- [HTTP Messaging Module](./messaging/http)
- [NATS Messaging Module](./messaging/nats)
- [PostgreSQL Unit of Work](./uow/postgres)
