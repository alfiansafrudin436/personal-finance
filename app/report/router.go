package report

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all reporting routes under the given group.
// The group is expected to already carry the JWT middleware.
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.GET("/dashboard", u.GetDashboard)
	g.GET("/summary", u.GetSummary)
	g.GET("/by-category", u.GetByCategory)
	g.GET("/trend", u.GetTrend)
	g.GET("/balances", u.GetAccountBalances)
}
