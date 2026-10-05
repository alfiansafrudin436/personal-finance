package category

import (
	"time"

	"personal-finance/app/category/repository"
	"personal-finance/utils"
)

// CategoryResponse is the wire shape of a category. Nullable columns are
// pointers so they serialize as JSON null rather than as the struct that
// sql.NullString marshals to.
type CategoryResponse struct {
	ID        string                  `json:"id"`
	UserID    *string                 `json:"userId"`
	Name      string                  `json:"name"`
	Type      repository.CategoryType `json:"type"`
	Icon      *string                 `json:"icon"`
	Color     *string                 `json:"color"`
	IsGlobal  bool                    `json:"isGlobal"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}

func toCategoryResponse(c repository.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID.String(),
		UserID:    utils.NullUUID(c.UserID),
		Name:      c.Name,
		Type:      c.Type,
		Icon:      utils.NullString(c.Icon),
		Color:     utils.NullString(c.Color),
		IsGlobal:  !c.UserID.Valid,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponseFromRow(c repository.ListCategoriesRow) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID.String(),
		UserID:    utils.NullUUID(c.UserID),
		Name:      c.Name,
		Type:      c.Type,
		Icon:      utils.NullString(c.Icon),
		Color:     utils.NullString(c.Color),
		IsGlobal:  c.IsGlobal,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponseFromDetail(c repository.GetCategoryByIDRow) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID.String(),
		UserID:    utils.NullUUID(c.UserID),
		Name:      c.Name,
		Type:      c.Type,
		Icon:      utils.NullString(c.Icon),
		Color:     utils.NullString(c.Color),
		IsGlobal:  c.IsGlobal,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponses(rows []repository.ListCategoriesRow) []CategoryResponse {
	out := make([]CategoryResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, toCategoryResponseFromRow(r))
	}
	return out
}
