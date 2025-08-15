package message

import "time"

// BaseMessageOption is a functional option for creating new messages.
type BaseMessageOption func(*Base)

// WithCorrelationID sets the correlation ID for the message.
func WithCorrelationID(correlationID string) BaseMessageOption {
	return func(b *Base) {
		b.correlationID = correlationID
	}
}

// WithCausationID sets the causation ID for the message.
func WithCausationID(causationID string) BaseMessageOption {
	return func(b *Base) {
		b.causationID = causationID
	}
}

// WithTimestamp sets the timestamp for the message.
func WithTimestamp(timestamp time.Time) BaseMessageOption {
	return func(b *Base) {
		b.timestamp = timestamp
	}
}

// WithMetadata sets the metadata for the message.
func WithMetadata(metadata map[string]string) BaseMessageOption {
	return func(b *Base) {
		b.metadata = metadata
	}
}

// WithName sets the name for the message.
func WithName(name string) BaseMessageOption {
	return func(b *Base) {
		b.name = name
	}
}

// WithID sets the ID for the message.
func WithID(id string) BaseMessageOption {
	return func(b *Base) {
		b.id = id
	}
}

// WithMetadataKeyValue sets a key-value pair in the metadata for the message.
func WithMetadataKeyValue(key, value string) BaseMessageOption {
	return func(b *Base) {
		if b.metadata == nil {
			b.metadata = make(map[string]string)
		}
		b.metadata[key] = value
	}
}
