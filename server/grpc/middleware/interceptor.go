package middleware

import (
	"context"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/log/hooks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type (
	// private type used to define context keys
	ctxKey int
)

const (
	// AddressKey is the ip address of the connected peer
	AddressKey ctxKey = iota + 1
)

func init() {
	log.AddHook(hooks.NewContextHook(AddressKey, func(entry *log.Entry, ctxValue interface{}) {
		entry.Data["address"] = ctxValue
	}))
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer errors.Recover(func(recErr error) {
			errors.Log(recErr)
			err = status.Error(codes.Internal, "internal server error")
		})

		if peer, ok := peer.FromContext(ctx); ok {
			ctx = context.WithValue(ctx, AddressKey, peer.Addr.String())
		}

		log.WithContext(ctx).Debugf("Start unary")
		defer log.WithContext(ctx).Debugf("End unary")

		resp, err = handler(ctx, req)
		if err != nil {
			log.WithContext(ctx).WithError(err).Debug("Unary handler returned an error")
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

		ctx := ss.Context()
		if peer, ok := peer.FromContext(ctx); ok {
			ctx = context.WithValue(ctx, AddressKey, peer.Addr.String())

			ss = &WrappedServerStream{
				ServerStream: ss,
				ctx:          ctx,
			}
		}

		log.WithContext(ctx).Debugf("Start stream")
		defer log.WithContext(ctx).Debugf("End stream")

		err = handler(srv, ss)
		if err != nil {
			log.WithContext(ss.Context()).WithError(err).Debug("Stream handler returned an error")
		}
		return err
	})
}
