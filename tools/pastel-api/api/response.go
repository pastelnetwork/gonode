package api

import (
	"bytes"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Error represents an API error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response represents an API response.
type Response struct {
	*Error `json:"error"`
	ID     interface{} `json:"id"`
	Result interface{} `json:"result"`
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

func newErrorResponse(id interface{}, code int, message string) *Response {
	return &Response{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}
