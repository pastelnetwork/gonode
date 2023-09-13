// Code generated by goa v3.13.0, DO NOT EDIT.
//
// nft WebSocket client streaming
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"io"

	"github.com/gorilla/websocket"
	nft "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "nft" service.
type ConnConfigurer struct {
	RegisterTaskStateFn goahttp.ConnConfigureFunc
	NftSearchFn         goahttp.ConnConfigureFunc
}

// RegisterTaskStateClientStream implements the
// nft.RegisterTaskStateClientStream interface.
type RegisterTaskStateClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NftSearchClientStream implements the nft.NftSearchClientStream interface.
type NftSearchClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "nft" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		RegisterTaskStateFn: fn,
		NftSearchFn:         fn,
	}
}

// Recv reads instances of "nft.TaskState" from the "registerTaskState"
// endpoint websocket connection.
func (s *RegisterTaskStateClientStream) Recv() (*nft.TaskState, error) {
	var (
		rv   *nft.TaskState
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

// Recv reads instances of "nft.NftSearchResult" from the "nftSearch" endpoint
// websocket connection.
func (s *NftSearchClientStream) Recv() (*nft.NftSearchResult, error) {
	var (
		rv   *nft.NftSearchResult
		body NftSearchResponseBody
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
	err = ValidateNftSearchResponseBody(&body)
	if err != nil {
		return rv, err
	}
	res := NewNftSearchResultOK(&body)
	return res, nil
}
