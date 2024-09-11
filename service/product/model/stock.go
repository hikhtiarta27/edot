package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Stock struct {
	ID             ulid.ULID
	ProductID      ulid.ULID
	AvailableStock uint64
	ReservedStock  uint64
	CreatedAt      time.Time
}

type Stocks []Stock

func (m Stocks) MapByProductID() map[ulid.ULID]*Stock {

	res := make(map[ulid.ULID]*Stock)

	for _, stock := range m {
		res[stock.ProductID] = &stock
	}

	return res
}
