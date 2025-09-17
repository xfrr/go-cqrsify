package messagingnats

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// SubjectBuilder defines a function that builds a NATS subject from an event name.
type SubjectBuilder func(m messaging.Message) string

// DefaultSubjectBuilder is the default implementation of SubjectBuilder.
// It returns the event name as the subject.
func DefaultSubjectBuilder(m messaging.Message) string {
	return m.MessageType()
}
