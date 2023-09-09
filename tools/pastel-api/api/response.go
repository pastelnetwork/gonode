package api

import (
	"bytes"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Response represents an API response.
type Response struct {
	ID     interface{} `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

// Bytes returns the response in json string format.
func (resp *Response) Bytes() ([]byte, error) {
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, errors.New(err)
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, data); err != nil {
		return nil, errors.New(err)
	}
	data = append(buffer.Bytes(), []byte("\n")...)

	return data, nil
}

func newResponse(id interface{}, result interface{}) *Response {
	return &Response{
		ID:     id,
		Result: result,
	}
}

func newErrorResponse(id interface{}, err error) *Response {
	return &Response{
		ID:    id,
		Error: err,
	}
}
