package grpc

import (
	"context"

	"github.com/pastelnetwork/go-commons/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recovery() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer errors.Recover(func(recErr error) {
			errors.Log(recErr)
			err = status.Error(codes.Internal, "internal server error")
		})
		resp, err = handler(ctx, req)
		return resp, err
	})
}
