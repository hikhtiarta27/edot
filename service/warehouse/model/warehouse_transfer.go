package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type WarehouseTransfer struct {
	ID              ulid.ULID
	FromWarehouseID ulid.ULID
	ToWarehouseID   ulid.ULID
	ProductID       ulid.ULID
	Stock           uint64
	CreatedAt       time.Time
}

func (WarehouseTransfer) TableName() string {
	return "warehouse_transfers"
}

func NewWarehouseTransfer(
	fromID ulid.ULID,
	toID ulid.ULID,
	productID ulid.ULID,
	stock uint64,
) (*WarehouseTransfer, error) {
	return &WarehouseTransfer{
		ID:              ulid.Make(),
		FromWarehouseID: fromID,
		ToWarehouseID:   toID,
		ProductID:       productID,
		Stock:           stock,
	}, nil
}
