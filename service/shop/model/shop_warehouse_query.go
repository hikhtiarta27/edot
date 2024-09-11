package model

import "github.com/oklog/ulid/v2"

type GetShopWarehouse struct {
	ID          ulid.ULID
	ShopID      ulid.ULID
	WarehouseID ulid.ULID
}

type SelectShopWarehouse struct {
	ShopID ulid.ULID
}
