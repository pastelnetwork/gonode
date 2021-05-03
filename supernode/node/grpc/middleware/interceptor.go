package middleware

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer errors.Recover(func(recErr error) {
			errors.Log(recErr)
			err = status.Error(codes.Internal, "internal server error")
		})

		ctx = WithRequestID(ctx)
		peer, _ := peer.FromContext(ctx)

		log.WithContext(ctx).
			WithField("address", peer.Addr.String()).
			WithField("method", info.FullMethod).
			Debugf("Start unary")

		defer log.WithContext(ctx).Debugf("End unary")

		resp, err = handler(ctx, req)
		if err != nil && !errors.IsContextCanceled(err) {
			if errors.IsContextCanceled(err) {
				log.WithContext(ctx).Debug("Closed by peer")
			} else {
				log.WithContext(ctx).WithError(err).Warn("Unary error")
			}
		}
		return resp, err
	})
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer errors.Recover(func(recErr error) {
			errors.Log(recErr)
			err = status.Error(codes.Internal, "internal server error")
		})

		ctx := WithRequestID(ss.Context())
		ss = &WrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}
		peer, _ := peer.FromContext(ctx)

		log.WithContext(ctx).
			WithField("address", peer.Addr.String()).
			WithField("method", info.FullMethod).
			Debugf("Start stream")

		defer log.WithContext(ctx).Debugf("End stream")

		err = handler(srv, ss)
		if err != nil {
			if errors.IsContextCanceled(err) {
				log.WithContext(ctx).Debug("Closed by peer")
			} else {
				log.WithContext(ctx).WithError(err).Warn("Stream error")
			}
		}
		return err
	})
}
