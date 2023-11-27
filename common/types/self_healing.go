package types

import (
	"database/sql"
	"time"
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
