package uow

import "context"

// Tx is a generic transaction. Backends wrap their native tx.
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
