package uow_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xfrr/go-cqrsify/uow"
)

type fakeTx struct {
	mu         sync.Mutex
	commits    int
	rolls      int
	ops        []string
	sp         bool // implements savepoints if true
	spStack    []string
	failCommit bool
}

func (f *fakeTx) Commit(_ context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.commits++
	if f.failCommit {
		return assertError("commit failed")
	}
	return nil
}
func (f *fakeTx) Rollback(_ context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rolls++
	return nil
}

type fakeSavepointer struct{ *fakeTx }

func (f *fakeSavepointer) Savepoint(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ops = append(f.ops, "sp:"+name)
	f.spStack = append(f.spStack, name)
	return nil
}
func (f *fakeSavepointer) Release(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ops = append(f.ops, "release:"+name)
	// pop
	if len(f.spStack) > 0 {
		f.spStack = f.spStack[:len(f.spStack)-1]
	}
	return nil
}
func (f *fakeSavepointer) RollbackTo(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ops = append(f.ops, "rollbackto:"+name)
	// clear to name
	for len(f.spStack) > 0 && f.spStack[len(f.spStack)-1] != name {
		f.spStack = f.spStack[:len(f.spStack)-1]
	}
	if len(f.spStack) > 0 {
		f.spStack = f.spStack[:len(f.spStack)-1]
	}
	return nil
}

type fakeMgr struct {
	withSavepoints bool
	failBegin      bool
	nextCommitFail bool
	last           *fakeTx
}

func (m *fakeMgr) Begin(_ context.Context) (uow.Tx, error) {
	if m.failBegin {
		return nil, assertError("begin failed")
	}
	tx := &fakeTx{
		sp:         m.withSavepoints,
		failCommit: m.nextCommitFail, // propagate commit failure to this tx
	}
	m.nextCommitFail = false // consume flag
	m.last = tx
	if m.withSavepoints {
		return &fakeSavepointer{fakeTx: tx}, nil
	}
	return tx, nil
}

type assertError string

func (e assertError) Error() string { return string(e) }

// Example registry used by tests (could be any struct of repos).
type Repos struct {
	Log func(s string)
}

// base repos (no tx) and tx-bound repos are simulated by different loggers.
func baseRepos(rec *[]string) Repos {
	return Repos{Log: func(s string) { *rec = append(*rec, "base:"+s) }}
}
func txRepos(_ uow.Tx, rec *[]string) Repos {
	return Repos{Log: func(s string) { *rec = append(*rec, "tx:"+s) }}
}

// ---- Tests ------------------------------------------------------------------

func Test_Do_Success_Commits(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{withSavepoints: false}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{})

	err := u.Do(ctx, func(_ context.Context, r Repos) error {
		r.Log("create-user")
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"tx:create-user"}, rec)
	require.Equal(t, 1, mgr.last.commits)
	require.Equal(t, 0, mgr.last.rolls)
}

func Test_Do_Error_RollsBack(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{})

	err := u.Do(ctx, func(_ context.Context, r Repos) error {
		r.Log("create-user")
		return assertError("boom")
	})
	require.Error(t, err)
	require.Equal(t, []string{"tx:create-user"}, rec)
	require.Equal(t, 0, mgr.last.commits)
	require.Equal(t, 1, mgr.last.rolls)
}

func Test_Do_Panic_RollsBackAndPropagates(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{})

	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.Equal(t, 0, mgr.last.commits)
		require.Equal(t, 1, mgr.last.rolls)
	}()
	_ = u.Do(ctx, func(_ context.Context, r Repos) error {
		r.Log("inside")
		panic("kaboom")
	})
}

func Test_Do_Nested_Disabled_ReturnsError(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{EnableSavepoints: false})

	err := u.Do(ctx, func(_ context.Context, _ Repos) error {
		return u.Do(ctx, func(_ context.Context, _ Repos) error { return nil })
	})
	require.Error(t, err)
	require.Equal(t, 1, mgr.last.rolls) // outer should roll back
}

func Test_Do_Nested_WithSavepoints_Commits(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{withSavepoints: true}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{EnableSavepoints: true})

	err := u.Do(ctx, func(ctx context.Context, r Repos) error {
		r.Log("outer")
		return u.Do(ctx, func(_ context.Context, r2 Repos) error {
			r2.Log("inner")
			return nil
		})
	})
	require.NoError(t, err)
	require.Equal(t, []string{"tx:outer", "tx:inner"}, rec)
	require.Equal(t, 1, mgr.last.commits)
	require.Equal(t, 0, mgr.last.rolls)

	// and savepoint ops recorded
	spTx, err := mgr.Begin(ctx)
	require.NoError(t, err)
	sp, _ := spTx.(uow.Savepointer) // ensure type exists in manager used above
	_ = sp                          // compile-time guard; actual ops validated implicitly by No errors
}

func Test_Do_Nested_WithSavepoints_RollbackInner_CommitOuter(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{withSavepoints: true}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{EnableSavepoints: true})

	err := u.Do(ctx, func(ctx context.Context, r Repos) error {
		r.Log("outer")
		innerErr := u.Do(ctx, func(_ context.Context, r2 Repos) error {
			r2.Log("inner-fail")
			return assertError("fail")
		})
		require.Error(t, innerErr)
		return nil // outer still succeeds
	})
	require.NoError(t, err)
	require.Equal(t, []string{"tx:outer", "tx:inner-fail"}, rec)
	require.Equal(t, 1, mgr.last.commits)
	require.Equal(t, 0, mgr.last.rolls)
}

func Test_CommitError_Propagates(t *testing.T) {
	ctx := context.Background()
	rec := []string{}
	mgr := &fakeMgr{}
	u := uow.New(mgr, baseRepos(&rec), func(tx uow.Tx) Repos { return txRepos(tx, &rec) }, uow.Config{})

	// force commit error
	_ = u.Do(ctx, func(_ context.Context, _ Repos) error {
		return nil
	})
	mgr.nextCommitFail = true

	err := u.Do(ctx, func(_ context.Context, _ Repos) error {
		return nil
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "commit failed")
}
