package usecase

import (
	"context"
	"log/slog"

	"github.com/tyzerrr/aws-log-practice/server/internal/domain"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context) ([]*entity.Product, error)
}

type productUsecase struct {
	logger            *slog.Logger
	productRepository domain.ProductRepository
}

func NewProductUsecase(
	logger *slog.Logger,
	productRepository domain.ProductRepository,
) ProductUsecase {
	return &productUsecase{
		logger:            logger,
		productRepository: productRepository,
	}
}

func (uc *productUsecase) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	return nil, nil
}
