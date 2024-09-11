package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type OrderDetail struct {
	ID          ulid.ULID
	OrderID     ulid.ULID
	ProductID   ulid.ULID
	ProductName string
	Qty         uint64
	Price       uint64
	TotalPrice  uint64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt
}

func (OrderDetail) TableName() string {
	return "order_details"
}

func (m OrderDetail) GetTotalPrice() uint64 {
	return m.Qty * m.Price
}
