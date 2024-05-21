package storagechallenge

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"testing"

	fuzz "github.com/google/gofuzz"
	json "github.com/json-iterator/go"
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
			ticket.RegTicketData.NFTTicketData.AppTicket = base64.RawStdEncoding.EncodeToString(b)

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
				SCService:     NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, nil, nil),
				storage:       common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir, nil),
			}

			if err := task.GenerateStorageChallenges(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("SCTask.GenerateStorageChallenges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
