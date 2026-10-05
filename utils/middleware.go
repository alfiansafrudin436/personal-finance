package utils

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// contextUserIDKey is the echo context key holding the authenticated user ID
const contextUserIDKey = "authUserID"

// JWTMiddleware validates the Authorization Bearer token and stores the
// authenticated user ID on the request context for downstream handlers.
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, ResponseError("Token otorisasi tidak ditemukan"))
			}

			tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, ResponseError("Token otorisasi harus diisi"))
			}

			claims, err := VerifyToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ResponseError(err.Error()))
			}

			userID, err := uuid.Parse(claims.ID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ResponseError("Token tidak valid"))
			}

			c.Set(contextUserIDKey, userID)
			return next(c)
		}
	}
}

// GetUserID returns the authenticated user ID placed on the context by
// JWTMiddleware. The second return value is false when the request was not
// authenticated, which should not happen on a route behind the middleware.
func GetUserID(c echo.Context) (uuid.UUID, bool) {
	id, ok := c.Get(contextUserIDKey).(uuid.UUID)
	return id, ok
}
