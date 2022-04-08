package pastel

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pastelnetwork/gonode/common/b85"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Refer https://pastel.wiki/en/Architecture/Components/PastelOpenAPITicketStructures

const (
	// ActionTypeSense indicates a sense action request
	ActionTypeSense = "sense"
	// ActionTypeCascade indicates a cascade action request
	ActionTypeCascade = "cascade"
)

// ActionTicketDatas is a collection of ActionTicketData (Note type here)
type ActionTicketDatas []ActionTicketData

// ActionTicketData is Pastel Action ticket structure
type ActionTicketData struct {
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	ActionType string                 `json:"action_type"`
	Signatures ActionTicketSignatures `json:"signatures"`
	Key1       string                 `json:"key1"`
	Key2       string                 `json:"key2"`
	// is used to check if the SNs that created this ticket was indeed top SN
	// when that action call was made
	CalledAt int `json:"called_at"` // block at which action was requested,
	// fee in PSL
	StorageFee int `json:"storage_fee"`

	ActionTicket     []byte       `json:"action_ticket"`
	ActionTicketData ActionTicket `json:"-"`
}

// ActionTicket is define a action ticket's details
type ActionTicket struct {
	Version       int         `json:"action_ticket_version"`
	Caller        string      `json:"caller"` // PastelID of the caller
	BlockNum      int         `json:"blocknum"`
	BlockHash     string      `json:"block_hash"`
	ActionType    string      `json:"action_type"`
	APITicket     string      `json:"api_ticket"` // as ascii85(api_ticket)
	APITicketData interface{} `json:"-"`
}

// APISenseTicket returns APITicketData as *APISenseTicket if ActionType is ActionTypeSense
func (ticket *ActionTicket) APISenseTicket() (*APISenseTicket, error) {
	if ticket.APITicketData == nil {
		return nil, errors.New("nil APITicketData")
	}

	if ticket.ActionType != ActionTypeSense {
		return nil, fmt.Errorf("invalid action type: %s", ticket.ActionType)
	}

	switch data := ticket.APITicketData.(type) {
	case *APISenseTicket:
		return data, nil
	default:
		return nil, fmt.Errorf("invalid type of api data, got type: %s", reflect.TypeOf(ticket.APITicketData))
	}
}

// APICascadeTicket returns APITicketData as *APICascadeTicket if ActionType is ActionTypeCascade
func (ticket *ActionTicket) APICascadeTicket() (*APICascadeTicket, error) {
	if ticket.APITicketData == nil {
		return nil, errors.New("nil APITicketData")
	}

	if ticket.ActionType != ActionTypeCascade {
		return nil, fmt.Errorf("invalid action type: %s", ticket.ActionType)
	}

	switch data := ticket.APITicketData.(type) {
	case *APICascadeTicket:
		return data, nil
	default:
		return nil, fmt.Errorf("invalid type of api data, got type: %s", reflect.TypeOf(ticket.APITicketData))
	}
}

// APISenseTicket represents pastel API Sense ticket.
type APISenseTicket struct {
	DataHash []byte `json:"data_hash"`

	DDAndFingerprintsIc  uint32   `json:"dd_and_fingerprints_ic"`
	DDAndFingerprintsMax uint32   `json:"dd_and_fingerprints_max"`
	DDAndFingerprintsIDs []string `json:"dd_and_fingerprints_ids"`
}

// APICascadeTicket represents pastel API Cascade ticket.
type APICascadeTicket struct {
	DataHash []byte `json:"data_hash"`

	RQIc  uint32   `json:"rq_ic"`
	RQMax uint32   `json:"rq_max"`
	RQIDs []string `json:"rq_ids"`
	RQOti []byte   `json:"rq_oti"`
}

// EncodeActionTicket encodes  ActionTicket into byte array
func EncodeActionTicket(ticket *ActionTicket) ([]byte, error) {
	appTicketBytes, err := json.Marshal(ticket.APITicketData)
	if err != nil {
		return nil, errors.Errorf("marshal api ticket: %w", err)
	}

	appTicket := b85.Encode(appTicketBytes)
	ticket.APITicket = appTicket

	b, err := json.Marshal(ticket)
	if err != nil {
		return nil, errors.Errorf("marshal action ticket: %w", err)
	}

	return b, nil
}

// DecodeActionTicket decoded byte array into ArtTicket
func DecodeActionTicket(b []byte) (*ActionTicket, error) {
	res := &ActionTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("unmarshal nft ticket: %w", err)
	}

	appDecodedBytes, err := b85.Decode(res.APITicket)
	if err != nil {
		return nil, fmt.Errorf("b85 decode: %v", err)
	}

	switch res.ActionType {
	case ActionTypeSense:
		apiTicket := &APISenseTicket{}
		err = json.Unmarshal(appDecodedBytes, apiTicket)
		if err != nil {
			return nil, errors.Errorf("unmarshal api sense ticket: %w", err)
		}
		res.APITicketData = apiTicket
	case ActionTypeCascade:
		apiTicket := &APICascadeTicket{}
		err = json.Unmarshal(appDecodedBytes, apiTicket)
		if err != nil {
			return nil, errors.Errorf("unmarshal api cascade ticket: %w", err)
		}
		res.APITicketData = apiTicket
	default:
		return nil, fmt.Errorf("invalid action type, got type is %s", res.ActionType)
	}

	return res, nil
}

// ActionTicketSignatures represents signatures from parties
type ActionTicketSignatures struct {
	Principal map[string]string `json:"principal,omitempty"`
	Mn1       map[string]string `json:"mn1,omitempty"`
	Mn2       map[string]string `json:"mn2,omitempty"`
	Mn3       map[string]string `json:"mn3,omitempty"`
}

// EncodeActionSignatures encodes ActionTicketSignatures into byte array
func EncodeActionSignatures(signatures ActionTicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	signatures.Mn1 = nil

	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("marshal signatures: %w", err)
	}

	return b, nil
}

// RegisterActionRequest represents request to register an action ticket
// tickets register action "action_ticket" "{signatures}" "pastelid of MN 1" "passphrase" "key1" "key2" fee
type RegisterActionRequest struct {
	Ticket      *ActionTicket
	Signatures  *ActionTicketSignatures
	Mn1PastelID string
	Passphrase  string
	Key1        string
	Key2        string
	Fee         int64
}

// ActivateActionRequest - represents request to activate an action ticket
type ActivateActionRequest struct {
	RegTxID    string
	BlockNum   int
	Fee        int64
	PastelID   string
	Passphrase string
}
