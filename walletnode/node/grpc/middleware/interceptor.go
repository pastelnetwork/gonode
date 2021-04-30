package middleware

import (
	"context"
	"strings"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, RequestIDKey, reqID)

		log.WithContext(ctx).WithField("method", strings.TrimPrefix(method, "/walletnode.")).Debugf("Start unary")

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		reqID, _ := random.String(8, random.Base62Chars)
		ctx = context.WithValue(ctx, RequestIDKey, reqID)

		log.WithContext(ctx).WithField("method", strings.TrimPrefix(method, "/walletnode.")).Debugf("Start stream")

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return clientStream, err
	}
}
