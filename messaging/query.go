package messaging

// Query represents an action or intent to change the state of the system.
type Query interface {
	Message

	QueryID() string
}

type BaseQueryModifier = BaseMessageModifier

func WithQueryID(queryID string) BaseQueryModifier {
	return func(b *BaseMessage) {
		b.id = queryID
	}
}

// BaseQuery provides a basic implementation of the Query interface.
type BaseQuery = BaseMessage

func (q BaseQuery) QueryID() string {
	return q.id
}

// NewBaseQuery creates a new BaseQuery with the given name and payload.
func NewBaseQuery(queryType string, modifiers ...BaseQueryModifier) BaseQuery {
	return NewMessage(
		queryType,
		modifiers...,
	)
}

// NewQueryFromJSON creates a BaseQuery from a JSONMessage.
func NewQueryFromJSON[P any](jsonMsg JSONMessage[P]) BaseQuery {
	return BaseQuery{
		id:        jsonMsg.ID,
		_type:     jsonMsg.Type,
		schema:    jsonMsg.SchemaURI,
		source:    jsonMsg.Source,
		timestamp: jsonMsg.Timestamp,
		metadata:  jsonMsg.Metadata,
	}
}
