package model

import "github.com/oklog/ulid/v2"

type ReserveReleaseStock struct {
	ProductID ulid.ULID
	Qty       uint64
	Action    StockAction
}

type StockAction string

const (
	StockRelease StockAction = "release"
	StockReserve StockAction = "reserve"
)
