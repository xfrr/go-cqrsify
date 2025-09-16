package messaging

// Query represents an action or intent to change the state of the system.
type Query interface {
	Message
}

// BaseQuery provides a basic implementation of the Query interface.
type BaseQuery struct {
	BaseMessage
}

// NewBaseQuery creates a new BaseQuery with the given name and payload.
func NewBaseQuery(queryType string, modifiers ...BaseMessageModifier) BaseQuery {
	return BaseQuery{
		BaseMessage: NewBaseMessage(
			queryType,
			modifiers...,
		),
	}
}
