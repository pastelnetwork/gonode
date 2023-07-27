package collectionregister

import (
	"context"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testCreatorPastelID = "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW"
)

func TestTaskRun(t *testing.T) {
	t.Parallel()

	type fields struct {
		Request *common.CollectionRegistrationRequest
	}

	type args struct {
		taskID         string
		ctx            context.Context
		networkFee     float64
		masterNodes    pastel.MasterNodes
		getTopMNsReply *pb.GetTopMNsReply
		primarySessID  string
		pastelIDS      []string
		signature      []byte
		returnErr      error
	}

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr error
	}{
		"success": {
			fields: fields{
				Request: &common.CollectionRegistrationRequest{
					BurnTXID:              "txid",
					AppPastelID:           testCreatorPastelID,
					AppPastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				taskID:     "1",
				ctx:        log.ContextWithServer(context.Background(), "test-ip"),
				networkFee: 0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
				},
				getTopMNsReply: &pb.GetTopMNsReply{MnTopList: []string{"127.0.0.1:4444", "127.0.0.1:4446", "127.0.0.1:4447"}},
				primarySessID:  "sesid1",
				pastelIDS:      []string{"2", "3"},
				signature:      []byte("sign"),
				returnErr:      nil,
			},
		},

		"failure": {
			wantErr: errors.New("test"),
			fields: fields{
				Request: &common.CollectionRegistrationRequest{
					BurnTXID:              "txid",
					AppPastelID:           testCreatorPastelID,
					AppPastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3"},
				signature:     []byte("sign"),
				returnErr:     errors.New("test"),
			},
		},
	}

	for name, tc := range tests {
		testCase := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
				ListenOnDone().
				//ListenOnSendTicketForSignature(100, nil).
				ListenOnClose(nil).ListenOnSendActionAct(nil)
			nodeClient.RegisterCollectionInterface.
				On("SendTicketForSignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("", nil).Times(1)
			nodeClient.RegisterCollectionInterface.
				On("SendTicketForSignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("100", nil).Times(1)
			nodeClient.RegisterCollectionInterface.
				On("SendTicketForSignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("", nil).Times(1)

			nodeClient.ConnectionInterface.On("RegisterCollection").Return(nodeClient.RegisterCollectionInterface)
			nodeClient.RegisterCollectionInterface.On("GetTopMNs", mock.Anything, mock.Anything).Return(testCase.args.getTopMNsReply, nil)
			nodeClient.RegisterCollectionInterface.On("MeshNodes", mock.Anything, mock.Anything).Return(nil)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
				ListenOnSignCollectionTicket([]byte(testCase.args.signature), testCase.args.returnErr).
				ListenOnGetBlockCount(100, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Hash: "abc123", Height: 100}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{PqKey: ""}}, nil).
				ListenOnSendFromAddress("pre-burnt-txid", nil).
				ListenOnGetRawTransactionVerbose1(&pastel.GetRawTransactionVerbose1Result{Confirmations: 12}, nil).
				ListenOnVerifyCollectionTicket(true, nil).ListenOnGetBalance(10, nil).
				ListenOnActivateCollectionTicket("txid", nil)

			service := NewService(NewConfig(), pastelClientMock, nodeClient)
			service.config.WaitTxnValidInterval = 1

			go service.Run(testCase.args.ctx)

			Request := testCase.fields.Request
			task := NewCollectionRegistrationTask(service, Request)

			//create context with timeout to automatically end process after 5 sec
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 5*time.Second)
			defer cancel()

			err := task.Run(ctx)
			if testCase.wantErr != nil {
				assert.True(t, task.Status().IsFailure())
			} else {
				task.Status().Is(common.StatusTaskCompleted)
				assert.Nil(t, err)
			}
		})
	}
}
