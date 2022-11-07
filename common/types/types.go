package types

import (
	"encoding/json"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
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
	Details   *Details
}

//Fields represents status log
type Fields map[string]interface{}

//Details represents status log details with additional fields
type Details struct {
	Message string
	Fields  Fields
}

//Stringify convert the Details' struct to stringify json
func (d *Details) Stringify() string {
	details, err := json.Marshal(&d)
	if err != nil {
		log.WithError(err).Error("unable to marshal task history details")
		return ""
	}

	return string(details)
}

//NewDetails initializes and return the valid detail object
func NewDetails(msg string, fields Fields) *Details {
	return &Details{
		Message: msg,
		Fields:  fields,
	}
}

//IsValid checks if the status log map is not empty
func (f Fields) IsValid() bool {
	return len(f) != 0
}
