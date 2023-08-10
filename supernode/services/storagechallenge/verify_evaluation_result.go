package storagechallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// VerifyEvaluationResult process the evaluation report from the challenger
func (task *SCTask) VerifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message) (types.Message, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingEvaluationResult.ChallengeID)
	logger.Debug("Start verifying storage challenge") // Incoming challenge message validation

	return types.Message{}, nil
}
