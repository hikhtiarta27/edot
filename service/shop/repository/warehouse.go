package repository

import (
	"context"
	"proto_buffer/warehouse"
	"shop/model"
	"time"

	"github.com/oklog/ulid/v2"
)

type Warehouse interface {
	Get(ctx context.Context, param *model.GetWarehouse) (*model.Warehouse, error)
}

type warehouseRepo struct {
	warehouseGrpc warehouse.WarehouseServiceClient
}

func NewWarehouse(
	warehouseGrpc warehouse.WarehouseServiceClient,
) Warehouse {
	return &warehouseRepo{
		warehouseGrpc: warehouseGrpc,
	}
}

func (r warehouseRepo) Get(ctx context.Context, param *model.GetWarehouse) (*model.Warehouse, error) {
	warehouseModel, err := r.warehouseGrpc.Get(ctx, &warehouse.GetRequest{
		Id: param.ID.String(),
	})

	if err != nil {
		return nil, err
	}

	return &model.Warehouse{
		ID:        ulid.MustParse(warehouseModel.Id),
		Name:      warehouseModel.Name,
		Status:    warehouseModel.Status,
		CreatedAt: time.Unix(warehouseModel.CreatedAt, 0),
	}, nil
}
