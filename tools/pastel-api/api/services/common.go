package services

import (
	"embed"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Common contains commonly used methods.
type Common struct{}

// RoutePath returns the given method and params, as one underscore-delimited string.
func (service *Common) RoutePath(method string, params []string) string {
	method = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(method, "")
	return strings.Join(append([]string{method}, params...), "_")
}

// UnmarshalFile reads a json string from the file, by the given `fs` and `filename` and unmarshal the data to the given `model` interface.
func (service *Common) UnmarshalFile(fs embed.FS, filename string, model interface{}) error {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return errors.New(err)
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return errors.New(err)
	}
	return nil
}

// NewCommon returns a new Common instance.
func NewCommon() *Common {
	return &Common{}
}
