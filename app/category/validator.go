package category

import (
	"database/sql"
	"regexp"
	"strings"

	"personal-finance/app/category/repository"
	"personal-finance/utils"

	"github.com/labstack/echo/v4"
)

// validCategoryTypes mirrors the category_type enum in the database.
var validCategoryTypes = map[repository.CategoryType]bool{
	repository.CategoryTypeIncome:   true,
	repository.CategoryTypeExpense:  true,
	repository.CategoryTypeTransfer: true,
}

// hexColorPattern matches a 7-character hex color such as #16a34a.
var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// categoryRequest is the request body for creating or updating a category
type categoryRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// parsedCategory holds the validated and normalized category payload
type parsedCategory struct {
	Name  string
	Type  repository.CategoryType
	Icon  sql.NullString
	Color sql.NullString
}

// ParseCategoryType validates an optional ?type= query filter
func ParseCategoryType(raw string) (repository.NullCategoryType, bool) {
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

// ValidateCategoryInput validates the create/update category request body
func ValidateCategoryInput(c echo.Context) (*utils.NetworkAPIError, *parsedCategory) {
	var req categoryRequest
	if err := c.Bind(&req); err != nil {
		e := utils.NewError("Request Body tidak valid",
			utils.WithLocation("body"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		e := utils.NewError("Nama kategori harus diisi",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("required"),
		)
		return &e, nil
	}
	if len(name) > 50 {
		e := utils.NewError("Nama kategori maksimal 50 karakter",
			utils.WithLocation("body"),
			utils.WithPath("name"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	catType := repository.CategoryType(strings.TrimSpace(req.Type))
	if !validCategoryTypes[catType] {
		e := utils.NewError("Tipe kategori tidak valid (income, expense, transfer)",
			utils.WithLocation("body"),
			utils.WithPath("type"),
			utils.WithType("invalid"),
			utils.WithValue(req.Type),
		)
		return &e, nil
	}

	icon := strings.TrimSpace(req.Icon)
	if len(icon) > 50 {
		e := utils.NewError("Ikon maksimal 50 karakter",
			utils.WithLocation("body"),
			utils.WithPath("icon"),
			utils.WithType("invalid"),
		)
		return &e, nil
	}

	color := strings.TrimSpace(req.Color)
	if color != "" && !hexColorPattern.MatchString(color) {
		e := utils.NewError("Warna harus berupa kode hex, contoh #16a34a",
			utils.WithLocation("body"),
			utils.WithPath("color"),
			utils.WithType("invalid"),
			utils.WithValue(req.Color),
		)
		return &e, nil
	}

	return nil, &parsedCategory{
		Name:  name,
		Type:  catType,
		Icon:  sql.NullString{String: icon, Valid: icon != ""},
		Color: sql.NullString{String: color, Valid: color != ""},
	}
}
