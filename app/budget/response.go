package budget

import (
	"time"

	"personal-finance/app/budget/repository"
	"personal-finance/utils"
)

// BudgetResponse is the wire shape of a budget together with how much of it
// has already been spent, so the client can render progress without a second
// request or any arithmetic of its own.
type BudgetResponse struct {
	ID            string    `json:"id"`
	CategoryID    string    `json:"categoryId"`
	CategoryName  string    `json:"categoryName"`
	CategoryIcon  *string   `json:"categoryIcon"`
	CategoryColor *string   `json:"categoryColor"`
	Amount        string    `json:"amount"`
	Spent         string    `json:"spent"`
	Remaining     string    `json:"remaining"`
	Period        string    `json:"period"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func toBudgetResponse(b repository.ListBudgetsRow) BudgetResponse {
	return BudgetResponse{
		ID:            b.ID.String(),
		CategoryID:    b.CategoryID.String(),
		CategoryName:  b.CategoryName,
		CategoryIcon:  utils.NullString(b.CategoryIcon),
		CategoryColor: utils.NullString(b.CategoryColor),
		Amount:        b.Amount,
		Spent:         b.Spent,
		Remaining:     utils.SubAmount(b.Amount, b.Spent),
		Period:        b.PeriodMonth.UTC().Format("2006-01"),
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}

func toBudgetResponses(rows []repository.ListBudgetsRow) []BudgetResponse {
	out := make([]BudgetResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, toBudgetResponse(r))
	}
	return out
}

// WrittenBudgetResponse is returned after a create or update.
type WrittenBudgetResponse struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"categoryId"`
	Amount     string    `json:"amount"`
	Period     string    `json:"period"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func toWrittenBudgetResponse(b repository.Budget) WrittenBudgetResponse {
	return WrittenBudgetResponse{
		ID:         b.ID.String(),
		CategoryID: b.CategoryID.String(),
		Amount:     b.Amount,
		Period:     b.PeriodMonth.UTC().Format("2006-01"),
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}
