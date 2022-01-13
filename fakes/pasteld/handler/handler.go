package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitwurx/jrpc2"
	"github.com/pastelnetwork/gonode/fakes/pasteld/storage"
)

// Handler handles rpc requests to this fake pasteld server
type Handler interface {
	HandleMasternode(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
}

type rpcHandler struct {
	storage.Store
}

// New returns new instance of rpcHandler
func New(store storage.Store) Handler {
	return &rpcHandler{
		Store: store,
	}
}

func (h *rpcHandler) handle(method string, params []string) (interface{}, error) {
	key := method + "*" + strings.Join(params, "*")
	data, err := h.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %w", err)
	}

	toRet := make(map[string]interface{})
	if err := json.Unmarshal(data, &toRet); err != nil {
		return nil, fmt.Errorf("unable to decode data: %w", err)
	}

	return toRet, nil
}

// HandleMasternode handles requests on method 'masternode'
func (h *rpcHandler) HandleMasternode(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Params[0] != "top" {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command not provided",
		}
	}

	data, err := h.handle("masternode", p.Params)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}
