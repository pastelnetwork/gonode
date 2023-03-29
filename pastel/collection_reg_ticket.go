package pastel

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
	MaxCollectionEntries                    uint      `json:"max_collection_entries"`
	CollectionItemCopyCount                 uint      `json:"collection_item_copy_count"`
	Royalty                                 float64   `json:"royalty"`
	Green                                   bool      `json:"green"`
	AppTicketData                           AppTicket `json:"app_ticket"`
}
