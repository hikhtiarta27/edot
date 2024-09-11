package grpc

import (
	"context"
	"proto_buffer/shop"
	"shop/model"
	"shop/repository"

	"github.com/oklog/ulid/v2"
)

type ShopGrpc struct {
	shop.UnimplementedShopServiceServer
	shopRepo          repository.Shop
	shopWarehouseRepo repository.ShopWarehouse
}

func NewShop(
	shopRepo repository.Shop,
	shopWarehouseRepo repository.ShopWarehouse,
) *ShopGrpc {
	return &ShopGrpc{
		shopRepo:          shopRepo,
		shopWarehouseRepo: shopWarehouseRepo,
	}
}

func (s ShopGrpc) Get(ctx context.Context, param *shop.GetRequest) (*shop.Shop, error) {

	id, err := ulid.Parse(param.Id)
	if err != nil {
		return nil, model.ErrInvalidUlid
	}

	shopModel, err := s.shopRepo.Get(ctx, &model.GetShop{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	if shopModel == nil {
		return nil, model.ErrShopNotFound
	}

	res := &shop.Shop{
		Id:        shopModel.ID.String(),
		Name:      shopModel.Name,
		CreatedAt: shopModel.CreatedAt.Unix(),
	}

	shopWarehouses, err := s.shopWarehouseRepo.Select(ctx, &model.SelectShopWarehouse{
		ShopID: id,
	})

	if err != nil {
		return nil, err
	}

	for _, shopWarehouse := range shopWarehouses {
		res.Warehouse = append(res.Warehouse, shopWarehouse.WarehouseID.String())
	}

	return res, nil
}
