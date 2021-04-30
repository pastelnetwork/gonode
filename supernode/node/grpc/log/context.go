package log

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
	log.AddHook(hooks.NewContextHook(AddressKey, func(entry *log.Entry, ctxValue interface{}) {
		entry.Data["address"] = ctxValue
	}))

	log.AddHook(hooks.NewContextHook(MethodKey, func(entry *log.Entry, ctxValue interface{}) {
		entry.Data["method"] = ctxValue
	}))

	log.AddHook(hooks.NewContextHook(RequestIDKey, func(entry *log.Entry, ctxValue interface{}) {
		entry.Message = fmt.Sprintf("[%v] %s", ctxValue, entry.Message)
	}))

}
