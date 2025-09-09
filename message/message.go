package message

import "time"

type Message interface {
	// ID returns the unique identifier of the message.
	// It can be optionally used to correlate messages.
	ID() string

	// Schema returns the schema URI of the message.
	// It can be optionally used to specify the format of the message.
	// Must be a valid URI reference.
	Schema() string

	// Source returns the source of the message.
	// It can be optionally used to specify the source of the message.
	Source() string

	// Timestamp returns the creation time of the message in UTC.
	Timestamp() time.Time

	// Metadata returns the metadata of the message.
	Metadata() map[string]string
}
