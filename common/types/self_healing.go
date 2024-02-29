package types

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"
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
	// SelfHealingAcknowledgementMessage represents the acknowledgement message
	SelfHealingAcknowledgementMessage
)

func (s SelfHealingMessageType) String() string {
	messages := [...]string{"", "challenge", "response", "verification", "completion", "acknowledgement"}
	if s < 1 || int(s) >= len(messages) {
		return "unknown"
	}

	return messages[s]
}

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

func (t TicketType) String() string {
	tickets := [...]string{"", "cascade", "sense", "nft"}
	if t < 1 || int(t) >= len(tickets) {
		return "unknown"
	}

	return tickets[t]
}

// PingInfo represents the structure of data to be inserted into the ping_history table
type PingInfo struct {
	ID                               int          `db:"id"`
	SupernodeID                      string       `db:"supernode_id"`
	IPAddress                        string       `db:"ip_address"`
	TotalPings                       int          `db:"total_pings"`
	TotalSuccessfulPings             int          `db:"total_successful_pings"`
	AvgPingResponseTime              float64      `db:"avg_ping_response_time"`
	IsOnline                         bool         `db:"is_online"`
	IsOnWatchlist                    bool         `db:"is_on_watchlist"`
	IsAdjusted                       bool         `db:"is_adjusted"`
	CumulativeResponseTime           float64      `db:"cumulative_response_time"`
	CreatedAt                        time.Time    `db:"created_at"`
	UpdatedAt                        time.Time    `db:"updated_at"`
	LastSeen                         sql.NullTime `db:"last_seen"`
	MetricsLastBroadcastAt           sql.NullTime `db:"metrics_last_broadcast_at"`
	GenerationMetricsLastBroadcastAt sql.NullTime `db:"generation_metrics_last_broadcast_at"`
	ExecutionMetricsLastBroadcastAt  sql.NullTime `db:"execution_metrics_last_broadcast_at"`
	LastResponseTime                 float64      `db:"-"`
}

// PingInfos represents array of ping info
type PingInfos []PingInfo

// SelfHealingReports represents the self-healing metrics for each challenge
type SelfHealingReports map[string]SelfHealingReport

// SelfHealingReport represents the self-healing challenges
type SelfHealingReport map[string]SelfHealingMessages

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
	Error                    string     `json:"error"`
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
	TxID                             string     `json:"tx_id"`
	TicketType                       TicketType `json:"ticket_type"`
	MissingKeys                      []string   `json:"missing_keys"`
	ReconstructedFileHash            []byte     `json:"reconstructed_file_hash"`
	IsReconstructionRequired         bool       `json:"is_reconstruction_required"`
	IsReconstructionRequiredByHealer bool       `json:"is_reconstruction_required_by_healer"`
	RaptorQSymbols                   []byte     `json:"raptor_q_symbols"`
	FileIDs                          []string   `json:"sense_file_ids"`
	IsVerified                       bool       `json:"is_verified"`
	Message                          string     `json:"message"`
}

// SelfHealingVerificationData represents the verification data for self-healing challenge
type SelfHealingVerificationData struct {
	NodeID         string            `json:"node_id"`
	NodeAddress    string            `json:"node_address"`
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

// SelfHealingChallengeEvent represents the challenge event that needs to be healed.
type SelfHealingChallengeEvent struct {
	ID          int64
	TriggerID   string
	ChallengeID string
	TicketID    string
	Data        []byte
	SenderID    string
	IsProcessed bool
	ExecMetric  SelfHealingExecutionMetric
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Hash returns the hash of the self-healing challenge reports
func (s SelfHealingReports) Hash() string {
	data, _ := json.Marshal(s)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}
