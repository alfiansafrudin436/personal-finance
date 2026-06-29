package repository

import (
	"context"
	"personal-finance/config"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	return r.query.GetUserByEmail(ctx, email)
}

func (r *Repository) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	return r.query.CreateUser(ctx, arg)
}
