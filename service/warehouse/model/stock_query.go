package model

import "github.com/oklog/ulid/v2"

type GetStock struct {
	ID        ulid.ULID
	ProductID ulid.ULID
}

type SelectStock struct {
	ProductIDs []ulid.ULID
}

type CreateStock struct {
	*Stock
	*WarehouseTransfer
}

type ReserveReleaseStock struct {
	*Stock
	Action StockAction
	Qty    uint64
}
