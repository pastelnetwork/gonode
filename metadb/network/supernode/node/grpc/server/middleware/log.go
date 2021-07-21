package middleware

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/metadb/network/proto"
	"google.golang.org/grpc/metadata"
)

type (
	// private type used to define context keys.
	ctxKey int
)

const (
	// RequestIDKey is unique numeric for every request.
	RequestIDKey ctxKey = iota

	// SessIDKey is unique numeric for every regiration process.
	SessIDKey
)

func init() {
	log.AddHook(hooks.NewContextHook(RequestIDKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["reqID"] = ctxValue
		return msg, fields
	}))
	log.AddHook(hooks.NewContextHook(SessIDKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["SessID"] = ctxValue
		return msg, fields
	}))
}

// WithRequestID returns a context with RequestID value.
func WithRequestID(ctx context.Context) context.Context {
	reqID, _ := random.String(8, random.Base62Chars)
	return log.ContextWithPrefix(ctx, fmt.Sprintf("server-%s", reqID))
}

// WithSessID returns a context with SessID value.
func WithSessID(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	mdVals := md.Get(proto.MetadataKeySessID)
	if len(mdVals) == 0 {
		return ctx
	}
	return context.WithValue(ctx, SessIDKey, mdVals[0])
}
