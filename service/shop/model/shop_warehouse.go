package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type ShopWarehouse struct {
	ID          ulid.ULID
	ShopID      ulid.ULID
	WarehouseID ulid.ULID
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt
}

func (ShopWarehouse) TableName() string {
	return "shop_warehouses"
}

func NewShopWarehouse(
	shopID ulid.ULID,
	warehouseID ulid.ULID,
) (*ShopWarehouse, error) {
	instance := &ShopWarehouse{
		ID:          ulid.Make(),
		WarehouseID: warehouseID,
		ShopID:      shopID,
	}

	return instance, nil
}
