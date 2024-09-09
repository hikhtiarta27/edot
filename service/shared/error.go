package shared

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	HttpStatusCode int
	GrpcStatusCode codes.Code
	Message        string
	Data           interface{}
}

func (e Error) Error() string {
	return e.Message
}

func GrpcUnaryParser() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if err != nil {
			if customErr, ok := (err).(*Error); ok {
				err = status.Error(customErr.GrpcStatusCode, customErr.Message)
			} else if grpcErr, _ := status.FromError(err); grpcErr != nil {
				err = status.Error(grpcErr.Code(), grpcErr.Message())
			}

			return
		}

		return
	}
}

func GrpcStreamParser() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		err = handler(srv, ss)

		if err != nil {
			if customErr, ok := (err).(*Error); ok {
				err = status.Error(customErr.GrpcStatusCode, customErr.Message)
			} else if grpcErr, _ := status.FromError(err); grpcErr != nil {
				err = status.Error(grpcErr.Code(), grpcErr.Message())
			}

			return
		}

		return
	}
}
