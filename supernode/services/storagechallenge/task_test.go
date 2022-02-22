package storagechallenge

import (
	"context"
	"encoding/json"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/pastelnetwork/gonode/common/b85"
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
		SuperNodeTask           *common.SuperNodeTask
		StorageChallengeService *StorageChallengeService
		storage                 *common.StorageHandler
	}
	type args struct {
		ctx context.Context
	}
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"success": {
			args:    args{},
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

			pMock := pastelMock.NewMockClient(t)
			pMock.ListenOnRegTickets(pastel.RegTickets{
				ticket,
			}, nil).ListenOnGetBlockCount(1, nil).ListenOnGetBlockVerbose1(nil, nil)

			p2pClientMock := p2pMock.NewMockClient(t).ListenOnRetrieve(nil, nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			// rqClientMock.ListenOnConnect(tt.args.connectErr)

			clientMock := sctest.NewMockClient(t)
			clientMock.ListenOnConnect("", nil)

			fsMock := storageMock.NewMockFileStorage()
			// storage := files.NewStorage(fsMock)

			testConfig := NewConfig()

			task := StorageChallengeTask{
				SuperNodeTask:           tt.fields.SuperNodeTask,
				StorageChallengeService: NewService(testConfig, fsMock, pMock.Client, clientMock, p2pClientMock, rqClientMock, nil),
				storage:                 common.NewStorageHandler(p2pClientMock, rqClientMock, testConfig.RaptorQServiceAddress, testConfig.RqFilesDir),
			}

			if err := task.GenerateStorageChallenges(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("StorageChallengeTask.GenerateStorageChallenges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
