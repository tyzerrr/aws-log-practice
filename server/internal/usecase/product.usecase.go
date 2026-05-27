package usecase

import (
	"context"
	"log/slog"

	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/transaction"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
)

type ProductUsecase interface {
	ListProducts(ctx context.Context) ([]*entity.Product, error)
}

type ProductRepositoryFactory func(db sqlc.DBTX) domain.ProductRepository

type productUsecase struct {
	logger    *slog.Logger
	newRepo   ProductRepositoryFactory
	txManager transaction.Tx
}

func NewProductUsecase(
	logger *slog.Logger,
	newRepo ProductRepositoryFactory,
	txManager transaction.Tx,
) ProductUsecase {
	return &productUsecase{
		logger:    logger,
		newRepo:   newRepo,
		txManager: txManager,
	}
}

func (uc *productUsecase) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	// NOTE: memory allocation strategy is not good...
	var products []*entity.Product
	if err := uc.txManager.ReadOnlyTransaction(ctx, func(ctx context.Context, tx sqlc.DBTX) error {
		var err error
		products, err = uc.newRepo(tx).FindAll(ctx)
		return err
	}); err != nil {
		return nil, err
	}
	return products, nil
}
