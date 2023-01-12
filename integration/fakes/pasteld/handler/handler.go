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

	if !(p.Params[0] == "getnetworkfee" || p.Params[0] == "getactionfees") {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command not provided",
		}
	}

	if p.Params[0] == "getactionfees" {
		p.Params = p.Params[:1]
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
	fmt.Println("got params: ", p.Params)
	if p.Params[0] == "find" && p.Params[1] == "id" && !isValidMNPastelID(p.Params[2]) {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "invalid command provided",
		}
	} else if (p.Params[0] == "tools" && p.Params[1] == "gettotalstoragefee") ||
		(p.Params[0] == "register") || (p.Params[0] == "activate") ||
		(p.Params[0] == "find" && p.Params[1] == "action-act") {
		p.Params = p.Params[:2]
	} else if len(p.Params) > 4 && p.Params[0] == "tools" && p.Params[1] == "validateownership" {
		p.Params = p.Params[:5]
	} else if len(p.Params) > 2 && p.Params[0] == "findbylabel" && p.Params[1] == "nft" {
		key := "tickets" + "*" + strings.Join(p.Params[:2], "*")
		data, err := h.store.Get(key)
		if err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
				Data:    "fetch failure",
			}
		}
		toRet := RegTickets{}
		if err := json.Unmarshal(data, &toRet); err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
				Data:    "unmarshal failure",
			}
		}

		return toRet, nil
	} else if len(p.Params) > 2 && p.Params[0] == "findbylabel" && p.Params[1] == "action" {
		key := "tickets" + "*" + strings.Join(p.Params[:2], "*")
		data, err := h.store.Get(key)
		if err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
				Data:    "fetch failure",
			}
		}
		toRet := ActTickets{}
		if err := json.Unmarshal(data, &toRet); err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
				Data:    "unmarshal failure",
			}
		}

		return toRet, nil
	} else if len(p.Params) >= 2 && p.Params[0] == "list" &&
		((p.Params[1] == "trade" && p.Params[2] == "available") || p.Params[1] == "act" || p.Params[1] == "nft") {
		if p.Params[1] == "act" {
			p.Params = p.Params[:2]

			key := "tickets" + "*" + strings.Join(p.Params, "*")
			data, err := h.store.Get(key)
			if err != nil {
				return nil, &jrpc2.ErrorObject{
					Code:    jrpc2.InternalErrorCode,
					Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
					Data:    "fetch failure",
				}
			}
			toRet := ActTickets{}
			if err := json.Unmarshal(data, &toRet); err != nil {
				return nil, &jrpc2.ErrorObject{
					Code:    jrpc2.InternalErrorCode,
					Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
					Data:    "unmarshal failure",
				}
			}

			return toRet, nil
		} else if p.Params[1] == "nft" {
			p.Params = p.Params[:2]

			key := "tickets" + "*" + strings.Join(p.Params, "*")
			data, err := h.store.Get(key)
			if err != nil {
				return nil, &jrpc2.ErrorObject{
					Code:    jrpc2.InternalErrorCode,
					Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
					Data:    "fetch failure",
				}
			}
			toRet := RegTickets{}
			if err := json.Unmarshal(data, &toRet); err != nil {
				return nil, &jrpc2.ErrorObject{
					Code:    jrpc2.InternalErrorCode,
					Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
					Data:    "unmarshal failure",
				}
			}

			return toRet, nil
		}

		key := "tickets" + "*" + strings.Join(p.Params, "*")
		data, err := h.store.Get(key)
		if err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
				Data:    "fetch failure",
			}
		}
		toRet := []TradeTicket{}
		if err := json.Unmarshal(data, &toRet); err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
				Data:    "unmarshal failure",
			}
		}

		return toRet, nil
	} else if p.Params[0] == "get" {
		fmt.Println("params: ", p.Params[len(p.Params)-1], p.Params[len(p.Params)-2])
		key := "tickets" + "*" + strings.Join(p.Params, "*")
		data, err := h.store.Get(key)
		if err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to fetch " + err.Error()),
				Data:    "fetch failure",
			}
		}
		if p.Params[len(p.Params)-1] == testconst.TestCascadeRegTXID ||
			p.Params[len(p.Params)-1] == testconst.TestSenseRegTXID {
			toRet := ActionRegTicket{}
			if err := json.Unmarshal(data, &toRet); err != nil {
				return nil, &jrpc2.ErrorObject{
					Code:    jrpc2.InternalErrorCode,
					Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
					Data:    "unmarshal failure",
				}
			}

			return toRet, nil
		}
		toRet := RegTicket{}
		if err := json.Unmarshal(data, &toRet); err != nil {
			return nil, &jrpc2.ErrorObject{
				Code:    jrpc2.InternalErrorCode,
				Message: jrpc2.ErrorMsg("unable to unmarsal " + err.Error()),
				Data:    "unmarshal failure",
			}
		}

		return toRet, nil
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

type TradeTicket struct {
	Height int             `json:"height"`
	TXID   string          `json:"txid"`
	Ticket TradeTicketData `json:"ticket"`
}

// TradeTicketData represents pastel Trade ticket data
type TradeTicketData struct {
	Type             string `json:"type"`              // "nft-trade"
	Version          int    `json:"version"`           // version
	PastelID         string `json:"pastelID"`          // PastelID of the buyer
	SellTXID         string `json:"sell_txid"`         // txid with sale ticket
	BuyTXID          string `json:"buy_txid"`          // txid with buy ticket
	NFTTXID          string `json:"nft_txid"`          // txid with either 1) NFT activation ticket or 2) trade ticket in it
	RegistrationTXID string `json:"registration_txid"` // txid with registration ticket
	Price            string `json:"price"`
	CopySerialNR     string `json:"copy_serial_nr"`
	Signature        string `json:"signature"`
}

// RegTickets is a collection of RegTicket
type RegTickets []RegTicket

// RegTicket represents pastel registration ticket.
type RegTicket struct {
	Height        int           `json:"height"`
	TXID          string        `json:"txid"`
	RegTicketData RegTicketData `json:"ticket"`
}

// RegTicketSignatures represents signatures from parties
type RegTicketSignatures struct {
	Creator map[string]string `json:"creator,omitempty"`
	Mn1     map[string]string `json:"mn1,omitempty"`
	Mn2     map[string]string `json:"mn2,omitempty"`
	Mn3     map[string]string `json:"mn3,omitempty"`
}

// ActionTicketSignatures represents signatures from parties
type ActionTicketSignatures struct {
	Principal map[string]string `json:"principal,omitempty"`
	Mn1       map[string]string `json:"mn1,omitempty"`
	Mn2       map[string]string `json:"mn2,omitempty"`
	Mn3       map[string]string `json:"mn3,omitempty"`
}

// RegTicketData is Pastel Registration ticket structure
type RegTicketData struct {
	Type           string              `json:"type"`
	Version        int                 `json:"version"`
	Signatures     RegTicketSignatures `json:"signatures"`
	Label          string              `json:"label"`
	CreatorHeight  int                 `json:"creator_height"`
	TotalCopies    int                 `json:"total_copies"`
	Royalty        float64             `json:"royalty"`
	RoyaltyAddress string              `json:"royalty_address"`
	Green          bool                `json:"green"`
	GreenAddress   string              `json:"green_address"`
	StorageFee     int                 `json:"storage_fee"`
	NFTTicket      []byte              `json:"nft_ticket"`
	//NFTTicketData  NFTTicket           `json:"-"`
}

// ActTickets is a collection of ActTicket
type ActTickets []ActTicket

// ActTicket represents pastel activation ticket.
type ActTicket struct {
	Height        int           `json:"height"`
	TXID          string        `json:"txid"`
	ActTicketData ActTicketData `json:"ticket"`
}

// ActTicketData represents activation ticket properties
type ActTicketData struct {
	PastelID      string `json:"pastelID"`
	Signature     string `json:"signature"`
	Type          string `json:"type"`
	CreatorHeight int    `json:"creator_height"`
	RegTXID       string `json:"reg_txid"`
	StorageFee    int    `json:"storage_fee"`
	Version       int    `json:"version"`
}

// ActionRegTicket represents pastel registration ticket.
type ActionRegTicket struct {
	Height           int              `json:"height"`
	TXID             string           `json:"txid"`
	ActionTicketData ActionTicketData `json:"ticket"`
}

// ActionTicketData is Pastel Action ticket structure
type ActionTicketData struct {
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	ActionType string                 `json:"action_type"`
	Signatures ActionTicketSignatures `json:"signatures"`
	Label      string                 `json:"label"`
	// is used to check if the SNs that created this ticket was indeed top SN
	// when that action call was made
	CalledAt int `json:"called_at"` // block at which action was requested,
	// fee in PSL
	StorageFee int `json:"storage_fee"`

	ActionTicket string `json:"action_ticket"`
	//ActionTicketData ActionTicket `json:"-"`
}
