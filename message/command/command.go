package command

import "github.com/xfrr/go-cqrsify/message"

// Command represents a command message.
type Command interface {
	message.Message
}

type Base struct {
	message.Base
}
