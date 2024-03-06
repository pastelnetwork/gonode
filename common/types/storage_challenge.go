package types

import (
	"encoding/json"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/errors"
)

// MessageType represents the type of message
type MessageType int

const (
	// ChallengeMessageType represents the challenge message
	ChallengeMessageType MessageType = iota + 1
	// ResponseMessageType represents the response message
	ResponseMessageType
	// EvaluationMessageType represents the evaluation message
	EvaluationMessageType
	// AffirmationMessageType represents the affirmation message
	AffirmationMessageType
	//BroadcastMessageType represents the message that needs to be broadcast
	BroadcastMessageType
)

// String returns the message string
func (m MessageType) String() string {
	switch m {
	case ChallengeMessageType:
		return "challenge"
	case ResponseMessageType:
		return "response"
	case EvaluationMessageType:
		return "evaluation"
	case AffirmationMessageType:
		return "affirmation"
	default:
		return "unknown"
	}
}

// MessageTypeFromString returns the message type from string
func MessageTypeFromString(str string) (MessageType, error) {
	switch str {
	case "challenge":
		return ChallengeMessageType, nil
	case "response":
		return ResponseMessageType, nil
	case "evaluation":
		return EvaluationMessageType, nil
	case "affirmation":
		return AffirmationMessageType, nil
	default:
		return 0, errors.New("invalid message type string")
	}
}

// StorageChallengeSignatures represents the signature struct for broadcasting
type StorageChallengeSignatures struct {
	Challenger map[string]string `json:"challenger,omitempty"`
	Recipient  map[string]string `json:"recipient,omitempty"`
	Obs        map[string]string `json:"obs,omitempty"`
}

// Message represents the storage challenge message
type Message struct {
	MessageType     MessageType `json:"message_type"`
	ChallengeID     string      `json:"challenge_id"`
	Data            MessageData `json:"data"`
	Sender          string      `json:"sender"`
	SenderSignature []byte      `json:"sender_signature"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// BroadcastMessage represents the storage challenge message that needs to be broadcast after evaluation
type BroadcastMessage struct {
	ChallengeID string
	Challenger  map[string][]byte
	Recipient   map[string][]byte
	Observers   map[string][]byte
}

type MessageDataList []MessageData

// MessageData represents the storage challenge message data
type MessageData struct {
	ChallengerID         string                 `json:"challenger_id"`
	Challenge            ChallengeData          `json:"challenge"`
	Observers            []string               `json:"observers"`
	RecipientID          string                 `json:"recipient_id"`
	Response             ResponseData           `json:"response"`
	ChallengerEvaluation EvaluationData         `json:"challenger_evaluation"`
	ObserverEvaluation   ObserverEvaluationData `json:"observer_evaluation"`
}

// ChallengeData represents the data of challenge
type ChallengeData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
	FileHash   string    `json:"file_hash"`
	StartIndex int       `json:"start_index"`
	EndIndex   int       `json:"end_index"`
}

// ResponseData represents the data of response
type ResponseData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Hash       string    `json:"hash"`
	Timestamp  time.Time `json:"timestamp"`
}

// EvaluationData represents the data of evaluation
type EvaluationData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
	Hash       string    `json:"hash"`
	IsVerified bool      `json:"is_verified"`
}

// ObserverEvaluationData represents the data of Observer's evaluation
type ObserverEvaluationData struct {
	Block                   int32     `json:"block"`
	Merkelroot              string    `json:"merkelroot"`
	IsChallengeTimestampOK  bool      `json:"is_challenge_timestamp_ok"`
	IsProcessTimestampOK    bool      `json:"is_process_timestamp_ok"`
	IsEvaluationTimestampOK bool      `json:"is_evaluation_timestamp_ok"`
	IsRecipientSignatureOK  bool      `json:"is_recipient_signature_ok"`
	IsChallengerSignatureOK bool      `json:"is_challenger_signature_ok"`
	IsEvaluationResultOK    bool      `json:"is_evaluation_result_ok"`
	Reason                  string    `json:"reason"`
	TrueHash                string    `json:"true_hash"`
	Timestamp               time.Time `json:"timestamp"`
}

// StorageChallengeLogMessage represents the message log to be stored in the DB
type StorageChallengeLogMessage struct {
	ID              int       `db:"id"`
	MessageType     int       `db:"message_type"`
	ChallengeID     string    `db:"challenge_id"`
	Data            []byte    `db:"data"`
	Sender          string    `db:"sender_id"`
	SenderSignature []byte    `db:"sender_signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// StorageChallengeMetric represents the metric log to be stored in the DB
type StorageChallengeMetric struct {
	ID          int    `db:"id"`
	MessageType int    `db:"message_type"`
	ChallengeID string `db:"challenge_id"`
	Data        []byte `db:"data"`
	SenderID    string `db:"sender_id"`
}

// BroadcastLogMessage represents the broadcast message log to be stored in the DB
type BroadcastLogMessage struct {
	ChallengeID string `db:"challenge_id"`
	Challenger  string `db:"challenger"`
	Recipient   string `db:"recipient"`
	Observers   string `db:"observers"`
	Data        []byte `db:"data"`
}

// BroadcastMessageMetrics is the struct for broadcast message metrics
type BroadcastMessageMetrics struct {
	ID          int       `db:"id"`
	ChallengeID string    `db:"challenge_id"`
	Challenger  string    `db:"challenger"`
	Recipient   string    `db:"recipient"`
	Observers   string    `db:"observers"`
	Data        []byte    `db:"data"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ProcessBroadcastChallengeMetricsRequest represents the request for broadcasting metrics
type ProcessBroadcastChallengeMetricsRequest struct {
	Data     []byte `json:"data"`
	SenderID string `json:"sender_id"`
}

type StorageChallengeMessages []Message

// Hash returns the hash of the storage-challenge challenge log data

func (mdl StorageChallengeMessages) Hash() string {
	data, _ := json.Marshal(mdl)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}

// NScMetric gets the latest challenge IDs from the DB
type NScMetric struct {
	Count       int
	ChallengeID string
	CreatedAt   string
}
