package storagechallenge

import (
	"bytes"
	"encoding/hex"
	"errors"
	"image"
	"image/png"
	"testing"

	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/sha3"
)

func Test_service_ProcessStorageChallenge(t *testing.T) {
	var (
		mockChallengeID       = "ch1"
		mockChallengingNodeID = "mn1"
		mockRespondingNodeID  = "mn2"
		mockChallengeFileHash = "key1"

		mockBlkHeight           int32  = 1
		mockChallengeStartPoint uint64 = 0
		mockChallengeEndPoint   uint64 = 20

		mockError = errors.New("mock error")
	)

	type mockFields struct {
		actor      func(t *testing.T) (actor messaging.Actor, assertion func())
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
		"success": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil).
						ListenOnMasterNodesExtra(pastel.MasterNodes{
							{
								ExtAddress: "127.0.0.1:14444",
								ExtKey:     mockChallengingNodeID,
							},
							{
								ExtAddress: "127.0.0.1:14445",
								ExtKey:     mockRespondingNodeID,
							},
						}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {

					m := test.NewMockActor(t).
						ListenOnSend(nil)

					return m, func() {

						m.AssertExpectations(t)

						m.AssertSendCall(1, mock.Anything, mock.Anything, mock.MatchedBy(func(msg *sendVerifyStorageChallengeMsg) bool {
							return assert.Equal(t, []string{"127.0.0.1:14445"}, msg.VerifierMasternodesAddr) &&
								assert.NotNil(t, msg.ChallengeMessage) &&
								assert.Equal(
									t,
									computeHashOfFileSlice(
										buff.Bytes(),
										mockChallengeStartPoint,
										mockChallengeEndPoint),
									msg.ChallengeResponseHash,
								)
						}))
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetSymbolFileByKeyCall(1, mock.Anything, mockChallengeFileHash, true).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "respond", mockChallengeID, mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:              storageChallengeIssuanceMessage,
				ChallengeStatus:          statusPending,
				ChallengingMasternodeID:  mockChallengingNodeID,
				RespondingMasternodeID:   mockRespondingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
				ChallengeID:              mockChallengeID,
			},
			assertion: assert.NoError,
		},
		"incorrect-status-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					return nil, func() {}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
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
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {
					return nil, func() {}
				},
			},
			inComingChallengeMessage: &ChallengeMessage{
				MessageType:     "mock-incorrect-msg-type",
				ChallengeStatus: statusPending,
			},
			assertion: assert.Error,
		},
		"get-symbol-file-by-key-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {
					return nil, func() {}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
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
				MessageType:         storageChallengeIssuanceMessage,
				ChallengeStatus:     statusPending,
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
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
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
				MessageType:         storageChallengeIssuanceMessage,
				ChallengeStatus:     statusPending,
				FileHashToChallenge: mockChallengeFileHash,
			},
			assertion: assert.Error,
		},
		"mn-extra-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil).
						ListenOnMasterNodesExtra(nil, mockError)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
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
				MessageType:              storageChallengeIssuanceMessage,
				ChallengeStatus:          statusPending,
				ChallengingMasternodeID:  mockChallengingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
			},
			assertion: assert.Error,
		},
		"mn-not-found-in-list-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil).
						ListenOnMasterNodesExtra(pastel.MasterNodes{
							{
								ExtAddress: "127.0.0.1:14444",
								ExtKey:     mockChallengingNodeID,
							},
						}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
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
				MessageType:              storageChallengeIssuanceMessage,
				ChallengeStatus:          statusPending,
				ChallengingMasternodeID:  mockChallengingNodeID,
				RespondingMasternodeID:   mockRespondingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
			},
			assertion: assert.Error,
		},
		"send-verify-msg-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(mockBlkHeight, nil).
						ListenOnMasterNodesExtra(pastel.MasterNodes{
							{
								ExtAddress: "127.0.0.1:14444",
								ExtKey:     mockChallengingNodeID,
							},
							{
								ExtAddress: "127.0.0.1:14445",
								ExtKey:     mockRespondingNodeID,
							},
						}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {

					m := test.NewMockActor(t).
						ListenOnSend(mockError)

					return m, func() {

						m.AssertExpectations(t)

						m.AssertSendCall(1, mock.Anything, mock.Anything, mock.MatchedBy(func(msg *sendVerifyStorageChallengeMsg) bool {
							return assert.Equal(t, []string{"127.0.0.1:14445"}, msg.VerifierMasternodesAddr) &&
								assert.NotNil(t, msg.ChallengeMessage) &&
								assert.Equal(
									t,
									computeHashOfFileSlice(
										buff.Bytes(),
										mockChallengeStartPoint,
										mockChallengeEndPoint),
									msg.ChallengeResponseHash,
								)
						}))
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
				MessageType:              storageChallengeIssuanceMessage,
				ChallengeStatus:          statusPending,
				ChallengingMasternodeID:  mockChallengingNodeID,
				RespondingMasternodeID:   mockRespondingNodeID,
				BlockNumChallengeSent:    mockBlkHeight,
				FileHashToChallenge:      mockChallengeFileHash,
				ChallengeSliceStartIndex: mockChallengeStartPoint,
				ChallengeSliceEndIndex:   mockChallengeEndPoint,
			},
			assertion: assert.Error,
		},
	}
	for name, tt := range tests {

		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actor, actorAssertion := tt.fields.actor(t)
			pclient, pclientAssertion := tt.fields.pclient(t)
			repository, repositoryAssertion := tt.fields.repository(t)

			s := &service{
				actor:      actor,
				pclient:    pclient,
				repository: repository,
			}

			tt.assertion(t, s.ProcessStorageChallenge(context.Background(), tt.inComingChallengeMessage))

			actorAssertion()
			pclientAssertion()
			repositoryAssertion()
		})
	}
}
