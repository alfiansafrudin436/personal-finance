package user

import (
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// updateUserRequest is the request body for updating a user
type updateUserRequest struct {
	Name string `json:"name"`
}

// ValidateUpdateInput validates the update user request body
func ValidateUpdateInput(c echo.Context) (*utils.NetworkAPIError, *updateUserRequest) {
	var req updateUserRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}
	if req.Name == "" {
		e := utils.NewError("Nama harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	return nil, &req
}

// ValidateBearerToken extracts and validates the Authorization Bearer token
func ValidateBearerToken(c echo.Context) (*utils.NetworkAPIError, string) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		e := utils.NewError("Header otorisasi tidak valid",
			utils.WithLocation("header"),
			utils.WithPath("Authorization"),
			utils.WithType("invalid"),
		)
		return &e, ""
	}
	token := authHeader[7:]
	if token == "" {
		e := utils.NewError("Token otorisasi harus diisi",
			utils.WithLocation("header"),
			utils.WithPath("Authorization"),
			utils.WithType("required"),
		)
		return &e, ""
	}
	return nil, token
}
