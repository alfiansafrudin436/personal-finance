package transaction

import (
	"time"

	"personal-finance/app/transaction/repository"
	"personal-finance/utils"
)

// TransactionResponse is the wire shape of a transaction, with the account and
// category it points at joined in. Nullable columns are pointers so they
// serialize as JSON null, and the date is a plain YYYY-MM-DD string.
type TransactionResponse struct {
	ID              string                  `json:"id"`
	AccountID       string                  `json:"accountId"`
	AccountName     string                  `json:"accountName"`
	ToAccountID     *string                 `json:"toAccountId"`
	ToAccountName   *string                 `json:"toAccountName"`
	CategoryID      string                  `json:"categoryId"`
	CategoryName    string                  `json:"categoryName"`
	CategoryType    repository.CategoryType `json:"categoryType"`
	CategoryIcon    *string                 `json:"categoryIcon"`
	CategoryColor   *string                 `json:"categoryColor"`
	Amount          string                  `json:"amount"`
	TransactionDate string                  `json:"transactionDate"`
	Description     *string                 `json:"description"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

func toTransactionResponse(t repository.ListTransactionsRow) TransactionResponse {
	return TransactionResponse{
		ID:              t.ID.String(),
		AccountID:       t.AccountID.String(),
		AccountName:     t.AccountName,
		ToAccountID:     utils.NullUUID(t.ToAccountID),
		ToAccountName:   utils.NullString(t.ToAccountName),
		CategoryID:      t.CategoryID.String(),
		CategoryName:    t.CategoryName,
		CategoryType:    t.CategoryType,
		CategoryIcon:    utils.NullString(t.CategoryIcon),
		CategoryColor:   utils.NullString(t.CategoryColor),
		Amount:          t.Amount,
		TransactionDate: utils.Date(t.TransactionDate),
		Description:     utils.NullString(t.Description),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

func toTransactionResponseFromDetail(t repository.GetTransactionByIDRow) TransactionResponse {
	return toTransactionResponse(repository.ListTransactionsRow(t))
}

func toTransactionResponses(rows []repository.ListTransactionsRow) []TransactionResponse {
	out := make([]TransactionResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, toTransactionResponse(r))
	}
	return out
}

// WrittenTransactionResponse is returned after a create or update. It carries
// the stored row rather than the joined view, which the client already has.
type WrittenTransactionResponse struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"accountId"`
	ToAccountID     *string   `json:"toAccountId"`
	CategoryID      string    `json:"categoryId"`
	Amount          string    `json:"amount"`
	TransactionDate string    `json:"transactionDate"`
	Description     *string   `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func toWrittenTransactionResponse(t repository.Transaction) WrittenTransactionResponse {
	return WrittenTransactionResponse{
		ID:              t.ID.String(),
		AccountID:       t.AccountID.String(),
		ToAccountID:     utils.NullUUID(t.ToAccountID),
		CategoryID:      t.CategoryID.String(),
		Amount:          t.Amount,
		TransactionDate: utils.Date(t.TransactionDate),
		Description:     utils.NullString(t.Description),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
