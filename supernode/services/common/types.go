package common

import dto "github.com/pastelnetwork/gonode/proto/supernode"

var (
	StorageChallengeIssuanceMessage     = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE.String()
	StorageChallengeResponseMessage     = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE.String()
	StorageChallengeVerificationMessage = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE.String()

	StatusPending                 = dto.StorageChallengeData_Status_PENDING.String()
	StatusResponded               = dto.StorageChallengeData_Status_RESPONDED.String()
	StatusSucceeded               = dto.StorageChallengeData_Status_SUCCEEDED.String()
	StatusFailedTimeout           = dto.StorageChallengeData_Status_FAILED_TIMEOUT.String()
	StatusFailedIncorrectResponse = dto.StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE.String()
)

// ChallengeMessage struct
type ChallengeMessage struct {
	MessageID                    string
	MessageType                  string
	ChallengeStatus              string
	BlockNumChallengeSent        int32
	BlockNumChallengeRespondedTo int32
	BlockNumChallengeVerified    int32
	MerklerootWhenChallengeSent  string
	ChallengingSupernodeID       string
	RespondingSupernodeID        string
	FileHashToChallenge          string
	ChallengeSliceStartIndex     uint64
	ChallengeSliceEndIndex       uint64
	ChallengeSliceCorrectHash    string
	ChallengeResponseHash        string
	ChallengeID                  string
}

// Challenge struct
type Challenge struct {
	ChallengeID                    string
	ChallengeStatus                string
	BlockNumChallengeSent          int32
	BlockNumChallengeRespondedTo   int32
	BlockNumChallengeVerified      int32
	MerklerootWhenChallengeSent    string
	ChallengeResponseTimeInSeconds float64
	ChallengingSupernodeID         string
	RespondingSupernodeID          string
	FileHashToChallenge            string
	ChallengeSliceStartIndex       uint64
	ChallengeSliceEndIndex         uint64
	ChallengeSliceCorrectHash      string
	ChallengeResponseHash          string
}
