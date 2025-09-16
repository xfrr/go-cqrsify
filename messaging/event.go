package messaging

// Event represents a significant occurrence or change in state within the system.
type Event interface {
	Message
}

// BaseEvent provides a basic implementation of the Event interface.
type BaseEvent struct {
	BaseMessage
}

// NewBaseEvent creates a new BaseEvent with the given name and payload.
func NewBaseEvent(eventType string, modifiers ...BaseMessageModifier) BaseEvent {
	return BaseEvent{
		BaseMessage: NewBaseMessage(
			eventType,
			modifiers...,
		),
	}
}
