package messaging

import (
	"encoding/json"
	"fmt"
	"time"
)

var _ MessageSerializer = (*JSONSerializer)(nil)
var _ MessageDeserializer = (*JSONDeserializer)(nil)

type JSONMessage[P any] struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Source    string            `json:"source,omitempty"`
	SchemaURI string            `json:"schema_uri,omitempty"`
	Payload   P                 `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// MessageSerializer defines a function type for serializing messages into byte slices.
type MessageSerializer interface {
	Serialize(msg Message) ([]byte, error)
}

// MessageDeserializer defines a function type for deserializing byte slices into messages.
type MessageDeserializer interface {
	Deserialize(msgData []byte) (Message, error)
}

// NoOpSerializer is a no-operation serializer that returns an empty byte slice.
type NoOpSerializer struct{}

// Serialize implements MessageSerializer.
func (s *NoOpSerializer) Serialize(_ Message) ([]byte, error) {
	return []byte{}, nil
}

// NoOpDeserializer is a no-operation deserializer that returns nil.
type NoOpDeserializer struct{}

// Deserialize implements MessageDeserializer.
func (d *NoOpDeserializer) Deserialize(_ []byte) (Message, error) {
	return nil, nil
}

// JSONSerializer is a JSON-based serializer.
type JSONSerializer struct {
	serializers map[string]func(msg Message) ([]byte, error)
}

// NewJSONSerializer creates a new JSONSerializer with the given serializers.
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{
		serializers: make(map[string]func(msg Message) ([]byte, error)),
	}
}

// RegisterSerializer registers a serializer function for the given message type.
func (s *JSONSerializer) RegisterSerializer(msgType string, serializer func(msg Message) ([]byte, error)) {
	s.serializers[msgType] = serializer
}

// Serialize implements MessageSerializer.
func (s *JSONSerializer) Serialize(msg Message) ([]byte, error) {
	serializer, ok := s.serializers[msg.MessageType()]
	if !ok {
		return nil, nil
	}
	return serializer(msg)
}

// JSONDeserializer is a JSON-based deserializer.
type JSONDeserializer struct {
	deserializers map[string]func(jsonMessage JSONMessage[json.RawMessage]) (Message, error)
}

// NewJSONDeserializer creates a new JSONDeserializer with the given deserializers.
func NewJSONDeserializer() *JSONDeserializer {
	return &JSONDeserializer{
		deserializers: make(map[string]func(jsonMessage JSONMessage[json.RawMessage]) (Message, error)),
	}
}

// RegisterDeserializer registers a deserializer function for the given message type.
func (d *JSONDeserializer) RegisterDeserializer(msgType string, deserializer func(jsonMessage JSONMessage[json.RawMessage]) (Message, error)) {
	d.deserializers[msgType] = deserializer
}

// Deserialize implements MessageDeserializer.
func (d *JSONDeserializer) Deserialize(msgData []byte) (Message, error) {
	var jsonMessage JSONMessage[json.RawMessage]
	if err := json.Unmarshal(msgData, &jsonMessage); err != nil {
		return nil, err
	}

	msgType := jsonMessage.Type
	deserializer, ok := d.deserializers[msgType]
	if !ok {
		return nil, nil
	}
	return deserializer(jsonMessage)
}

// RegisterJSONMessageSerializer is a helper function to register a payload serializer for a specific message type.
func RegisterJSONMessageSerializer[T Message, P any](s *JSONSerializer, msgType string, serializer func(e T) JSONMessage[P]) *JSONSerializer {
	s.RegisterSerializer(msgType, func(msg Message) ([]byte, error) {
		castMsg, ok := msg.(T)
		if !ok {
			return nil, InvalidMessageTypeError{
				Actual:   msgType,
				Expected: fmt.Sprintf("%T", msg),
			}
		}

		jsonMessage := serializer(castMsg)
		return json.Marshal(jsonMessage)
	})
	return s
}

// RegisterJSONMessageDeserializer is a helper function to register a payload deserializer for a specific message type.
func RegisterJSONMessageDeserializer[T Message, P any](d *JSONDeserializer, msgType string, deserializer func(jsonMessage JSONMessage[P]) (T, error)) {
	d.RegisterDeserializer(msgType, func(jsonMessage JSONMessage[json.RawMessage]) (Message, error) {
		var payload P
		if err := json.Unmarshal(jsonMessage.Payload, &payload); err != nil {
			return nil, err
		}

		typedMessage := JSONMessage[P]{
			ID:        jsonMessage.ID,
			Type:      jsonMessage.Type,
			Source:    jsonMessage.Source,
			SchemaURI: jsonMessage.SchemaURI,
			Payload:   payload,
			Timestamp: jsonMessage.Timestamp,
			Metadata:  jsonMessage.Metadata,
		}

		return deserializer(typedMessage)
	})
}

func NewJSONMessage[P any](msg Message, payload P) JSONMessage[P] {
	return JSONMessage[P]{
		ID:        msg.MessageID(),
		Type:      msg.MessageType(),
		Source:    msg.MessageSource(),
		SchemaURI: msg.MessageSchemaURI(),
		Payload:   payload,
		Timestamp: msg.MessageTimestamp(),
		Metadata:  msg.MessageMetadata(),
	}
}
