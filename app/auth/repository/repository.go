package repository

import (
	"context"
	"database/sql"
	"personal-finance/config"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateUserResetToken(ctx context.Context, arg UpdateUserResetTokenParams) error
	GetUserByResetToken(ctx context.Context, token sql.NullString) (GetUserByResetTokenRow, error)
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

func (r *Repository) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	return r.query.UpdateUserPassword(ctx, arg)
}

func (r *Repository) UpdateUserResetToken(ctx context.Context, arg UpdateUserResetTokenParams) error {
	return r.query.UpdateUserResetToken(ctx, arg)
}

func (r *Repository) GetUserByResetToken(ctx context.Context, token sql.NullString) (GetUserByResetTokenRow, error) {
	return r.query.GetUserByResetToken(ctx, token)
}
