package repository

import (
	"context"
	"product/model"
	"proto_buffer/shop"
	"time"

	"github.com/oklog/ulid/v2"
)

type Shop interface {
	Get(ctx context.Context, param *model.GetShop) (*model.Shop, error)
}

type shopRepo struct {
	shopGrpc shop.ShopServiceClient
}

func NewShop(
	shopGrpc shop.ShopServiceClient,
) Shop {
	return &shopRepo{
		shopGrpc: shopGrpc,
	}
}

func (r shopRepo) Get(ctx context.Context, param *model.GetShop) (*model.Shop, error) {

	shopModel, err := r.shopGrpc.Get(ctx, &shop.GetRequest{
		Id: param.ShopID.String(),
	})

	if err != nil {
		return nil, err
	}

	res := &model.Shop{
		ID:        ulid.MustParse(shopModel.Id),
		Name:      shopModel.Name,
		CreatedAt: time.Unix(shopModel.CreatedAt, 0),
	}

	for _, sw := range shopModel.Warehouse {
		res.Warehouse = append(res.Warehouse, ulid.MustParse(sw))
	}

	return res, nil
}
