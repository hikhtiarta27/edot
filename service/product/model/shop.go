package model

import (
	"shared"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
)

var ErrShopWarehouseNotFound = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.NotFound,
	Message:        "shop warehouse not found",
}

type Shop struct {
	ID        ulid.ULID
	Name      string
	Warehouse []ulid.ULID
	CreatedAt time.Time
}
