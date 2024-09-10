package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Product struct {
	ID        ulid.ULID
	Slug      string
	Name      string
	Price     uint64
	CreatedAt time.Time
}
