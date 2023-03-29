package pastel

// CollectionRegTicket represents pastel collection ticket.
type CollectionRegTicket struct {
	Height                  int                     `json:"height"`
	TXID                    string                  `json:"txid"`
	CollectionRegTicketData CollectionRegTicketData `json:"ticket"`
}

// CollectionRegTicketData is Pastel collection ticket structure
type CollectionRegTicketData struct {
	Type                string              `json:"type"`
	Version             int                 `json:"version"`
	NftCollectionTicket string              `json:"nft_collection_ticket"`
	Signatures          RegTicketSignatures `json:"signatures"`
	PermittedUsers      []string            `json:"permitted_users"`
	Key                 string              `json:"key"`
	Label               string              `json:"label"`
	CreatorHeight       uint                `json:"creator_height"`
	ClosingHeight       uint                `json:"closing_height"`
	NftMaxCount         uint                `json:"nft_max_count"`
	NftCopyCount        uint                `json:"nft_copy_count"`
	Royalty             float64             `json:"royalty"`
	RoyaltyAddress      string              `json:"royalty_address"`
	Green               bool                `json:"green"`
	StorageFee          int64               `json:"storage_fee"`
	NFTCollectionTicket NFTCollectionTicket `json:"-"`
}

// NFTCollectionTicket represents json of NFTCollectionTicket inside Pastel's CollectionRegTicketData
type NFTCollectionTicket struct {
	NftCollectionTicketVersion int       `json:"nft_collection_ticket_version"`
	NftCollectionName          string    `json:"nft_collection_name"`
	Creator                    string    `json:"creator"`
	PermittedUsers             []string  `json:"permitted_users"`
	BlockNum                   uint      `json:"blocknum"`
	BlockHash                  string    `json:"block_hash"`
	ClosingHeight              uint      `json:"closing_height"`
	NftMaxCount                uint      `json:"nft_max_count"`
	NftCopyCount               uint      `json:"nft_copy_count"`
	Royalty                    float64   `json:"royalty"`
	Green                      bool      `json:"green"`
	AppTicket                  string    `json:"app_ticket"`
	AppTicketData              AppTicket `json:"-"`
}
