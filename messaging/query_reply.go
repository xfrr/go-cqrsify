package messaging

type QueryReply = Message

// BaseQueryReply is a base implementation of QueryReply.
type BaseQueryReply = BaseMessage

// NewBaseQueryReply creates a new BaseQueryReply with the given name and payload.
func NewBaseQueryReply(query Query, modifiers ...BaseMessageModifier) BaseQueryReply {
	return NewMessage(query.MessageType()+".reply", modifiers...)
}

// NewQueryReplyFromJSON creates a BaseQueryReply from a JSONMessage.
func NewQueryReplyFromJSON[P any](jsonMsg JSONMessage[P]) BaseQueryReply {
	return BaseQueryReply{
		id:        jsonMsg.ID,
		_type:     jsonMsg.Type,
		schema:    jsonMsg.SchemaURI,
		source:    jsonMsg.Source,
		timestamp: jsonMsg.Timestamp,
		metadata:  jsonMsg.Metadata,
	}
}
