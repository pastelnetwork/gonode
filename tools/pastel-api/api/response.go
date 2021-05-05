package api

import (
	"bytes"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	ID     string      `json:"id"`
	Error  *Error      `json:"error"`
	Result interface{} `json:"result"`
}

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

func newResponse(id string, result interface{}) *Response {
	return &Response{
		ID:     id,
		Result: result,
	}
}

func newErrorResponse(id string, code int, message string) *Response {
	return &Response{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}
