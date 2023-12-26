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
	//SelfHealingVerificationMessage represents the verification message
	SelfHealingVerificationMessage
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
	LastSeen               sql.NullTime `db:"last_seen"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
	LastResponseTime       float64      `db:"-"`
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
}

// ChallengeTicket represents the ticket details for self-healing challenge
type ChallengeTicket struct {
	TxID        string     `json:"tx_id"`
	TicketType  TicketType `json:"ticket_type"`
	MissingKeys []string   `json:"missing_keys"`
	DataHash    []byte     `json:"data_hash"`
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
	Block            int32             `json:"block"`
	Merkelroot       string            `json:"merkelroot"`
	Timestamp        time.Time         `json:"timestamp"`
	RespondedTickets []RespondedTicket `json:"responded_tickets"`
}

// VerifiedTicket represents the details of ticket verified in self-healing challenge
type VerifiedTicket struct {
	TxID                     string     `json:"tx_id"`
	TicketType               TicketType `json:"ticket_type"`
	MissingKeys              []string   `json:"missing_keys"`
	DataHash                 string     `json:"data_hash"`
	ReconstructedFileHash    []byte     `json:"reconstructed_file_hash"`
	IsReconstructionRequired bool       `json:"is_reconstruction_required"`
	IsVerified               bool       `json:"is_verified"`
	Message                  string     `json:"message"`
}

// SelfHealingVerificationData represents the verification data for self-healing challenge
type SelfHealingVerificationData struct {
	Block           int32            `json:"block"`
	Merkelroot      string           `json:"merkelroot"`
	Timestamp       time.Time        `json:"timestamp"`
	VerifiedTickets []VerifiedTicket `json:"verified_tickets"`
}

// SelfHealingLogMessage represents the message log to be stored in the DB
type SelfHealingLogMessage struct {
	ID              int       `db:"id"`
	MessageType     int       `db:"message_type"`
	ChallengeID     string    `db:"challenge_id"`
	Data            []byte    `db:"data"`
	SenderID        string    `db:"sender_id"`
	SenderSignature []byte    `db:"sender_signature"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
