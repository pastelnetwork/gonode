package storagechallenge

import (
	"bytes"
	"encoding/hex"
	"errors"
	"image"
	"image/png"
	"testing"

	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func Test_service_VerifyStorageChallenge(t *testing.T) {
	var (
		mockChallengeID       = "ch1"
		mockChallengingNodeID = "mn1"
		mockChallengeFileHash = "key1"

		mockChallengeStartPoint uint64 = 0
		mockChallengeEndPoint   uint64 = 20
		mockExpiredBlocks       int32  = 2
		mockBlkHeight           int32  = 1

		mockError = errors.New("mock error")
	)

	type mockFields struct {
		pclient    func(t *testing.T) (pclient pastel.Client, assertion func())
		repository func(t *testing.T) (repo Repository, assertion func())
	}

	testImage := image.NewRGBA(image.Rect(0, 0, 60, 90))
	var buff bytes.Buffer

	png.Encode(&buff, testImage)

	var computeHashOfFileSlice = func(fileData []byte, challengeSliceStartIndex, challengeSliceEndIndex uint64) string {
		challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
		algorithm := sha3.New256()
		algorithm.Write(challengeDataSlice)
		return hex.EncodeToString(algorithm.Sum(nil))
	}

	tests := map[string]struct {
		fields                   mockFields
		inComingChallengeMessage *ChallengeMessage
		assertion                assert.ErrorAssertionFunc
	}{
		"verify-success": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything)
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					testImage := image.NewRGBA(image.Rect(0, 0, 60, 90))
					var buff bytes.Buffer

					require.NoError(t, png.Encode(&buff, testImage), "could not create test image")

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "succeeded", mockChallengeID, mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:              storageChallengeResponseMessage,
				ChallengeStatus:          statusResponded,
				ChallengingMasternodeID:  mockChallengingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
				ChallengeResponseHash:    computeHashOfFileSlice(buff.Bytes(), mockChallengeStartPoint, mockChallengeEndPoint),
				ChallengeID:              mockChallengeID,
			},
			assertion: assert.NoError,
		},
		"verify-failed": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything)
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "failed", mockChallengeID, mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:              storageChallengeResponseMessage,
				ChallengeStatus:          statusResponded,
				ChallengingMasternodeID:  mockChallengingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
				ChallengeID:              mockChallengeID,
			},
			assertion: assert.NoError,
		},
		"verify-timeout": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight+3, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything)
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "timeout", mockChallengeID, mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:              storageChallengeResponseMessage,
				ChallengeStatus:          statusResponded,
				ChallengingMasternodeID:  mockChallengingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
				ChallengeResponseHash:    computeHashOfFileSlice(buff.Bytes(), mockChallengeStartPoint, mockChallengeEndPoint),
				ChallengeID:              mockChallengeID,
			},
			assertion: assert.NoError,
		},
		"incorrect-status-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {
					return nil, func() {}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				ChallengeStatus: "mock-incorrect-status",
			},
			assertion: assert.Error,
		},
		"incorrect-message-type-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {
					return nil, func() {}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:     "mock-incorrect-msg-type",
				ChallengeStatus: statusResponded,
			},
			assertion: assert.Error,
		},
		"get-symbol-file-by-key-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(nil, mockError)

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:         storageChallengeResponseMessage,
				ChallengeStatus:     statusResponded,
				FileHashToChallenge: mockChallengeFileHash,
			},
			assertion: assert.Error,
		},
		"get-block-count-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(0, mockError)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything)
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil)

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:         storageChallengeResponseMessage,
				ChallengeStatus:     statusResponded,
				FileHashToChallenge: mockChallengeFileHash,
			},
			assertion: assert.Error,
		},
	}
	for name, tt := range tests {

		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pclient, pclientAssertion := tt.fields.pclient(t)
			repository, repositoryAssertion := tt.fields.repository(t)

			s := &service{
				pclient:                       pclient,
				repository:                    repository,
				storageChallengeExpiredBlocks: mockExpiredBlocks,
			}

			tt.assertion(t, s.VerifyStorageChallenge(context.Background(), tt.inComingChallengeMessage))

			pclientAssertion()
			repositoryAssertion()
		})
	}
}
