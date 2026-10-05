package utils

import (
	"errors"

	"github.com/lib/pq"
)

// PostgreSQL error codes we react to specifically.
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
)

// IsUniqueViolation reports whether err is a Postgres unique constraint error,
// so handlers can answer 409 instead of a generic 500.
func IsUniqueViolation(err error) bool {
	return hasPQCode(err, pgUniqueViolation)
}

// IsForeignKeyViolation reports whether err is a Postgres foreign key error.
func IsForeignKeyViolation(err error) bool {
	return hasPQCode(err, pgForeignKeyViolation)
}

// IsCheckViolation reports whether err is a Postgres CHECK constraint error.
func IsCheckViolation(err error) bool {
	return hasPQCode(err, pgCheckViolation)
}

func hasPQCode(err error, code string) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == code
	}
	return false
}
