package jsonschema

// ValidatorOption defines a functional option for configuring the Validator.
type ValidatorOption func(*Validator)

// WithFilepathResolver sets the FilepathResolver for the Validator.
func WithFilepathResolver(resolver FilepathResolver) ValidatorOption {
	return func(v *Validator) {
		v.filepathResolver = resolver
	}
}

// WithProblemURL sets the base URL for problem types.
func WithProblemURL(url string) ValidatorOption {
	return func(v *Validator) {
		v.problemURL = url
	}
}
