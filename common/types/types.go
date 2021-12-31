package types

// MeshedSuperNode represents meshed sn
type MeshedSuperNode struct {
	SessID string
	NodeID string
}

// NftRegMetadata represents nft reg metadata
type NftRegMetadata struct {
	CreatorPastelID string
	BlockHash       string
}

// ActionRegMetadata represents action reg metadata
type ActionRegMetadata struct {
	CreatorPastelID string
	BlockHash       string
	BurnTxID        string
}
