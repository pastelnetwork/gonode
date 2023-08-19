package storagechallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/common/types"
)

// BroadcastStorageChallengeResult receives and process the storage challenge result
func (task *SCTask) BroadcastStorageChallengeResult(ctx context.Context, incomingEvaluationResult types.Message) (types.Message, error) {
	log.WithContext(ctx).WithField("method", "BroadcastStorageChallengeResult").
		WithField("challengeID", incomingEvaluationResult.ChallengeID).
		Debug("Start processing broadcasting message") // Incoming challenge message validation
	return types.Message{}, nil
}
