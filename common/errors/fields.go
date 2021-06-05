package errors

import (
	"fmt"
	"strings"
)

// Fields contains key/value data
type Fields map[string]interface{}

// Error returns the underlying error's message.
func (fields Fields) String() string {
	var strs []string
	for name, value := range fields {
		strs = append(strs, fmt.Sprintf("%s=%s", name, value))
	}
	return strings.Join(strs, " ")
}

// ExtractFields extracts fields from err
func ExtractFields(err error) Fields {
	if e, ok := err.(*Error); ok {
		return e.fields
	}
	return nil
}
