package services

import (
	"fmt"
)

// InnerError generates an error by http status code.
// `vals` is values for Sprintf, where the first index can be `format`, and the rest are values.
func InnerError(code int, vals ...interface{}) *struct {
	Code    int
	Message string
} {

	var format string
	if len(vals) > 0 {
		switch v := vals[0].(type) {
		case string:
			format = v
		case error:
			format = v.Error()
		}
		vals = vals[1:]
	}

	return &struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: fmt.Sprintf(format, vals...),
	}
}
