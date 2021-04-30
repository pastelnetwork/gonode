package middleware

import (
	"context"

	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/walletnode/node/grpc/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

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
