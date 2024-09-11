package model

import (
	"shared"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var ErrDuplicateStock = &shared.Error{
	HttpStatusCode: 422,
	GrpcStatusCode: codes.AlreadyExists,
	Message:        "duplicate stock of product",
}

var ErrInsufficientStock = &shared.Error{
	HttpStatusCode: 422,
	GrpcStatusCode: codes.Aborted,
	Message:        "insufficient stock",
}

type StockAction string

const (
	StockRelease = "release"
	StockReserve = "reserve"
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
