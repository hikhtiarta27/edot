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
	Detail(ctx context.Context, param *shop.DetailRequest) (*shop.ShopWarehouse, error)
}

type shopUsecase struct {
	shopRepo          repository.Shop
	shopWarehouseRepo repository.ShopWarehouse
	warehouseRepo     repository.Warehouse
}

func NewShop(
	shopRepo repository.Shop,
	shopWarehouseRepo repository.ShopWarehouse,
	warehouseRepo repository.Warehouse,
) Shop {
	return &shopUsecase{
		shopRepo:          shopRepo,
		warehouseRepo:     warehouseRepo,
		shopWarehouseRepo: shopWarehouseRepo,
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
		return nil, model.ErrShopNotFound
	}

	var shopWarehouses []model.ShopWarehouse

	// validate warehouse id
	for _, wrh := range param.WarehouseID {
		// TODO: change to outside for-loop
		_, err = s.warehouseRepo.Get(ctx, &model.GetWarehouse{
			ID: wrh,
		})

		if err != nil {
			return nil, err
		}

		// TODO: change to outside for-loop
		shopWarehouse, err := s.shopWarehouseRepo.Get(ctx, &model.GetShopWarehouse{
			ShopID:      param.ID,
			WarehouseID: wrh,
		})

		if err != nil {
			return nil, err
		}

		if shopWarehouse != nil {
			return nil, model.ErrDuplicateShopWarehouse
		}

		shopWarehouse, err = model.NewShopWarehouse(param.ID, wrh)
		if err != nil {
			return nil, err
		}

		shopWarehouses = append(shopWarehouses, *shopWarehouse)
	}

	err = s.shopWarehouseRepo.CreateBatch(ctx, shopWarehouses)
	if err != nil {
		return nil, err
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

func (s *shopUsecase) Detail(ctx context.Context, param *shop.DetailRequest) (*shop.ShopWarehouse, error) {

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
		return nil, model.ErrShopNotFound
	}

	res := &shop.ShopWarehouse{
		Shop: &shop.Shop{
			ID:        shopModel.ID,
			Name:      shopModel.Name,
			CreatedAt: shopModel.CreatedAt,
			UpdatedAt: shopModel.UpdatedAt,
		},
	}

	shopWarehouses, err := s.shopWarehouseRepo.Select(ctx, &model.SelectShopWarehouse{
		ShopID: param.ID,
	})

	if err != nil {
		return nil, err
	}

	for _, shopWarehouse := range shopWarehouses {
		res.Warehouse = append(res.Warehouse, shopWarehouse.WarehouseID)
	}

	return res, nil
}
