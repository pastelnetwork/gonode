package storagechallenge

import dto "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"

var (
	MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE     = dto.MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE.String()
	MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE     = dto.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE.String()
	MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE = dto.MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE.String()

	Status_PENDING                   = dto.Status_PENDING.String()
	Status_RESPONDED                 = dto.Status_RESPONDED.String()
	Status_SUCCEEDED                 = dto.Status_SUCCEEDED.String()
	Status_FAILED_TIMEOUT            = dto.Status_FAILED_TIMEOUT.String()
	Status_FAILED_INCORRECT_RESPONSE = dto.Status_FAILED_INCORRECT_RESPONSE.String()
)

type ChallengeAnalysisStatus int

const (
	ANALYSYS_STATUS_ISSUED ChallengeAnalysisStatus = iota
	ANALYSIS_STATUS_RESPONDED_TO
	ANALYSIS_STATUS_CORRECT
	ALALYSIS_STATUS_INCORRECT
	ALALYSIS_STATUS_TIMEOUT
)

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

type PastelBlock struct {
	BlockHash                       string
	BlockNumber                     uint
	TotalChallengesIssued           uint
	TotalChallengesRespondedTo      uint
	TotalChallengesCorrect          uint
	TotalChallengesIncorrect        uint
	TotalChallengeTimeout           uint
	ChallengeResponseSuccessRatePct float32 `gorm:"column:challenge_response_success_rate_pct"`
}

type Masternode struct {
	NodeID                          string
	MasternodeIPAddress             string
	TotalChallengesIssued           uint
	TotalChallengesRespondedTo      uint
	TotalChallengesCorrect          uint
	TotalChallengesIncorrect        uint
	TotalChallengeTimeout           uint
	ChallengeResponseSuccessRatePct float32 `gorm:"column:challenge_response_success_rate_pct"`
}

type SymbolFile struct {
	FileHash               string
	FileLengthInBytes      uint
	TotalChallengesForFile uint
	OriginalFilePath       string
}

type XORDistance struct {
	XORDistanceID  string `gorm:"column:xor_distance_id"`
	MasternodeID   string
	SymbolFileHash string
	XorDistance    uint64
	Masternode     *Masternode `gorm:"foreignKey:MasternodeID;references:NodeID"`
	SymbolFile     *SymbolFile `gorm:"foreignKey:SymbolFileHash;references:FileHash"`
}

func (XORDistance) TableName() string {
	return "xor_distances"
}
