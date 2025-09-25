package uow

import (
	"context"
	"errors"
	"fmt"
)

// Config tunes UoW behavior.
type Config struct {
	// EnableSavepoints allows nested Do() calls when the underlying Tx implements Savepointer.
	EnableSavepoints bool
}

// BindFn is a function that binds a Tx to a repository registry T.
type BindFn[T any] func(tx Tx) T

// UnitOfWork is a storage-agnostic, generic UoW that wires a repository registry T.
type UnitOfWork[T any] struct {
	mgr     TransactionManager
	cfg     Config
	depth   int // 0 => not in tx; >=1 => inside a tx
	tx      Tx
	baseR   T         // repositories bound to "no tx" access
	bind    BindFn[T] // binds a Tx to a T
	current T         // current repos (base or tx-bound)
	closed  bool
}

// New constructs a new generic UoW.
func New[T any](mgr TransactionManager, baseRepos T, bind BindFn[T], cfg Config) *UnitOfWork[T] {
	return &UnitOfWork[T]{
		mgr:     mgr,
		cfg:     cfg,
		baseR:   baseRepos,
		current: baseRepos,
		bind:    bind,
	}
}

// Repos returns the currently active repository registry.
// Outside transactions this is base repos; inside Do it's tx-bound repos.
func (u *UnitOfWork[T]) Repos() T { return u.current }

// Do executes fn within a transaction boundary (or savepoint if nested).
// It provides the tx-bound repository registry to fn.
func (u *UnitOfWork[T]) Do(ctx context.Context, fn func(ctx context.Context, repos T) error) error {
	if u.closed {
		return errors.New("uow: already closed")
	}

	if u.depth > 0 {
		return u.doNested(ctx, fn)
	}
	return u.doTransaction(ctx, fn)
}

func (u *UnitOfWork[T]) doNested(ctx context.Context, fn func(ctx context.Context, repos T) error) (err error) {
	if !u.cfg.EnableSavepoints {
		return errors.New("uow: nested transaction attempted without savepoints enabled")
	}

	sp, ok := u.tx.(Savepointer)
	if !ok {
		return errors.New("uow: nested transaction requested but backend does not support savepoints")
	}

	spName := fmt.Sprintf("sp_%d", u.depth+1)
	if err = sp.Savepoint(ctx, spName); err != nil {
		return fmt.Errorf("uow: savepoint: %w", err)
	}

	u.depth++
	defer func(prev T) {
		u.depth--
		u.current = prev // restore previous repos binding

		if p := recover(); p != nil {
			_ = sp.RollbackTo(ctx, spName)
			panic(p)
		} else if err != nil {
			_ = sp.RollbackTo(ctx, spName)
		} else {
			if rerr := sp.Release(ctx, spName); rerr != nil {
				err = fmt.Errorf("uow: release savepoint: %w", rerr)
			}
		}
	}(u.current)

	// Run with current tx-bound repos unchanged
	return fn(ctx, u.current)
}

func (u *UnitOfWork[T]) doTransaction(ctx context.Context, fn func(ctx context.Context, repos T) error) (err error) {
	tx, err := u.mgr.Begin(ctx)
	if err != nil {
		return fmt.Errorf("uow: begin: %w", err)
	}
	u.tx = tx
	u.depth = 1

	// Bind tx repositories
	prevRepos := u.current
	u.current = u.bind(tx)

	defer func() {
		u.depth = 0
		u.tx = nil
		u.current = prevRepos // restore base repos after tx ends

		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			if cerr := tx.Commit(ctx); cerr != nil {
				err = fmt.Errorf("uow: commit: %w", cerr)
			}
		}
	}()

	return fn(ctx, u.current)
}

// Close marks the UoW unusable.
func (u *UnitOfWork[T]) Close() { u.closed = true }
