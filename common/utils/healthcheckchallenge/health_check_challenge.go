package storagechallenge

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/utils"
)

// HCSummaryStatsRes is the struct for metrics
type HCSummaryStatsRes struct {
	HCSummaryStats HCSummaryStats
}

type HCSummaryStats struct {
	TotalChallenges                      int
	TotalChallengesProcessed             int
	TotalChallengesEvaluatedByChallenger int
	TotalChallengesVerified              int
	SlowResponsesObservedByObservers     int
	InvalidSignaturesObservedByObservers int
	InvalidEvaluationObservedByObservers int
}

// Hash returns the hash of the SCSummaryStats
func (ss *HCSummaryStatsRes) Hash() string {
	data, _ := json.Marshal(ss)
	hash, _ := utils.Sha3256hash(data)

	return string(hash)
}
