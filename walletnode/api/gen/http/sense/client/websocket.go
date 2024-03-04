// Code generated by goa v3.15.0, DO NOT EDIT.
//
// sense WebSocket client streaming
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"io"

	"github.com/gorilla/websocket"
	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "sense" service.
type ConnConfigurer struct {
	RegisterTaskStateFn goahttp.ConnConfigureFunc
}

// RegisterTaskStateClientStream implements the
// sense.RegisterTaskStateClientStream interface.
type RegisterTaskStateClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "sense" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		RegisterTaskStateFn: fn,
	}
}

// Recv reads instances of "sense.TaskState" from the "registerTaskState"
// endpoint websocket connection.
func (s *RegisterTaskStateClientStream) Recv() (*sense.TaskState, error) {
	var (
		rv   *sense.TaskState
		body RegisterTaskStateResponseBody
		err  error
	)
	err = s.conn.ReadJSON(&body)
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		s.conn.Close()
		return rv, io.EOF
	}
	if err != nil {
		return rv, err
	}
	err = ValidateRegisterTaskStateResponseBody(&body)
	if err != nil {
		return rv, err
	}
	res := NewRegisterTaskStateTaskStateOK(&body)
	return res, nil
}
