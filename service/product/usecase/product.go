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
	shopRepo    repository.Shop
}

func NewProduct(
	productRepo repository.Product,
	stockRepo repository.Stock,
	shopRepo repository.Shop,
) Product {
	return &productUsecase{
		productRepo: productRepo,
		stockRepo:   stockRepo,
		shopRepo:    shopRepo,
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
			Price:     prd.Price,
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

	stockByID := stocks.MapByProductID()

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

	shop, err := s.shopRepo.Get(ctx, &model.GetShop{
		ShopID: param.ShopID,
	})

	if err != nil {
		return nil, err
	}

	isShopWarehouseExist := false

	for _, sw := range shop.Warehouse {
		if sw == param.WarehouseID {
			isShopWarehouseExist = true
			break
		}
	}

	if !isShopWarehouseExist {
		return nil, model.ErrShopWarehouseNotFound
	}

	prd, err := model.NewProduct(param.Name, param.Price, param.ShopID)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.Create(ctx, prd)
	if err != nil {
		return nil, err
	}

	_, err = s.stockRepo.Create(ctx, &model.CreateStock{
		ProductID:   prd.ID,
		WarehouseID: param.WarehouseID,
		Stock:       param.Stock,
	})

	if err != nil {
		return nil, err
	}

	return &product.Product{
		ID:        prd.ID,
		Slug:      prd.Slug,
		Name:      prd.Name,
		Price:     prd.Price,
		CreatedAt: prd.CreatedAt,
	}, nil
}
