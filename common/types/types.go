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

// StorageChallengeStatus represents possible storage challenge statuses
type StorageChallengeStatus int

const (
	//UndefinedStorageChallengeStatus represents invalid storage challenge type
	UndefinedStorageChallengeStatus StorageChallengeStatus = iota
	//GeneratedStorageChallengeStatus represents when the challenge is stored after generation
	GeneratedStorageChallengeStatus
	//ProcessedStorageChallengeStatus represents when the challenge is stored after processing
	ProcessedStorageChallengeStatus
	//VerifiedStorageChallengeStatus represents when the challenge is stored after verification
	VerifiedStorageChallengeStatus
)

// StorageChallenge represents storage challenge
type StorageChallenge struct {
	ID              int64
	ChallengeID     string
	FileHash        string
	ChallengingNode string
	RespondingNode  string
	VerifyingNodes  string
	GeneratedHash   string
	Status          StorageChallengeStatus
	StartingIndex   int
	EndingIndex     int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SelfHealingStatus represents possible self-healing statuses of failed challenge
type SelfHealingStatus int

const (
	//UndefinedSelfHealingStatus represents invalid status for self-healing operation
	UndefinedSelfHealingStatus SelfHealingStatus = iota
	//CreatedSelfHealingStatus represents when the failed challenge gets stored in DB
	CreatedSelfHealingStatus
	//InProgressSelfHealingStatus represents when the challenge is retrieved for self-healing
	InProgressSelfHealingStatus
	//FailedSelfHealingStatus represents when the reconstruction has been completed
	FailedSelfHealingStatus
	//CompletedSelfHealingStatus represents when the reconstruction has been completed
	CompletedSelfHealingStatus
)

// selfHealingStatusMap represents a map of self-healing statuses
var selfHealingStatusMap = map[SelfHealingStatus]string{
	CreatedSelfHealingStatus:    "Created",
	InProgressSelfHealingStatus: "InProgress",
	FailedSelfHealingStatus:     "Failed",
	CompletedSelfHealingStatus:  "Completed",
}

// String returns the string for self-healing status
func (s SelfHealingStatus) String() string {
	if status, ok := selfHealingStatusMap[s]; ok {
		return status
	}

	return "Undefined"
}

// SelfHealingChallenge represents self-healing challenge
type SelfHealingChallenge struct {
	ID                    int64
	ChallengeID           string
	MerkleRoot            string
	FileHash              string
	ChallengingNode       string
	RespondingNode        string
	VerifyingNode         string
	ReconstructedFileHash []byte
	Status                SelfHealingStatus
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// Fields represents status log
type Fields map[string]interface{}

// Details represents status log details with additional fields
type Details struct {
	Message string
	Fields  Fields
}

// Stringify convert the Details' struct to stringify json
func (d *Details) Stringify() string {
	details, err := json.Marshal(&d)
	if err != nil {
		log.WithError(err).Error("unable to marshal task history details")
		return ""
	}

	return string(details)
}

// NewDetails initializes and return the valid detail object
func NewDetails(msg string, fields Fields) *Details {
	return &Details{
		Message: msg,
		Fields:  fields,
	}
}

// IsValid checks if the status log map is not empty
func (f Fields) IsValid() bool {
	return len(f) != 0
}
