package domain

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductRepository interface {
	CreateOne(ctx context.Context, product *entity.Product) (*entity.Product, error)
	FindAllActiveProducts(ctx context.Context) ([]*entity.Product, error)
}
