package infrastructure

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductRepository struct {
	db *db.DBPool
}

func NewProductRepository(db *db.DBPool) domain.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	return nil, nil
}
