package repository

import (
	"context"
	"warehouse/model"

	"gorm.io/gorm"
)

type WarehouseTransfer interface {
	Select(ctx context.Context, param *model.SelectWarehouseTransfer) ([]model.WarehouseTransfer, error)
	Create(ctx context.Context, warehouseTransfer *model.WarehouseTransfer) error
}

type warehouseTransferRepo struct {
	db *gorm.DB
}

func NewWarehouseTransfer(
	db *gorm.DB,
) WarehouseTransfer {
	return &warehouseTransferRepo{
		db: db,
	}
}

func (r warehouseTransferRepo) Select(ctx context.Context, param *model.SelectWarehouseTransfer) ([]model.WarehouseTransfer, error) {

	var warehouseTransfers []model.WarehouseTransfer

	q := r.db.
		WithContext(ctx)

	err := q.
		Find(&warehouseTransfers).
		Error

	if err != nil {
		return nil, err
	}

	return warehouseTransfers, nil
}

func (r warehouseTransferRepo) Create(ctx context.Context, warehouseTransfer *model.WarehouseTransfer) error {

	err := r.db.
		WithContext(ctx).
		Create(warehouseTransfer).
		Error

	if err != nil {
		return err
	}

	return nil
}
