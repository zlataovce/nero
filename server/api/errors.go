package api

// HTTPError is an error wrapper carrying an HTTP status code.
type HTTPError struct {
	// Err is the wrapped error.
	Err error
	// Status is the HTTP status code.
	Status int
	// Type is the OpenAPI error type, may be empty.
	Type string
}

// Error returns the string representation of the error.
func (he *HTTPError) Error() string {
	return he.Err.Error()
}

// Unwrap returns the wrapped error.
func (he *HTTPError) Unwrap() error {
	return he.Err
}
