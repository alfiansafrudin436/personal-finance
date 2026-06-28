package repository

import (
	"context"
	"personal-finance/config"
)

type AuthRepository interface {
	GetUserByEmailOrUsername(ctx context.Context, email string) (GetUserByEmailOrUsernameRow, error)
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

func (r *Repository) GetUserByEmailOrUsername(ctx context.Context, email string) (GetUserByEmailOrUsernameRow, error) {
	return r.query.GetUserByEmailOrUsername(ctx, email)
}
