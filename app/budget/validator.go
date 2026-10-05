package budget

import (
	"strings"
	"time"

	"personal-finance/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// budgetRequest is the request body for creating or updating a budget
type budgetRequest struct {
	CategoryID string `json:"categoryId"`
	Amount     string `json:"amount"`
	Period     string `json:"period"`
}

// parsedBudget holds the validated and normalized budget payload
type parsedBudget struct {
	CategoryID uuid.UUID
	Amount     string
	Period     time.Time
}

// ValidateBudgetInput validates the create/upsert budget request body
func ValidateBudgetInput(c echo.Context) (*utils.NetworkAPIError, *parsedBudget) {
	var req budgetRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
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

	period, err := utils.ParseMonth(strings.TrimSpace(req.Period))
	if err != nil {
		e := utils.NewError(err.Error(),
			utils.WithLocation("body"),
			utils.WithPath("period"),
			utils.WithType("invalid"),
			utils.WithValue(req.Period),
		)
		return &e, nil
	}

	return nil, &parsedBudget{
		CategoryID: categoryID,
		Amount:     amount,
		Period:     period,
	}
}

// ValidateAmountInput validates a body that only carries an amount
func ValidateAmountInput(c echo.Context) (*utils.NetworkAPIError, string) {
	var req budgetRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, ""
	}

	amount, err := utils.ParseAmount(req.Amount)
	if err != nil {
		e := utils.NewError(err.Error(),
			utils.WithLocation("body"),
			utils.WithPath("amount"),
			utils.WithType("invalid"),
			utils.WithValue(req.Amount),
		)
		return &e, ""
	}

	return nil, amount
}
