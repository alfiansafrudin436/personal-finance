package repository

import (
	"context"
	"personal-finance/config"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	ListCategories(ctx context.Context, arg ListCategoriesParams) ([]ListCategoriesRow, error)
	GetCategoryByID(ctx context.Context, arg GetCategoryByIDParams) (GetCategoryByIDRow, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	DeleteCategory(ctx context.Context, arg DeleteCategoryParams) (int64, error)
	CountCategoryTransactions(ctx context.Context, categoryID uuid.UUID) (int64, error)
}

type Repository struct {
	query *Queries
}

func NewRepository() *Repository {
	db := config.Application.DB
	return &Repository{
		query: New(db),
	}
}

func (r *Repository) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	return r.query.CreateCategory(ctx, arg)
}

func (r *Repository) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]ListCategoriesRow, error) {
	return r.query.ListCategories(ctx, arg)
}

func (r *Repository) GetCategoryByID(ctx context.Context, arg GetCategoryByIDParams) (GetCategoryByIDRow, error) {
	return r.query.GetCategoryByID(ctx, arg)
}

func (r *Repository) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	return r.query.UpdateCategory(ctx, arg)
}

func (r *Repository) DeleteCategory(ctx context.Context, arg DeleteCategoryParams) (int64, error) {
	return r.query.DeleteCategory(ctx, arg)
}

func (r *Repository) CountCategoryTransactions(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	return r.query.CountCategoryTransactions(ctx, categoryID)
}
