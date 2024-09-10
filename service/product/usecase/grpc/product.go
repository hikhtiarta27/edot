package grpc

import (
	"context"
	"product/model"
	"product/repository"
	"proto_buffer/product"

	"github.com/oklog/ulid/v2"
)

type ProductGrpc struct {
	product.UnimplementedProductServiceServer
	productRepo repository.Product
}

func NewProduct(
	productRepo repository.Product,
) *ProductGrpc {
	return &ProductGrpc{
		productRepo: productRepo,
	}
}

func (s ProductGrpc) Get(ctx context.Context, param *product.GetRequest) (*product.Product, error) {

	id, err := ulid.Parse(param.Id)

	if err != nil {
		return nil, model.ErrInvalidID
	}

	productModel, err := s.productRepo.Get(ctx, &model.GetProduct{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	if productModel == nil {
		return nil, model.ErrProductNotFound
	}

	return &product.Product{
		Id:        productModel.ID.String(),
		Slug:      productModel.Slug,
		Name:      productModel.Name,
		Price:     productModel.Price,
		CreatedAt: productModel.CreatedAt.Unix(),
	}, nil
}
