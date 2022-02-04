package storagechallenge

import (
	"bytes"
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
	"github.com/stretchr/testify/require"
)

func Test_service_GenerateStorageChallenges(t *testing.T) {
	var (
		mockBlkHeight                 int32 = 1
		mockNumberOfChallengeReplicas       = 3

		mockMerkleRoot                = "1"
		mockChallengingNodeID         = "mn1"
		mockRespondingNodeID          = "mn2"
		mockCurrentNodeID             = ""
		mockChallengeRQSymbolFileHash = "key1"

		mockMasternodeList     = []string{mockChallengingNodeID, mockRespondingNodeID, mockCurrentNodeID}
		mockRQSymbolFileHashes = []string{mockChallengeRQSymbolFileHash}

		mockError = errors.New("mock error")
	)

	type mockFields struct {
		actor      func(t *testing.T) (messaging.Actor, func())
		pclient    func(t *testing.T) (pastel.Client, func())
		repository func(t *testing.T) (Repository, func())
	}

	tests := map[string]struct {
		fields    mockFields
		assertion assert.ErrorAssertionFunc
	}{
		"success": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil).
						ListenOnMasterNodesExtra(pastel.MasterNodes{
							{
								ExtAddress: "127.0.0.1:14444",
								ExtKey:     mockChallengingNodeID,
							},
							{
								ExtAddress: "127.0.0.1:14445",
								ExtKey:     mockRespondingNodeID,
							},
							{
								ExtAddress: "127.0.0.1:14446",
								ExtKey:     mockCurrentNodeID,
							},
						}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {

					m := test.NewMockActor(t).
						ListenOnSend(nil)

					return m, func() {

						m.AssertExpectations(t)

						m.AssertSendCall(1, mock.Anything, mock.Anything, mock.Anything)
					}
				},
				repository: func(t *testing.T) (Repository, func()) {

					testImage := image.NewRGBA(image.Rect(0, 0, 60, 90))
					var buff bytes.Buffer

					require.NoError(t, png.Encode(&buff, testImage), "could not create test image")

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(mockRQSymbolFileHashes, nil).
						ListenOnGetNClosestFileHashesToAGivenComparisonString(mockRQSymbolFileHashes).
						ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia([]string{mockRespondingNodeID}).
						ListenOnGetNClosestMasternodeIDsToComparisionString([]string{mockRespondingNodeID}).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(2, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything).
							AssertGetNClosestFileHashesToAGivenComparisonStringCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockRQSymbolFileHashes).
							AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(1, mock.Anything, 3, mock.AnythingOfType("string"), mockCurrentNodeID).
							AssertGetNClosestMasternodeIDsToComparisionStringCall(2, mock.Anything, 1, mock.AnythingOfType("string"), []string{mockRespondingNodeID}).
							AssertGetSymbolFileByKeyCall(2, mock.Anything, mockChallengeRQSymbolFileHash, false).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "sent", mock.AnythingOfType("string"), mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			assertion: assert.NoError,
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
					return nil, func() {}
				},
			},
			assertion: assert.Error,
		},
		"get-block-verbose-1-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(nil, mockError)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {
					return nil, func() {}
				},
			},
			assertion: assert.Error,
		},
		"get-mn-list-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(nil, mockError)

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything)
					}
				},
			},
			assertion: assert.Error,
		},
		"get-challenging-mn-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{})

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID)
					}
				},
			},
			assertion: assert.Error,
		},
		"list-rq-symbol-file-hash-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(nil, mockError)

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything)
					}
				},
			},
			assertion: assert.Error,
		},
		"get-symbol-file-bey-key-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(mockRQSymbolFileHashes, nil).
						ListenOnGetSymbolFileByKey(nil, mockError).
						ListenOnGetNClosestFileHashesToAGivenComparisonString(mockRQSymbolFileHashes).
						ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia([]string{mockRespondingNodeID}).
						ListenOnGetNClosestMasternodeIDsToComparisionString([]string{mockRespondingNodeID})

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(2, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything).
							AssertGetSymbolFileByKeyCall(2, mock.Anything, mockChallengeRQSymbolFileHash, false).
							AssertGetNClosestFileHashesToAGivenComparisonStringCall(1, mock.Anything, 1, mock.AnythingOfType("string"), []string{}).
							AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(1, mock.Anything, 3, mock.AnythingOfType("string"), mockCurrentNodeID).
							AssertGetNClosestMasternodeIDsToComparisionStringCall(2, mock.Anything, 1, mock.AnythingOfType("string"), []string{mockRespondingNodeID})
					}
				},
			},
			assertion: assert.NoError,
		},
		"empty-rq-symbol-file-value-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(mockRQSymbolFileHashes, nil).
						ListenOnGetSymbolFileByKey([]byte{}, nil).
						ListenOnGetNClosestFileHashesToAGivenComparisonString(mockRQSymbolFileHashes).
						ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia([]string{mockRespondingNodeID}).
						ListenOnGetNClosestMasternodeIDsToComparisionString([]string{mockRespondingNodeID})

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(2, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything).
							AssertGetSymbolFileByKeyCall(2, mock.Anything, mockChallengeRQSymbolFileHash, false).
							AssertGetNClosestFileHashesToAGivenComparisonStringCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockRQSymbolFileHashes).
							AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(1, mock.Anything, 3, mock.AnythingOfType("string"), mockCurrentNodeID).
							AssertGetNClosestMasternodeIDsToComparisionStringCall(2, mock.Anything, 1, mock.AnythingOfType("string"), []string{mockRespondingNodeID})
					}
				},
			},
			assertion: assert.NoError,
		},
		"mn-extra-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil).
						ListenOnMasterNodesExtra(nil, mockError)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					testImage := image.NewRGBA(image.Rect(0, 0, 60, 90))
					var buff bytes.Buffer

					require.NoError(t, png.Encode(&buff, testImage), "could not create test image")

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(mockRQSymbolFileHashes, nil).
						ListenOnGetNClosestFileHashesToAGivenComparisonString(mockRQSymbolFileHashes).
						ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia([]string{mockChallengingNodeID, mockRespondingNodeID}).
						ListenOnGetNClosestMasternodeIDsToComparisionString([]string{mockRespondingNodeID}).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(2, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything).
							AssertGetNClosestFileHashesToAGivenComparisonStringCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockRQSymbolFileHashes).
							AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(1, mock.Anything, 3, mock.AnythingOfType("string"), mockCurrentNodeID).
							AssertGetNClosestMasternodeIDsToComparisionStringCall(2, mock.Anything, 1, mock.AnythingOfType("string"), []string{mockRespondingNodeID}).
							AssertGetSymbolFileByKeyCall(2, mock.Anything, mockChallengeRQSymbolFileHash, false).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "sent", mock.AnythingOfType("string"), mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			assertion: assert.NoError,
		},
		"mn-not-found-in-list-err": {
			fields: mockFields{
				pclient: func(t *testing.T) (pastel.Client, func()) {

					p := test.NewMockPastelClient(t).
						ListenOnGetBlockCount(1, nil).
						ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Height: int64(mockBlkHeight), MerkleRoot: mockMerkleRoot}, nil).
						ListenOnMasterNodesExtra(pastel.MasterNodes{
							{
								ExtAddress: "127.0.0.1:14445",
								ExtKey:     mockRespondingNodeID,
							},
							{
								ExtAddress: "127.0.0.1:14446",
								ExtKey:     mockCurrentNodeID,
							},
						}, nil)

					return p.Client, func() {

						p.AssertExpectations(t)

						p.AssertGetBlockCountCall(1, mock.Anything).
							AssertGetBlockVerbose1Call(1, mock.Anything, mockBlkHeight).
							AssertMasterNodesExtraCall(1, mock.Anything)
					}
				},
				actor: func(t *testing.T) (messaging.Actor, func()) {
					return nil, func() {}
				},
				repository: func(t *testing.T) (Repository, func()) {

					testImage := image.NewRGBA(image.Rect(0, 0, 60, 90))
					var buff bytes.Buffer

					require.NoError(t, png.Encode(&buff, testImage), "could not create test image")

					r := test.NewMockRepository(t).
						ListenOnGetListOfMasternode(mockMasternodeList, nil).
						ListenOnGetNClosestMasternodeIDsToComparisionStringWithIgnore([]string{mockChallengingNodeID}).
						ListenOnListSymbolFileKeysFromNFTTicket(mockRQSymbolFileHashes, nil).
						ListenOnGetNClosestFileHashesToAGivenComparisonString(mockRQSymbolFileHashes).
						ListenOnGetNClosestMasternodesToAGivenFileUsingKademlia([]string{mockRespondingNodeID}).
						ListenOnGetNClosestMasternodeIDsToComparisionString([]string{mockRespondingNodeID}).
						ListenOnGetSymbolFileByKey(buff.Bytes(), nil).
						ListenOnSaveChallengMessageState()

					return r, func() {

						r.AssertExpectations(t)

						r.AssertGetListOfMasternodeCall(1, mock.Anything).
							AssertGetNClosestMasternodeIDsToComparisionStringWithIgnoreCall(2, mock.Anything, 1, mock.AnythingOfType("string"), mockMasternodeList, mockCurrentNodeID).
							AssertListSymbolFileKeysFromNFTTicketCall(1, mock.Anything).
							AssertGetNClosestFileHashesToAGivenComparisonStringCall(1, mock.Anything, 1, mock.AnythingOfType("string"), mockRQSymbolFileHashes).
							AssertGetNClosestMasternodesToAGivenFileUsingKademliaCall(1, mock.Anything, 3, mock.AnythingOfType("string"), mockCurrentNodeID).
							AssertGetNClosestMasternodeIDsToComparisionStringCall(2, mock.Anything, 1, mock.AnythingOfType("string"), []string{mockRespondingNodeID}).
							AssertGetSymbolFileByKeyCall(2, mock.Anything, mockChallengeRQSymbolFileHash, false).
							AssertSaveChallengMessageStateCall(1, mock.Anything, "sent", mock.AnythingOfType("string"), mockChallengingNodeID, mockBlkHeight)
					}
				},
			},
			assertion: assert.NoError,
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
				actor:                     actor,
				pclient:                   pclient,
				repository:                repository,
				nodeID:                    mockCurrentNodeID,
				numberOfChallengeReplicas: mockNumberOfChallengeReplicas,
			}

			tt.assertion(t, s.GenerateStorageChallenges(context.Background(), 1))

			actorAssertion()
			pclientAssertion()
			repositoryAssertion()

		})
	}
}
