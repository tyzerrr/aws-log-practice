package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (tx *TransactionManager) ReadOnlyTransaction(
	ctx context.Context,
	fn func(context.Context, sqlc.DBTX) error,
) error {
	return tx.run(ctx, fn, pgx.TxOptions{AccessMode: pgx.ReadOnly})
}

func (tx *TransactionManager) ReadWriteTransaction(
	ctx context.Context,
	fn func(context.Context, sqlc.DBTX) error,
) error {
	return tx.run(ctx, fn, pgx.TxOptions{AccessMode: pgx.ReadWrite})
}

func (txm *TransactionManager) run(
	ctx context.Context,
	fn func(context.Context, sqlc.DBTX) error,
	opts pgx.TxOptions,
) error {
	tx, err := txm.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	if err = fn(ctx, tx); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true
	return nil
}
