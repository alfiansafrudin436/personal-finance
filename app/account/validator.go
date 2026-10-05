package account

import (
	"strings"

	"personal-finance/app/account/repository"
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// validAccountTypes mirrors the account_type enum in the database.
var validAccountTypes = map[repository.AccountType]bool{
	repository.AccountTypeCash:       true,
	repository.AccountTypeBank:       true,
	repository.AccountTypeEWallet:    true,
	repository.AccountTypeCreditCard: true,
	repository.AccountTypeInvestment: true,
}

// accountRequest is the request body for creating or updating an account
type accountRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

// parsedAccount holds the validated and normalized account payload
type parsedAccount struct {
	Name     string
	Type     repository.AccountType
	Balance  string
	Currency string
}

// ValidateAccountInput validates the create/update account request body
func ValidateAccountInput(c echo.Context) (*utils.NetworkAPIError, *parsedAccount) {
	var req accountRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		e := utils.NewError("Nama akun harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	if len(name) > 50 {
		e := utils.NewError("Nama akun maksimal 50 karakter",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	accType := repository.AccountType(strings.TrimSpace(req.Type))
	if !validAccountTypes[accType] {
		e := utils.NewError("Tipe akun tidak valid (cash, bank, e_wallet, credit_card, investment)",
			utils.WithLocation("body"),
			utils.WithPath("type"),
			utils.WithType("invalid"),
			utils.WithValue(req.Type),
		)
		return &e, nil
	}

	// Opening balance may legitimately be zero or negative (e.g. a credit card).
	balance, err := utils.ParseSignedAmount(req.Balance)
	if err != nil {
		e := utils.NewError(err.Error(),
			utils.WithLocation("body"),
			utils.WithPath("balance"),
			utils.WithType("invalid"),
			utils.WithValue(req.Balance),
		)
		return &e, nil
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "IDR"
	}
	if len(currency) != 3 {
		e := utils.NewError("Mata uang harus 3 karakter (contoh: IDR)",
			utils.WithLocation("body"),
			utils.WithPath("currency"),
			utils.WithType("invalid"),
			utils.WithValue(req.Currency),
		)
		return &e, nil
	}

	return nil, &parsedAccount{
		Name:     name,
		Type:     accType,
		Balance:  balance,
		Currency: currency,
	}
}
