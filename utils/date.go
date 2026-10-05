package utils

import (
	"errors"
	"time"
)

// DateLayout is the wire format for all date-only fields.
const DateLayout = "2006-01-02"

// ParseDate parses a YYYY-MM-DD date in UTC, so a DATE column round-trips
// without the local timezone shifting the day.
func ParseDate(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.New("tanggal harus diisi")
	}
	t, err := time.ParseInLocation(DateLayout, raw, time.UTC)
	if err != nil {
		return time.Time{}, errors.New("format tanggal harus YYYY-MM-DD")
	}
	return t, nil
}

// ParseMonth parses a YYYY-MM month and returns its first day in UTC.
func ParseMonth(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.New("periode harus diisi")
	}
	t, err := time.ParseInLocation("2006-01", raw, time.UTC)
	if err != nil {
		return time.Time{}, errors.New("format periode harus YYYY-MM")
	}
	return t, nil
}

// StartOfMonth returns the first day of the month that t falls in, in UTC.
func StartOfMonth(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// EndOfMonth returns the last day of the month that t falls in, in UTC.
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, -1)
}
