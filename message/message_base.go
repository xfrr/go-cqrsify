package message

import (
	"time"
)

// Base is the base message implementation.
type Base struct {
	id        string
	schema    string
	source    string
	timestamp time.Time
	metadata  map[string]string
}

func NewBase(modifiers ...BaseModifier) Base {
	b := Base{
		id:        "",
		schema:    "",
		source:    "",
		timestamp: time.Now().UTC(),
		metadata:  map[string]string{},
	}
	for _, o := range modifiers {
		o(&b)
	}
	return b
}

func (b Base) ID() string                  { return b.id }
func (b Base) Schema() string              { return b.schema }
func (b Base) Source() string              { return b.source }
func (b Base) Timestamp() time.Time        { return b.timestamp }
func (b Base) Metadata() map[string]string { return b.metadata }
