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
				log.WithContext(ctx).WithError(err).Error("Handler error")
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
				log.WithContext(ctx).WithError(err).Error("Handler error")
			}
		}
		return err
	})
}

// UnaryClientInterceptor returns a new unary client interceptor.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = WithRequestID(ctx)

		log.WithContext(ctx).WithField("method", method).Debugf("Start unary client")
		defer log.WithContext(ctx).Debugf("End unary client")

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil && !errors.IsContextCanceled(err) {
			log.WithContext(ctx).WithError(err).Error("Handler error")
		}
		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = WithRequestID(ctx)

		log.WithContext(ctx).WithField("method", method).Debugf("Start stream")
		defer log.WithContext(ctx).Debugf("End stream")

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil && !errors.IsContextCanceled(err) {
			log.WithContext(ctx).WithError(err).Error("Handler error")
		}
		return clientStream, err
	}
}
