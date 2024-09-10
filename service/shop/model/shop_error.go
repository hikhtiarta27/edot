package model

import (
	"shared"

	"google.golang.org/grpc/codes"
)

var ErrInvalidShopID = &shared.Error{
	HttpStatusCode: 400,
	GrpcStatusCode: codes.InvalidArgument,
	Message:        "invalid shop id ",
}
