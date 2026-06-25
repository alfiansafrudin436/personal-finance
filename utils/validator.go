package utils

// TCustomAPIError mimics Express-like validation error shape
type TCustomAPIError struct {
	Location string `json:"location,omitempty"`
	Msg      string `json:"msg"`
	Path     string `json:"path,omitempty"`
	Type     string `json:"type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// NetworkAPIError wraps multiple TCustomAPIError entries
type NetworkAPIError struct {
	Errors []TCustomAPIError `json:"errors"`
}

// NewError returns a NetworkAPIError with a single error entry
func NewError(msg string, opts ...func(*TCustomAPIError)) NetworkAPIError {
	e := TCustomAPIError{Msg: msg}
	for _, opt := range opts {
		opt(&e)
	}
	return NetworkAPIError{Errors: []TCustomAPIError{e}}
}

// NewErrors returns a NetworkAPIError with multiple errors
func NewErrors(errs ...TCustomAPIError) NetworkAPIError {
	return NetworkAPIError{Errors: errs}
}

// WithLocation sets the location field on a TCustomAPIError
func WithLocation(loc string) func(*TCustomAPIError) {
	return func(e *TCustomAPIError) {
		e.Location = loc
	}
}

// WithPath sets the path field on a TCustomAPIError
func WithPath(path string) func(*TCustomAPIError) {
	return func(e *TCustomAPIError) {
		e.Path = path
	}
}

// WithType sets the type field on a TCustomAPIError
func WithType(t string) func(*TCustomAPIError) {
	return func(e *TCustomAPIError) {
		e.Type = t
	}
}

// WithValue sets the value field on a TCustomAPIError
func WithValue(val string) func(*TCustomAPIError) {
	return func(e *TCustomAPIError) {
		e.Value = val
	}
}
