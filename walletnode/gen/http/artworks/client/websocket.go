// Code generated by goa v3.5.2, DO NOT EDIT.
//
// artworks WebSocket client streaming
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"io"

	"github.com/gorilla/websocket"
	artworks "github.com/pastelnetwork/gonode/walletnode/gen/artworks"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "artworks" service.
type ConnConfigurer struct {
	RegisterTaskStateFn goahttp.ConnConfigureFunc
	ArtSearchFn         goahttp.ConnConfigureFunc
}

// RegisterTaskStateClientStream implements the
// artworks.RegisterTaskStateClientStream interface.
type RegisterTaskStateClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// ArtSearchClientStream implements the artworks.ArtSearchClientStream
// interface.
type ArtSearchClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "artworks" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		RegisterTaskStateFn: fn,
		ArtSearchFn:         fn,
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

// Recv reads instances of "artworks.ArtworkSearchResult" from the "artSearch"
// endpoint websocket connection.
func (s *ArtSearchClientStream) Recv() (*artworks.ArtworkSearchResult, error) {
	var (
		rv   *artworks.ArtworkSearchResult
		body ArtSearchResponseBody
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
	err = ValidateArtSearchResponseBody(&body)
	if err != nil {
		return rv, err
	}
	res := NewArtSearchArtworkSearchResultOK(&body)
	return res, nil
}
