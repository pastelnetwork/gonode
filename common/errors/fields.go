package errors

import (
	"fmt"
)

// Fields contains key/value data
type Fields map[string]interface{}

// Error returns the underlying error's message.
func (fields Fields) String() string {
	var str string
	for name, value := range fields {
		str += fmt.Sprintf("%s=%s ", name, value)
	}
	return str
}

// ExtractFields extracts fields from err
func ExtractFields(err error) Fields {
	if e, ok := err.(*Error); ok {
		return e.fields
	}
	return nil
}
