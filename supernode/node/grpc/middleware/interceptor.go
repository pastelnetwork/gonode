package middleware

import (
	"context"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/log"
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

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, log.RequestIDKey, reqID)

		method := strings.TrimPrefix(info.FullMethod, "/proto.")

		var address string
		if peer, ok := peer.FromContext(ctx); ok {
			address = peer.Addr.String()
		}

		log.WithContext(ctx).WithField("address", address).WithField("method", method).Debugf("Start unary")
		resp, err = handler(ctx, req)
		log.WithContext(ctx).WithError(err).Debugf("End unary")

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

		ctx := ss.Context()

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, log.RequestIDKey, reqID)

		ss = &WrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		method := strings.TrimPrefix(info.FullMethod, "/proto.")

		var address string
		if peer, ok := peer.FromContext(ctx); ok {
			address = peer.Addr.String()
		}

		log.WithContext(ctx).WithField("address", address).WithField("method", method).Debugf("Start stream")
		err = handler(srv, ss)
		log.WithContext(ctx).WithError(err).Debugf("End stream")

		return err
	})
}

// UnaryClientInterceptor returns a new unary client interceptor.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, log.RequestIDKey, reqID)

		var address string
		if peer, ok := peer.FromContext(ctx); ok {
			address = peer.Addr.String()
		}

		log.WithContext(ctx).WithField("address", address).WithField("method", method).Debugf("Start unary client")
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.WithContext(ctx).WithError(err).Debugf("End unary client")

		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, log.RequestIDKey, reqID)

		var address string
		if peer, ok := peer.FromContext(ctx); ok {
			address = peer.Addr.String()
		}

		log.WithContext(ctx).WithField("address", address).WithField("method", method).Debugf("Start stream")
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		log.WithContext(ctx).WithError(err).Debugf("End stream")

		return clientStream, err
	}
}
