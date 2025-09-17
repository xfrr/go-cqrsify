package messagingnats

// PubSubMessageBusOption defines a function type for configuring PubSubMessageBus options.
type PubSubMessageBusOption func(*PubSubMessageBus)

// WithSubjectBuilder sets a custom SubjectBuilder for the PubSubMessageBus.
func WithSubjectBuilder(sb SubjectBuilder) PubSubMessageBusOption {
	return func(p *PubSubMessageBus) {
		p.subjectBuilder = sb
	}
}
