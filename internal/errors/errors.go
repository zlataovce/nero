package errors

import "errors"

// ErrUnsupported is an alias for errors.ErrUnsupported.
var ErrUnsupported = errors.ErrUnsupported

// New redirects to the errors.New method.
func New(text string) error {
	return errors.New(text)
}

// Is redirects to the errors.Is method.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As redirects to the errors.As method.
func As(err error, target any) bool {
	//goland:noinspection GoErrorsAs
	return errors.As(err, target)
}
