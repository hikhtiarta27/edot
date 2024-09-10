package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Warehouse struct {
	ID        ulid.ULID
	Name      string
	Status    string
	CreatedAt time.Time
}
