package middleware

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer errors.Recover(func(recErr error) {
			err = status.Error(codes.Internal, "internal server error")
		})

		ctx = WithRequestID(ctx)
		ctx = WithSessID(ctx)

		return handler(ctx, req)
	})
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer errors.Recover(func(recErr error) {
			err = status.Error(codes.Internal, "internal server error")
		})

		ctx := ss.Context()
		ctx = WithRequestID(ctx)
		ctx = WithSessID(ctx)

		ss = &WrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}
		return handler(srv, ss)
	})
}
