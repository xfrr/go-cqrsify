package postgres

import (
	"context"
	"database/sql"

	"github.com/xfrr/go-cqrsify/uow"
)

// ensure compile-time interface conformance
var _ uow.TransactionManager = (*Manager)(nil)

type Manager struct{ db *sql.DB }

func NewManager(db *sql.DB) *Manager { return &Manager{db: db} }

func (m *Manager) Begin(ctx context.Context) (uow.Tx, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &txWrapper{tx: tx}, nil
}
