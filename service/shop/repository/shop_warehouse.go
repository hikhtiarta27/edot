package repository

import (
	"context"
	"shared"
	"shop/model"

	"gorm.io/gorm"
)

type ShopWarehouse interface {
	Get(ctx context.Context, param *model.GetShopWarehouse) (*model.ShopWarehouse, error)
	Select(ctx context.Context, param *model.SelectShopWarehouse) ([]model.ShopWarehouse, error)
	CreateBatch(ctx context.Context, shopWarehouses []model.ShopWarehouse) error
}

type shopWarehouseRepo struct {
	db *gorm.DB
}

func NewShopWarehouse(
	db *gorm.DB,
) ShopWarehouse {
	return &shopWarehouseRepo{
		db: db,
	}
}

func (r shopWarehouseRepo) Get(ctx context.Context, param *model.GetShopWarehouse) (*model.ShopWarehouse, error) {

	var shopWarehouse *model.ShopWarehouse

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ShopID) {
		q = q.Where("shop_id = ?", param.ShopID)
	}

	if !shared.IsZero(param.WarehouseID) {
		q = q.Where("warehouse_id = ?", param.WarehouseID)
	}

	err := q.
		First(&shopWarehouse).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return shopWarehouse, nil
}

func (r shopWarehouseRepo) Select(ctx context.Context, param *model.SelectShopWarehouse) ([]model.ShopWarehouse, error) {

	var shopWarehouses []model.ShopWarehouse

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ShopID) {
		q = q.Where("shop_id = ?", param.ShopID)
	}

	err := q.
		Find(&shopWarehouses).
		Error

	if err != nil {
		return nil, err
	}

	return shopWarehouses, nil
}

func (r shopWarehouseRepo) CreateBatch(ctx context.Context, shopWarehouses []model.ShopWarehouse) error {
	return r.db.
		WithContext(ctx).
		CreateInBatches(&shopWarehouses, 10).
		Error
}
