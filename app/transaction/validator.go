package transaction

import (
	"database/sql"
	"strings"
	"time"

	"personal-finance/app/transaction/repository"
	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// maxDescriptionLength caps the free-text description.
const maxDescriptionLength = 500

// transactionRequest is the request body for creating or updating a transaction
type transactionRequest struct {
	AccountID       string `json:"accountId"`
	CategoryID      string `json:"categoryId"`
	ToAccountID     string `json:"toAccountId"`
	Amount          string `json:"amount"`
	TransactionDate string `json:"transactionDate"`
	Description     string `json:"description"`
}

// parsedTransaction holds the validated and normalized transaction payload
type parsedTransaction struct {
	AccountID       uuid.UUID
	CategoryID      uuid.UUID
	ToAccountID     uuid.NullUUID
	Amount          string
	TransactionDate time.Time
	Description     sql.NullString
}

// ValidateTransactionInput validates the create/update transaction body.
// Whether to_account_id is required depends on the category type, which lives
// in the database, so that rule is enforced in the usecase instead.
func ValidateTransactionInput(c echo.Context) (*utils.NetworkAPIError, *parsedTransaction) {
	var req transactionRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	accountID, err := uuid.Parse(strings.TrimSpace(req.AccountID))
	if err != nil {
		e := utils.NewError("Akun harus dipilih",
			utils.WithLocation("body"),
			utils.WithPath("accountId"),
			utils.WithType("required"),
			utils.WithValue(req.AccountID),
		)
		return &e, nil
	}

	categoryID, err := uuid.Parse(strings.TrimSpace(req.CategoryID))
	if err != nil {
		e := utils.NewError("Kategori harus dipilih",
			utils.WithLocation("body"),
			utils.WithPath("categoryId"),
			utils.WithType("required"),
			utils.WithValue(req.CategoryID),
		)
		return &e, nil
	}

	var toAccountID uuid.NullUUID
	if raw := strings.TrimSpace(req.ToAccountID); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			e := utils.NewError("Akun tujuan tidak valid",
				utils.WithLocation("body"),
				utils.WithPath("toAccountId"),
				utils.WithType("invalid"),
				utils.WithValue(req.ToAccountID),
			)
			return &e, nil
		}
		toAccountID = uuid.NullUUID{UUID: parsed, Valid: true}
	}

	amount, err := utils.ParseAmount(req.Amount)
	if err != nil {
		e := utils.NewError(err.Error(),
			utils.WithLocation("body"),
			utils.WithPath("amount"),
			utils.WithType("invalid"),
			utils.WithValue(req.Amount),
		)
		return &e, nil
	}

	date, err := utils.ParseDate(strings.TrimSpace(req.TransactionDate))
	if err != nil {
		e := utils.NewError(err.Error(),
			utils.WithLocation("body"),
			utils.WithPath("transactionDate"),
			utils.WithType("invalid"),
			utils.WithValue(req.TransactionDate),
		)
		return &e, nil
	}

	description := strings.TrimSpace(req.Description)
	if len(description) > maxDescriptionLength {
		e := utils.NewError("Deskripsi maksimal 500 karakter",
			utils.WithLocation("body"),
			utils.WithPath("description"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	return nil, &parsedTransaction{
		AccountID:       accountID,
		CategoryID:      categoryID,
		ToAccountID:     toAccountID,
		Amount:          amount,
		TransactionDate: date,
		Description:     sql.NullString{String: description, Valid: description != ""},
	}
}

// validCategoryTypes mirrors the category_type enum in the database.
var validCategoryTypes = map[repository.CategoryType]bool{
	repository.CategoryTypeIncome:   true,
	repository.CategoryTypeExpense:  true,
	repository.CategoryTypeTransfer: true,
}

// ParseCategoryTypeFilter validates an optional ?type= query filter
func ParseCategoryTypeFilter(raw string) (repository.NullCategoryType, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return repository.NullCategoryType{}, true
	}
	t := repository.CategoryType(raw)
	if !validCategoryTypes[t] {
		return repository.NullCategoryType{}, false
	}
	return repository.NullCategoryType{CategoryType: t, Valid: true}, true
}
