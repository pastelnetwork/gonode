package storagechallenge

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/utils"
)

// SCSummaryStatsRes is the struct for metrics
type SCSummaryStatsRes struct {
	SCSummaryStats SCSummaryStats
}

type SCSummaryStats struct {
	TotalChallenges                      int
	TotalChallengesProcessed             int
	TotalChallengesVerifiedByChallenger  int
	TotalChallengesVerifiedByObservers   int
	SlowResponsesObservedByObservers     int
	InvalidSignaturesObservedByObservers int
	InvalidEvaluationObservedByObservers int
}

// Hash returns the hash of the SCSummaryStats
func (ss *SCSummaryStatsRes) Hash() string {
	data, _ := json.Marshal(ss)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}
