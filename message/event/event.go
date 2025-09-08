package event

import (
	"github.com/xfrr/go-cqrsify/message"
)

// Event represents a event message.
type Event interface {
	message.Message
}

type Base struct {
	message.Base
}

func New(modifiers ...message.BaseModifier) *Base {
	return &Base{
		Base: message.NewBase(modifiers...),
	}
}
