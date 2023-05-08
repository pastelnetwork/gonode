package domain

import "time"

// NodeReplicationInfo is the struct for replication info
type NodeReplicationInfo struct {
	LastReplicatedAt *time.Time `json:"last_replicated,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedAt        time.Time  `json:"created_at"`
	Active           bool       `json:"active"`
	IP               string     `json:"ip"`
	Port             int        `json:"port"`
	ID               []byte     `json:"id"`
	IsAdjusted       bool       `json:"is_adjusted"`
}
