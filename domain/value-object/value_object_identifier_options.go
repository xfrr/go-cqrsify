package valueobject

type IdentifierOption func(*IdentifierOptions)

type IdentifierOptions struct {
	customValidationFn func(id Identifier[any]) error
}

func (o *IdentifierOptions) Apply(opts ...IdentifierOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithCustomIdentifierValidation(fn func(id Identifier[any]) error) IdentifierOption {
	return func(o *IdentifierOptions) {
		o.customValidationFn = fn
	}
}
