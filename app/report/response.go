package report

import (
	"personal-finance/app/report/repository"
	"personal-finance/utils"
)

// CategoryBreakdownResponse is one slice of a per-category report.
type CategoryBreakdownResponse struct {
	CategoryID       string  `json:"categoryId"`
	CategoryName     string  `json:"categoryName"`
	CategoryIcon     *string `json:"categoryIcon"`
	CategoryColor    *string `json:"categoryColor"`
	Total            string  `json:"total"`
	TransactionCount int64   `json:"transactionCount"`
}

func toCategoryBreakdown(rows []repository.GetSpendingByCategoryRow) []CategoryBreakdownResponse {
	out := make([]CategoryBreakdownResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, CategoryBreakdownResponse{
			CategoryID:       r.CategoryID.String(),
			CategoryName:     r.CategoryName,
			CategoryIcon:     utils.NullString(r.CategoryIcon),
			CategoryColor:    utils.NullString(r.CategoryColor),
			Total:            r.Total,
			TransactionCount: r.TransactionCount,
		})
	}
	return out
}

// RecentTransactionResponse is a compact transaction row for the dashboard.
type RecentTransactionResponse struct {
	ID              string                  `json:"id"`
	AccountName     string                  `json:"accountName"`
	CategoryName    string                  `json:"categoryName"`
	CategoryType    repository.CategoryType `json:"categoryType"`
	CategoryIcon    *string                 `json:"categoryIcon"`
	CategoryColor   *string                 `json:"categoryColor"`
	Amount          string                  `json:"amount"`
	TransactionDate string                  `json:"transactionDate"`
	Description     *string                 `json:"description"`
}

func toRecentTransactions(rows []repository.GetRecentTransactionsRow) []RecentTransactionResponse {
	out := make([]RecentTransactionResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, RecentTransactionResponse{
			ID:              r.ID.String(),
			AccountName:     r.AccountName,
			CategoryName:    r.CategoryName,
			CategoryType:    r.CategoryType,
			CategoryIcon:    utils.NullString(r.CategoryIcon),
			CategoryColor:   utils.NullString(r.CategoryColor),
			Amount:          r.Amount,
			TransactionDate: utils.Date(r.TransactionDate),
			Description:     utils.NullString(r.Description),
		})
	}
	return out
}

// AccountBalanceResponse is one account on the balances list.
type AccountBalanceResponse struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     repository.AccountType `json:"type"`
	Balance  string                 `json:"balance"`
	Currency string                 `json:"currency"`
}

func toAccountBalances(rows []repository.GetAccountBalancesRow) []AccountBalanceResponse {
	out := make([]AccountBalanceResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, AccountBalanceResponse{
			ID:       r.ID.String(),
			Name:     r.Name,
			Type:     r.Type,
			Balance:  r.Balance,
			Currency: r.Currency,
		})
	}
	return out
}

// TrendPointResponse is one month on the income/expense trend.
type TrendPointResponse struct {
	Period       string `json:"period"`
	TotalIncome  string `json:"totalIncome"`
	TotalExpense string `json:"totalExpense"`
	Net          string `json:"net"`
}

func toTrend(rows []repository.GetMonthlyTrendRow) []TrendPointResponse {
	out := make([]TrendPointResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, TrendPointResponse{
			Period:       r.Period,
			TotalIncome:  r.TotalIncome,
			TotalExpense: r.TotalExpense,
			Net:          utils.SubAmount(r.TotalIncome, r.TotalExpense),
		})
	}
	return out
}
