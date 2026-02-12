# AI Coding Agent Instructions for go-cqrsify

## Project Overview

go-cqrsify is a lightweight Go library implementing Domain-Driven Design (DDD), Clean Architecture, Event-Sourcing, and Event-Driven systems. The API is pre-1.0 and unstable.

## Architecture Layers

### Domain Layer (`domain/`)

- **Aggregates**: Generic interfaces `Aggregate[ID]`, `VersionedAggregate[ID]`, `EventSourcedAggregate[ID]`
  - Implement with `*BaseAggregate[T]` for ID-typed aggregates
  - All aggregates must embed or implement these interfaces
  - Example: [aggregate/main.go](../examples/aggregate/main.go) shows `CustomAggregate` pattern
- **Events**: Implement `Event` interface or use `BaseEvent`
  - Events must have unique human-readable names via `Name()` method
  - Always include aggregate reference via `EventAggregateReference`
  - Use functional options pattern: `WithEventTimestamp(time.Time)`
- **Repositories**: Repository pattern with three variants
  - `Repository[T, ID]`: Basic get/save/exists
  - `VersionedRepository[T, ID]`: Adds version support
  - `EventSourcedRepository[T, ID]`: Event-based loading/saving
  - Implementations in `domain/inmemory/` for in-memory stores
- **Event Store**: Abstraction over event persistence
  - Implement `EventSaver`, `EventRetriever[ID]`, `EventSearcher`
  - Supports version-scoped retrieval and batch operations
- **Policies**: Domain policies for complex decisions
  - `PolicyEngine[T]` manages multiple policies
  - Policies return `Result` with allow/deny decisions
  - Use composition patterns in `policy/` subpackage

### Messaging Layer (`messaging/`)

- **Command Bus**: `CommandDispatcher` + `CommandConsumer`
  - `Dispatch()` provides at-least-once semantics
  - `Subscribe()` registers handlers by command name
  - In-memory impl: `CommandBusInMemory`
- **Event Bus**: `EventPublisher` + `EventConsumer`
  - `Publish()` provides at-least-once semantics
  - Subscribe to specific event names via `MessageHandler[Event]`
  - In-memory impl: `EventBusInMemory`
- **Query Bus**: Similar pattern for query/response
- **Message Handling**: Use `MessageHandlerFn` wrapper to adapt signatures
  - Generic helpers like `EventHandlerFn[E]` for type-safe casting
  - Handlers receive `context.Context` for cancellation and deadlines

### HTTP Integration (`messaging/http/`)

- HTTP wrappers for command/query buses
- Implements adapter pattern between HTTP transport and domain messaging

### Unit of Work (`uow/`)

- Generic transactional boundary: `UnitOfWork[T]`
  - `Do(ctx, fn)` executes code within a transaction
  - Automatically binds repositories to transaction (`BindFn[T]`)
  - Supports nested transactions via savepoints (when `EnableSavepoints=true`)
  - Pattern: instantiate with transaction manager and repository factory
  - Example: `uow/postgres/` shows persistence implementation

### Error Handling (`errors/`)

- Custom error classification: `Permanent`, `Temporary`, `Retryable`
- Classify operations for retry logic
- Use `errors.Unwrap()` to inspect error chains

## Critical Patterns

### Generic Type Parameters

- All core types are generics: `Aggregate[ID]`, `Repository[T, ID]`
- Go 1.22+ required
- When implementing repositories/stores, always specify concrete ID type
- Use `CastAggregate[OutID, InID]()` to safely cast aggregates between ID types

### Functional Options Pattern

Throughout codebase for flexible configuration:

```go
// Event creation
NewEvent(name, aggRef, WithEventTimestamp(time.Now()))

// Event retrieval
RetrieveMany(ctx, aggID, RetrieveEventsFromVersion(5), RetrieveEventsBatchSize(100))

// Message buses
NewEventBusInMemory(WithMiddlewares(...))
```

### Interface Segregation

All major components are small, composable interfaces:

- `EventStore` = `EventSaver` + `EventRetriever[ID]` + `EventSearcher`
- `CommandBus` = `CommandDispatcher` + `CommandConsumer`
- Prefer taking minimal interfaces in functions

### Handler Registration by String Name

Both message buses use string-based handler registration:

```go
bus.Subscribe(ctx, MessageHandlerFn[Command](func(ctx context.Context, cmd Command) error {
    // handler
}))
```

Handlers must handle type matching; use `EventHandlerFn[E]()` wrapper for type safety.

## Development Workflow

### Testing

- Use testify: `require`, `assert` in test files
- Pattern: `Test*Suite(t *testing.T)` functions
- Run: `make test` (verbose)
- Coverage: `make test-cover` generates `coverage.out`
- View: `make cover-html` opens browser with coverage report

### Code Quality

- Linting: `make lint` (golangci-lint, 5m timeout)
- Mocking: Uses `moq` with `//go:generate moq` directives
  - Run `go generate ./...` to update mocks in `mock/` subdirs

### Building & Benchmarks

- `make bench`: Runs benchmarks on all packages
- No build targets (pure library)

### Go Module Management

- Workspace: `go.work` manages monorepo of main + examples
- Run tests from workspace root: `go test ./...`
- Examples are separate modules with own `go.mod`

## Integration Points

### External Dependencies

- `github.com/stretchr/testify`: Testing only
- `gopkg.in/yaml.v3`, `gojsonschema`: Schema validation in `pkg/apix`
- NATS example in `examples/nats/` (optional integration)
- PostgreSQL example in `uow/postgres/` (optional integration)

### Search & Criteria

- `pkg/criteria/`: Flexible query building for repository searches
- `SearchCriteriaOptions` passed to `Search()` methods
- Filters in `domain/inmemory/filters.go` show how to evaluate criteria

### Retry & Backoff

- `pkg/retry/`: Retry mechanism with backoff strategies
- Used for resilience in event-driven scenarios

## File Organization Rules

1. **Interface definitions first**: Each package starts with interface definitions
2. **Base implementations**: `Base*` types (BaseAggregate, BaseEvent) are default impl
3. **Implementations by storage**: InMemory in `inmemory/`, PostgreSQL in `postgres/`
4. **Mock generation**: `//go:generate moq` comments on interfaces
5. **Test files alongside**: `*_test.go` in same package as implementation
6. **Examples in top-level `examples/`**: Not in `testdata`

## Things to Avoid

- **Don't bypass generic constraints**: Always specify concrete ID types
- **Don't register handlers without string names**: Handlers are name-indexed
- **Don't mix cleanup**: Use contexts for cancellation, not custom shutdown
- **Don't assume event ordering outside aggregate**: Events within an aggregate are ordered, across aggregates they may not be
- **Don't manually implement interfaces if Base\* exists**: Use composition with embedded `*BaseAggregate[ID]`
- **Don't call moq-generated mocks directly**: Regenerate with `go generate ./...` after interface changes
