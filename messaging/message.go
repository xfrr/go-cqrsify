package messaging

import "time"

var _ Message = (*baseMessage)(nil)

// Message represents a generic message.
// It can be used as a base interface for more specific message types like Event, Command or Query.
// Implementations must ensure that MessageID is unique for each message instance.
type Message interface {
	// MessageID returns the unique identifier of the message.
	// It can be optionally used to correlate messages.
	MessageID() string

	// MessageType returns the type of the message related to the originating occurrence.
	// It is used to route messages to the appropriate handler.
	//
	// - MUST be a non-empty string
	//
	// - SHOULD be a short, human-readable string that describes the purpose of the message.
	//
	// - SHOULD be unique within the context of the message source.
	//
	// - SHOULD be prefixed with a reverse-DNS name.
	// The prefixed domain dictates the organization which defines the semantics of this event type.
	//
	// - SHOULD be versioned using semantic versioning (e.g., "com.my_org.object.created.v1").
	MessageType() string

	// MessageSchemaURI returns the schema URI of the message.
	// It can be optionally used to specify the format of the message.
	// Must be a valid URI reference.
	MessageSchemaURI() string

	// MessageSource returns the source of the message.
	// It can be optionally used to specify the source of the message.
	MessageSource() string

	// MessageTimestamp returns the creation time of the message in UTC.
	MessageTimestamp() time.Time

	// MessageMetadata returns the metadata of the message.
	MessageMetadata() map[string]string
}

// baseMessage is the base message implementation.
type baseMessage struct {
	id        string
	_type     string
	schema    string
	source    string
	timestamp time.Time
	metadata  map[string]string
}

func (b baseMessage) MessageID() string                  { return b.id }
func (b baseMessage) MessageType() string                { return b._type }
func (b baseMessage) MessageSchemaURI() string           { return b.schema }
func (b baseMessage) MessageSource() string              { return b.source }
func (b baseMessage) MessageTimestamp() time.Time        { return b.timestamp }
func (b baseMessage) MessageMetadata() map[string]string { return b.metadata }

func newBaseMessage(msgType string, modifiers ...baseMessageModifier) baseMessage {
	b := baseMessage{
		_type:     msgType,
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
