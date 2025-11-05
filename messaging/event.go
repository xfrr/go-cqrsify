package messaging

// Event represents a significant occurrence or change in state within the system.
type Event interface {
	Message

	EventID() string
}

// BaseEvent provides a basic implementation of the Event interface.
type BaseEvent = BaseMessage

func (e BaseEvent) EventID() string {
	return e.id
}

type BaseEventModifier = BaseMessageModifier

// NewBaseEvent creates a new BaseEvent with the given name and payload.
func NewBaseEvent(eventType string, modifiers ...BaseEventModifier) BaseEvent {
	return NewMessage(
		eventType,
		modifiers...,
	)
}

// NewEventFromJSON creates a BaseEvent from a JSONMessage.
func NewEventFromJSON[P any](jsonMsg JSONMessage[P]) BaseEvent {
	return BaseEvent{
		id:        jsonMsg.ID,
		_type:     jsonMsg.Type,
		schema:    jsonMsg.SchemaURI,
		source:    jsonMsg.Source,
		timestamp: jsonMsg.Timestamp,
		metadata:  jsonMsg.Metadata,
	}
}
