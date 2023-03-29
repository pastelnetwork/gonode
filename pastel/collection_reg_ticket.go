package pastel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pastelnetwork/gonode/common/b85"
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
	CollectionTicket     []byte              `json:"collection_ticket"`
	CollectionTicketData CollectionTicket    `json:"-"`
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
	MaxCollectionEntries                    uint      `json:"max_collection_entries"`
	CollectionItemCopyCount                 uint      `json:"collection_item_copy_count"`
	Royalty                                 float64   `json:"royalty"`
	Green                                   bool      `json:"green"`
	AppTicket                               string    `json:"app_ticket"`
	AppTicketData                           AppTicket `json:"-"`
}

// DecodeCollectionTicket decoded byte array into ArtTicket
func DecodeCollectionTicket(b []byte) (*CollectionTicket, error) {
	res := CollectionTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("error unmarshalling collection ticket: %w", err)
	}

	appDecodedBytes, err := base64.RawStdEncoding.DecodeString(res.AppTicket)
	if err != nil {
		appDecodedBytes, err = b85.Decode(res.AppTicket)
		if err != nil {
			return nil, fmt.Errorf("b85 decode: %v", err)
		}
	}

	appTicket := AppTicket{}
	err = json.Unmarshal(appDecodedBytes, &appTicket)
	if err != nil {
		return nil, errors.Errorf("unmarshal app ticket: %w", err)
	}

	return &CollectionTicket{
		CollectionTicketVersion:                 res.CollectionTicketVersion,
		CollectionName:                          res.CollectionName,
		ItemType:                                res.ItemType,
		Creator:                                 res.Creator,
		ListOfPastelIDsOfAuthorizedContributors: res.ListOfPastelIDsOfAuthorizedContributors,
		BlockNum:                                res.BlockNum,
		BlockHash:                               res.BlockHash,
		MaxCollectionEntries:                    res.MaxCollectionEntries,
		CollectionItemCopyCount:                 res.CollectionItemCopyCount,
		Royalty:                                 res.Royalty,
		Green:                                   res.Green,
		AppTicketData:                           appTicket,
	}, nil
}
