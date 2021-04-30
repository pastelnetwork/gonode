package middleware

import (
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
)

type (
	// private type used to define context keys
	ctxKey int
)

const (
	// AddressKey is the ip address of the connected peer
	AddressKey ctxKey = iota + 1

	// MethodKey is proto rpc
	MethodKey

	// RequestIDKey is unique numeric for every request
	RequestIDKey
)

func init() {
	log.AddHook(hooks.NewContextHook(RequestIDKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		return fmt.Sprintf("%v: %s", ctxValue, msg), fields
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
