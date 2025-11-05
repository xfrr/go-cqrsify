package messagingnats

import (
	"fmt"
	"time"

	"github.com/xfrr/go-cqrsify/messaging"
)

// defaultReplySubjectBuilder is the default subject builder that uses the message type as the subject.
var defaultReplySubjectBuilder = NewMessageTypeReplySubjectBuilder()

// ReplySubjectBuilder builds a NATS subject for a given message.
type ReplySubjectBuilder interface {
	Build(m messaging.Message) string
}

// ReplySubjectBuilderFunc is a function type that implements the ReplySubjectBuilder interface.
type ReplySubjectBuilderFunc func(m messaging.Message) string

// Build builds a NATS subject for a given message.
func (f ReplySubjectBuilderFunc) Build(m messaging.Message) string {
	return f(m)
}

type MessageTypeReplySubjectBuilder struct{}

func NewMessageTypeReplySubjectBuilder() *MessageTypeReplySubjectBuilder {
	return &MessageTypeReplySubjectBuilder{}
}

// Build builds a NATS subject based on the message type.
func (b *MessageTypeReplySubjectBuilder) Build(m messaging.Message) string {
	baseSubject := m.MessageType() + ".reply"
	switch {
	case m.MessageID() != "":
		return baseSubject + "." + m.MessageID()
	case !m.MessageTimestamp().IsZero():
		return fmt.Sprintf("%s.%d", baseSubject, m.MessageTimestamp().UnixNano())
	default:
		return fmt.Sprintf("%s.%d", baseSubject, time.Now().UnixNano())
	}
}
