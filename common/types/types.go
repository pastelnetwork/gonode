package types

import (
	"time"
)

// MeshedSuperNode represents meshed sn
type MeshedSuperNode struct {
	SessID string
	NodeID string
}

// NftRegMetadata represents nft reg metadata
type NftRegMetadata struct {
	CreatorPastelID string
	BlockHash       string
	BlockHeight     string
	Timestamp       string
}

// ActionRegMetadata represents action reg metadata
type ActionRegMetadata struct {
	CreatorPastelID string
	BlockHash       string
	BurnTxID        string
	BlockHeight     string
	Timestamp       string
	OpenAPISubsetID string
	EstimatedFee    int64
}

// TaskHistory represents task history
type TaskHistory struct {
	ID        int
	TaskID    string
	CreatedAt time.Time
	Status    string
}
