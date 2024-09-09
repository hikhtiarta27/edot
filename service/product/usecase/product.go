package usecase

import (
	"context"
	"product/model"
	"product/repository"
	"product/usecase/product"

	"github.com/oklog/ulid/v2"
)

type Product interface {
	List(ctx context.Context, param *product.ListRequest) ([]product.Product, error)
	Create(ctx context.Context, param *product.CreateRequest) (*product.Product, error)
}

type productUsecase struct {
	productRepo repository.Product
	stockRepo   repository.Stock
}

func NewProduct(
	productRepo repository.Product,
	stockRepo repository.Stock,
) Product {
	return &productUsecase{
		productRepo: productRepo,
		stockRepo:   stockRepo,
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

	var (
		productRes = make([]product.Product, 0)
		productIDs []ulid.ULID
	)

	for _, prd := range products {

		productIDs = append(productIDs, prd.ID)

		productRes = append(productRes, product.Product{
			ID:        prd.ID,
			Slug:      prd.Slug,
			Name:      prd.Name,
			CreatedAt: prd.CreatedAt,
		})
	}

	stocks, err := s.stockRepo.Select(ctx, &model.SelectStock{
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, err
	}

	stockByID := stocks.ToMap()

	for i, prd := range productRes {

		stock, ok := stockByID[prd.ID]
		if !ok {
			continue
		}

		productRes[i].AvailableStock = stock.AvailableStock
		productRes[i].ReservedStock = stock.ReservedStock
	}

	return productRes, nil
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
		ID:        prd.ID,
		Slug:      prd.Slug,
		Name:      prd.Name,
		CreatedAt: prd.CreatedAt,
	}, nil
}
