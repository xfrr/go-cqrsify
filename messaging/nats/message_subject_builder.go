package messagingnats

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// defaultSubjectBuilder is the default subject builder that uses the message type as the subject.
var defaultSubjectBuilder = NewMessageTypeSubjectBuilder()

// SubjectBuilder builds a NATS subject for a given message.
type SubjectBuilder interface {
	Build(m messaging.Message) string
}

// SubjectBuilderFunc is a function type that implements the SubjectBuilder interface.
type SubjectBuilderFunc func(m messaging.Message) string

// Build builds a NATS subject for a given message.
func (f SubjectBuilderFunc) Build(m messaging.Message) string {
	return f(m)
}

type MessageTypeSubjectBuilder struct{}

func NewMessageTypeSubjectBuilder() *MessageTypeSubjectBuilder {
	return &MessageTypeSubjectBuilder{}
}

// Build builds a NATS subject based on the message type.
func (b *MessageTypeSubjectBuilder) Build(m messaging.Message) string {
	return m.MessageType()
}
