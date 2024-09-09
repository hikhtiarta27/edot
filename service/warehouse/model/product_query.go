package model

import (
	"github.com/oklog/ulid/v2"
)

type GetProduct struct {
	ID ulid.ULID
}

type UpdateStockProduct struct {
	ID             ulid.ULID
	AvailableStock uint64
}
