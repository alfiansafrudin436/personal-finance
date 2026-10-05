package report

import (
	"strings"
	"time"

	"personal-finance/app/report/repository"
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// dateRange is an inclusive reporting window.
type dateRange struct {
	Start time.Time
	End   time.Time
}

// ParseDateRange reads ?startDate= and ?endDate= and defaults to the current
// calendar month when either is absent.
func ParseDateRange(c echo.Context) (dateRange, error) {
	now := time.Now()
	r := dateRange{
		Start: utils.StartOfMonth(now),
		End:   utils.EndOfMonth(now),
	}

	if raw := strings.TrimSpace(c.QueryParam("startDate")); raw != "" {
		parsed, err := utils.ParseDate(raw)
		if err != nil {
			return r, err
		}
		r.Start = parsed
	}

	if raw := strings.TrimSpace(c.QueryParam("endDate")); raw != "" {
		parsed, err := utils.ParseDate(raw)
		if err != nil {
			return r, err
		}
		r.End = parsed
	}

	if r.End.Before(r.Start) {
		return r, errEndBeforeStart
	}

	return r, nil
}

// errEndBeforeStart is returned when the requested window is inverted.
var errEndBeforeStart = &rangeError{"endDate tidak boleh lebih awal dari startDate"}

type rangeError struct{ msg string }

func (e *rangeError) Error() string { return e.msg }

// ParseReportType reads ?type= and defaults to expense, the breakdown that
// matters most on a spending report.
func ParseReportType(raw string) (repository.CategoryType, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return repository.CategoryTypeExpense, true
	}
	t := repository.CategoryType(raw)
	switch t {
	case repository.CategoryTypeIncome, repository.CategoryTypeExpense, repository.CategoryTypeTransfer:
		return t, true
	default:
		return "", false
	}
}
