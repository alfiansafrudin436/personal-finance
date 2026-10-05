package repository

import (
	"context"
	"personal-finance/config"

	"github.com/google/uuid"
)

type ReportRepository interface {
	GetPeriodSummary(ctx context.Context, arg GetPeriodSummaryParams) (GetPeriodSummaryRow, error)
	GetSpendingByCategory(ctx context.Context, arg GetSpendingByCategoryParams) ([]GetSpendingByCategoryRow, error)
	GetMonthlyTrend(ctx context.Context, arg GetMonthlyTrendParams) ([]GetMonthlyTrendRow, error)
	GetAccountBalances(ctx context.Context, userID uuid.UUID) ([]GetAccountBalancesRow, error)
	GetTotalBalance(ctx context.Context, userID uuid.UUID) (string, error)
	GetRecentTransactions(ctx context.Context, arg GetRecentTransactionsParams) ([]GetRecentTransactionsRow, error)
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

func (r *Repository) GetPeriodSummary(ctx context.Context, arg GetPeriodSummaryParams) (GetPeriodSummaryRow, error) {
	return r.query.GetPeriodSummary(ctx, arg)
}

func (r *Repository) GetSpendingByCategory(ctx context.Context, arg GetSpendingByCategoryParams) ([]GetSpendingByCategoryRow, error) {
	return r.query.GetSpendingByCategory(ctx, arg)
}

func (r *Repository) GetMonthlyTrend(ctx context.Context, arg GetMonthlyTrendParams) ([]GetMonthlyTrendRow, error) {
	return r.query.GetMonthlyTrend(ctx, arg)
}

func (r *Repository) GetAccountBalances(ctx context.Context, userID uuid.UUID) ([]GetAccountBalancesRow, error) {
	return r.query.GetAccountBalances(ctx, userID)
}

func (r *Repository) GetTotalBalance(ctx context.Context, userID uuid.UUID) (string, error) {
	return r.query.GetTotalBalance(ctx, userID)
}

func (r *Repository) GetRecentTransactions(ctx context.Context, arg GetRecentTransactionsParams) ([]GetRecentTransactionsRow, error) {
	return r.query.GetRecentTransactions(ctx, arg)
}
