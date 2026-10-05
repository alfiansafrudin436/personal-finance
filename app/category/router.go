package category

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all category-related routes under the given group.
// The group is expected to already carry the JWT middleware.
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.GET("", u.GetAll)
	g.POST("", u.Create)
	g.GET("/:id", u.GetByID)
	g.PUT("/:id", u.Update)
	g.DELETE("/:id", u.Delete)
}
