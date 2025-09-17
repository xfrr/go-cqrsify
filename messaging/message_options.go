package messaging

import "time"

// baseMessageModifier is a functional option for creating new messages.
type baseMessageModifier func(*baseMessage)

// WithTimestamp sets the timestamp for the message.
func WithTimestamp(timestamp time.Time) baseMessageModifier {
	return func(b *baseMessage) {
		b.timestamp = timestamp
	}
}

// WithID sets the ID for the message.
func WithID(id string) baseMessageModifier {
	return func(b *baseMessage) {
		b.id = id
	}
}

// WithMetadata sets the metadata for the message.
func WithMetadata(metadata map[string]string) baseMessageModifier {
	return func(b *baseMessage) {
		b.metadata = metadata
	}
}

// WithMetadataKeyValue sets a key-value pair in the metadata for the message.
func WithMetadataKeyValue(key, value string) baseMessageModifier {
	return func(b *baseMessage) {
		if b.metadata == nil {
			b.metadata = make(map[string]string)
		}
		b.metadata[key] = value
	}
}

// WithSchema sets the schema for the message.
func WithSchema(schema string) baseMessageModifier {
	return func(b *baseMessage) {
		b.schema = schema
	}
}

// WithSource sets the source for the message.
func WithSource(source string) baseMessageModifier {
	return func(b *baseMessage) {
		b.source = source
	}
}
