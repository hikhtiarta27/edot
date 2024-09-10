package usecase

import (
	"context"
	"shop/model"
	"shop/repository"
	"shop/usecase/shop"
)

type Shop interface {
	List(ctx context.Context) ([]shop.Shop, error)
	Create(ctx context.Context, param *shop.CreateRequest) (*shop.Shop, error)
	AssignWarehouse(ctx context.Context, param *shop.AssignWarehouseRequest) (*shop.ShopWarehouse, error)
}

type shopUsecase struct {
	shopRepo      repository.Shop
	warehouseRepo repository.Warehouse
}

func NewShop(
	shopRepo repository.Shop,
	warehouseRepo repository.Warehouse,
) Shop {
	return &shopUsecase{
		shopRepo:      shopRepo,
		warehouseRepo: warehouseRepo,
	}
}

func (s *shopUsecase) List(ctx context.Context) ([]shop.Shop, error) {

	shopModels, err := s.shopRepo.Select(ctx)
	if err != nil {
		return nil, err
	}

	var shops = make([]shop.Shop, 0)

	for _, shp := range shopModels {
		shops = append(shops, shop.Shop{
			ID:        shp.ID,
			Name:      shp.Name,
			CreatedAt: shp.CreatedAt,
			UpdatedAt: shp.UpdatedAt,
		})
	}

	return shops, nil
}

func (s *shopUsecase) Create(ctx context.Context, param *shop.CreateRequest) (*shop.Shop, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	shopModel, err := model.NewShop(param.Name)
	if err != nil {
		return nil, err
	}

	err = s.shopRepo.Create(ctx, shopModel)
	if err != nil {
		return nil, err
	}

	return &shop.Shop{
		ID:        shopModel.ID,
		Name:      shopModel.Name,
		CreatedAt: shopModel.CreatedAt,
		UpdatedAt: shopModel.UpdatedAt,
	}, nil

}

func (s *shopUsecase) AssignWarehouse(ctx context.Context, param *shop.AssignWarehouseRequest) (*shop.ShopWarehouse, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	shopModel, err := s.shopRepo.Get(ctx, &model.GetShop{
		ID: param.ID,
	})

	if err != nil {
		return nil, err
	}

	if shopModel == nil {
		return nil, model.ErrInvalidShopID
	}

	// validate warehouse id
	for _, wrh := range param.WarehouseID {
		_, err = s.warehouseRepo.Get(ctx, &model.GetWarehouse{
			ID: wrh,
		})

		if err != nil {
			return nil, err
		}
	}

	return &shop.ShopWarehouse{
		Shop: &shop.Shop{
			ID:        shopModel.ID,
			Name:      shopModel.Name,
			CreatedAt: shopModel.CreatedAt,
			UpdatedAt: shopModel.UpdatedAt,
		},
		Warehouse: param.WarehouseID,
	}, nil
}
