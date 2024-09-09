package repository

import (
	"context"
	"database/sql"
	"warehouse/model"

	"gorm.io/gorm"
)

type WarehouseTransfer interface {
	Select(ctx context.Context, param *model.SelectWarehouseTransfer) ([]model.WarehouseTransfer, error)
	Create(ctx context.Context, warehouseTransfer *model.WarehouseTransfer) error
}

type warehouseTransferRepo struct {
	db          *gorm.DB
	productRepo Product
}

func NewWarehouseTransfer(
	db *gorm.DB,
	productRepo Product,
) WarehouseTransfer {
	return &warehouseTransferRepo{
		db:          db,
		productRepo: productRepo,
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

	tx := r.db.Begin(&sql.TxOptions{})

	err := tx.
		WithContext(ctx).
		Create(warehouseTransfer).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = r.productRepo.UpdateStock(ctx, &model.UpdateStockProduct{
		ID:             warehouseTransfer.ProductID,
		AvailableStock: warehouseTransfer.Stock,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
