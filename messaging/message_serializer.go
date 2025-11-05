package messaging

import (
	"encoding/json"
	"fmt"
	"time"
)

var _ MessageSerializer = (*JSONSerializer)(nil)
var _ MessageDeserializer = (*JSONDeserializer)(nil)

type JSONMessage[P any] struct {
	ID        string            `json:"id,omitempty"`
	Type      string            `json:"type"`
	Source    string            `json:"source,omitempty"`
	SchemaURI string            `json:"schemaUri,omitempty"`
	Payload   P                 `json:"payload,omitempty"`
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

type JSONMessageEncoder[T Message, P any] func(msg T) JSONMessage[P]
type JSONMessageDecoder[T Message, P any] func(jsonMessage JSONMessage[P]) (T, error)

// JSONSerializer is a JSON-based serializer.
type JSONSerializer struct {
	encoders map[string]func(msg Message) ([]byte, error)
}

// NewJSONSerializer creates a new JSONSerializer with the given serializers.
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{
		encoders: make(map[string]func(msg Message) ([]byte, error)),
	}
}

// RegisterEncoder registers a serializer function for the given message type.
func (s *JSONSerializer) RegisterEncoder(msgType string, encoder func(msg Message) ([]byte, error)) {
	s.encoders[msgType] = encoder
}

// Serialize implements MessageSerializer.
func (s *JSONSerializer) Serialize(msg Message) ([]byte, error) {
	encoder, ok := s.encoders[msg.MessageType()]
	if !ok {
		return nil, fmt.Errorf("no encoder registered for message type %s", msg.MessageType())
	}
	return encoder(msg)
}

// JSONDeserializer holds JSON decoders for different message types.
type JSONDeserializer struct {
	decoders map[string]func([]byte) (Message, error)
}

// NewJSONDeserializer creates a new JSONDeserializer with the given deserializers.
func NewJSONDeserializer() *JSONDeserializer {
	return &JSONDeserializer{
		decoders: make(map[string]func([]byte) (Message, error)),
	}
}

// RegisterDecoder registers a deserializer function for the given message type.
func (d *JSONDeserializer) RegisterDecoder(msgType string, decoder func([]byte) (Message, error)) {
	d.decoders[msgType] = decoder
}

// Deserialize implements MessageDeserializer.
func (d *JSONDeserializer) Deserialize(msgData []byte) (Message, error) {
	decoder, ok := d.decoders[string(msgData)]
	if !ok {
		return nil, fmt.Errorf("no decoder registered for message type in data: %s", string(msgData))
	}
	return decoder(msgData)
}

// RegisterJSONMessageSerializer is a helper function to register a payload serializer for a specific message type.
func RegisterJSONMessageSerializer[T Message, P any](s *JSONSerializer, msgType string, encoder JSONMessageEncoder[T, P]) *JSONSerializer {
	s.RegisterEncoder(msgType, func(msg Message) ([]byte, error) {
		castMsg, ok := msg.(T)
		if !ok {
			return nil, InvalidMessageTypeError{
				Actual:   fmt.Sprintf("%T", msg),
				Expected: fmt.Sprintf("%T", castMsg),
			}
		}

		jsonMessage := encoder(castMsg)
		return json.Marshal(jsonMessage)
	})
	return s
}

// RegisterJSONMessageDeserializer is a helper function to register a payload deserializer for a specific message type.
func RegisterJSONMessageDeserializer[T Message, P any](d *JSONDeserializer, msgType string, decoder JSONMessageDecoder[T, P]) {
	d.RegisterDecoder(msgType, func(msgData []byte) (Message, error) {
		var jsonMessage JSONMessage[json.RawMessage]
		if err := json.Unmarshal(msgData, &jsonMessage); err != nil {
			return nil, err
		}

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

		return decoder(typedMessage)
	})
}

func NewJSONMessage[P any](msg Message, payload P) JSONMessage[P] {
	var id string
	if bmsg, ok := msg.(BaseMessage); ok {
		id = bmsg.id
	}

	return JSONMessage[P]{
		ID:        id,
		Type:      msg.MessageType(),
		Source:    msg.MessageSource(),
		SchemaURI: msg.MessageSchemaURI(),
		Payload:   payload,
		Timestamp: msg.MessageTimestamp(),
		Metadata:  msg.MessageMetadata(),
	}
}
