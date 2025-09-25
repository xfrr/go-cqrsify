package uow

import "context"

// Savepointer is an optional capability for nested transactions.
type Savepointer interface {
	Savepoint(ctx context.Context, name string) error
	Release(ctx context.Context, name string) error
	RollbackTo(ctx context.Context, name string) error
}
