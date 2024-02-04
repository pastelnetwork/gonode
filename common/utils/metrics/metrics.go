package metrics

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/utils"
)

// Metrics is the struct for metrics
type Metrics struct {
	SCMetrics          []byte
	SHTriggerMetrics   SHTriggerMetrics
	SHExecutionMetrics SHExecutionMetrics
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

func (m *Metrics) Hash() string {
	data, _ := json.Marshal(m)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}
