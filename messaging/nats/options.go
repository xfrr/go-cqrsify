package messagingnats

// MessageBusOption defines a function type for configuring PubSubMessageBus options.
type MessageBusOption func(*PubSubMessageBus)

// WithSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageBus.
func WithSubjectBuilder(sb SubjectBuilder) MessageBusOption {
	return func(p *PubSubMessageBus) {
		p.subjectBuilder = sb
	}
}
