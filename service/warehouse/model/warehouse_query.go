package model

import "github.com/oklog/ulid/v2"

type SelectWarehouse struct {
}

type GetWarehouse struct {
	ID ulid.ULID
}
