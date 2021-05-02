package log

import (
	"github.com/pastelnetwork/gonode/common/log/hooks"
)

type (
	// private type used to define context keys
	ctxKey int
)

const (
	// PrefixKey is the prefix of the log record
	PrefixKey ctxKey = iota + 1
)

func init() {
	AddHook(hooks.NewContextHook(PrefixKey, func(ctxValue interface{}, msg string, fields hooks.ContextHookFields) (string, hooks.ContextHookFields) {
		fields["prefix"] = ctxValue
		return msg, fields
	}))
}
