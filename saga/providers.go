package saga

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type UUIDProvider interface{ New() string }

type UUIDProviderFunc func() string

func (f UUIDProviderFunc) New() string {
	return f()
}

type TimeProvider interface{ Now() time.Time }

type TimeProviderFunc func() time.Time

func (f TimeProviderFunc) Now() time.Time {
	return f()
}

var (
	DefaultTimeProvider = TimeProviderFunc(time.Now)
	DefaultUUIDProvider = UUIDProviderFunc(func() string {
		uuid, _ := uuid.NewV4()
		return uuid.String()
	})
)
