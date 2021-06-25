// Code generated by goa v3.3.1, DO NOT EDIT.
//
// artworks WebSocket client streaming
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"io"

	"github.com/gorilla/websocket"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "artworks" service.
type ConnConfigurer struct {
	RegisterTaskStateFn         goahttp.ConnConfigureFunc
	DownloadTaskStateEndpointFn goahttp.ConnConfigureFunc
}

// RegisterTaskStateClientStream implements the
// artworks.RegisterTaskStateClientStream interface.
type RegisterTaskStateClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// DownloadTaskStateClientStream implements the
// artworks.DownloadTaskStateEndpointClientStream interface.
type DownloadTaskStateClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "artworks" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		RegisterTaskStateFn:         fn,
		DownloadTaskStateEndpointFn: fn,
	}
}

// Recv reads instances of "artworks.TaskState" from the "registerTaskState"
// endpoint websocket connection.
func (s *RegisterTaskStateClientStream) Recv() (*artworks.TaskState, error) {
	var (
		rv   *artworks.TaskState
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

// Recv reads instances of "artworks.DownloadTaskState" from the
// "downloadTaskState" endpoint websocket connection.
func (s *DownloadTaskStateClientStream) Recv() (*artworks.DownloadTaskState, error) {
	var (
		rv   *artworks.DownloadTaskState
		body DownloadTaskStateResponseBody
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
	err = ValidateDownloadTaskStateResponseBody(&body)
	if err != nil {
		return rv, err
	}
	res := NewDownloadTaskStateOK(&body)
	return res, nil
}
