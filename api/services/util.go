package services

import (
	"fmt"
)

var InnerError = func(code int, val ...interface{}) *struct {
	Code    int
	Message string
} {

	var format string
	if len(val) > 0 {
		switch v := val[0].(type) {
		case string:
			format = v
		case error:
			format = v.Error()
		}
		val = val[1:]
	}

	return &struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: fmt.Sprintf(format, val...),
	}
}
