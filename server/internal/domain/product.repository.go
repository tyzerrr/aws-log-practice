package domain

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]*entity.Product, error)
}
