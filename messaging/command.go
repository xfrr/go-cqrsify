package messaging

// Command represents an action or intent to change the state of the system.
type Command interface {
	Message

	CommandID() string
}

// BaseCommand provides a basic implementation of the Command interface.
type BaseCommand struct {
	BaseMessage
}

func (c BaseCommand) CommandID() string {
	return c.id
}

type BaseCommandModifier = baseMessageModifier

// NewBaseCommand creates a new BaseCommand with the given name and payload.
func NewBaseCommand(commandType string, modifiers ...BaseCommandModifier) BaseCommand {
	return BaseCommand{
		BaseMessage: NewMessage(
			commandType,
			modifiers...,
		),
	}
}

// NewCommandFromJSON creates a BaseCommand from a JSONMessage.
func NewCommandFromJSON[P any](jsonMsg JSONMessage[P]) BaseCommand {
	return BaseCommand{
		BaseMessage: BaseMessage{
			id:        jsonMsg.ID,
			_type:     jsonMsg.Type,
			schema:    jsonMsg.SchemaURI,
			source:    jsonMsg.Source,
			timestamp: jsonMsg.Timestamp,
			metadata:  jsonMsg.Metadata,
		},
	}
}
