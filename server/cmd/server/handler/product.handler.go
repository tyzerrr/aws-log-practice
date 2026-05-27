package handler

import (
	"context"
	"log/slog"

	productv1 "github.com/tyzerrr/aws-log-practice/server/gen/product/v1"
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
	return nil, nil
}
