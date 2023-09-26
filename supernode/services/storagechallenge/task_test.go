package storagechallenge

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/stretchr/testify/assert"

	fuzz "github.com/google/gofuzz"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	sctest "github.com/pastelnetwork/gonode/supernode/node/test/storage_challenge"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

func TestTaskGenerateStorageChallenges(t *testing.T) {
	type fields struct {
		SuperNodeTask *common.SuperNodeTask
		SCService     *SCService
	}
	type args struct {
		MerkleRoot string
		PastelID   string
	}
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"my_node_not_challenger": {
			args: args{
				MerkleRoot: "aaaaaaaaa",
				PastelID:   "A",
			},
			wantErr: false,
		},
		"my_node_is_challenger": {
			args: args{
				MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
				PastelID:   hex.EncodeToString([]byte("PrimaryID")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ticket := pastel.RegTicket{}
			f := fuzz.New()
			f.Fuzz(&ticket)
			ticket.Height = 2

			b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicketData.AppTicket = base64.StdEncoding.EncodeToString(b)

			b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicket = b

			nodes := pastel.MasterNodes{}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "C"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "D"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "5072696d6172794944"})

			pMock := pastelMock.NewMockClient(t)
			pMock.ListenOnRegTickets(pastel.RegTickets{
				ticket,
			}, nil).ListenOnActionTickets(nil, nil).ListenOnGetBlockCount(1, nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil).ListenOnSign([]byte{}, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(retrieveValue, nil).
				ListenOnNClosestNodes(closestNodes[0:2], nil).ListenOnDisableKey(nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().ListenOnProcessStorageChallengeFunc(nil).ConnectionInterface.On("Close").Return(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, defaultChallengeStateLogging{}, nil),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}

			if err := task.GenerateStorageChallenges(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.GenerateStorageChallenges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//{MessageId: "0edd3e067b93ff597fd7e0d83ddd05161ba9b678b01c4c2fae7874b9274e6181", MessageType: StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE (1), ChallengeStatus: StorageChallengeData_Status_PENDING (1), BlockNumChallengeSent: 1, BlockNumChallengeRespondedTo: 0, BlockNumChallengeVerified: 0, MerklerootWhenChallengeSent: "5072696d6172794944", ChallengingMasternodeId: "5072696d6172794944", RespondingMasternodeId: "B", ChallengeFile: *github.com/pastelnetwork/gonode/proto/supernode.StorageChallengeDataChallengeFile {state: (*"google.golang.org/protobuf/internal/impl.MessageState")(0xc0002ae230), sizeCache: 0, unknownFields: []uint8 len: 0, cap: 0, nil, FileHashToChallenge: "亁zȲǘ", ChallengeSliceStartIndex: 0, ChallengeSliceEndIndex: 22}, ChallengeSliceCorrectHash: "", ChallengeResponseHash: "", ChallengeId: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639"}

func TestTaskProcessStorageChallenge(t *testing.T) {
	type fields struct {
		SuperNodeTask *common.SuperNodeTask
		SCService     *SCService
	}
	type args struct {
		incomingChallengeMessage types.Message
		PastelID                 string
		MerkleRoot               string
	}
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"success": {
			args: args{
				incomingChallengeMessage: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9669",
					MessageType: types.ChallengeMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						RecipientID: "B",
					},
				},

				PastelID:   "B",
				MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
			},
			wantErr: false,
		},
		// "my_node_is_challenger": {
		// 	args: args{
		// 		MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
		// 		PastelID:   hex.EncodeToString([]byte("PrimaryID")),
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ticket := pastel.RegTicket{}
			f := fuzz.New()
			f.Fuzz(&ticket)
			ticket.Height = 2

			b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicketData.AppTicket = base64.StdEncoding.EncodeToString(b)

			b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicket = b

			nodes := pastel.MasterNodes{}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "C"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "D"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "5072696d6172794944"})

			pMock := pastelMock.NewMockClient(t)
			pMock.ListenOnRegTickets(pastel.RegTickets{
				ticket,
			}, nil).ListenOnActionTickets(nil, nil).ListenOnGetBlockCount(1, nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil).ListenOnVerify(true, nil).ListenOnSign([]byte{}, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(retrieveValue, nil).ListenOnNClosestNodes(closestNodes[0:2], nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().
				ListenOnVerifyStorageChallengeFunc(nil).ConnectionInterface.On("Close").Return(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, defaultChallengeStateLogging{}, nil),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}

			if resp, err := task.ProcessStorageChallenge(context.Background(), tt.args.incomingChallengeMessage); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.ProcessStorageChallenge() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(resp)
			}
		})
	}
}

//{ MessageId: "be0771c56e1cb07748550f7b4650f7dba23c5af2f20f71a679eb217ddc88f7c4", MessageType: StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE (2), ChallengeStatus: StorageChallengeData_Status_RESPONDED (2), BlockNumChallengeSent: 1, BlockNumChallengeRespondedTo: 1, BlockNumChallengeVerified: 0, MerklerootWhenChallengeSent: "5072696d6172794944", ChallengingMasternodeId: "5072696d6172794944", RespondingMasternodeId: "B", ChallengeFile: *github.com/pastelnetwork/gonode/proto/supernode.StorageChallengeDataChallengeFile {state: (*"google.golang.org/protobuf/internal/impl.MessageState")(0xc0002b20a0), sizeCache: 0, unknownFields: []uint8 len: 0, cap: 0, nil, FileHashToChallenge: "亁zȲǘ", ChallengeSliceStartIndex: 0, ChallengeSliceEndIndex: 22}, ChallengeSliceCorrectHash: "", ChallengeResponseHash: "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977", ChallengeId: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639"}

func createTempDirInPath(basePath string) (string, error) {
	dir, err := os.MkdirTemp(basePath, ".pastel")
	if err != nil {
		return "", err
	}
	return dir, nil
}

func TestVerifyStorageChallenge(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".pastel")

	// Make sure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		t.Fatalf("Failed to create base path: %v", err)
	}

	tempDir, err := createTempDirInPath(basePath)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up after test

	type fields struct {
		SuperNodeTask *common.SuperNodeTask
		SCService     *SCService
	}
	type args struct {
		incomingChallengeMessage types.Message
		verifyEvaluationResult1  types.Message
		verifyEvaluationResult2  types.Message
		PastelID                 string
		MerkleRoot               string
		currentBlockCount        int
	}
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"success": {
			args: args{
				incomingChallengeMessage: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.ResponseMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "B",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult1: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "D",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult2: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "E",
					SenderSignature: []byte{1, 3, 4},
				},
				PastelID:          "5072696d6172794944",
				MerkleRoot:        hex.EncodeToString([]byte("PrimaryID")),
				currentBlockCount: 1,
			},
			wantErr: false,
		},
		"tooManyBlocksPassed": {
			args: args{
				incomingChallengeMessage: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.ResponseMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      -4,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "B",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult1: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "D",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult2: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "E",
					SenderSignature: []byte{1, 3, 4},
				},
				PastelID:          "5072696d6172794944",
				MerkleRoot:        hex.EncodeToString([]byte("PrimaryID")),
				currentBlockCount: 10,
			},
			wantErr: false,
		},
		"badBlockHash": {
			args: args{
				incomingChallengeMessage: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.ResponseMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E"},
					},
					Sender:          "B",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult1: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "D",
					SenderSignature: []byte{1, 3, 4},
				},
				verifyEvaluationResult2: types.Message{
					ChallengeID: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType: types.AffirmationMessageType,
					Data: types.MessageData{
						ChallengerID: "5072696d6172794944",
						Challenge: types.ChallengeData{
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
							FileHash:   "亁zȲǘ",
							StartIndex: 0,
							EndIndex:   22,
						},
						Response: types.ResponseData{
							Hash:       "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Block:      1,
							Merkelroot: "5072696d6172794944",
							Timestamp:  time.Now().UTC(),
						},
						ObserverEvaluation: types.ObserverEvaluationData{
							IsChallengeTimestampOK:  true,
							IsProcessTimestampOK:    true,
							IsEvaluationTimestampOK: true,
							IsChallengerSignatureOK: true,
							IsRecipientSignatureOK:  true,
							IsEvaluationResultOK:    true,
							TrueHash:                "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
							Timestamp:               time.Now().UTC(),
						},
						RecipientID: "B",
						Observers:   []string{"D", "E", "F"},
					},
					Sender:          "E",
					SenderSignature: []byte{1, 3, 4},
				},
				PastelID:          "5072696d6172794944",
				MerkleRoot:        hex.EncodeToString([]byte("PrimaryID")),
				currentBlockCount: 1,
			},
			wantErr: false,
		},
		// "my_node_is_challenger": {
		// 	args: args{
		// 		MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
		// 		PastelID:   hex.EncodeToString([]byte("PrimaryID")),
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket := pastel.RegTicket{}
			f := fuzz.New()
			f.Fuzz(&ticket)
			ticket.Height = 2

			b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicketData.AppTicket = base64.StdEncoding.EncodeToString(b)

			b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
			if err != nil {
				t.Fatalf("faied to marshal, err: %s", err)
			}
			ticket.RegTicketData.NFTTicket = b

			nodes := pastel.MasterNodes{}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "B", ExtAddress: "B"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "C"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "D", ExtAddress: "D"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "E", ExtAddress: "E"})
			nodes = append(nodes, pastel.MasterNode{ExtKey: "5072696d6172794944"})

			pMock := pastelMock.NewMockClient(t)
			pMock.ListenOnRegTickets(pastel.RegTickets{
				ticket,
			}, nil).ListenOnActionTickets(nil, nil).ListenOnGetBlockCount(int32(tt.args.currentBlockCount), nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil).ListenOnVerify(true, nil).
				ListenOnSign([]byte{}, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).
				ListenOnRetrieve(retrieveValue, nil).
				ListenOnNClosestNodes(closestNodes[0:2], nil).ListenOnEnableKey(nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("D", nil).ListenOnStorageChallengeInterface().
				ListenOnVerifyEvaluationResultFunc(tt.args.verifyEvaluationResult1.Sender, tt.args.verifyEvaluationResult1, nil).ConnectionInterface.On("Close").
				Return(nil)

			clientMock.ListenOnConnect("E", nil).ListenOnStorageChallengeInterface().
				ListenOnVerifyEvaluationResultFunc(tt.args.verifyEvaluationResult2.Sender, tt.args.verifyEvaluationResult2, nil).ConnectionInterface.On("Close").
				Return(nil)

			clientMock.ListenOnConnect("B", nil).ListenOnStorageChallengeInterface().
				ListenOnVerifyEvaluationResultFunc(tt.args.verifyEvaluationResult2.Sender, tt.args.verifyEvaluationResult2, nil).ConnectionInterface.On("Close").
				Return(nil)

			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().
				ListenOnBroadcastStorageChallengeResultFunc(nil).ConnectionInterface.On("Close").
				Return(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, defaultChallengeStateLogging{}, nil),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}
			task.config.IsTestConfig = true

			err = task.StoreChallengeMessage(context.Background(),
				types.Message{
					ChallengeID:     "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
					MessageType:     types.ChallengeMessageType,
					Sender:          "5072696d6172794944",
					SenderSignature: []byte{1, 2, 3},
				})
			assert.NoError(t, err)

			if resp, err := task.VerifyStorageChallenge(context.Background(), tt.args.incomingChallengeMessage); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.VerifyStorageChallenge() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(resp)
			}
		})
	}

	defer func() {
		store, err := local.OpenHistoryDB()
		assert.NoError(t, err)

		store.CleanupStorageChallenges()
	}()
}
