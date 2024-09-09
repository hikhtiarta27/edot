package model

import (
	"shared"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var ErrInvalidUlid = &shared.Error{
	HttpStatusCode: 400,
	Message:        "invalid id",
}

var ErrWarehouseNotFound = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.NotFound,
	Message:        "invalid id",
}

var ErrWarehouseInactive = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.NotFound,
	Message:        "warehouse inactive",
}

type WarehouseStatus string

var (
	WarehouseActive   WarehouseStatus = "ACTIVE"
	WarehouseInactive WarehouseStatus = "INACTIVE"
)

type Warehouse struct {
	ID        ulid.ULID
	Name      string
	Status    WarehouseStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}

func (Warehouse) TableName() string {
	return "warehouses"
}

func NewWarehouse(
	name string,
) (*Warehouse, error) {
	return &Warehouse{
		ID:     ulid.Make(),
		Name:   strings.TrimSpace(name),
		Status: WarehouseActive,
	}, nil
}

func (m *Warehouse) Activate() {
	m.Status = WarehouseActive
}

func (m *Warehouse) Deactivate() {
	m.Status = WarehouseInactive
}
