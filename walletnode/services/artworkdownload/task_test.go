package artworkdownload

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"
	taskMock "github.com/pastelnetwork/gonode/common/service/task/test"
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

	// t.Parallel()

	type args struct {
		ctx                context.Context
		returnErr          error
		ticketOwnershipErr error
		masterNodesTopErr  error
		connectErr         error
		downloadErr        error
		signErr            error
		file               []byte
		ttxid              string
		signature          []byte
		returnMn           pastel.MasterNodes
		taskID             string
	}

	type fields struct {
		Ticket *Ticket
	}

	testCases := []struct {
		args               args
		fields             fields
		assertion          assert.ErrorAssertionFunc
		numUpdateStatus    int
		numTicketOwnership int
		numMasterNodesTop  int
		numSign            int
		numConnect         int
		numDownload        int
		numDownLoadArtwork int
		numDone            int
	}{
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:                context.Background(),
				returnErr:          nil,
				ticketOwnershipErr: nil,
				signErr:            nil,
				masterNodesTopErr:  nil,
				connectErr:         nil,
				downloadErr:        nil,
				file:               []byte("testfile"),
				ttxid:              "ttxid",
				signature:          []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:          assert.NoError,
			numUpdateStatus:    3,
			numTicketOwnership: 1,
			numMasterNodesTop:  1,
			numSign:            1,
			numConnect:         3,
			numDownload:        3,
			numDownLoadArtwork: 3,
			numDone:            3,
		},
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:                context.Background(),
				returnErr:          nil,
				ticketOwnershipErr: fmt.Errorf("failed to get ticket ownership"),
				signErr:            nil,
				masterNodesTopErr:  nil,
				connectErr:         nil,
				downloadErr:        nil,
				file:               []byte("testfile"),
				ttxid:              "ttxid",
				signature:          []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:          assert.NoError,
			numUpdateStatus:    1,
			numTicketOwnership: 1,
			numMasterNodesTop:  0,
			numSign:            0,
			numConnect:         0,
			numDownload:        0,
			numDownLoadArtwork: 0,
			numDone:            0,
		},
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:                context.Background(),
				returnErr:          nil,
				ticketOwnershipErr: nil,
				signErr:            fmt.Errorf("failed to sign data"),
				masterNodesTopErr:  nil,
				connectErr:         nil,
				downloadErr:        nil,
				file:               []byte("testfile"),
				ttxid:              "ttxid",
				signature:          []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:          assert.NoError,
			numUpdateStatus:    1,
			numTicketOwnership: 1,
			numSign:            1,
			numMasterNodesTop:  0,
			numConnect:         0,
			numDownload:        0,
			numDownLoadArtwork: 0,
			numDone:            0,
		},
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:                context.Background(),
				returnErr:          nil,
				ticketOwnershipErr: nil,
				signErr:            nil,
				masterNodesTopErr:  fmt.Errorf("failed to get top masternodes"),
				connectErr:         nil,
				downloadErr:        nil,
				file:               []byte("testfile"),
				ttxid:              "ttxid",
				signature:          []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:          assert.NoError,
			numUpdateStatus:    1,
			numTicketOwnership: 1,
			numSign:            1,
			numMasterNodesTop:  1,
			numConnect:         0,
			numDownLoadArtwork: 0,
			numDownload:        0,
			numDone:            0,
		},
		{
			fields: fields{Ticket: &Ticket{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
			args: args{
				ctx:                context.Background(),
				returnErr:          nil,
				ticketOwnershipErr: nil,
				signErr:            nil,
				masterNodesTopErr:  nil,
				connectErr:         fmt.Errorf("failed to dial supernoded address"),
				downloadErr:        nil,
				file:               []byte("testfile"),
				ttxid:              "ttxid",
				signature:          []byte("sign"),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
				},
				taskID: "downloadtask",
			},
			assertion:          assert.NoError,
			numUpdateStatus:    2,
			numTicketOwnership: 1,
			numSign:            1,
			numMasterNodesTop:  1,
			numConnect:         3,
			numDownLoadArtwork: 0,
			numDownload:        0,
			numDone:            0,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			nodeClient := test.NewMockClient(t)
			if testCase.numConnect > 0 {
				nodeClient.ListenOnConnect(testCase.args.connectErr)
			}
			if testCase.numDownLoadArtwork > 0 {
				nodeClient.ListenOnDownloadArtwork()
			}
			if testCase.numDownload > 0 {
				nodeClient.ListenOnDownload(testCase.args.file, testCase.args.downloadErr)
			}
			if testCase.numDone > 0 {
				nodeClient.ListenOnDone()
			}

			pastelClient := pastelMock.NewMockClient(t)
			if testCase.numTicketOwnership > 0 {
				pastelClient.ListenOnTicketOwnership(testCase.args.ttxid, testCase.args.ticketOwnershipErr)
			}
			if testCase.numMasterNodesTop > 0 {
				pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.masterNodesTopErr)
			}
			if testCase.numSign > 0 {
				pastelClient.ListenOnSign(testCase.args.signature, testCase.args.signErr)
			}

			service := &Service{
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
				config:       NewConfig(),
			}

			taskClient := taskMock.NewMockTask(t)
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

			// pastelClient mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertTicketOwnershipCall(testCase.numTicketOwnership, mock.Anything,
				testCase.fields.Ticket.Txid,
				testCase.fields.Ticket.PastelID,
				testCase.fields.Ticket.PastelIDPassphrase,
			)
			pastelClient.AssertMasterNodesTopCall(testCase.numMasterNodesTop, mock.Anything)
			pastelClient.AssertSignCall(testCase.numSign, mock.Anything, mock.Anything,
				testCase.fields.Ticket.PastelID, testCase.fields.Ticket.PastelIDPassphrase)

			// nodeClient mock assertion
			nodeClient.Connection.AssertExpectations(t)
			nodeClient.Client.AssertExpectations(t)
			nodeClient.DownloadArtwork.AssertExpectations(t)
			nodeClient.AssertConnectCall(testCase.numConnect, mock.Anything, mock.Anything)
			nodeClient.AssertDownloadArtworkCall(testCase.numDownLoadArtwork)
			nodeClient.AssertDownloadCall(testCase.numDownload, mock.Anything, testCase.fields.Ticket.Txid,
				mock.Anything, string(testCase.args.signature), testCase.args.ttxid)
			nodeClient.AssertDoneCall(testCase.numDone)
		})
	}
}
