package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
)

type DBPool struct {
	pool    *pgxpool.Pool
	logger  *slog.Logger
	Querier sqlc.Querier
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
		pool:    p,
		logger:  logger,
		Querier: sqlc.New(p),
	}, nil
}

func (db *DBPool) Close() {
	db.pool.Close()
}
