package types

import "time"

// HealthCheckMessageType represents the type of message sent in the health-check process
type HealthCheckMessageType int

const (
	// HealthCheckChallengeMessageType represents the challenge message
	HealthCheckChallengeMessageType HealthCheckMessageType = iota + 1
	// HealthCheckResponseMessageType represents the response message
	HealthCheckResponseMessageType
	// HealthCheckEvaluationMessageType represents the evaluation message
	HealthCheckEvaluationMessageType
	// HealthCheckVerificationMessageType represents the verification message
	HealthCheckVerificationMessageType
	// HealthCheckCompletionMessageType represents the challenge message processed successfully
	HealthCheckCompletionMessageType
	// HealthCheckBroadcastMessageType represents the broadcast message processed successfully
	HealthCheckBroadcastMessageType
	// HealthCheckAffirmationMessageType represents the affirmation message
	HealthCheckAffirmationMessageType
)

// String returns the message string
func (hcm HealthCheckMessageType) String() string {
	switch hcm {
	case HealthCheckChallengeMessageType:
		return "challenge"
	case HealthCheckResponseMessageType:
		return "response"
	case HealthCheckEvaluationMessageType:
		return "evaluation"
	case HealthCheckAffirmationMessageType:
		return "affirmation"
	default:
		return "unknown"
	}
}

// BroadcastHealthCheckMessage represents the healthcheck challenge message that needs to be broadcast after evaluation
type BroadcastHealthCheckMessage struct {
	ChallengeID string
	Challenger  map[string][]byte
	Recipient   map[string][]byte
	Observers   map[string][]byte
}

// BroadcastHealthCheckLogMessage represents the broadcast message log to be stored in the DB
type BroadcastHealthCheckLogMessage struct {
	ChallengeID string `db:"challenge_id"`
	Challenger  string `db:"challenger"`
	Recipient   string `db:"recipient"`
	Observers   string `db:"observers"`
	Data        []byte `db:"data"`
}

// HealthCheckChallengeData represents the data of challenge
type HealthCheckChallengeData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
}

// HealthCheckResponseData represents the data of response
type HealthCheckResponseData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
}

// HealthCheckEvaluationData represents the data of evaluation
type HealthCheckEvaluationData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
	IsVerified bool      `json:"is_verified"`
}

// HealthCheckObserverEvaluationData represents the data of Observer's evaluation
type HealthCheckObserverEvaluationData struct {
	Block                   int32     `json:"block"`
	Merkelroot              string    `json:"merkelroot"`
	IsChallengeTimestampOK  bool      `json:"is_challenge_timestamp_ok"`
	IsProcessTimestampOK    bool      `json:"is_process_timestamp_ok"`
	IsEvaluationTimestampOK bool      `json:"is_evaluation_timestamp_ok"`
	IsRecipientSignatureOK  bool      `json:"is_recipient_signature_ok"`
	IsChallengerSignatureOK bool      `json:"is_challenger_signature_ok"`
	IsEvaluationResultOK    bool      `json:"is_evaluation_result_ok"`
	Timestamp               time.Time `json:"timestamp"`
}

// HealthCheckMessageData represents the health check challenge message data
type HealthCheckMessageData struct {
	ChallengerID         string                            `json:"challenger_id"`
	Challenge            HealthCheckChallengeData          `json:"challenge"`
	Observers            []string                          `json:"observers"`
	RecipientID          string                            `json:"recipient_id"`
	Response             HealthCheckResponseData           `json:"response"`
	ChallengerEvaluation HealthCheckEvaluationData         `json:"challenger_evaluation"`
	ObserverEvaluation   HealthCheckObserverEvaluationData `json:"observer_evaluation"`
}

// HealthCheckMessage represents the healthcheck challenge message
type HealthCheckMessage struct {
	MessageType                HealthCheckMessageType `json:"message_type"`
	ChallengeID                string                 `json:"challenge_id"`
	Data                       HealthCheckMessageData `json:"data"`
	Sender                     string                 `json:"sender"`
	SenderSignature            []byte                 `json:"sender_signature"`
	StorageChallengeSignatures StorageChallengeStatus
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// HealthCheckChallengeMetric represents the metric log to be stored in the DB
type HealthCheckChallengeMetric struct {
	ID          int    `db:"id"`
	MessageType int    `db:"message_type"`
	ChallengeID string `db:"challenge_id"`
	Data        []byte `db:"data"`
	SenderID    string `db:"sender_id"`
}

// HealthCheckChallengeLogMessage represents the message log to be stored in the DB
type HealthCheckChallengeLogMessage struct {
	ID              int       `db:"id"`
	MessageType     int       `db:"message_type"`
	ChallengeID     string    `db:"challenge_id"`
	Data            []byte    `db:"data"`
	Sender          string    `db:"sender_id"`
	SenderSignature []byte    `db:"sender_signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// BroadcastHealthCheckMessageMetrics is the struct for broadcast message metrics
type BroadcastHealthCheckMessageMetrics struct {
	ID          int       `db:"id"`
	ChallengeID string    `db:"challenge_id"`
	Challenger  string    `db:"challenger"`
	Recipient   string    `db:"recipient"`
	Observers   string    `db:"observers"`
	Data        []byte    `db:"data"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ProcessBroadcastHealthCheckChallengeMetricsRequest represents the request for broadcasting metrics
type ProcessBroadcastHealthCheckChallengeMetricsRequest struct {
	Data     []byte `json:"data"`
	SenderID string `json:"sender_id"`
}
