package handler

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	productv1 "github.com/tyzerrr/aws-log-practice/server/gen/product/v1"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/valueobject"
	"github.com/tyzerrr/aws-log-practice/server/internal/usecase"
)

type ProductHandler struct {
	logger  *slog.Logger
	usecase usecase.ProductUsecase
}

func NewProductHandler(
	logger *slog.Logger,
	usecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *ProductHandler) RegisterProduct(ctx context.Context, req *productv1.RegisterProductRequest) (*productv1.RegisterProductResponse, error) {
	if req == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("RegisterProducts req cannot be nil"))
	}
	name, err := valueobject.NewProductName(req.Product.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	priceAmount, err := valueobject.NewProductAmount(req.Product.PriceAmount)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	currencyCode, err := valueobject.NewCurrency(req.Product.CurrencyCode)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	newProduct, err := h.usecase.CreateProduct(ctx, &entity.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: req.Product.Description,
		PriceAmount: priceAmount,
		Currency:    currencyCode,
		Archived:    false,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &productv1.RegisterProductResponse{
		Product: toProtoProduct(newProduct),
	}, nil
}

func (h *ProductHandler) GetActiveProductsCount(ctx context.Context, req *productv1.GetActiveProductsCountRequest) (*productv1.GetActiveProductsCountResponse, error) {
	if req == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("GetActiveProductsCount request cannot be nil"))
	}

	count, err := h.usecase.GetActiveProductsCount(ctx)
	if err != nil {
		h.logger.Error("failed to get active products count", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get active products count"))
	}
	return &productv1.GetActiveProductsCountResponse{
		Counts: count,
	}, nil
}

func (h *ProductHandler) ListActiveProducts(ctx context.Context, req *productv1.ListActiveProductsRequest) (*productv1.ListActiveProductsResponse, error) {
	if req == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("ListActiveProducts request is nil"))
	}

	products, err := h.usecase.ListActiveProducts(ctx)
	if err != nil {
		h.logger.Error("failed to list products", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to list products"))
	}

	out := make([]*productv1.Product, 0, len(products))
	for _, p := range products {
		out = append(out, toProtoProduct(p))
	}
	return &productv1.ListActiveProductsResponse{Products: out}, nil
}

func toProtoProduct(p *entity.Product) *productv1.Product {
	return &productv1.Product{
		Id:           p.ID.String(),
		Name:         p.Name.Value(),
		Description:  p.Description,
		PriceAmount:  int32(p.PriceAmount.Value()),
		CurrencyCode: p.Currency.Value(),
		IsActive:     !p.Archived,
	}
}
