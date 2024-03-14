package metrics

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/utils"
)

// Metrics is the struct for metrics
type Metrics struct {
	SCMetrics          SCMetrics
	SHTriggerMetrics   SHTriggerMetrics
	SHExecutionMetrics SHExecutionMetrics
}

type SCMetrics struct {
	TotalChallenges                      int
	TotalChallengesProcessed             int
	TotalChallengesEvaluatedByChallenger int
	TotalChallengesVerified              int
	SlowResponsesObservedByObservers     int
	InvalidSignaturesObservedByObservers int
	InvalidEvaluationObservedByObservers int
}

type HCMetrics struct {
	TotalChallenges                      int
	TotalChallengesProcessed             int
	TotalChallengesEvaluatedByChallenger int
	TotalChallengesVerified              int
	SlowResponsesObservedByObservers     int
	InvalidSignaturesObservedByObservers int
	InvalidEvaluationObservedByObservers int
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
	TotalChallengesIssued       int `json:"total_challenges_issued"`       // sum of challenge tickets in trigger table message_type = 1
	TotalChallengesAcknowledged int `json:"total_challenges_acknowledged"` // healer node accepted or rejected -- message_type = 2 && is_reconsrcution_req = false
	TotalChallengesAccepted     int `json:"total_challenges_accepted"`     // healer node accepted the challenge -- message_type = 2 && is_reconsrcution_req = true
	TotalChallengesRejected     int `json:"total_challenges_rejected"`     // message_type = 3 and array of self healing message shows < 3 is_verified = true selfhealig verification messages

	TotalChallengeEvaluationsVerified      int `json:"total_challenge_evaluations_verified"`        // message_type = 3 and array of self healing message shows >= 3 is_verified = true selfhealig verification messages
	TotalReconstructionsApproved           int `json:"total_reconstructions_approved"`              // message_type = 3 and array of self healing message shows >= 3 is_verified = true selfhealig verification messages
	TotalReconstructionsNotRquiredApproved int `json:"total_reconstructions_not_required_approved"` // message_type = 3 and array of self healing message shows >= 3 is_verified = true selfhealig verification messages

	TotalChallengeEvaluationsUnverified                  int `json:"total_challenges_unverified"`                                // message_type = 3 and array of self healing message shows < 3 is_verified = false selfhealig verification messages
	TotalReconstructionsNotApproved                      int `json:"total_reconstructions_not_approved"`                         // message_type = 3 and array of self healing message shows < 3 is_verified = false selfhealig verification messages
	TotalReconstructionsNotRequiredEvaluationNotApproved int `json:"total_reconstructions_not_required_evaluation_not_approved"` // message_type = 3 and array of self healing message shows < 3 is_verified = false selfhealig verification messages

	TotalFilesHealed       int `json:"total_files_healed"`        // message_type = 4
	TotalFileHealingFailed int `json:"total_file_healing_failed"` // message_type !=4

	TotalReconstructionRequiredHashMismatch int `json:"total_reconstruction_required_hash_mismatch"` // message_type = 3 and is_reconsrcution_req = true
}

// Hash returns the hash of the metrics
func (m *Metrics) Hash() string {
	data, _ := json.Marshal(m)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}
