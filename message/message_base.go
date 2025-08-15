package message

import (
	"time"
)

// Base is the base message implementation.
type Base struct {
	id            string
	name          string
	correlationID string
	causationID   string
	timestamp     time.Time
	metadata      map[string]string
}

func NewBase(name string, opts ...BaseMessageOption) Base {
	b := Base{
		id:            "",
		correlationID: "",
		causationID:   "",
		name:          name,
		timestamp:     time.Now().UTC(),
		metadata:      map[string]string{},
	}
	for _, o := range opts {
		o(&b)
	}
	return b
}

func (b Base) Name() string                { return b.name }
func (b Base) MessageID() string           { return b.id }
func (b Base) CorrelationID() string       { return b.correlationID }
func (b Base) CausationID() string         { return b.causationID }
func (b Base) Timestamp() time.Time        { return b.timestamp }
func (b Base) Metadata() map[string]string { return b.metadata }
