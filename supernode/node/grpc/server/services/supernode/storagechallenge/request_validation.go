package storagechallenge

import (
	"fmt"

	dto "github.com/pastelnetwork/gonode/proto/supernode"
)

func validateGenerateStorageChallengeData(req *dto.GenerateStorageChallengesRequest) validationErrorStack {
	var key string
	var es validationErrorStack = make([]*validationError, 0)

	if req.GetChallengesPerMasternodePerBlock() <= 0 {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "ChallengesPerMasternodePerBlock")},
			reason: "invalid negative or zero value number",
		})
	}

	return es
}

func validateStorageChallengeData(req *dto.StorageChallengeData, key string) validationErrorStack {
	var es validationErrorStack = make([]*validationError, 0)

	if req.GetChallengeId() == "" {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "ChallengeId")},
			reason: reasonInvalidEmptyValue,
		})
	}

	if req.GetChallengingMasternodeId() == "" {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "ChallengingMasternodeId")},
			reason: reasonInvalidEmptyValue,
		})
	}

	if req.GetMessageId() == "" {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "MessageId")},
			reason: reasonInvalidEmptyValue,
		})
	}

	if req.GetMessageType() == 0 {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "MessageType")},
			reason: "invalid unknown value",
		})
	}

	es = append(es, validateStorageChallengeDataChallengeFile(req.GetChallengeFile(), fmt.Sprintf("%s.%s", key, "ChallengeFile"))...)
	return es
}

func validateStorageChallengeDataChallengeFile(req *dto.StorageChallengeDataChallengeFile, key string) validationErrorStack {
	var es validationErrorStack = make([]*validationError, 0)

	if req.GetFileHashToChallenge() == "" {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "FileHashToChallenge")},
			reason: reasonInvalidEmptyValue,
		})
	}

	if req.GetChallengeSliceStartIndex() < 0 {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "ChallengeSliceStartIndex")},
			reason: "start index could not be negative",
		})
	}

	if req.GetChallengeSliceStartIndex() >= req.GetChallengeSliceEndIndex() {
		es = append(es, &validationError{
			keys:   []string{joinKeyPart(key, "ChallengeSliceStartIndex"), joinKeyPart(key, "ChallengeSliceEndIndex")},
			reason: "ChallengeSliceStartIndex could not be larger than or equal to ChallengeSliceEndIndex",
		})
	}

	return es
}
