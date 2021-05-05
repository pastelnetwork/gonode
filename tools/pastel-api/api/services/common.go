package services

import (
	"embed"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

type Common struct{}

func (service *Common) RoutePath(method string, params []string) string {
	method = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(method, "")
	return strings.Join(append([]string{method}, params...), "_")
}

func (service *Common) UnmarshalFile(fs embed.FS, filename string, model interface{}) error {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return errors.New(err)
	}

	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}
	return nil
}

func NewCommon() *Common {
	return &Common{}
}
