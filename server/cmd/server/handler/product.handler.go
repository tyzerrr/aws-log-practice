package handler

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	productv1 "github.com/tyzerrr/aws-log-practice/server/gen/product/v1"
	"github.com/tyzerrr/aws-log-practice/server/internal/domain/entity"
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

func (h *ProductHandler) RegisterProducts(ctx context.Context, req *productv1.RegisterProductsRequest) (*productv1.RegisterProductsResponse, error) {
	return nil, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	if req == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request is nil"))
	}

	products, err := h.usecase.ListProducts(ctx)
	if err != nil {
		h.logger.Error("failed to list products", slog.String("error", err.Error()))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to list products"))
	}

	out := make([]*productv1.Product, 0, len(products))
	for _, p := range products {
		out = append(out, toProtoProduct(p))
	}
	return &productv1.ListProductsResponse{Products: out}, nil
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
