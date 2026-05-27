package infrastructure

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductRepository struct {
	querier *sqlc.Queries
}

func NewProductRepository(db sqlc.DBTX) domain.ProductRepository {
	return &ProductRepository{
		querier: sqlc.New(db),
	}
}

func (pr *ProductRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	return nil, nil
}
