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
	LastSeen         *time.Time `json:"last_seen"`
}

// ToRepKeys is the list for replication keys that need to be replicated
type ToRepKeys []ToRepKey

// ToRepKey is the struct for replication keys that need to be replicated
type ToRepKey struct {
	Key       []byte    `json:"key"`
	UpdatedAt time.Time `json:"updatedAt"`
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	ID        string    `json:"id"`
	Attempts  int       `json:"attempts"`
}

// DisabledKeys is a disabled key
type DisabledKey struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"createdAt"`
}

// DisabledKeys is the list for disabled keys
type DisabledKeys []DisabledKey
