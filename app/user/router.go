package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all user-related routes under the given group
// Add JWT middleware here if needed:
// g.Use(middleware.JWT([]byte(config.Application.JWT.SecretKey)))
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.GET("", u.GetAll)
	g.GET("/:id", u.GetByID)
	g.PUT("/:id", u.Update)
	g.DELETE("/:id", u.Delete)
}
