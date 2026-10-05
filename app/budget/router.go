package budget

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all budget-related routes under the given group.
// The group is expected to already carry the JWT middleware.
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.GET("", u.GetAll)
	g.POST("", u.Create)
	g.PUT("/:id", u.Update)
	g.DELETE("/:id", u.Delete)
}
