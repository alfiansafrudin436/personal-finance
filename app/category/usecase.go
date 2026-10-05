package category

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"personal-finance/app/category/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Usecase handles all category business logic
type Usecase struct {
	repo   repository.CategoryRepository
	appCfg *config.App
}

// NewUsecase creates a new category Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// GetAll returns the categories of the authenticated user together with the
// global ones, optionally filtered by ?type=income|expense|transfer.
func (u *Usecase) GetAll(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	catType, valid := ParseCategoryType(c.QueryParam("type"))
	if !valid {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("Tipe kategori tidak valid (income, expense, transfer)"))
	}

	ctx := c.Request().Context()
	categories, err := u.repo.ListCategories(ctx, repository.ListCategoriesParams{
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
		Type:   catType,
	})
	if err != nil {
		log.Println("GetAll - ListCategories error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data kategori"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toCategoryResponses(categories)))
}

// GetByID returns one category, which may be a global category
func (u *Usecase) GetByID(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	ctx := c.Request().Context()
	cat, err := u.repo.GetCategoryByID(ctx, repository.GetCategoryByIDParams{
		ID:     id,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Kategori tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toCategoryResponseFromDetail(cat)))
}

// Create creates a category owned by the authenticated user
func (u *Usecase) Create(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	apiErr, req := ValidateCategoryInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	cat, err := u.repo.CreateCategory(ctx, repository.CreateCategoryParams{
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
		Name:   req.Name,
		Type:   req.Type,
		Icon:   req.Icon,
		Color:  req.Color,
	})
	if err != nil {
		if utils.IsUniqueViolation(err) {
			return c.JSON(http.StatusConflict, utils.ResponseError("Kategori dengan nama dan tipe tersebut sudah ada"))
		}
		log.Println("Create - CreateCategory error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal membuat kategori"))
	}

	return c.JSON(http.StatusCreated, utils.ResponseOK(toCategoryResponse(cat)))
}

// Update updates a category. Global categories are shared by every user, so
// the query is scoped to user_id and editing one answers 403 rather than 404.
func (u *Usecase) Update(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	apiErr, req := ValidateCategoryInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	existing, err := u.repo.GetCategoryByID(ctx, repository.GetCategoryByIDParams{
		ID:     id,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Kategori tidak ditemukan"))
	}
	if existing.IsGlobal {
		return c.JSON(http.StatusForbidden, utils.ResponseError("Kategori bawaan tidak dapat diubah"))
	}

	cat, err := u.repo.UpdateCategory(ctx, repository.UpdateCategoryParams{
		ID:     id,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
		Name:   req.Name,
		Type:   req.Type,
		Icon:   req.Icon,
		Color:  req.Color,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, utils.ResponseError("Kategori tidak ditemukan"))
		}
		if utils.IsUniqueViolation(err) {
			return c.JSON(http.StatusConflict, utils.ResponseError("Kategori dengan nama dan tipe tersebut sudah ada"))
		}
		log.Println("Update - UpdateCategory error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate kategori"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toCategoryResponse(cat)))
}

// Delete removes a category owned by the user. A category still referenced by
// transactions is kept, because transactions.category_id is ON DELETE RESTRICT.
func (u *Usecase) Delete(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	ctx := c.Request().Context()
	existing, err := u.repo.GetCategoryByID(ctx, repository.GetCategoryByIDParams{
		ID:     id,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Kategori tidak ditemukan"))
	}
	if existing.IsGlobal {
		return c.JSON(http.StatusForbidden, utils.ResponseError("Kategori bawaan tidak dapat dihapus"))
	}

	count, err := u.repo.CountCategoryTransactions(ctx, id)
	if err != nil {
		log.Println("Delete - CountCategoryTransactions error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus kategori"))
	}
	if count > 0 {
		return c.JSON(http.StatusConflict, utils.ResponseError("Kategori masih digunakan oleh transaksi"))
	}

	rows, err := u.repo.DeleteCategory(ctx, repository.DeleteCategoryParams{
		ID:     id,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		log.Println("Delete - DeleteCategory error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus kategori"))
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Kategori tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Kategori berhasil dihapus"))
}
