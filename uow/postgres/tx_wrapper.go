package postgres

import (
	"context"
	"database/sql"

	"github.com/xfrr/go-cqrsify/uow"
)

// ensure compile-time interface conformance
var _ interface {
	uow.Tx
	uow.Savepointer
} = (*txWrapper)(nil)

type txWrapper struct {
	tx *sql.Tx
}

func (t *txWrapper) Commit(_ context.Context) error   { return t.tx.Commit() }
func (t *txWrapper) Rollback(_ context.Context) error { return t.tx.Rollback() }

func (t *txWrapper) Savepoint(ctx context.Context, name string) error {
	_, err := t.tx.ExecContext(ctx, "SAVEPOINT "+name)
	return err
}
func (t *txWrapper) Release(ctx context.Context, name string) error {
	_, err := t.tx.ExecContext(ctx, "RELEASE SAVEPOINT "+name)
	return err
}
func (t *txWrapper) RollbackTo(ctx context.Context, name string) error {
	_, err := t.tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+name)
	return err
}

// Unwrap is a helper to unwrap *sql.Tx when binding tx-scoped repos
func Unwrap(tx uow.Tx) (*sql.Tx, bool) {
	tw, ok := tx.(*txWrapper)
	if !ok {
		return nil, false
	}
	return tw.tx, true
}

func BindTx[T any](tx uow.Tx) (T, bool) {
	sqltx, ok := Unwrap(tx)
	if !ok {
		var zero T
		return zero, false
	}

	return any(sqltx).(T), true
}
