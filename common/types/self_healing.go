package types

import (
	"database/sql"
	"time"
)

// SelfHealingMessageType represents the type of message sent in the self-healing process
type SelfHealingMessageType int

const (
	// SelfHealingChallengeMessage represents the challenge message
	SelfHealingChallengeMessage SelfHealingMessageType = iota + 1
	// SelfHealingResponseMessage represents the response message
	SelfHealingResponseMessage
	// SelfHealingVerificationMessage represents the verification message
	SelfHealingVerificationMessage
	// SelfHealingCompletionMessage represents the challenge message processed successfully
	SelfHealingCompletionMessage
)

// TicketType represents the type of ticket; nft, cascade, sense
type TicketType int

const (
	// TicketTypeCascade represents the cascade ticket type
	TicketTypeCascade TicketType = iota + 1
	// TicketTypeSense represents the sense ticket type
	TicketTypeSense
	// TicketTypeNFT represents the NFT ticket type
	TicketTypeNFT
)

// PingInfo represents the structure of data to be inserted into the ping_history table
type PingInfo struct {
	ID                     int          `db:"id"`
	SupernodeID            string       `db:"supernode_id"`
	IPAddress              string       `db:"ip_address"`
	TotalPings             int          `db:"total_pings"`
	TotalSuccessfulPings   int          `db:"total_successful_pings"`
	AvgPingResponseTime    float64      `db:"avg_ping_response_time"`
	IsOnline               bool         `db:"is_online"`
	IsOnWatchlist          bool         `db:"is_on_watchlist"`
	IsAdjusted             bool         `db:"is_adjusted"`
	CumulativeResponseTime float64      `db:"cumulative_response_time"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
	LastSeen               sql.NullTime `db:"last_seen"`
	MetricsLastBroadcastAt sql.NullTime `json:"metrics_last_broadcast_at"`
	LastResponseTime       float64      `db:"-"`
}

// PingInfos represents array of ping info
type PingInfos []PingInfo

// SelfHealingMessages represents the self-healing metrics for each challenge = message_type = 3
type SelfHealingMessages []SelfHealingMessage

// SelfHealingMessage represents the self-healing message
type SelfHealingMessage struct {
	TriggerID              string                 `json:"trigger_id"`
	MessageType            SelfHealingMessageType `json:"message_type"`
	SelfHealingMessageData SelfHealingMessageData `json:"data"`
	SenderID               string                 `json:"sender_id"`
	SenderSignature        []byte                 `json:"sender_signature"`
}

// SelfHealingMessageData represents the self-healing message data == message_type = 2
type SelfHealingMessageData struct {
	ChallengerID string                      `json:"challenger_id"`
	RecipientID  string                      `json:"recipient_id"`
	Challenge    SelfHealingChallengeData    `json:"challenge"`
	Response     SelfHealingResponseData     `json:"response"`
	Verification SelfHealingVerificationData `json:"verification"`
}

// SelfHealingChallengeData represents the challenge data for self-healing sent by the challenger
type SelfHealingChallengeData struct {
	Block            int32             `json:"block"`
	Merkelroot       string            `json:"merkelroot"`
	Timestamp        time.Time         `json:"timestamp"`
	ChallengeTickets []ChallengeTicket `json:"challenge_tickets"`
	NodesOnWatchlist string            `json:"nodes_on_watchlist"`
}

// ChallengeTicket represents the ticket details for self-healing challenge
type ChallengeTicket struct {
	TxID        string     `json:"tx_id"`
	TicketType  TicketType `json:"ticket_type"`
	MissingKeys []string   `json:"missing_keys"`
	DataHash    []byte     `json:"data_hash"`
	Recipient   string     `json:"recipient"`
}

// RespondedTicket represents the details of ticket responded in a self-healing challenge
type RespondedTicket struct {
	TxID                     string     `json:"tx_id"`
	TicketType               TicketType `json:"ticket_type"`
	MissingKeys              []string   `json:"missing_keys"`
	ReconstructedFileHash    []byte     `json:"reconstructed_file_hash"`
	FileIDs                  []string   `json:"sense_file_ids"`
	RaptorQSymbols           []byte     `json:"raptor_q_symbols"`
	IsReconstructionRequired bool       `json:"is_reconstruction_required"`
}

// SelfHealingResponseData represents the response data for self-healing sent by the recipient
type SelfHealingResponseData struct {
	ChallengeID     string          `json:"challenge_id"`
	Block           int32           `json:"block"`
	Merkelroot      string          `json:"merkelroot"`
	Timestamp       time.Time       `json:"timestamp"`
	RespondedTicket RespondedTicket `json:"responded_ticket"`
	Verifiers       []string        `json:"verifiers"`
}

// VerifiedTicket represents the details of ticket verified in self-healing challenge
type VerifiedTicket struct {
	TxID                     string     `json:"tx_id"`
	TicketType               TicketType `json:"ticket_type"`
	MissingKeys              []string   `json:"missing_keys"`
	ReconstructedFileHash    []byte     `json:"reconstructed_file_hash"`
	IsReconstructionRequired bool       `json:"is_reconstruction_required"`
	RaptorQSymbols           []byte     `json:"raptor_q_symbols"`
	FileIDs                  []string   `json:"sense_file_ids"`
	IsVerified               bool       `json:"is_verified"`
	Message                  string     `json:"message"`
}

// SelfHealingVerificationData represents the verification data for self-healing challenge
type SelfHealingVerificationData struct {
	ChallengeID    string            `json:"challenge_id"`
	Block          int32             `json:"block"`
	Merkelroot     string            `json:"merkelroot"`
	Timestamp      time.Time         `json:"timestamp"`
	VerifiedTicket VerifiedTicket    `json:"verified_ticket"`
	VerifiersData  map[string][]byte `json:"verifiers_data"`
}

// SelfHealingGenerationMetric represents the self-healing generation metrics for trigger events
type SelfHealingGenerationMetric struct {
	ID              int       `db:"id"`
	TriggerID       string    `db:"trigger_id"`
	MessageType     int       `db:"message_type"`
	Data            []byte    `db:"data"`
	SenderID        string    `db:"sender_id"`
	SenderSignature []byte    `db:"sender_signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// CombinedSelfHealingMetrics represents the combination of generation and execution metrics
type CombinedSelfHealingMetrics struct {
	GenerationMetrics []SelfHealingGenerationMetric
	ExecutionMetrics  []SelfHealingExecutionMetric
}

// SelfHealingExecutionMetric represents the self-healing execution metrics for trigger events
type SelfHealingExecutionMetric struct {
	ID              int       `db:"id"`
	TriggerID       string    `db:"trigger_id"`
	ChallengeID     string    `db:"challenge_id"`
	MessageType     int       `db:"message_type"`
	Data            []byte    `db:"data"`
	SenderID        string    `db:"sender_id"`
	SenderSignature []byte    `db:"sender_signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// SelfHealingMetricType represents the type of self-healing metric
type SelfHealingMetricType int

const (
	// GenerationSelfHealingMetricType represents the generation metric for self-healing
	GenerationSelfHealingMetricType SelfHealingMetricType = 1
	// ExecutionSelfHealingMetricType represents the execution metric for self-healing
	ExecutionSelfHealingMetricType SelfHealingMetricType = 2
)

// ProcessBroadcastMetricsRequest represents the request for broadcasting metrics
type ProcessBroadcastMetricsRequest struct {
	Data            []byte                `json:"data"`
	Type            SelfHealingMetricType `json:"type"`
	SenderID        string                `json:"sender_id"`
	SenderSignature []byte                `json:"sender_signature"`
}

// SelfHealingMetrics represents the self-healing metrics for each challenge
type SelfHealingMetrics struct {
	ChallengeID                   string `db:"challenge_id"`
	SentTicketsForSelfHealing     int    `db:"sent_tickets_for_self_healing"`
	EstimatedMissingKeys          int    `db:"estimated_missing_keys"`
	TicketsInProgress             int    `db:"tickets_in_progress"`
	TicketsRequiredSelfHealing    int    `db:"tickets_required_self_healing"`
	SuccessfullySelfHealedTickets int    `db:"successfully_self_healed_tickets"`
	SuccessfullyVerifiedTickets   int    `db:"successfully_verified_tickets"`
}
