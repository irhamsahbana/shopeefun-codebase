package adapter

// import "codebase-app/internal/pkg/validator"
func WithValidator(v Validator) Option {
	return func(a *Adapter) {
		a.Validator = v
	}
}

/*
func WithValidator(v *validator.Validator) Option {
	return func(a *Adapter) {
		a.Validator = v
	}
}

*/
