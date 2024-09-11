package model

import (
	"shared"

	"google.golang.org/grpc/codes"
)

var ErrInvalidUlid = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.InvalidArgument,
	Message:        "invalid id",
}

var ErrShopNotFound = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.NotFound,
	Message:        "shop not found",
}

var ErrDuplicateShopWarehouse = &shared.Error{
	HttpStatusCode: 422,
	GrpcStatusCode: codes.AlreadyExists,
	Message:        "shop warehouse already exist",
}
