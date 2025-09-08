package message

import "time"

type Message interface {
	ID() string

	Schema() string

	Source() string

	Timestamp() time.Time

	Metadata() map[string]string
}
