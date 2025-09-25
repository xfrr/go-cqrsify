package uow

import "context"

// TransactionManager starts new transactions for a backend.
type TransactionManager interface {
	Begin(ctx context.Context) (Tx, error)
}
