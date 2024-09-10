package repository

import (
	"context"
	"shared"
	"time"
	"warehouse/model"

	"gorm.io/gorm"
)

type Warehouse interface {
	Get(ctx context.Context, param *model.GetWarehouse) (*model.Warehouse, error)
	Select(ctx context.Context, param *model.SelectWarehouse) ([]model.Warehouse, error)
	Create(ctx context.Context, warehouse *model.Warehouse) error
	Update(ctx context.Context, warehouse *model.Warehouse) error
}

type warehouseRepo struct {
	db *gorm.DB
}

func NewWarehouse(
	db *gorm.DB,
) Warehouse {
	return &warehouseRepo{
		db: db,
	}
}

func (r warehouseRepo) Select(ctx context.Context, param *model.SelectWarehouse) ([]model.Warehouse, error) {

	var warehouses []model.Warehouse

	q := r.db.
		WithContext(ctx)

	err := q.
		Find(&warehouses).
		Error

	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r warehouseRepo) Create(ctx context.Context, warehouse *model.Warehouse) error {
	return r.db.
		WithContext(ctx).
		Create(&warehouse).
		Error
}

func (r warehouseRepo) Get(ctx context.Context, param *model.GetWarehouse) (*model.Warehouse, error) {

	warehouse := &model.Warehouse{}

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ID) {
		q = q.Where("id = ?", param.ID)
	}

	if param.Status != "" {
		q = q.Where("status = ?", param.Status)
	}

	err := q.
		First(warehouse).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (r warehouseRepo) Update(ctx context.Context, warehouse *model.Warehouse) error {

	now := time.Now()
	warehouse.UpdatedAt = &now

	return r.db.
		WithContext(ctx).
		Updates(warehouse).
		Error
}
