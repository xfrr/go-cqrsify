package messagingnats

import (
	"github.com/xfrr/go-cqrsify/messaging"
)

// SubjectBuilderFunc defines a function that builds a NATS subject from an event name.
type SubjectBuilderFunc func(m messaging.Message) string

// DefaultSubjectBuilder is the default implementation of SubjectBuilder.
// It returns the event name as the subject.
func DefaultSubjectBuilder(m messaging.Message) string {
	return m.MessageType()
}
