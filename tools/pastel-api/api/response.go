package api

import (
	"bytes"
	"encoding/json"
	"net/http"

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

func (resp *Response) Write(w http.ResponseWriter, code int) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return errors.New(err)
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, data); err != nil {
		return errors.New(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	res := append(buffer.Bytes(), []byte("\n")...)
	if _, err := w.Write(res); err != nil {
		return errors.New(err)
	}
	return nil
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
