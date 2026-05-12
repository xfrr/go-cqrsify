package saga

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type UUIDProvider interface{ New() (string, error) }

type UUIDProviderFunc func() (string, error)

func (f UUIDProviderFunc) New() (string, error) {
	return f()
}

type TimeProvider interface{ Now() time.Time }

type TimeProviderFunc func() time.Time

func (f TimeProviderFunc) Now() time.Time {
	return f()
}

var (
	DefaultTimeProvider = TimeProviderFunc(time.Now)
	DefaultUUIDProvider = UUIDProviderFunc(func() (string, error) {
		uuid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		return uuid.String(), nil
	})
)
