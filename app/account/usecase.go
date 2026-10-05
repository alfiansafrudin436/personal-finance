package account

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"personal-finance/app/account/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Usecase handles all account business logic
type Usecase struct {
	repo   repository.AccountRepository
	appCfg *config.App
}

// NewUsecase creates a new account Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// GetAll returns every account owned by the authenticated user.
// Archived accounts are hidden unless ?includeArchived=true.
func (u *Usecase) GetAll(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	includeArchived, _ := strconv.ParseBool(c.QueryParam("includeArchived"))

	ctx := c.Request().Context()
	accounts, err := u.repo.ListAccounts(ctx, repository.ListAccountsParams{
		UserID:          userID,
		IncludeArchived: sql.NullBool{Bool: includeArchived, Valid: true},
	})
	if err != nil {
		log.Println("GetAll - ListAccounts error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data akun"))
	}

	if accounts == nil {
		accounts = []repository.Account{}
	}

	summary, err := u.repo.GetAccountsSummary(ctx, userID)
	if err != nil {
		log.Println("GetAll - GetAccountsSummary error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil ringkasan akun"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"accounts":      accounts,
		"totalBalance":  summary.TotalBalance,
		"totalAccounts": summary.TotalAccounts,
	}))
}

// GetByID returns a single account owned by the authenticated user
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
	acc, err := u.repo.GetAccountByID(ctx, repository.GetAccountByIDParams{ID: id, UserID: userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Akun tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(acc))
}

// Create creates a new account with an optional opening balance
func (u *Usecase) Create(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	apiErr, req := ValidateAccountInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	acc, err := u.repo.CreateAccount(ctx, repository.CreateAccountParams{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	})
	if err != nil {
		if utils.IsUniqueViolation(err) {
			return c.JSON(http.StatusConflict, utils.ResponseError("Akun dengan nama tersebut sudah ada"))
		}
		log.Println("Create - CreateAccount error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal membuat akun"))
	}

	return c.JSON(http.StatusCreated, utils.ResponseOK(acc))
}

// Update updates the name, type and currency of an account. The balance is
// derived from transactions, so it is deliberately not editable here.
func (u *Usecase) Update(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	apiErr, req := ValidateAccountInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()
	acc, err := u.repo.UpdateAccount(ctx, repository.UpdateAccountParams{
		ID:       id,
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		Currency: req.Currency,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, utils.ResponseError("Akun tidak ditemukan"))
		}
		if utils.IsUniqueViolation(err) {
			return c.JSON(http.StatusConflict, utils.ResponseError("Akun dengan nama tersebut sudah ada"))
		}
		log.Println("Update - UpdateAccount error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate akun"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(acc))
}

// Archive hides an account without touching its transaction history
func (u *Usecase) Archive(c echo.Context) error {
	return u.setArchived(c, true)
}

// Unarchive restores a previously archived account
func (u *Usecase) Unarchive(c echo.Context) error {
	return u.setArchived(c, false)
}

func (u *Usecase) setArchived(c echo.Context, archived bool) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	ctx := c.Request().Context()
	if _, err := u.repo.GetAccountByID(ctx, repository.GetAccountByIDParams{ID: id, UserID: userID}); err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Akun tidak ditemukan"))
	}

	if err := u.repo.SetAccountArchived(ctx, repository.SetAccountArchivedParams{
		ID:         id,
		UserID:     userID,
		IsArchived: archived,
	}); err != nil {
		log.Println("setArchived - SetAccountArchived error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengubah status akun"))
	}

	if archived {
		return c.JSON(http.StatusOK, utils.ResponseOK("Akun berhasil diarsipkan"))
	}
	return c.JSON(http.StatusOK, utils.ResponseOK("Akun berhasil diaktifkan kembali"))
}

// Delete permanently removes an account. An account that still has
// transactions must be archived instead, so history is never silently lost.
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
	count, err := u.repo.CountAccountTransactions(ctx, repository.CountAccountTransactionsParams{
		UserID:    userID,
		AccountID: id,
	})
	if err != nil {
		log.Println("Delete - CountAccountTransactions error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus akun"))
	}
	if count > 0 {
		return c.JSON(http.StatusConflict, utils.ResponseError("Akun masih memiliki transaksi, arsipkan akun ini sebagai gantinya"))
	}

	rows, err := u.repo.DeleteAccount(ctx, repository.DeleteAccountParams{ID: id, UserID: userID})
	if err != nil {
		log.Println("Delete - DeleteAccount error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus akun"))
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Akun tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Akun berhasil dihapus"))
}
