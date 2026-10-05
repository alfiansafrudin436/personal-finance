package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all user-related routes under the given group.
// The group is expected to already carry the JWT middleware, so /me resolves
// the caller from the token.
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.GET("/me", u.Me)
	g.PUT("/me", u.UpdateMe)
	g.GET("", u.GetAll)
	g.GET("/:id", u.GetByID)
	g.PUT("/:id", u.Update)
	g.DELETE("/:id", u.Delete)
}
