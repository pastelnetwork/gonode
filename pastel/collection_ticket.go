package pastel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
)

// CollectionRegTicket represents pastel collection ticket.
type CollectionRegTicket struct {
	Height                  int                     `json:"height"`
	TXID                    string                  `json:"txid"`
	CollectionRegTicketData CollectionRegTicketData `json:"ticket"`
}

// CollectionRegTicketData is Pastel collection ticket structure
type CollectionRegTicketData struct {
	Type                 string              `json:"type"`
	Version              int                 `json:"version"`
	Signatures           RegTicketSignatures `json:"signatures"`
	PermittedUsers       []string            `json:"permitted_users"`
	Key                  string              `json:"key"`
	Label                string              `json:"label"`
	CreatorHeight        uint                `json:"creator_height"`
	RoyaltyAddress       string              `json:"royalty_address"`
	StorageFee           int64               `json:"storage_fee"`
	CollectionTicketData CollectionTicket    `json:"collection_ticket"`
}

// CollectionTicket is the json representation of CollectionTicket inside Pastel's CollectionRegTicketData
type CollectionTicket struct {
	CollectionTicketVersion                 int       `json:"collection_ticket_version"`
	CollectionName                          string    `json:"collection_name"`
	ItemType                                string    `json:"item_type"`
	Creator                                 string    `json:"creator"`
	ListOfPastelIDsOfAuthorizedContributors []string  `json:"list_of_pastelids_of_authorized_contributors"`
	BlockNum                                uint      `json:"blocknum"`
	BlockHash                               string    `json:"block_hash"`
	CollectionFinalAllowedBlockHeight       uint      `json:"collection_final_allowed_block_height"`
	MaxCollectionEntries                    uint      `json:"max_collection_entries"`
	CollectionItemCopyCount                 uint      `json:"collection_item_copy_count"`
	Royalty                                 float64   `json:"royalty"`
	Green                                   bool      `json:"green"`
	AppTicket                               string    `json:"app_ticket"`
	AppTicketData                           AppTicket `json:"-"`
}

// CollectionActTickets is an array of CollectionActTicket
type CollectionActTickets []CollectionActTicket

// CollectionActTicket represents pastel collection-activation ticket.
type CollectionActTicket struct {
	Height                  int                     `json:"height"`
	TXID                    string                  `json:"txid"`
	CollectionActTicketData CollectionActTicketData `json:"ticket"`
}

// CollectionActTicketData represents activation collection ticket properties
type CollectionActTicketData struct {
	PastelID          string `json:"pastelID"`
	Signature         string `json:"signature"`
	Type              string `json:"type"`
	CreatorHeight     int    `json:"creator_height"`
	RegTXID           string `json:"reg_txid"`
	StorageFee        int    `json:"storage_fee"`
	Version           int    `json:"version"`
	CalledAt          int    `json:"called_at"`
	IsFull            bool   `json:"is_full"`
	IsExpiredByHeight bool   `json:"is_expired_by_height"`
}

// EncodeCollectionTicket encodes CollectionTicket to bytes
func EncodeCollectionTicket(ticket *CollectionTicket) ([]byte, error) {
	appTicketBytes, err := json.Marshal(ticket.AppTicketData)
	if err != nil {
		return nil, errors.Errorf("marshal api ticket: %w", err)
	}

	appTicket := base64.RawStdEncoding.EncodeToString(appTicketBytes)
	ticket.AppTicket = appTicket

	b, err := json.Marshal(ticket)
	if err != nil {
		return nil, errors.Errorf("marshal action ticket: %w", err)
	}

	return b, nil
}

// DecodeCollectionTicket decodes collection ticket from bytes
func DecodeCollectionTicket(b []byte) (*CollectionTicket, error) {
	res := &CollectionTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("unmarshal collection ticket: %w", err)
	}

	appDecodedBytes, err := base64.RawStdEncoding.DecodeString(res.AppTicket)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %v", err)
	}

	appTicket := &AppTicket{}
	err = json.Unmarshal(appDecodedBytes, appTicket)
	if err != nil {
		return nil, errors.Errorf("unmarshal api sense ticket: %w", err)
	}
	res.AppTicketData = *appTicket

	return res, nil
}

// DecodeCollectionAppTicket decodes collection app ticket from base64 string
func (regTicket *CollectionRegTicket) DecodeCollectionAppTicket() error {
	appDecodedBytes, err := base64.RawStdEncoding.DecodeString(regTicket.CollectionRegTicketData.CollectionTicketData.AppTicket)
	if err != nil {
		return fmt.Errorf("base64 decode: %v", err)
	}

	appTicket := &AppTicket{}
	err = json.Unmarshal(appDecodedBytes, appTicket)
	if err != nil {
		return errors.Errorf("unmarshal api sense ticket: %w", err)
	}
	regTicket.CollectionRegTicketData.CollectionTicketData.AppTicketData = *appTicket

	return nil
}

// CollectionTicketSignatures represents signatures from parties
type CollectionTicketSignatures struct {
	Principal map[string]string `json:"principal,omitempty"`
	Mn2       map[string]string `json:"mn2,omitempty"`
	Mn3       map[string]string `json:"mn3,omitempty"`
}

// EncodeCollectionTicketSignatures encodes CollectionTicketSignatures into byte array
func EncodeCollectionTicketSignatures(signatures CollectionTicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("marshal signatures: %w", err)
	}

	return b, nil
}

// RegisterCollectionRequest represents request to register a collection ticket
type RegisterCollectionRequest struct {
	Ticket      *CollectionTicket
	Signatures  *CollectionTicketSignatures
	Mn1PastelID string
	Passphrase  string
	Label       string
	Fee         int64
}

// ActivateCollectionRequest - represents request to activate an action ticket
type ActivateCollectionRequest struct {
	RegTxID    string
	BlockNum   int
	Fee        int64
	PastelID   string
	Passphrase string
}
