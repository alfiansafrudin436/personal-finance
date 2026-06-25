package auth

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all auth-related routes under the given group
func RegisterRoutes(g *echo.Group) {
	u := NewUsecase()
	g.POST("/login", u.Login)
	g.POST("/register", u.Register)
	g.GET("/validate", u.ValidateToken)
	g.GET("/refresh", u.RefreshToken)
	g.POST("/forgot-password", u.ForgotPassword)
	g.POST("/reset-password", u.ResetPassword)
}
