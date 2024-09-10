package grpc

import (
	"context"
	"proto_buffer/warehouse"
	"warehouse/model"
	"warehouse/repository"

	"github.com/oklog/ulid/v2"
)

type WarehouseGrpc struct {
	warehouse.UnimplementedWarehouseServiceServer
	warehouseRepo repository.Warehouse
}

func NewWarehouse(
	warehouseRepo repository.Warehouse,
) *WarehouseGrpc {
	return &WarehouseGrpc{
		warehouseRepo: warehouseRepo,
	}
}

func (s WarehouseGrpc) Get(ctx context.Context, param *warehouse.GetRequest) (*warehouse.Warehouse, error) {

	id, err := ulid.Parse(param.Id)
	if err != nil {
		return nil, model.ErrInvalidUlid
	}

	wrh, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID: id,
	})

	if wrh == nil {
		return nil, model.ErrWarehouseNotFound
	}

	if err != nil {
		return nil, err
	}

	return &warehouse.Warehouse{
		Id:        wrh.ID.String(),
		Name:      wrh.Name,
		Status:    string(wrh.Status),
		CreatedAt: wrh.CreatedAt.Unix(),
	}, nil
}
