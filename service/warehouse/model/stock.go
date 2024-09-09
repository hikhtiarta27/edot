package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Stock struct {
	ID             ulid.ULID
	ProductID      ulid.ULID
	AvailableStock uint64
	ReservedStock  uint64
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	DeletedAt      gorm.DeletedAt
}

func (Stock) TableName() string {
	return "stocks"
}

func NewStock(
	productID ulid.ULID,
	stock uint64,
) (*Stock, error) {
	return &Stock{
		ID:             ulid.Make(),
		ProductID:      productID,
		AvailableStock: stock,
	}, nil
}
