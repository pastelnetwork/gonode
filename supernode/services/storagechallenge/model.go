package storagechallenge

import dto "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"

var (
	storageChallengeIssuanceMessage     = dto.MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE.String()
	storageChallengeResponseMessage     = dto.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE.String()
	storageChallengeVerificationMessage = dto.MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE.String()

	statusPending                 = dto.Status_PENDING.String()
	statusResponded               = dto.Status_RESPONDED.String()
	statusSucceeded               = dto.Status_SUCCEEDED.String()
	statusFailedTimeout           = dto.Status_FAILED_TIMEOUT.String()
	statusFailedIncorrectResponse = dto.Status_FAILED_INCORRECT_RESPONSE.String()
)

// ChallengeMessage struct
type ChallengeMessage struct {
	MessageID                     string
	MessageType                   string
	ChallengeStatus               string
	TimestampChallengeSent        int64
	TimestampChallengeRespondedTo int64
	TimestampChallengeVerified    int64
	BlockHashWhenChallengeSent    string
	ChallengingMasternodeID       string
	RespondingMasternodeID        string
	FileHashToChallenge           string
	ChallengeSliceStartIndex      uint64
	ChallengeSliceEndIndex        uint64
	ChallengeSliceCorrectHash     string
	ChallengeResponseHash         string
	ChallengeID                   string
}

// Challenge struct
type Challenge struct {
	ChallengeID                    string
	ChallengeStatus                string
	TimestampChallengeSent         int64
	TimestampChallengeRespondedTo  int64
	TimestampChallengeVerified     int64
	BlockHashWhenChallengeSent     string
	ChallengeResponseTimeInSeconds float64
	ChallengingMasternodeID        string
	RespondingMasternodeID         string
	FileHashToChallenge            string
	ChallengeSliceStartIndex       uint64
	ChallengeSliceEndIndex         uint64
	ChallengeSliceCorrectHash      string
	ChallengeResponseHash          string
}

// type PastelBlock struct {
// 	BlockHash                       string
// 	BlockNumber                     uint
// 	TotalChallengesIssued           uint
// 	TotalChallengesRespondedTo      uint
// 	TotalChallengesCorrect          uint
// 	TotalChallengesIncorrect        uint
// 	TotalChallengeTimeout           uint
// 	ChallengeResponseSuccessRatePct float32 `gorm:"column:challenge_response_success_rate_pct"`
// }

// type Masternode struct {
// 	NodeID                          string
// 	MasternodeIPAddress             string
// 	TotalChallengesIssued           uint
// 	TotalChallengesRespondedTo      uint
// 	TotalChallengesCorrect          uint
// 	TotalChallengesIncorrect        uint
// 	TotalChallengeTimeout           uint
// 	ChallengeResponseSuccessRatePct float32 `gorm:"column:challenge_response_success_rate_pct"`
// }
