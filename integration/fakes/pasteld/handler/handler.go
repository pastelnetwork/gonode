package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"

	"github.com/bitwurx/jrpc2"
	"github.com/pastelnetwork/gonode/integration/fakes/common/storage"
)

// Handler handles rpc requests to this fake pasteld server
type Handler interface {
	HandleMasternode(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleStorageFee(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleTickets(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleGetBlockCount(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleGetBlock(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleZGetBalance(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandlePastelid(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleSendMany(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleGetOperationStatus(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
	HandleGetRawTransaction(params json.RawMessage) (interface{}, *jrpc2.ErrorObject)
}

type rpcHandler struct {
	store storage.Store
}

// New returns new instance of rpcHandler
func New(store storage.Store) Handler {
	return &rpcHandler{
		store: store,
	}
}

func (h *rpcHandler) handle(method string, params []string) (interface{}, error) {
	key := method + "*" + strings.Join(params, "*")
	data, err := h.store.Get(key)
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
		log.Println("failed to parse params")
		return nil, err
	}

	if !(p.Params[0] == "top" || p.Params[0] == "list") {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command provided",
		}
	}

	data, err := h.handle("masternode", p.Params)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

// HandleMasternode handles requests on method 'masternode'
func (h *rpcHandler) HandleGetRawTransaction(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		log.Println("failed to parse params")
		return nil, err
	}

	data, err := h.handle("getrawtransaction", []string{})
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

// HandleStorageFee handles requests on method 'StorageFee'
func (h *rpcHandler) HandleStorageFee(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Params[0] != "getnetworkfee" {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command not provided",
		}
	}

	data, err := h.handle("storagefee", p.Params)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

func (h *rpcHandler) HandleTickets(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Params[0] == "find" && p.Params[1] == "id" && !isValidMNPastelID(p.Params[2]) {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command provided",
		}
	} else if (p.Params[0] == "tools" && p.Params[1] == "gettotalstoragefee") ||
		(p.Params[0] == "register") {
		p.Params = p.Params[:2]
	}

	data, err := h.handle("tickets", p.Params)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

func (h *rpcHandler) HandleGetBlock(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	data, err := h.handle("getblock", p.Params[:1])
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

func (h *rpcHandler) HandleZGetBalance(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Params[0] != testconst.RegSpendableAddress {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.ErrorMsg("unexpected address"),
			Data:    "unexpected address",
		}
	}

	key := "z_getbalance" + "*" + p.Params[0]
	data, err := h.store.Get(key)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	numStr := string(data)

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return num, nil
}

func (h *rpcHandler) HandleGetBlockCount(_ json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	key := "getblockcount" + "*"
	data, err := h.store.Get(key)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	numStr := string(data)

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return num, nil
}

func (h *rpcHandler) HandlePastelid(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if !(p.Params[0] == "sign" || p.Params[0] == "verify") {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command provided",
		}
	}

	data, err := h.handle("pastelid", p.Params[:1])
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(fmt.Sprintf("unable to fetch: %s", err.Error())),
			Data:    fmt.Sprintf("unable to fetch: %s", err.Error()),
		}
	}

	return data, nil
}

func (h *rpcHandler) HandleGetOperationStatus(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	key := "z_getoperationstatus" + "*" + p.Params[0]
	data, err := h.store.Get(key)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.ErrorMsg("not found"),
			Data:    "not found",
		}
	}

	toRet := []GetOperationStatusResult{}
	if err := json.Unmarshal(data, &toRet); err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
			Data:    "unmarshal failure",
		}
	}

	return toRet, nil
}

func (h *rpcHandler) HandleSendMany(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RpcParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.Params[0] != testconst.RegSpendableAddress {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.ErrorMsg("unexpected address"),
			Data:    "unexpected address",
		}
	}

	key := "z_sendmanywithchangetosender" + "*" + p.Params[0]
	data, err := h.store.Get(key)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.ErrorMsg("not found"),
			Data:    "not found",
		}
	}

	return string(data), nil
}
