package pastel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
)

// ContractTicket defines the contract ticket
type ContractTicket struct {
	ContractTicketData string `json:"contract_ticket"`
	Key                string `json:"key"`
	SecondaryKey       string `json:"secondary_key"`
	SubType            string `json:"sub_type"`
	Timestamp          int64  `json:"timestamp"`
	Type               string `json:"type"`
	Version            int    `json:"version"`
}

// TxInfo defines the transaction information
type TxInfo struct {
	CompressedSize       int    `json:"compressed_size"`
	CompressionRatio     string `json:"compression_ratio"`
	IsCompressed         bool   `json:"is_compressed"`
	MultisigOutputsCount int    `json:"multisig_outputs_count"`
	MultisigTxTotalFee   int    `json:"multisig_tx_total_fee"`
	UncompressedSize     int    `json:"uncompressed_size"`
}

// Contract defines the contract
type Contract struct {
	Height int            `json:"height"`
	Ticket ContractTicket `json:"ticket"`
	TxInfo TxInfo         `json:"tx_info"`
	TxID   string         `json:"txid"`
}

// -------------------------------------- Contract Types -------------------------------------------//

// CascadeMultiVolumeTicket defines the cascade multi volume ticket contract type
type CascadeMultiVolumeTicket struct {
	NameOfOriginalFile        string         `json:"name_of_original_file"`
	SizeOfOriginalFileMB      int            `json:"size_of_original_file_mb"`
	SHA3256HashOfOriginalFile string         `json:"sha3_256_hash_of_original_file"`
	Volumes                   map[int]string `json:"volumes"` // key (int): index of the volume, value (string): txid of the volume
}

func (c *Contract) GetCascadeMultiVolumeMetadataTicket() (t CascadeMultiVolumeTicket, err error) {
	if c.Ticket.SubType != string(CascadeMultiVolumeMetadata) {
		return t, errors.New("contract is not of type cascade_multi_volume_metadata")
	}

	data, err := base64.StdEncoding.DecodeString(c.Ticket.ContractTicketData)
	if err != nil {
		return t, fmt.Errorf("unable to b64 decode contract: %w", err)
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return t, fmt.Errorf("unable to decode contract: %w", err)
	}

	return t, nil
}
