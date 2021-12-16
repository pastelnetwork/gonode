package storagechallenge

import (
	dto "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

func mapChallengeMessage(dto *dto.StorageChallengeData) *storagechallenge.ChallengeMessage {
	return &storagechallenge.ChallengeMessage{
		MessageID:                    dto.GetMessageId(),
		MessageType:                  dto.GetMessageType().String(),
		ChallengeStatus:              dto.GetChallengeStatus().String(),
		BlockNumChallengeSent:        dto.GetBlockNumChallengeSent(),
		BlockNumChallengeRespondedTo: dto.GetBlockNumChallengeRespondedTo(),
		BlockNumChallengeVerified:    dto.GetBlockNumChallengeVerified(),
		MerklerootWhenChallengeSent:  dto.GetMerklerootWhenChallengeSent(),
		ChallengingMasternodeID:      dto.GetChallengingMasternodeId(),
		RespondingMasternodeID:       dto.GetRespondingMasternodeId(),
		FileHashToChallenge:          dto.GetChallengeFile().GetFileHashToChallenge(),
		ChallengeSliceStartIndex:     uint64(dto.GetChallengeFile().GetChallengeSliceStartIndex()),
		ChallengeSliceEndIndex:       uint64(dto.GetChallengeFile().GetChallengeSliceEndIndex()),
		ChallengeSliceCorrectHash:    dto.GetChallengeSliceCorrectHash(),
		ChallengeResponseHash:        dto.GetChallengeResponseHash(),
		ChallengeID:                  dto.GetChallengeId(),
	}
}
