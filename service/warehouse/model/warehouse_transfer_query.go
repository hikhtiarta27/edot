package model

import "github.com/oklog/ulid/v2"

type SelectWarehouseTransfer struct {
}

type SumStockWarehouseTransferRequest struct {
	ProductIDs []ulid.ULID
}

type SumStockWarehouseTransferResult map[ulid.ULID]uint64
