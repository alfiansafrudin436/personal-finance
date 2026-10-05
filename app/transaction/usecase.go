package transaction

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"personal-finance/app/transaction/repository"
	"personal-finance/config"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Pagination bounds for the transaction list.
const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// Usecase handles all transaction business logic
type Usecase struct {
	repo   repository.TransactionRepository
	appCfg *config.App
}

// NewUsecase creates a new transaction Usecase
func NewUsecase() *Usecase {
	return &Usecase{
		repo:   repository.NewRepository(),
		appCfg: config.Application,
	}
}

// balanceEffects returns the account balance adjustments a transaction implies.
// Income credits the account, expense debits it, and a transfer moves the
// amount from the source account to the destination account.
func balanceEffects(
	catType repository.CategoryType,
	accountID uuid.UUID,
	toAccountID uuid.NullUUID,
	amount string,
) []repository.BalanceEffect {
	switch catType {
	case repository.CategoryTypeIncome:
		return []repository.BalanceEffect{{AccountID: accountID, Delta: amount}}
	case repository.CategoryTypeExpense:
		return []repository.BalanceEffect{{AccountID: accountID, Delta: utils.NegateAmount(amount)}}
	case repository.CategoryTypeTransfer:
		effects := []repository.BalanceEffect{{AccountID: accountID, Delta: utils.NegateAmount(amount)}}
		if toAccountID.Valid {
			effects = append(effects, repository.BalanceEffect{AccountID: toAccountID.UUID, Delta: amount})
		}
		return effects
	default:
		return nil
	}
}

// reverseEffects negates a set of balance effects, used to undo a stored
// transaction before an update or a delete.
func reverseEffects(effects []repository.BalanceEffect) []repository.BalanceEffect {
	reversed := make([]repository.BalanceEffect, 0, len(effects))
	for _, e := range effects {
		reversed = append(reversed, repository.BalanceEffect{
			AccountID: e.AccountID,
			Delta:     utils.NegateAmount(e.Delta),
		})
	}
	return reversed
}

// refFailure is a rejected payload: a status code and a message to return.
// validateReferences reports failures this way rather than returning the
// result of c.JSON, which is nil when the response is written successfully.
type refFailure struct {
	status  int
	message string
}

// validateReferences checks that the account, the optional destination account
// and the category all belong to the user, and that the shape of the payload
// matches the category type. It returns the resolved category type.
func (u *Usecase) validateReferences(
	c echo.Context,
	userID uuid.UUID,
	req *parsedTransaction,
) (repository.CategoryType, *refFailure) {
	ctx := c.Request().Context()

	account, err := u.repo.GetAccountForTransaction(ctx, repository.GetAccountForTransactionParams{
		ID:     req.AccountID,
		UserID: userID,
	})
	if err != nil {
		return "", &refFailure{http.StatusBadRequest, "Akun tidak ditemukan"}
	}
	if account.IsArchived {
		return "", &refFailure{http.StatusBadRequest, "Akun sudah diarsipkan dan tidak dapat digunakan"}
	}

	cat, err := u.repo.GetCategoryForTransaction(ctx, repository.GetCategoryForTransactionParams{
		ID:     req.CategoryID,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		return "", &refFailure{http.StatusBadRequest, "Kategori tidak ditemukan"}
	}

	if cat.Type == repository.CategoryTypeTransfer {
		if !req.ToAccountID.Valid {
			return "", &refFailure{http.StatusBadRequest, "Akun tujuan harus dipilih untuk transaksi transfer"}
		}
		if req.ToAccountID.UUID == req.AccountID {
			return "", &refFailure{http.StatusBadRequest, "Akun tujuan harus berbeda dengan akun sumber"}
		}
		toAccount, err := u.repo.GetAccountForTransaction(ctx, repository.GetAccountForTransactionParams{
			ID:     req.ToAccountID.UUID,
			UserID: userID,
		})
		if err != nil {
			return "", &refFailure{http.StatusBadRequest, "Akun tujuan tidak ditemukan"}
		}
		if toAccount.IsArchived {
			return "", &refFailure{http.StatusBadRequest, "Akun tujuan sudah diarsipkan dan tidak dapat digunakan"}
		}
	} else if req.ToAccountID.Valid {
		return "", &refFailure{http.StatusBadRequest, "Akun tujuan hanya berlaku untuk transaksi transfer"}
	}

	return cat.Type, nil
}

// GetAll returns a paginated, filtered list of transactions.
// Supported query params: page, pageSize, startDate, endDate, accountId,
// categoryId, type and search.
func (u *Usecase) GetAll(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var startDate sql.NullTime
	if raw := strings.TrimSpace(c.QueryParam("startDate")); raw != "" {
		parsed, err := utils.ParseDate(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError("Format startDate harus YYYY-MM-DD"))
		}
		startDate = sql.NullTime{Time: parsed, Valid: true}
	}

	var endDate sql.NullTime
	if raw := strings.TrimSpace(c.QueryParam("endDate")); raw != "" {
		parsed, err := utils.ParseDate(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError("Format endDate harus YYYY-MM-DD"))
		}
		endDate = sql.NullTime{Time: parsed, Valid: true}
	}

	if startDate.Valid && endDate.Valid && endDate.Time.Before(startDate.Time) {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("endDate tidak boleh lebih awal dari startDate"))
	}

	var accountID uuid.NullUUID
	if raw := strings.TrimSpace(c.QueryParam("accountId")); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError("accountId tidak valid"))
		}
		accountID = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	var categoryID uuid.NullUUID
	if raw := strings.TrimSpace(c.QueryParam("categoryId")); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.ResponseError("categoryId tidak valid"))
		}
		categoryID = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	catType, valid := ParseCategoryTypeFilter(c.QueryParam("type"))
	if !valid {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("Tipe transaksi tidak valid (income, expense, transfer)"))
	}

	var search sql.NullString
	if raw := strings.TrimSpace(c.QueryParam("search")); raw != "" {
		search = sql.NullString{String: raw, Valid: true}
	}

	ctx := c.Request().Context()

	total, err := u.repo.CountTransactions(ctx, repository.CountTransactionsParams{
		UserID:     userID,
		StartDate:  startDate,
		EndDate:    endDate,
		AccountID:  accountID,
		CategoryID: categoryID,
		Type:       catType,
		Search:     search,
	})
	if err != nil {
		log.Println("GetAll - CountTransactions error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data transaksi"))
	}

	items, err := u.repo.ListTransactions(ctx, repository.ListTransactionsParams{
		UserID:     userID,
		StartDate:  startDate,
		EndDate:    endDate,
		AccountID:  accountID,
		CategoryID: categoryID,
		Type:       catType,
		Search:     search,
		RowLimit:   int32(pageSize),
		RowOffset:  int32((page - 1) * pageSize),
	})
	if err != nil {
		log.Println("GetAll - ListTransactions error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengambil data transaksi"))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return c.JSON(http.StatusOK, utils.ResponseOK(map[string]interface{}{
		"items": toTransactionResponses(items),
		"pagination": map[string]interface{}{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	}))
}

// GetByID returns one transaction with its account and category joined in
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
	trx, err := u.repo.GetTransactionByID(ctx, repository.GetTransactionByIDParams{ID: id, UserID: userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Transaksi tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toTransactionResponseFromDetail(trx)))
}

// Create records a transaction and applies it to the account balances
func (u *Usecase) Create(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	apiErr, req := ValidateTransactionInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	catType, fail := u.validateReferences(c, userID, req)
	if fail != nil {
		return c.JSON(fail.status, utils.ResponseError(fail.message))
	}

	ctx := c.Request().Context()
	created, err := u.repo.CreateWithBalance(ctx, repository.CreateTransactionParams{
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		TransactionDate: req.TransactionDate,
		Description:     req.Description,
	}, balanceEffects(catType, req.AccountID, req.ToAccountID, req.Amount))
	if err != nil {
		log.Println("Create - CreateWithBalance error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menyimpan transaksi"))
	}

	return c.JSON(http.StatusCreated, utils.ResponseOK(toWrittenTransactionResponse(created)))
}

// Update rewrites a transaction. The stored version is first reversed out of
// the account balances, then the new version is applied, both in one database
// transaction, so editing the amount, the account or the category stays exact.
func (u *Usecase) Update(c echo.Context) error {
	userID, ok := utils.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ResponseError("Tidak terautentikasi"))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ResponseError("ID tidak valid"))
	}

	apiErr, req := ValidateTransactionInput(c)
	if apiErr != nil {
		return c.JSON(http.StatusBadRequest, apiErr)
	}

	ctx := c.Request().Context()

	existing, err := u.repo.GetTransactionRaw(ctx, repository.GetTransactionRawParams{ID: id, UserID: userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Transaksi tidak ditemukan"))
	}

	oldCat, err := u.repo.GetCategoryForTransaction(ctx, repository.GetCategoryForTransactionParams{
		ID:     existing.CategoryID,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		log.Println("Update - GetCategoryForTransaction (existing) error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate transaksi"))
	}

	newCatType, fail := u.validateReferences(c, userID, req)
	if fail != nil {
		return c.JSON(fail.status, utils.ResponseError(fail.message))
	}

	effects := reverseEffects(balanceEffects(
		oldCat.Type, existing.AccountID, existing.ToAccountID, existing.Amount,
	))
	effects = append(effects, balanceEffects(
		newCatType, req.AccountID, req.ToAccountID, req.Amount,
	)...)

	updated, err := u.repo.UpdateWithBalance(ctx, repository.UpdateTransactionParams{
		ID:              id,
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		ToAccountID:     req.ToAccountID,
		Amount:          req.Amount,
		TransactionDate: req.TransactionDate,
		Description:     req.Description,
	}, effects)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, utils.ResponseError("Transaksi tidak ditemukan"))
		}
		log.Println("Update - UpdateWithBalance error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal mengupdate transaksi"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK(toWrittenTransactionResponse(updated)))
}

// Delete removes a transaction and reverses its effect on the balances
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

	existing, err := u.repo.GetTransactionRaw(ctx, repository.GetTransactionRawParams{ID: id, UserID: userID})
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Transaksi tidak ditemukan"))
	}

	cat, err := u.repo.GetCategoryForTransaction(ctx, repository.GetCategoryForTransactionParams{
		ID:     existing.CategoryID,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		log.Println("Delete - GetCategoryForTransaction error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus transaksi"))
	}

	effects := reverseEffects(balanceEffects(
		cat.Type, existing.AccountID, existing.ToAccountID, existing.Amount,
	))

	rows, err := u.repo.DeleteWithBalance(ctx, repository.DeleteTransactionParams{ID: id, UserID: userID}, effects)
	if err != nil {
		log.Println("Delete - DeleteWithBalance error:", err)
		return c.JSON(http.StatusInternalServerError, utils.ResponseError("Gagal menghapus transaksi"))
	}
	if rows == 0 {
		return c.JSON(http.StatusNotFound, utils.ResponseError("Transaksi tidak ditemukan"))
	}

	return c.JSON(http.StatusOK, utils.ResponseOK("Transaksi berhasil dihapus"))
}
