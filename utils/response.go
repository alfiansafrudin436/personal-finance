package utils

// ResponseStatusOK is a generic success response
type ResponseStatusOK[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

// ResponseStatusError is a generic error response
type ResponseStatusError struct {
	Errors []TCustomAPIError `json:"errors"`
}

// ResponseOK creates a success response (status=200)
func ResponseOK[T any](data T) ResponseStatusOK[T] {
	return ResponseStatusOK[T]{
		Status: 200,
		Data:   data,
	}
}

// ResponseError creates an error response with a simple message
func ResponseError(message string) ResponseStatusError {
	return ResponseStatusError{
		Errors: []TCustomAPIError{
			{Msg: message},
		},
	}
}
