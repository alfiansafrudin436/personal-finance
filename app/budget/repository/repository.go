package repository

import (
	"context"
	"personal-finance/config"
)

type BudgetRepository interface {
	UpsertBudget(ctx context.Context, arg UpsertBudgetParams) (Budget, error)
	ListBudgets(ctx context.Context, arg ListBudgetsParams) ([]ListBudgetsRow, error)
	GetBudgetByID(ctx context.Context, arg GetBudgetByIDParams) (Budget, error)
	UpdateBudget(ctx context.Context, arg UpdateBudgetParams) (Budget, error)
	DeleteBudget(ctx context.Context, arg DeleteBudgetParams) (int64, error)
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

func (r *Repository) UpsertBudget(ctx context.Context, arg UpsertBudgetParams) (Budget, error) {
	return r.query.UpsertBudget(ctx, arg)
}

func (r *Repository) ListBudgets(ctx context.Context, arg ListBudgetsParams) ([]ListBudgetsRow, error) {
	return r.query.ListBudgets(ctx, arg)
}

func (r *Repository) GetBudgetByID(ctx context.Context, arg GetBudgetByIDParams) (Budget, error) {
	return r.query.GetBudgetByID(ctx, arg)
}

func (r *Repository) UpdateBudget(ctx context.Context, arg UpdateBudgetParams) (Budget, error) {
	return r.query.UpdateBudget(ctx, arg)
}

func (r *Repository) DeleteBudget(ctx context.Context, arg DeleteBudgetParams) (int64, error) {
	return r.query.DeleteBudget(ctx, arg)
}
