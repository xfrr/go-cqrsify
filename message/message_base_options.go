package message

import "time"

// BaseModifier is a functional option for creating new messages.
type BaseModifier func(*Base)

// WithTimestamp sets the timestamp for the message.
func WithTimestamp(timestamp time.Time) BaseModifier {
	return func(b *Base) {
		b.timestamp = timestamp
	}
}

// WithID sets the ID for the message.
func WithID(id string) BaseModifier {
	return func(b *Base) {
		b.id = id
	}
}

// WithMetadata sets the metadata for the message.
func WithMetadata(metadata map[string]string) BaseModifier {
	return func(b *Base) {
		b.metadata = metadata
	}
}

// WithMetadataKeyValue sets a key-value pair in the metadata for the message.
func WithMetadataKeyValue(key, value string) BaseModifier {
	return func(b *Base) {
		if b.metadata == nil {
			b.metadata = make(map[string]string)
		}
		b.metadata[key] = value
	}
}

// WithSchema sets the schema for the message.
func WithSchema(schema string) BaseModifier {
	return func(b *Base) {
		b.schema = schema
	}
}

// WithSource sets the source for the message.
func WithSource(source string) BaseModifier {
	return func(b *Base) {
		b.source = source
	}
}
