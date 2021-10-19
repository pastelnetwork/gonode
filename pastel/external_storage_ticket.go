package pastel

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

// RegExternalStorageTickets is a collection of RegTicket
type RegExternalStorageTickets []RegExternalStorageTicket

// RegExternalStorageTicket represents pastel registration ticket.
type RegExternalStorageTicket struct {
	Height                       int                          `json:"height"`
	TXID                         string                       `json:"txid"`
	RegExternalStorageTicketData RegExternalStorageTicketData `json:"ticket"`
}

// RegExternalStorageTicketData represents external dupe detection registration ticket
type RegExternalStorageTicketData struct {
	Type                      string                `json:"type"`
	Version                   int                   `json:"version"`
	Signatures                TicketSignatures      `json:"signatures"`
	Key1                      string                `json:"key1"`
	Key2                      string                `json:"key2"`
	CreatorHeight             int                   `json:"creator_height"`
	ExternalStorageTicket     []byte                `json:"external_storage_ticket"`
	ExternalStorageTicketData ExternalStorageTicket `json:"-"`
}

// ExternalStorageTicket represents external dupe detection ticket properties
type ExternalStorageTicket struct {
	Version                                 int                 `json:"storage_version"`
	BlockHash                               string              `json:"block_hash"`
	BlockNum                                int                 `json:"block_num"`
	Identifier                              string              `json:"storage_request_identifier"`
	FileSizeInMegabytes                     float64             `json:"file_size"`
	MaximumFee                              float64             `json:"maximum_psl_fee_for_storage"`
	StorageFeeMultiplierFactor              float64             `json:"psl_storage_fee_multiplier_factor"`
	DesiredRQSymbolFilesToOriginalDataRatio float64             `json:"desired_rq_symbol_files_to_original_data_ratio"`
	DesiredRedundancyPerRQSymbolFile        int                 `json:"desired_redundancy_per_rq_symbol_file"`
	SendingAddress                          string              `json:"storage_sending_address"`
	EffectiveTotalFee                       float64             `json:"effective_total_fee_paid_by_user_in_psl"`
	KamedilaJSONListHash                    string              `json:"kademlia_hash_list_json_sha3_256_hash"`
	SupernodesSignature                     []map[string]string `json:"supernode_signatures_on_external_dupe_request"`
	UserSignature                           []map[string]string `json:"user_signature_on_external_dupe_request"`
}

// RegisterExternalStorageRequest represents request to register an dupe detection request
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "key1", "key2", 100)")
type RegisterExternalStorageRequest struct {
	Ticket      *ExternalStorageTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
}

// GetRegisterExternalStorageFeeRequest represents a request to get registration fee
type GetRegisterExternalStorageFeeRequest struct {
	Ticket      *ExternalStorageTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Passphrase  string
	Key1        string
	Key2        string
	Fee         int64
	ImgSizeInMb int64
}

// EncodeExternalStorageTicket encodes  ExternalStorageTicket into byte array
func EncodeExternalStorageTicket(ticket *ExternalStorageTicket) ([]byte, error) {
	b, err := json.Marshal(ticket)
	if err != nil {
		return nil, errors.Errorf("failed to marshal edd ticket: %w", err)
	}

	return b, nil
}

// DecodeExternalStorageTicket decoded byte array into ExternalStorageTicket
func DecodeExternalStorageTicket(b []byte) (*ExternalStorageTicket, error) {
	res := ExternalStorageTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("failed to unmarshal edd ticket: %w", err)
	}

	return &res, nil
}
