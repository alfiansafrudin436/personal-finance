package repository

import (
	"context"
	"fmt"

	"personal-finance/config"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BalanceEffect is a single signed adjustment to an account balance.
// Delta is a decimal string so it reaches NUMERIC(15,2) without going
// through float64; a debit is expressed as a negative Delta.
type BalanceEffect struct {
	AccountID uuid.UUID
	Delta     string
}

type TransactionRepository interface {
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]ListTransactionsRow, error)
	CountTransactions(ctx context.Context, arg CountTransactionsParams) (int64, error)
	GetTransactionByID(ctx context.Context, arg GetTransactionByIDParams) (GetTransactionByIDRow, error)
	GetTransactionRaw(ctx context.Context, arg GetTransactionRawParams) (GetTransactionRawRow, error)
	GetAccountForTransaction(ctx context.Context, arg GetAccountForTransactionParams) (GetAccountForTransactionRow, error)
	GetCategoryForTransaction(ctx context.Context, arg GetCategoryForTransactionParams) (GetCategoryForTransactionRow, error)

	// CreateWithBalance, UpdateWithBalance and DeleteWithBalance each run the
	// write and its balance adjustments inside one database transaction, so an
	// account balance can never drift from the transactions behind it.
	CreateWithBalance(ctx context.Context, arg CreateTransactionParams, effects []BalanceEffect) (Transaction, error)
	UpdateWithBalance(ctx context.Context, arg UpdateTransactionParams, effects []BalanceEffect) (Transaction, error)
	DeleteWithBalance(ctx context.Context, arg DeleteTransactionParams, effects []BalanceEffect) (int64, error)
}

type Repository struct {
	db    *sqlx.DB
	query *Queries
}

func NewRepository() *Repository {
	db := config.Application.DB
	return &Repository{
		db:    db,
		query: New(db),
	}
}

func (r *Repository) ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]ListTransactionsRow, error) {
	return r.query.ListTransactions(ctx, arg)
}

func (r *Repository) CountTransactions(ctx context.Context, arg CountTransactionsParams) (int64, error) {
	return r.query.CountTransactions(ctx, arg)
}

func (r *Repository) GetTransactionByID(ctx context.Context, arg GetTransactionByIDParams) (GetTransactionByIDRow, error) {
	return r.query.GetTransactionByID(ctx, arg)
}

func (r *Repository) GetTransactionRaw(ctx context.Context, arg GetTransactionRawParams) (GetTransactionRawRow, error) {
	return r.query.GetTransactionRaw(ctx, arg)
}

func (r *Repository) GetAccountForTransaction(ctx context.Context, arg GetAccountForTransactionParams) (GetAccountForTransactionRow, error) {
	return r.query.GetAccountForTransaction(ctx, arg)
}

func (r *Repository) GetCategoryForTransaction(ctx context.Context, arg GetCategoryForTransactionParams) (GetCategoryForTransactionRow, error) {
	return r.query.GetCategoryForTransaction(ctx, arg)
}

// withTx runs fn against a transactional Queries, committing on success and
// rolling back on any error or panic.
func (r *Repository) withTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := fn(r.query.WithTx(tx.Tx)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

// applyEffects adjusts every account balance named in effects.
func applyEffects(ctx context.Context, q *Queries, userID uuid.UUID, effects []BalanceEffect) error {
	for _, e := range effects {
		if err := q.AdjustAccountBalance(ctx, AdjustAccountBalanceParams{
			ID:     e.AccountID,
			UserID: userID,
			Delta:  e.Delta,
		}); err != nil {
			return fmt.Errorf("adjust balance for account %s: %w", e.AccountID, err)
		}
	}
	return nil
}

func (r *Repository) CreateWithBalance(ctx context.Context, arg CreateTransactionParams, effects []BalanceEffect) (Transaction, error) {
	var created Transaction
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		created, err = q.CreateTransaction(ctx, arg)
		if err != nil {
			return err
		}
		return applyEffects(ctx, q, arg.UserID, effects)
	})
	return created, err
}

func (r *Repository) UpdateWithBalance(ctx context.Context, arg UpdateTransactionParams, effects []BalanceEffect) (Transaction, error) {
	var updated Transaction
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		updated, err = q.UpdateTransaction(ctx, arg)
		if err != nil {
			return err
		}
		return applyEffects(ctx, q, arg.UserID, effects)
	})
	return updated, err
}

func (r *Repository) DeleteWithBalance(ctx context.Context, arg DeleteTransactionParams, effects []BalanceEffect) (int64, error) {
	var rows int64
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		rows, err = q.DeleteTransaction(ctx, arg)
		if err != nil {
			return err
		}
		if rows == 0 {
			return nil
		}
		return applyEffects(ctx, q, arg.UserID, effects)
	})
	return rows, err
}
