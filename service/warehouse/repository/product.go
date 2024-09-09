package repository

import (
	"context"
	"proto_buffer/product"
	"time"
	"warehouse/model"

	"github.com/oklog/ulid/v2"
)

type Product interface {
	Get(ctx context.Context, param *model.GetProduct) (*model.Product, error)
	UpdateStock(ctx context.Context, param *model.UpdateStockProduct) (*model.Product, error)
}

type productRepo struct {
	productGrpc product.ProductServiceClient
}

func NewProduct(
	productGrpc product.ProductServiceClient,
) Product {
	return &productRepo{
		productGrpc: productGrpc,
	}
}

func (r productRepo) Get(ctx context.Context, param *model.GetProduct) (*model.Product, error) {
	productModel, err := r.productGrpc.Get(ctx, &product.GetRequest{
		Id: param.ID.String(),
	})

	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:             ulid.MustParse(productModel.Id),
		Slug:           productModel.Slug,
		Name:           productModel.Name,
		AvailableStock: productModel.AvailableStock,
		ReservedStock:  productModel.ReservedStock,
		CreatedAt:      time.Unix(productModel.CreatedAt, 0),
	}, nil
}

func (r productRepo) UpdateStock(ctx context.Context, param *model.UpdateStockProduct) (*model.Product, error) {
	productModel, err := r.productGrpc.UpdateStock(ctx, &product.UpdateStockRequest{
		Id:             param.ID.String(),
		AvailableStock: param.AvailableStock,
	})

	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:             ulid.MustParse(productModel.Id),
		Slug:           productModel.Slug,
		Name:           productModel.Name,
		AvailableStock: productModel.AvailableStock,
		ReservedStock:  productModel.ReservedStock,
		CreatedAt:      time.Unix(productModel.CreatedAt, 0),
	}, nil
}
