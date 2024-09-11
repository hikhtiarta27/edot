package model

import "github.com/oklog/ulid/v2"

type GetOrderDetail struct {
	ID ulid.ULID
}

type SelectOrderDetail struct {
	OrderID ulid.ULID
}
