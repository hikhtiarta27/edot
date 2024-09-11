package model

import "github.com/oklog/ulid/v2"

type SelectStock struct {
	ProductIDs []ulid.ULID
}

type CreateStock struct {
	ProductID   ulid.ULID
	WarehouseID ulid.ULID
	Stock       uint64
}
