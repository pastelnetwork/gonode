package storagechallenge

import (
	dto "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

func mapChallengeMessage(dto *dto.StorageChallengeData) *storagechallenge.ChallengeMessage {
	return &storagechallenge.ChallengeMessage{
		MessageID:                     dto.GetMessageId(),
		MessageType:                   dto.GetMessageType().String(),
		ChallengeStatus:               dto.GetChallengeStatus().String(),
		TimestampChallengeSent:        dto.GetTimestampChallengeSent(),
		TimestampChallengeRespondedTo: dto.GetTimestampChallengeRespondedTo(),
		TimestampChallengeVerified:    dto.GetTimestampChallengeVerified(),
		BlockHashWhenChallengeSent:    dto.GetBlockHashWhenChallengeSent(),
		ChallengingMasternodeID:       dto.GetChallengingMasternodeId(),
		RespondingMasternodeID:        dto.GetRespondingMasternodeId(),
		FileHashToChallenge:           dto.GetChallengeFile().GetFileHashToChallenge(),
		ChallengeSliceStartIndex:      uint64(dto.GetChallengeFile().GetChallengeSliceStartIndex()),
		ChallengeSliceEndIndex:        uint64(dto.GetChallengeFile().GetChallengeSliceEndIndex()),
		ChallengeSliceCorrectHash:     dto.GetChallengeSliceCorrectHash(),
		ChallengeResponseHash:         dto.GetChallengeResponseHash(),
		ChallengeID:                   dto.GetChallengeId(),
	}
}
