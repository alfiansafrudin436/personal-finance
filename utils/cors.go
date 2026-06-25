package utils

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
)

// GetMiddleWareConfig returns CORS middleware configuration
func GetMiddleWareConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}
}
