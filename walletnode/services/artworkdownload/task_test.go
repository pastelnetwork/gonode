package artworkdownload

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"
	staskMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNode(address, pastelID string) *node.Node {
	return node.NewNode(nil, address, pastelID)
}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *Service
		Ticket  *Ticket
	}

	service := &Service{}
	ticket := &Ticket{}

	testCases := []struct {
		args args
		want *Task
	}{
		{
			args: args{
				service: service,
				Ticket:  ticket,
			},
			want: &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
				Ticket:  ticket,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			task := NewTask(testCase.args.service, testCase.args.Ticket)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Ticket, task.Ticket)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}

func TestTaskPastelTopNodes(t *testing.T) {
	t.Parallel()

	type fields struct {
		Task   task.Task
		Ticket *Ticket
	}

	type args struct {
		ctx       context.Context
		returnMn  pastel.MasterNodes
		returnErr error
	}

	testCases := []struct {
		fields    fields
		args      args
		want      node.List
		assertion assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				Ticket: &Ticket{},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4444", "1"),
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			args: args{
				ctx:       context.Background(),
				returnMn:  nil,
				returnErr: fmt.Errorf("connection timeout"),
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClient := pastelMock.NewMockClient(t)
			pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClient.Client,
			}

			task := &Task{
				Task:    testCase.fields.Task,
				Service: service,
				Ticket:  testCase.fields.Ticket,
			}
			got, err := task.pastelTopNodes(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertMasterNodesTopCall(1, mock.Anything)
		})
	}

}

func TestRun(t *testing.T) {
	// t.Skip()

	t.Parallel()

	type args struct {
		ctx       context.Context
		returnErr error
		file      []byte
		ttxid     string
		signature []byte
		returnMn  pastel.MasterNodes
		taskID    string
	}

	type fields struct {
		Ticket *Ticket
	}

	testCases := []struct {
		args            args
		fields          fields
		assertion       assert.ErrorAssertionFunc
		numUpdateStatus int
	}{
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:       context.Background(),
				returnErr: nil,
				file:      []byte("testfile"),
				ttxid:     "ttxid",
				signature: []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:       assert.NoError,
			numUpdateStatus: 3,
		},
	}

	t.Run("group", func(t *testing.T) {

		for i, testCase := range testCases {
			testCase := testCase

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				nodeClient := test.NewMockClient(t)
				nodeClient.
					ListenOnConnect(testCase.args.returnErr).
					ListenOnDownloadArtwork().
					ListenOnDownload(testCase.args.file, testCase.args.returnErr).
					ListenOnDone()

				pastelClientMock := pastelMock.NewMockClient(t)
				pastelClientMock.
					ListenOnTicketOwnership(testCase.args.ttxid, testCase.args.returnErr).
					ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr).
					ListenOnSign(testCase.args.signature, testCase.args.returnErr)

				service := &Service{
					pastelClient: pastelClientMock.Client,
					nodeClient:   nodeClient.Client,
					config:       NewConfig(),
				}

				taskClient := staskMock.NewMockTask(t)
				taskClient.
					ListenOnID(testCase.args.taskID).
					ListenOnUpdateStatus().
					ListenOnSetStatusNotifyFunc()

				task := &Task{
					Task:    taskClient.Task,
					Service: service,
					Ticket:  testCase.fields.Ticket,
				}

				//create context with timeout to automatically end process after 1 sec
				ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
				defer cancel()
				err := task.Run(ctx)

				testCase.assertion(t, err)

				// taskClient mock assertion
				taskClient.AssertExpectations(t)
				taskClient.AssertIDCall(1)
				taskClient.AssertUpdateStatusCall(testCase.numUpdateStatus, mock.Anything)
				taskClient.AssertSetStatusNotifyFuncCall(1, mock.Anything)
			})
		}
	})
}
