package types

import (
	"database/sql"
	"time"
)

// SelfHealingMessageType represents the type of message sent in the self-healing process
type SelfHealingMessageType int

const (
	//SelfHealingChallengeMessage represents the challenge message
	SelfHealingChallengeMessage SelfHealingMessageType = iota + 1
	//SelfHealingResponseMessage represents the response message
	SelfHealingResponseMessage
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
	CumulativeResponseTime float64      `json:"cumulative_response_time"`
	LastSeen               sql.NullTime `db:"last_seen"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
}

// PingInfos represents array of ping info
type PingInfos []PingInfo

// SelfHealingMessage represents the self-healing message
type SelfHealingMessage struct {
	ChallengeID            string                 `json:"challenge_id"`
	MessageType            SelfHealingMessageType `json:"message_type"`
	SelfHealingMessageData SelfHealingMessageData `json:"data"`
	SenderID               string                 `json:"sender_id"`
	SenderSignature        []byte                 `json:"sender_signature"`
}

// SelfHealingMessageData represents the self-healing message data
type SelfHealingMessageData struct {
	ChallengerID string                   `json:"challenger_id"`
	RecipientID  string                   `json:"recipient_id"`
	Challenge    SelfHealingChallengeData `json:"challenge"`
	Response     SelfHealingResponseData  `json:"response"`
}

// SelfHealingChallengeData represents the challenge data for self-healing sent by the challenger
type SelfHealingChallengeData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
	Tickets    []Ticket  `json:"tickets"`
}

// Ticket represents the ticket details that require self-healing
type Ticket struct {
	TxID                  string     `json:"tx_id"`
	TicketType            TicketType `json:"ticket_type"`
	MissingKeys           []string   `json:"missing_keys"`
	DataHash              string     `json:"data_hash"`
	ReconstructedFileHash string     `json:"reconstructed_file_hash"`
}

// SelfHealingResponseData represents the response data for self-healing sent by the recipient
type SelfHealingResponseData struct {
	Block      int32     `json:"block"`
	Merkelroot string    `json:"merkelroot"`
	Timestamp  time.Time `json:"timestamp"`
	Tickets    []Ticket  `json:"tickets"`
}
