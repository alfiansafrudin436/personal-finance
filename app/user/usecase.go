package user

import (
	"log"
	"net/http"
	"personal-finance/app/user/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Usecase handles all user business logic
type Usecase struct {
	repo   repository.UserRepository
	appCfg *config.App
}

// NewUsecase creates a new user Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// GetAll returns all active users
func (u *Usecase) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		log.Println("GetAll - GetAllUsers error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data pengguna"))
	}

	if users == nil {
		users = []repository.GetAllUsersRow{}
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(users))
}

// GetByID returns a single user by ID
func (u *Usecase) GetByID(c echo.Context) error {
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	ctx := c.Request().Context()
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Pengguna tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(user))
}

// Update updates a user's name
func (u *Usecase) Update(c echo.Context) error {
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	apiErr, req := ValidateUpdateInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()

	// Check if user exists
	if _, err := u.repo.GetUserByID(ctx, id); err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Pengguna tidak ditemukan"))
	}

	if err := u.repo.UpdateUserUsername(ctx, repository.UpdateUserUsernameParams{
		ID:       id,
		Username: req.Name,
	}); err != nil {
		log.Println("Update - UpdateUserName error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate pengguna"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Pengguna berhasil diupdate"))
}

// Delete deactivates a user (soft delete)
func (u *Usecase) Delete(c echo.Context) error {
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	ctx := c.Request().Context()

	// Check if user exists
	if _, err := u.repo.GetUserByID(ctx, id); err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Pengguna tidak ditemukan"))
	}

	if err := u.repo.DeactivateUser(ctx, id); err != nil {
		log.Println("Delete - DeactivateUser error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menonaktifkan pengguna"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Pengguna berhasil dinonaktifkan"))
}

// Me returns the profile of the authenticated user, resolved from the JWT so
// the client never has to pass its own ID.
func (u *Usecase) Me(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	ctx := c.Request().Context()
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Pengguna tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(user))
}

// UpdateMe updates the username of the authenticated user
func (u *Usecase) UpdateMe(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	apiErr, req := ValidateUpdateInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	if err := u.repo.UpdateUserUsername(ctx, repository.UpdateUserUsernameParams{
		ID:       userID,
		Username: req.Name,
	}); err != nil {
		log.Println("UpdateMe - UpdateUserUsername error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate profil"))
	}

	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusOK, utils.ResponseOK("Profil berhasil diupdate"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(user))
}
