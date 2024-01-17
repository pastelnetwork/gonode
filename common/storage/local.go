package storage

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/types"
)

// Metrics is the struct for metrics
type Metrics struct {
	SCMetrics          []byte
	SHTriggerMetrics   []byte
	SHExecutionMetrics []byte
}

// SHTriggerMetrics represents the self-healing trigger metrics
type SHTriggerMetrics []SHTriggerMetric

// SHTriggerMetric represents the self-healing execution metrics
type SHTriggerMetric struct {
	TriggerID              string `json:"trigger_id"`
	NodesOffline           int    `json:"nodes_offline"`
	ListOfNodes            string `json:"list_of_nodes"`
	TotalFilesIdentified   int    `json:"total_files_identified"`
	TotalTicketsIdentified int    `json:"total_tickets_identified"`
}

// SHExecutionMetrics represents the self-healing execution metrics
type SHExecutionMetrics struct {
	TotalChallengesIssued     int `json:"total_challenges_issued"`     // sum of challenge tickets in trigger table message_type = 1
	TotalChallengesRejected   int `json:"total_challenges_rejected"`   // healer node accepted or rejected -- message_type = 2 && is_reconsrcution_req = false
	TotalChallengesAccepted   int `json:"total_challenges_verified"`   // healer node accepted the challenge -- message_type = 2 && is_reconsrcution_req = true
	TotalChallengesFailed     int `json:"total_challenges_failed"`     // message_type = 3 and array of self healing message shows < 3 is_verified = true selfhealig verification messages
	TotalChallengesSuccessful int `json:"total_challenges_successful"` // message_type = 3 and array of self healing message shows >= 3 is_verified = true selfhealig verification messages
	TotalFilesHealed          int `json:"total_files_healed"`          // message_type = 4
	TotalFileHealingFailed    int `json:"total_file_healing_failed"`   // message_type !=4
}

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
	InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error
	InsertBroadcastMessage(challenge types.BroadcastLogMessage) error
	QueryStorageChallengeMessage(challengeID string, messageType int) (challenge types.StorageChallengeLogMessage, err error)
	CleanupStorageChallenges() (err error)
	CleanupSelfHealingChallenges() (err error)
	InsertSelfHealingGenerationMetrics(metrics types.SelfHealingGenerationMetric) error
	InsertSelfHealingExecutionMetrics(metrics types.SelfHealingExecutionMetric) error
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)
	UpsertPingHistory(pingInfo types.PingInfo) error
	GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error)
	GetAllPingInfos() (types.PingInfos, error)
	GetSelfHealingGenerationMetrics(timestamp time.Time) ([]types.SelfHealingGenerationMetric, error)
	GetSelfHealingExecutionMetrics(timestamp time.Time) ([]types.SelfHealingExecutionMetric, error)
	GetWatchlistPingInfo() ([]types.PingInfo, error)
	GetAllPingInfoForOnlineNodes() (types.PingInfos, error)
	UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error
	CloseHistoryDB(ctx context.Context)
	QueryMetrics(ctx context.Context, from time.Time, to *time.Time) (m Metrics, err error)
}
