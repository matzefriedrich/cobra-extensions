package types

// CobraXError represents a custom error type with support for attaching a cause and chaining errors.
type CobraXError struct {
	msg   string
	cause error
}

// Error returns the error message associated with the CobraXError.
func (c *CobraXError) Error() string {
	return c.msg
}

// Unwrap returns the underlying cause of the CobraXError, allowing error unwrapping as part of Go's error chain.
func (c *CobraXError) Unwrap() error {
	return c.cause
}

// Is checks if the target error is equivalent to the current CobraXError by comparing their error messages.
func (c *CobraXError) Is(target error) bool {
	return c.Error() == target.Error()
}

var _ error = (*CobraXError)(nil)

type CobraXErrorOption func(target *CobraXError)

// NewCobraXError creates a new instance of CobraXError with a message and optional configurations passed as options.
func NewCobraXError(msg string, options ...CobraXErrorOption) *CobraXError {
	err := &CobraXError{
		msg: msg,
	}
	for _, option := range options {
		option(err)
	}
	return err
}

// WithCause sets the underlying cause for a CobraXError using the provided error as a CobraXErrorOption.
func WithCause(cause error) CobraXErrorOption {
	return func(target *CobraXError) {
		target.cause = cause
	}
}
