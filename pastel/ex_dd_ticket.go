package pastel

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

// RegExDDTickets is a collection of RegTicket
type RegExDDTickets []RegExDDTicket

// RegExDDTicket represents pastel registration ticket.
type RegExDDTicket struct {
	Height            int               `json:"height"`
	TXID              string            `json:"txid"`
	RegExDDTicketData RegExDDTicketData `json:"ticket"`
}

// RegExDDTicketData represents external dupe detection registration ticket
type RegExDDTicketData struct {
	Type           string           `json:"type"`
	Version        int              `json:"version"`
	Signatures     TicketSignatures `json:"signatures"`
	Key1           string           `json:"key1"`
	Key2           string           `json:"key2"`
	CreatorHeight  int              `json:"creator_height"`
	ExDDTicket     []byte           `json:"edd_ticket"`
	ExDDTicketData ExDDTicket       `json:"-"`
}

// ExDDTicket represents external dupe detection ticket properties
type ExDDTicket struct {
	Version              int                 `json:"external_dupe_detection_version"`
	BlockHash            string              `json:"block_hash"`
	BlockNum             int                 `json:"block_num"`
	Identifier           string              `json:"external_dupe_detection_request_identifier"`
	ImageHash            string              `json:"image_file_sha3_256_hash"`
	MaximumFee           float64             `json:"maximum_psl_fee_for_external_dupe_detection"`
	SendingAddress       string              `json:"external_dupe_detection_sending_address"`
	EffectiveTotalFee    float64             `json:"effective_total_fee_paid_by_user_in_psl"`
	KamedilaJSONListHash string              `json:"kademlia_hash_list_json_sha3_256_hash"`
	SupernodesSignature  []map[string]string `json:"supernode_signatures_on_external_dupe_request"`
	UserSignature        []map[string]string `json:"user_signature_on_external_dupe_request"`
}

// RegisterExDDRequest represents request to register an dupe detection request
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "key1", "key2", 100)")
type RegisterExDDRequest struct {
	Ticket      *ExDDTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
}

// GetRegisterExDDFeeRequest represents a request to get registration fee
type GetRegisterExDDFeeRequest struct {
	Ticket      *ExDDTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Passphrase  string
	Key1        string
	Key2        string
	Fee         int64
	ImgSizeInMb int64
}

// EncodeExDDTicket encodes  ExDDTicket into byte array
func EncodeExDDTicket(ticket *ExDDTicket) ([]byte, error) {
	b, err := json.Marshal(ticket)
	if err != nil {
		return nil, errors.Errorf("failed to marshal edd ticket: %w", err)
	}

	return b, nil
}

// DecodeExDDTicket decoded byte array into ExDDTicket
func DecodeExDDTicket(b []byte) (*ExDDTicket, error) {
	res := ExDDTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("failed to unmarshal edd ticket: %w", err)
	}

	return &res, nil
}
