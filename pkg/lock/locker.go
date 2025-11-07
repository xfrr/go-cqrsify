package lock

import (
	"context"
	"time"
)

// Locker represents a lock mechanism.
type Locker interface {
	// TryLock tries to acquire a lock for key with TTL, returns true if acquired.
	TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	// Unlock releases the lock for key previously acquired.
	Unlock(ctx context.Context, key string) error
}

// Renewer represents a lock renewer mechanism.
type Renewer interface {
	// Renew renews the lock for key with new TTL.
	Renew(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
