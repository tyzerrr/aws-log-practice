package infrastructure

import (
	"context"

	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db/sqlc"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/valueobject"
)

type ProductRepository struct {
	querier *sqlc.Queries
}

func NewProductRepository(db sqlc.DBTX) domain.ProductRepository {
	return &ProductRepository{
		querier: sqlc.New(db),
	}
}

func (pr *ProductRepository) CreateOne(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	model, err := pr.querier.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:         product.Name.Value(),
		Description:  product.Description,
		PriceAmount:  product.PriceAmount.Value(),
		CurrencyCode: product.Currency.Value(),
		IsActive:     !product.Archived,
	})
	if err != nil {
		return nil, err
	}
	return restoreProduct(&model)
}

func (pr *ProductRepository) FindAllActiveProducts(ctx context.Context) ([]*entity.Product, error) {
	return nil, nil
}

func restoreProduct(model *sqlc.Product) (*entity.Product, error) {
	name, err := valueobject.NewProductName(model.Name)
	if err != nil {
		return nil, err
	}

	priceAmount, err := valueobject.NewProductAmount(model.PriceAmount)
	if err != nil {
		return nil, err
	}

	currency, err := valueobject.NewCurrency(model.CurrencyCode)
	if err != nil {
		return nil, err
	}

	return entity.NewProduct(
		model.ID.Bytes,
		name,
		model.Description,
		priceAmount,
		currency,
		!model.IsActive,
	), nil
}
