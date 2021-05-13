package middleware

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/random"
)

type (
	// private type used to define context keys
	ctxKey int
)

const (
	// AddressKey is the ip address of the connected peer
	AddressKey ctxKey = iota

	// MethodKey is proto rpc
	MethodKey

	// RequestIDKey is unique numeric for every request
	RequestIDKey
)

func init() {
	log.AddHook(hooks.NewContextHook(RequestIDKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["reqID"] = ctxValue
		return msg, fields
	}))

	log.AddHook(hooks.NewContextHook(AddressKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["address"] = ctxValue
		return msg, fields
	}))

	log.AddHook(hooks.NewContextHook(MethodKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["method"] = ctxValue
		return msg, fields
	}))
}

// WithRequestID returns a context with RequestID value.
func WithRequestID(ctx context.Context) context.Context {
	reqID, _ := random.String(8, random.Base62Chars)
	ctx = context.WithValue(ctx, RequestIDKey, reqID)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("server-%s", reqID))

	return ctx
}
