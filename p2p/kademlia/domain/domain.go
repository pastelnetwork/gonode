package domain

import (
	"sort"
	"time"
)

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
	Key       string    `db:"key"`
	UpdatedAt time.Time `db:"updatedAt"`
	IP        string    `db:"ip"`
	Port      int       `db:"port"`
	ID        string    `db:"id"`
	Attempts  int       `db:"attempts"`
}

// DisabledKey is a disabled key
type DisabledKey struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"createdAt"`
}

// DisabledKeys is the list for disabled keys
type DisabledKeys []DisabledKey

// KeysWithTimestamp is the list for keys with timestamp
type KeysWithTimestamp []KeyWithTimestamp

// KeyWithTimestamp is a key with timestamp
type KeyWithTimestamp struct {
	Key       string    `db:"key" `
	CreatedAt time.Time `db:"createdAt"`
}

// FindFirstAfter finds the index of the first KeyWithTimestamp that has a CreatedAt timestamp after the specified time.
func (kwt KeysWithTimestamp) FindFirstAfter(t time.Time) int {
	index := sort.Search(len(kwt), func(i int) bool {
		return kwt[i].CreatedAt.After(t)
	})

	if index < len(kwt) && kwt[index].CreatedAt.After(t) {
		return index
	}

	// return -1 if no such element is found
	return -1
}

// DelKey is a key that could be deleted
type DelKey struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"createdAt"`
	Nodes     string    `json:"nodes"`
	Count     int       `json:"count"`
}

// DelKeys is the list for disabled keys
type DelKeys []DelKey
