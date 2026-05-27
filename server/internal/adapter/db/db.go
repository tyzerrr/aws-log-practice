package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewDBPool(ctx context.Context, logger *slog.Logger, connectionString string) (*DBPool, error) {
	p, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		logger.Error("failed to initialize db pool", slog.String("error", err.Error()))
		return nil, err
	}
	if err := p.Ping(ctx); err != nil {
		p.Close()
		logger.Error("failed to ping to db pool, so close db pool", slog.String("error", err.Error()))
		return nil, err
	}

	return &DBPool{
		Pool:   p,
		logger: logger,
	}, nil
}

func (db *DBPool) Close() {
	db.Pool.Close()
}
