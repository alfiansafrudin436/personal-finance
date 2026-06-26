package repository

import (
	"context"
	"personal-finance/config"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	GetUserByID(ctx context.Context, id int32) (GetUserByIDRow, error)
	UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) error
	DeactivateUser(ctx context.Context, id int32) error
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

func (r *Repository) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	return r.query.GetAllUsers(ctx)
}

func (r *Repository) GetUserByID(ctx context.Context, id int32) (GetUserByIDRow, error) {
	return r.query.GetUserByID(ctx, id)
}

func (r *Repository) UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) error {
	return r.query.UpdateUserUsername(ctx, arg)
}

func (r *Repository) DeactivateUser(ctx context.Context, id int32) error {
	return r.query.DeactivateUser(ctx, id)
}
