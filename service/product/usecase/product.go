package usecase

import (
	"context"
	"product/model"
	"product/repository"
	"product/usecase/product"
)

type Product interface {
	List(ctx context.Context, param *product.ListRequest) ([]product.Product, error)
	Create(ctx context.Context, param *product.CreateRequest) (*product.Product, error)
}

type productUsecase struct {
	productRepo repository.Product
}

func NewProduct(
	productRepo repository.Product,

) Product {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (s productUsecase) List(ctx context.Context, param *product.ListRequest) ([]product.Product, error) {

	products, err := s.productRepo.Select(ctx, &model.SelectProduct{
		UseSearchEngine: true,
		Keyword:         param.Keyword,
	})
	if err != nil {
		return nil, err
	}

	var res = make([]product.Product, 0)

	for _, prd := range products {
		res = append(res, product.Product{
			ID:             prd.ID,
			Slug:           prd.Slug,
			Name:           prd.Name,
			AvailableStock: prd.AvailableStock,
			ReservedStock:  prd.ReservedStock,
			CreatedAt:      prd.CreatedAt,
		})
	}

	return res, nil
}

func (s productUsecase) Create(ctx context.Context, param *product.CreateRequest) (*product.Product, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	prd, err := model.NewProduct(param.Name)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.Create(ctx, prd)
	if err != nil {
		return nil, err
	}

	return &product.Product{
		ID:             prd.ID,
		Slug:           prd.Slug,
		Name:           prd.Name,
		AvailableStock: prd.AvailableStock,
		ReservedStock:  prd.ReservedStock,
		CreatedAt:      prd.CreatedAt,
	}, nil
}
