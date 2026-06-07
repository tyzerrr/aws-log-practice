package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	ReadOnlyTransaction(context.Context, func(context.Context, pgx.Tx) error) error
	ReadWriteTransaction(context.Context, func(context.Context, pgx.Tx) error) error
	Close()
}

type db struct {
	pool *pgxpool.Pool
}

func newDB(ctx context.Context, databaseURL string) (DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &db{
		pool: pool,
	}, nil
}

func (d *db) ReadOnlyTransaction(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	return d.run(ctx, fn, pgx.TxOptions{AccessMode: pgx.ReadOnly})
}

func (d *db) ReadWriteTransaction(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	return d.run(ctx, fn, pgx.TxOptions{AccessMode: pgx.ReadWrite})
}

func (d *db) run(ctx context.Context, fn func(context.Context, pgx.Tx) error, opts pgx.TxOptions) error {
	tx, err := d.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true
	return nil
}

func (d *db) Close() {
	d.pool.Close()
}
