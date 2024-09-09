package model

import "github.com/oklog/ulid/v2"

type SelectProduct struct {
	UseSearchEngine bool
	Keyword         string
}

type GetProduct struct {
	ID ulid.ULID
}

type UpdateStockProduct struct {
	ID    ulid.ULID
	Stock uint64
}
