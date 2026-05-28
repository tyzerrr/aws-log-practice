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
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetActiveProductsCount(ctx context.Context) (int64, error)
	ListActiveProducts(ctx context.Context) ([]*entity.Product, error)
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

func (uc *productUsecase) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var newProduct *entity.Product
	if err := uc.txManager.ReadWriteTransaction(ctx, func(ctx context.Context, tx sqlc.DBTX) error {
		var err error
		newProduct, err = uc.newRepo(tx).CreateOne(ctx, product)
		return err
	}); err != nil {
		return nil, err
	}
	return newProduct, nil
}

func (uc *productUsecase) GetActiveProductsCount(ctx context.Context) (int64, error) {
	var count int64
	if err := uc.txManager.ReadOnlyTransaction(ctx, func(ctx context.Context, tx sqlc.DBTX) error {
		products, err := uc.newRepo(tx).FindAllActiveProducts(ctx)
		count = int64(len(products))
		return err
	}); err != nil {
		return -1, err
	}
	return count, nil
}

func (uc *productUsecase) ListActiveProducts(ctx context.Context) ([]*entity.Product, error) {
	// NOTE: memory allocation strategy is not good...
	var products []*entity.Product
	if err := uc.txManager.ReadOnlyTransaction(ctx, func(ctx context.Context, tx sqlc.DBTX) error {
		var err error
		products, err = uc.newRepo(tx).FindAllActiveProducts(ctx)
		return err
	}); err != nil {
		return nil, err
	}
	return products, nil
}
