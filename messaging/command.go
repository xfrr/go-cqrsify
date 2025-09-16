package messaging

// Command represents an action or intent to change the state of the system.
type Command interface {
	Message
}

// BaseCommand provides a basic implementation of the Command interface.
type BaseCommand struct {
	BaseMessage
}

// NewBaseCommand creates a new BaseCommand with the given name and payload.
func NewBaseCommand(commandType string, modifiers ...BaseMessageModifier) BaseCommand {
	return BaseCommand{
		BaseMessage: NewBaseMessage(
			commandType,
			modifiers...,
		),
	}
}
