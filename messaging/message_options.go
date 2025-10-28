package messaging

import "time"

// BaseMessageModifier is a functional option for creating new messages.
type BaseMessageModifier func(*BaseMessage)

// WithTimestamp sets the timestamp for the message.
func WithTimestamp(timestamp time.Time) BaseMessageModifier {
	return func(b *BaseMessage) {
		b.timestamp = timestamp
	}
}

// WithID sets the ID for the message.
func WithID(id string) BaseMessageModifier {
	return func(b *BaseMessage) {
		b.id = id
	}
}

// WithMetadata sets the metadata for the message.
func WithMetadata(metadata map[string]string) BaseMessageModifier {
	return func(b *BaseMessage) {
		b.metadata = metadata
	}
}

// WithMetadataKeyValue sets a key-value pair in the metadata for the message.
func WithMetadataKeyValue(key, value string) BaseMessageModifier {
	return func(b *BaseMessage) {
		if b.metadata == nil {
			b.metadata = make(map[string]string)
		}
		b.metadata[key] = value
	}
}

// WithSchema sets the schema for the message.
func WithSchema(schema string) BaseMessageModifier {
	return func(b *BaseMessage) {
		b.schema = schema
	}
}

// WithSource sets the source for the message.
func WithSource(source string) BaseMessageModifier {
	return func(b *BaseMessage) {
		b.source = source
	}
}
