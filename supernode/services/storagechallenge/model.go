package storagechallenge

import dto "github.com/pastelnetwork/gonode/proto/supernode"

var (
	storageChallengeIssuanceMessage     = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE.String()
	storageChallengeResponseMessage     = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE.String()
	storageChallengeVerificationMessage = dto.StorageChallengeData_MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE.String()

	statusPending                 = dto.StorageChallengeData_Status_PENDING.String()
	statusResponded               = dto.StorageChallengeData_Status_RESPONDED.String()
	statusSucceeded               = dto.StorageChallengeData_Status_SUCCEEDED.String()
	statusFailedTimeout           = dto.StorageChallengeData_Status_FAILED_TIMEOUT.String()
	statusFailedIncorrectResponse = dto.StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE.String()
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
	ChallengingMasternodeID      string
	RespondingMasternodeID       string
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
	ChallengingMasternodeID        string
	RespondingMasternodeID         string
	FileHashToChallenge            string
	ChallengeSliceStartIndex       uint64
	ChallengeSliceEndIndex         uint64
	ChallengeSliceCorrectHash      string
	ChallengeResponseHash          string
}
