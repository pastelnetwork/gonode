package data

import (
	"embed"

	"github.com/pastelnetwork/gonode/common/errors"
)

//go:embed *.json
var files embed.FS

func ReadFile(filename string) ([]byte, error) {
	data, err := files.ReadFile(filename)
	if err != nil {
		return nil, errors.New(err)
	}
	return data, nil
}
