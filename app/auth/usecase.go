package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"personal-finance/app/auth/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// Usecase handles all auth business logic
type Usecase struct {
	repo   repository.AuthRepository
	appCfg *config.App
}

// NewUsecase creates a new auth Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// Login handles user login and returns a JWT token
func (u *Usecase) Login(c echo.Context) error {
	apiErr, req := ValidateLoginInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Println("Login - GetUserByEmail error:", err)
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Email atau password salah"))
	}

	if !utils.CompareHashedPassword(user.PasswordHash, req.Password) {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Email atau password salah"))
	}

	token, err := utils.ParseToken(utils.TokenParams{
		ID: user.ID.String(),
		// Role: "user",
	})
	if err != nil {
		log.Println("Login - ParseToken error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal membuat token"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Username,
			"email": user.Email,
			// "role":  user.Role,
		},
	}))
}

// // Register handles new user registration
func (u *Usecase) Register(c echo.Context) error {
	apiErr, req := ValidateRegisterInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()

	// Check if email is already taken
	_, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusConflict, utils.ResponseError("Email sudah digunakan"))
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Register - HashPassword error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal memproses password"))
	}

	newID, err := u.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		log.Println("Register - CreateUser error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mendaftarkan pengguna"))
	}

	return c.JSON(http.StatusCreated, utils.ResponseOK(map[string]interface{}{
		"id":       newID,
		"email":    req.Email,
		"username": req.Username,
	}))
}

// // ValidateToken validates a Bearer token from the Authorization header
// func (u *Usecase) ValidateToken(c echo.Context) error {
// 	apiErr, tokenStr := ValidateBearerToken(c)
// 	if apiErr != nil {
// 		return c.JSON(http.StatusUnauthorized, apiErr)
// 	}

// 	claims, err := utils.VerifyToken(tokenStr)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError(err.Error()))
// 	}

// 	ctx := c.Request().Context()
// 	userID, err := uuid.Parse(claims.ID)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Token tidak valid"))
// 	}

// 	user, err := u.repo.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Pengguna tidak ditemukan"))
// 	}

// 	return c.JSON(http.StatusOK, utils.ResponseOK(user))
// }

// // RefreshToken issues a new token for the authenticated user
// func (u *Usecase) RefreshToken(c echo.Context) error {
// 	apiErr, tokenStr := ValidateBearerToken(c)
// 	if apiErr != nil {
// 		return c.JSON(http.StatusUnauthorized, apiErr)
// 	}

// 	claims, err := utils.VerifyToken(tokenStr)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError(err.Error()))
// 	}

// 	ctx := c.Request().Context()
// 	userID, err := uuid.Parse(claims.ID)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Token tidak valid"))
// 	}

// 	user, err := u.repo.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Pengguna tidak ditemukan"))
// 	}

// 	newToken, err := utils.ParseToken(utils.TokenParams{
// 		ID:   user.ID.String(),
// 		Role: user.Role,
// 	})
// 	if err != nil {
// 		log.Println("RefreshToken - ParseToken error:", err)
// 		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal membuat token baru"))
// 	}

// 	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]string{
// 		"token": newToken,
// 	}))
// }

// // ForgotPassword sends a password reset email
// func (u *Usecase) ForgotPassword(c echo.Context) error {
// 	apiErr, email := ValidateForgotPasswordInput(c)
// 	if apiErr != nil {
// 		return c.JSON(http.StatusBadRequest, apiErr)
// 	}

// 	ctx := c.Request().Context()
// 	user, err := u.repo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		// Return success even if user not found to prevent email enumeration
// 		return c.JSON(http.StatusOK, utils.ResponseOK("Jika email terdaftar, link reset password telah dikirim"))
// 	}

// 	// Generate a short-lived reset token (1 hour)
// 	resetToken, err := utils.ParseTokenWithExpiration(utils.TokenParams{ID: user.ID.String()}, 3600*1e9)
// 	if err != nil {
// 		log.Println("ForgotPassword - ParseTokenWithExpiration error:", err)
// 		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal membuat token reset"))
// 	}

// 	// Save reset token to DB
// 	if err := u.repo.UpdateUserResetToken(ctx, repository.UpdateUserResetTokenParams{
// 		ID:                 user.ID,
// 		ResetPasswordToken: sql.NullString{String: resetToken, Valid: true},
// 	}); err != nil {
// 		log.Println("ForgotPassword - UpdateUserResetToken error:", err)
// 		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menyimpan token reset"))
// 	}

// 	// Send reset email
// 	resetLink := u.appCfg.FEUrl + "/reset-password?token=" + resetToken
// 	emailBody := "<p>Klik link berikut untuk reset password Anda:</p><a href=\"" + resetLink + "\">" + resetLink + "</a>"

// 	if err := u.appCfg.SendEmail(email, "Reset Password", emailBody, nil); err != nil {
// 		log.Println("ForgotPassword - SendEmail error:", err)
// 		// Don't expose internal email errors to client
// 	}

// 	return c.JSON(http.StatusOK, utils.ResponseOK("Jika email terdaftar, link reset password telah dikirim"))
// }

// // ResetPassword resets the user's password using a valid reset token
// func (u *Usecase) ResetPassword(c echo.Context) error {
// 	apiErr, tokenStr, newPassword := ValidateResetPasswordInput(c)
// 	if apiErr != nil {
// 		return c.JSON(http.StatusBadRequest, apiErr)
// 	}

// 	ctx := c.Request().Context()

// 	user, err := u.repo.GetUserByResetToken(ctx, sql.NullString{String: tokenStr, Valid: true})
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, utils.ResponseError("Token tidak valid atau sudah kadaluarsa"))
// 	}

// 	hashedPassword, err := utils.HashPassword(newPassword)
// 	if err != nil {
// 		log.Println("ResetPassword - HashPassword error:", err)
// 		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal memproses password"))
// 	}

// 	if err := u.repo.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
// 		ID:       user.ID,
// 		Password: hashedPassword,
// 	}); err != nil {
// 		log.Println("ResetPassword - UpdateUserPassword error:", err)
// 		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate password"))
// 	}

// 	return c.JSON(http.StatusOK, utils.ResponseOK("Password berhasil diubah"))
// }
