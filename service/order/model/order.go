package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type OrderStatus string

var (
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusPaid    OrderStatus = "PAID"
	OrderStatusExpired OrderStatus = "EXPIRED"
)

type Order struct {
	ID         ulid.ULID
	TotalItem  uint64
	TotalPrice uint64
	Status     OrderStatus
	ExpiredAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  gorm.DeletedAt

	Detail []OrderDetail `gorm:"-"`
}

func (Order) TableName() string {
	return "orders"
}

func NewOrder(
	totalItem uint64,
) (*Order, error) {
	instance := &Order{
		ID:        ulid.Make(),
		TotalItem: totalItem,
		Status:    OrderStatusPending,
	}

	return instance, nil
}

func (m *Order) AddDetail(prd *Product, qty uint64) {

	orderDetail := OrderDetail{
		ID:          ulid.Make(),
		OrderID:     m.ID,
		ProductID:   prd.ID,
		ProductName: prd.Name,
		Qty:         qty,
		Price:       prd.Price,
	}

	orderDetail.TotalPrice = orderDetail.GetTotalPrice()

	m.TotalPrice += orderDetail.GetTotalPrice()

	m.Detail = append(m.Detail, orderDetail)
}
