package repository

import (
	"context"
	"personal-finance/config"

	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	GetAccountByID(ctx context.Context, arg GetAccountByIDParams) (Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	SetAccountArchived(ctx context.Context, arg SetAccountArchivedParams) error
	DeleteAccount(ctx context.Context, arg DeleteAccountParams) (int64, error)
	CountAccountTransactions(ctx context.Context, arg CountAccountTransactionsParams) (int64, error)
	GetAccountsSummary(ctx context.Context, userID uuid.UUID) (GetAccountsSummaryRow, error)
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

func (r *Repository) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	return r.query.CreateAccount(ctx, arg)
}

func (r *Repository) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	return r.query.ListAccounts(ctx, arg)
}

func (r *Repository) GetAccountByID(ctx context.Context, arg GetAccountByIDParams) (Account, error) {
	return r.query.GetAccountByID(ctx, arg)
}

func (r *Repository) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	return r.query.UpdateAccount(ctx, arg)
}

func (r *Repository) SetAccountArchived(ctx context.Context, arg SetAccountArchivedParams) error {
	return r.query.SetAccountArchived(ctx, arg)
}

func (r *Repository) DeleteAccount(ctx context.Context, arg DeleteAccountParams) (int64, error) {
	return r.query.DeleteAccount(ctx, arg)
}

func (r *Repository) CountAccountTransactions(ctx context.Context, arg CountAccountTransactionsParams) (int64, error) {
	return r.query.CountAccountTransactions(ctx, arg)
}

func (r *Repository) GetAccountsSummary(ctx context.Context, userID uuid.UUID) (GetAccountsSummaryRow, error) {
	return r.query.GetAccountsSummary(ctx, userID)
}
