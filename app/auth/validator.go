package auth

import (
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// loginRequest is the request body for login
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// registerRequest is the request body for registration
type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"name"`
	Password string `json:"password"`
}

// // forgotPasswordRequest is the request body for forgot password
// type forgotPasswordRequest struct {
// 	Email string `json:"email"`
// }

// // resetPasswordRequest is the request body for reset password
// type resetPasswordRequest struct {
// 	Token       string `json:"token"`
// 	NewPassword string `json:"new_password"`
// }

// ValidateLoginInput validates the login request body
func ValidateLoginInput(c echo.Context) (*utils.NetworkAPIError, *loginRequest) {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}
	if req.Email == "" || !utils.IsValidEmail(req.Email) {
		e := utils.NewError("Email tidak valid atau tidak diisi",
			utils.WithLocation("body"),
			utils.WithPath("email"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	if req.Password == "" {
		e := utils.NewError("Password harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("password"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	return nil, &req
}

// ValidateRegisterInput validates the registration request body
func ValidateRegisterInput(c echo.Context) (*utils.NetworkAPIError, *registerRequest) {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}
	if req.Email == "" || !utils.IsValidEmail(req.Email) {
		e := utils.NewError("Email tidak valid atau tidak diisi",
			utils.WithLocation("body"),
			utils.WithPath("email"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	if req.Username == "" {
		e := utils.NewError("Nama harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	if req.Password == "" {
		e := utils.NewError("Password harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("password"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	return nil, &req
}

// // ValidateForgotPasswordInput validates the forgot password request body
// func ValidateForgotPasswordInput(c echo.Context) (*utils.NetworkAPIError, string) {
// 	var req forgotPasswordRequest
// 	if err := c.Bind(&req); err != nil {
// 		e := utils.NewError("Request Body tidak valid",
// 			utils.WithLocation("body"),
// 			utils.WithType("invalid"),
// 		)
// 		return &e, ""
// 	}
// 	if req.Email == "" || !utils.IsValidEmail(req.Email) {
// 		e := utils.NewError("Email tidak valid atau tidak diisi",
// 			utils.WithLocation("body"),
// 			utils.WithPath("email"),
// 			utils.WithType("required"),
// 		)
// 		return &e, ""
// 	}
// 	return nil, req.Email
// }

// // ValidateResetPasswordInput validates the reset password request body
// func ValidateResetPasswordInput(c echo.Context) (*utils.NetworkAPIError, string, string) {
// 	var req resetPasswordRequest
// 	if err := c.Bind(&req); err != nil {
// 		e := utils.NewError("Request Body tidak valid",
// 			utils.WithLocation("body"),
// 			utils.WithType("invalid"),
// 		)
// 		return &e, "", ""
// 	}
// 	if req.Token == "" {
// 		e := utils.NewError("Token harus diisi",
// 			utils.WithLocation("body"),
// 			utils.WithPath("token"),
// 			utils.WithType("required"),
// 		)
// 		return &e, "", ""
// 	}
// 	if req.NewPassword == "" {
// 		e := utils.NewError("Password baru harus diisi",
// 			utils.WithLocation("body"),
// 			utils.WithPath("new_password"),
// 			utils.WithType("required"),
// 		)
// 		return &e, "", ""
// 	}
// 	return nil, req.Token, req.NewPassword
// }

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
