package lock

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryLockerRenewExistingKey(t *testing.T) {
	t.Parallel()

	l := NewInMemoryLocker()
	ctx := context.Background()

	locked, err := l.TryLock(ctx, "k", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLock returned error: %v", err)
	}
	if !locked {
		t.Fatalf("TryLock did not lock key")
	}

	renewed, err := l.Renew(ctx, "k", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Renew returned error: %v", err)
	}
	if !renewed {
		t.Fatalf("Renew returned false for existing key")
	}
}

func TestInMemoryLockerRenewMissingKey(t *testing.T) {
	t.Parallel()

	l := NewInMemoryLocker()
	renewed, err := l.Renew(context.Background(), "missing", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Renew returned error: %v", err)
	}
	if renewed {
		t.Fatalf("Renew returned true for missing key")
	}
}

func TestInMemoryLockerRefreshCompatibilityAlias(t *testing.T) {
	t.Parallel()

	l := NewInMemoryLocker()
	ctx := context.Background()
	locked, err := l.TryLock(ctx, "k", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("TryLock returned error: %v", err)
	}
	if !locked {
		t.Fatalf("TryLock did not lock key")
	}

	refreshed, err := l.Refresh(ctx, "k", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}
	if !refreshed {
		t.Fatalf("Refresh returned false for existing key")
	}
}
