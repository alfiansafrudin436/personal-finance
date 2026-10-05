package budget

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"personal-finance/app/budget/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Usecase handles all budget business logic
type Usecase struct {
	repo   repository.BudgetRepository
	appCfg *config.App
}

// NewUsecase creates a new budget Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// GetAll returns the budgets for one month together with the amount already
// spent per category. The month comes from ?period=YYYY-MM and defaults to
// the current month.
func (u *Usecase) GetAll(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	period := utils.StartOfMonth(time.Now())
	if raw := strings.TrimSpace(c.QueryParam("period")); raw != "" {
		parsed, err := utils.ParseMonth(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		}
		period = parsed
	}

	ctx := c.Request().Context()
	budgets, err := u.repo.ListBudgets(ctx, repository.ListBudgetsParams{
		UserID:      userID,
		PeriodMonth: period,
	})
	if err != nil {
		log.Println("GetAll - ListBudgets error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data anggaran"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"period":  period.Format("2006-01"),
		"budgets": toBudgetResponses(budgets),
	}))
}

// Create sets the budget for a category and month. Setting it again for the
// same category and month replaces the amount rather than failing, which is
// what the single budget row per category per month is for.
func (u *Usecase) Create(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	apiErr, req := ValidateBudgetInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	b, err := u.repo.UpsertBudget(ctx, repository.UpsertBudgetParams{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		PeriodMonth: req.Period,
	})
	if err != nil {
		if utils.IsForeignKeyViolation(err) {
			return c.JSON(http.StatusBadRequest, utils.ResponseError("Kategori tidak ditemukan"))
		}
		log.Println("Create - UpsertBudget error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menyimpan anggaran"))
	}

	return c.JSON(http.StatusCreated, utils.ResponseOK(toWrittenBudgetResponse(b)))
}

// Update changes the amount of an existing budget
func (u *Usecase) Update(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	apiErr, amount := ValidateAmountInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	b, err := u.repo.UpdateBudget(ctx, repository.UpdateBudgetParams{
		ID:     id,
		UserID: userID,
		Amount: amount,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, utils.ResponseError("Anggaran tidak ditemukan"))
		}
		log.Println("Update - UpdateBudget error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate anggaran"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toWrittenBudgetResponse(b)))
}

// Delete removes a budget
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
	rows, err := u.repo.DeleteBudget(ctx, repository.DeleteBudgetParams{ID: id, UserID: userID})
	if err != nil {
		log.Println("Delete - DeleteBudget error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus anggaran"))
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Anggaran tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Anggaran berhasil dihapus"))
}
