package message

import (
	"encoding/json"
	"time"
)

// Ensure Envelope implements json.Marshaler and json.Unmarshaler interfaces.
var _ json.Marshaler = (*Envelope)(nil)
var _ json.Unmarshaler = (*Envelope)(nil)

type envelopeJSONAlias[T any] struct {
	ID        string            `json:"id"`
	Schema    string            `json:"schema"`
	Source    string            `json:"source"`
	Metadata  map[string]string `json:"metadata"`
	Payload   T                 `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewEnvelope creates a new Message Envelope instance.
func NewEnvelope(msgBase Base, payload any) Envelope {
	return Envelope{
		Base:    msgBase,
		payload: payload,
	}
}

// Envelope is a wrapper for messages that adds metadata.
type Envelope struct {
	Base
	payload any
}

// Payload returns the payload of the message envelope.
func (e *Envelope) Payload() any {
	return e.payload
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Envelope) UnmarshalJSON(data []byte) error {
	var aux envelopeJSONAlias[any]
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Base = Base{
		id:        aux.ID,
		schema:    aux.Schema,
		source:    aux.Source,
		timestamp: aux.Timestamp,
		metadata:  aux.Metadata,
	}

	e.payload = aux.Payload
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (e *Envelope) MarshalJSON() ([]byte, error) {
	aux := envelopeJSONAlias[any]{
		ID:        e.Base.ID(),
		Schema:    e.Base.Schema(),
		Source:    e.Base.Source(),
		Metadata:  e.Base.Metadata(),
		Payload:   e.Payload(),
		Timestamp: e.Base.Timestamp(),
	}
	return json.Marshal(aux)
}

// UnmarshalEnvelope is a helper function to unmarshal JSON data into an Envelope.
func UnmarshalEnvelopeJSON[T any](data []byte) (Envelope, error) {
	var aux envelopeJSONAlias[T]
	if err := json.Unmarshal(data, &aux); err != nil {
		return Envelope{}, err
	}

	envelope := Envelope{
		Base: Base{
			id:        aux.ID,
			schema:    aux.Schema,
			source:    aux.Source,
			timestamp: aux.Timestamp,
			metadata:  aux.Metadata,
		},
		payload: aux.Payload,
	}

	return envelope, nil
}

// MarshalEnvelopeJSON is a helper function to marshal an Envelope into JSON data.
func MarshalEnvelopeJSON(e Envelope) ([]byte, error) {
	aux := envelopeJSONAlias[any]{
		ID:        e.Base.ID(),
		Schema:    e.Base.Schema(),
		Source:    e.Base.Source(),
		Metadata:  e.Base.Metadata(),
		Payload:   e.Payload(),
		Timestamp: e.Base.Timestamp(),
	}
	return json.Marshal(aux)
}

// EnvelopePayloadAs attempts to cast the payload of the envelope to the specified type T.
// It returns the casted payload and a boolean indicating whether the cast was successful.
func EnvelopePayloadAs[T any](e Envelope) (T, bool) {
	payload, ok := e.Payload().(T)
	return payload, ok
}
