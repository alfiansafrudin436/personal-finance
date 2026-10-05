package utils

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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

// NullString unwraps a sql.NullString into a plain pointer, so a NULL column
// serializes as JSON null instead of the {"String":"","Valid":false} shape
// that sql.NullString marshals to.
func NullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

// NullUUID unwraps a uuid.NullUUID into a plain pointer, for the same reason.
func NullUUID(v uuid.NullUUID) *string {
	if !v.Valid {
		return nil
	}
	s := v.UUID.String()
	return &s
}

// Date formats a DATE column as YYYY-MM-DD, dropping the zero time that a
// full timestamp would otherwise carry into the response.
func Date(t time.Time) string {
	return t.UTC().Format(DateLayout)
}
