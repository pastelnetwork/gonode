package storagechallenge

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/pastelnetwork/gonode/common/b85"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	sctest "github.com/pastelnetwork/gonode/supernode/node/test/storage_challenge"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

func TestTaskGenerateStorageChallenges(t *testing.T) {
	type fields struct {
		SuperNodeTask *common.SuperNodeTask
		SCService     *SCService
		storage       *common.StorageHandler
	}
	type args struct {
		ctx        context.Context
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
			ticket.RegTicketData.NFTTicketData.AppTicket = b85.Encode(b)

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
			}, nil).ListenOnGetBlockCount(1, nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(retrieveValue, nil).ListenOnNClosestNodes(closestNodes[0:2], nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().ListenOnProcessStorageChallengeFunc(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, rqClientMock, defaultChallengeStateLogging{}),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}

			if err := task.GenerateStorageChallenges(tt.args.ctx); (err != nil) != tt.wantErr {
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
		storage       *common.StorageHandler
	}
	type args struct {
		ctx                      context.Context
		incomingChallengeMessage *pb.StorageChallengeData
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
				incomingChallengeMessage: &pb.StorageChallengeData{
					MessageId:                    "0edd3e067b93ff597fd7e0d83ddd05161ba9b678b01c4c2fae7874b9274e6181",
					MessageType:                  pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE,
					ChallengeStatus:              pb.StorageChallengeData_Status_PENDING,
					BlockNumChallengeSent:        1,
					BlockNumChallengeRespondedTo: 0,
					BlockNumChallengeVerified:    0,
					MerklerootWhenChallengeSent:  "5072696d6172794944",
					ChallengingMasternodeId:      "5072696d6172794944",
					RespondingMasternodeId:       "B",
					ChallengeFile:                &pb.StorageChallengeDataChallengeFile{FileHashToChallenge: "亁zȲǘ", ChallengeSliceStartIndex: 0, ChallengeSliceEndIndex: 22},
					ChallengeSliceCorrectHash:    "",
					ChallengeResponseHash:        "",
					ChallengeId:                  "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
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
			ticket.RegTicketData.NFTTicketData.AppTicket = b85.Encode(b)

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
			}, nil).ListenOnGetBlockCount(1, nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(retrieveValue, nil).ListenOnNClosestNodes(closestNodes[0:2], nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().ListenOnVerifyStorageChallengeFunc(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, rqClientMock, defaultChallengeStateLogging{}),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}

			if resp, err := task.ProcessStorageChallenge(tt.args.ctx, tt.args.incomingChallengeMessage); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.ProcessStorageChallenge() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(resp)
			}
		})
	}
}

//{ MessageId: "be0771c56e1cb07748550f7b4650f7dba23c5af2f20f71a679eb217ddc88f7c4", MessageType: StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE (2), ChallengeStatus: StorageChallengeData_Status_RESPONDED (2), BlockNumChallengeSent: 1, BlockNumChallengeRespondedTo: 1, BlockNumChallengeVerified: 0, MerklerootWhenChallengeSent: "5072696d6172794944", ChallengingMasternodeId: "5072696d6172794944", RespondingMasternodeId: "B", ChallengeFile: *github.com/pastelnetwork/gonode/proto/supernode.StorageChallengeDataChallengeFile {state: (*"google.golang.org/protobuf/internal/impl.MessageState")(0xc0002b20a0), sizeCache: 0, unknownFields: []uint8 len: 0, cap: 0, nil, FileHashToChallenge: "亁zȲǘ", ChallengeSliceStartIndex: 0, ChallengeSliceEndIndex: 22}, ChallengeSliceCorrectHash: "", ChallengeResponseHash: "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977", ChallengeId: "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639"}

func TestVerifyStorageChallenge(t *testing.T) {
	type fields struct {
		SuperNodeTask *common.SuperNodeTask
		SCService     *SCService
		storage       *common.StorageHandler
	}
	type args struct {
		ctx                      context.Context
		incomingChallengeMessage *pb.StorageChallengeData
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
				incomingChallengeMessage: &pb.StorageChallengeData{
					MessageId:                    "be0771c56e1cb07748550f7b4650f7dba23c5af2f20f71a679eb217ddc88f7c4",
					MessageType:                  pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
					ChallengeStatus:              pb.StorageChallengeData_Status_RESPONDED,
					BlockNumChallengeSent:        1,
					BlockNumChallengeRespondedTo: 1,
					BlockNumChallengeVerified:    0,
					MerklerootWhenChallengeSent:  "5072696d6172794944",
					ChallengingMasternodeId:      "5072696d6172794944",
					RespondingMasternodeId:       "B",
					ChallengeFile: &pb.StorageChallengeDataChallengeFile{
						FileHashToChallenge:      "亁zȲǘ",
						ChallengeSliceStartIndex: 0,
						ChallengeSliceEndIndex:   22,
					},
					ChallengeSliceCorrectHash: "",
					ChallengeResponseHash:     "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
					ChallengeId:               "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
				},
				PastelID:          "C",
				MerkleRoot:        hex.EncodeToString([]byte("PrimaryID")),
				currentBlockCount: 1,
			},
			wantErr: false,
		},
		"tooManyBlocksPassed": {
			args: args{
				incomingChallengeMessage: &pb.StorageChallengeData{
					MessageId:                    "be0771c56e1cb07748550f7b4650f7dba23c5af2f20f71a679eb217ddc88f7c4",
					MessageType:                  pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
					ChallengeStatus:              pb.StorageChallengeData_Status_RESPONDED,
					BlockNumChallengeSent:        1,
					BlockNumChallengeRespondedTo: 1,
					BlockNumChallengeVerified:    0,
					MerklerootWhenChallengeSent:  "5072696d6172794944",
					ChallengingMasternodeId:      "5072696d6172794944",
					RespondingMasternodeId:       "B",
					ChallengeFile: &pb.StorageChallengeDataChallengeFile{
						FileHashToChallenge:      "亁zȲǘ",
						ChallengeSliceStartIndex: 0,
						ChallengeSliceEndIndex:   22,
					},
					ChallengeSliceCorrectHash: "",
					ChallengeResponseHash:     "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca977",
					ChallengeId:               "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
				},
				PastelID:          "C",
				MerkleRoot:        hex.EncodeToString([]byte("PrimaryID")),
				currentBlockCount: 10,
			},
			wantErr: false,
		},
		"badBlockHash": {
			args: args{
				incomingChallengeMessage: &pb.StorageChallengeData{
					MessageId:                    "be0771c56e1cb07748550f7b4650f7dba23c5af2f20f71a679eb217ddc88f7c4",
					MessageType:                  pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
					ChallengeStatus:              pb.StorageChallengeData_Status_RESPONDED,
					BlockNumChallengeSent:        1,
					BlockNumChallengeRespondedTo: 1,
					BlockNumChallengeVerified:    0,
					MerklerootWhenChallengeSent:  "5072696d6172794944",
					ChallengingMasternodeId:      "5072696d6172794944",
					RespondingMasternodeId:       "B",
					ChallengeFile: &pb.StorageChallengeDataChallengeFile{
						FileHashToChallenge:      "亁zȲǘ",
						ChallengeSliceStartIndex: 0,
						ChallengeSliceEndIndex:   22,
					},
					ChallengeSliceCorrectHash: "",
					ChallengeResponseHash:     "b7f16a4620c7e1b2bc49658f7466b321bede60a7a3ec94f00dff5b3af79ca978",
					ChallengeId:               "40fb87182c3d3643837d9e8590365f5f227088f828ad448dd18fb717231d9639",
				},
				PastelID:          "C",
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
			ticket.RegTicketData.NFTTicketData.AppTicket = b85.Encode(b)

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
			}, nil).ListenOnGetBlockCount(int32(tt.args.currentBlockCount), nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: tt.args.MerkleRoot,
			}, nil).ListenOnMasterNodesExtra(nodes, nil)

			closestNodes := []string{"A", "B", "C", "D"}
			retrieveValue := []byte("I retrieved this result")
			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(retrieveValue, nil).ListenOnNClosestNodes(closestNodes[0:2], nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil).ListenOnStorageChallengeInterface().ListenOnVerifyStorageChallengeFunc(nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()
			testConfig.PastelID = tt.args.PastelID

			task := SCTask{
				SuperNodeTask: tt.fields.SuperNodeTask,
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, rqClientMock, defaultChallengeStateLogging{}),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
				stateStorage:  defaultChallengeStateLogging{},
			}

			if resp, err := task.VerifyStorageChallenge(tt.args.ctx, tt.args.incomingChallengeMessage); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.ProcessStorageChallenge() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(resp)
			}
		})
	}
}
