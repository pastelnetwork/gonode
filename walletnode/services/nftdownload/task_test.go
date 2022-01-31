package nftdownload

//
//import (
//	"context"
//	"fmt"
//	"github.com/pastelnetwork/gonode/walletnode/services/common"
//	"testing"
//	"time"
//
//	"github.com/pastelnetwork/gonode/common/service/task"
//	taskMock "github.com/pastelnetwork/gonode/common/service/task/test"
//	"github.com/pastelnetwork/gonode/pastel"
//	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
//	test "github.com/pastelnetwork/gonode/walletnode/node/test/nft_download"
//	"github.com/pastelnetwork/gonode/walletnode/services/nftdownload/node"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//func newTestNode(address, pastelID string) *node.NftDownloadNodeClient {
//	return node.NewNode(nil, address, pastelID)
//}
//
//func TestNewTask(t *testing.T) {
//	// t.Parallel()
//
//	type args struct {
//		service *NftDownloadingService
//		Ticket  *NftDownloadingRequest
//	}
//
//	service := &NftDownloadingService{}
//	ticket := &NftDownloadingRequest{}
//
//	testCases := []struct {
//		args args
//		want *NftDownloadingTask
//	}{
//		{
//			args: args{
//				service: service,
//				Ticket:  ticket,
//			},
//			want: &NftDownloadingTask{
//				WalletNodeTask:     common.NewWalletNodeTask(logPrefix),
//				NftDownloadingService: service,
//				Request:            ticket,
//			},
//		},
//	}
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			// t.Parallel()
//
//			task := NewNftDownloadTask(testCase.args.service, testCase.args.Ticket)
//			assert.Equal(t, testCase.want.NftDownloadingService, task.NftDownloadingService)
//			assert.Equal(t, testCase.want.Request, task.Request)
//			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
//		})
//	}
//}
//
//func TestTaskPastelTopNodes(t *testing.T) {
//	// t.Parallel()
//
//	type fields struct {
//		Task   task.Task
//		Ticket *NftDownloadingRequest
//	}
//
//	type args struct {
//		ctx       context.Context
//		returnMn  pastel.MasterNodes
//		returnErr error
//	}
//
//	testCases := []struct {
//		fields    fields
//		args      args
//		want      node.List
//		assertion assert.ErrorAssertionFunc
//	}{
//		{
//			fields: fields{
//				Ticket: &NftDownloadingRequest{},
//			},
//			args: args{
//				ctx: context.Background(),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//				},
//				returnErr: nil,
//			},
//			want: node.List{
//				newTestNode("127.0.0.1:4444", "1"),
//				newTestNode("127.0.0.1:4445", "2"),
//			},
//			assertion: assert.NoError,
//		}, {
//			args: args{
//				ctx:       context.Background(),
//				returnMn:  nil,
//				returnErr: fmt.Errorf("connection timeout"),
//			},
//			want:      nil,
//			assertion: assert.Error,
//		},
//	}
//
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			// t.Parallel()
//
//			//create new mock service
//			pastelClient := pastelMock.NewMockClient(t)
//			pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr)
//			service := &NftDownloadingService{
//				pastelClient: pastelClient.Client,
//			}
//
//			task := &NftDownloadingTask{
//				WalletNodeTask: &common.WalletNodeTask{
//					Task:      testCase.fields.Task,
//					LogPrefix: logPrefix,
//				},
//				NftDownloadingService: service,
//				Request:            testCase.fields.Ticket,
//			}
//			got, err := task.pastelTopNodes(testCase.args.ctx)
//			testCase.assertion(t, err)
//			assert.Equal(t, testCase.want, got)
//
//			//mock assertion
//			pastelClient.AssertExpectations(t)
//			pastelClient.AssertMasterNodesTopCall(1, mock.Anything)
//		})
//	}
//
//}
//
//func TestTaskRun(t *testing.T) {
//	t.Parallel()
//
//	type args struct {
//		ctx                context.Context
//		returnErr          error
//		ticketOwnershipErr error
//		masterNodesTopErr  error
//		connectErr         error
//		downloadErr        error
//		closeErr           error
//		signErr            error
//		file               []byte
//		ttxid              string
//		signature          []byte
//		returnMn           pastel.MasterNodes
//		taskID             string
//	}
//
//	type fields struct {
//		Ticket *NftDownloadingRequest
//	}
//
//	testCases := []struct {
//		args               args
//		fields             fields
//		assertion          assert.ErrorAssertionFunc
//		numUpdateStatus    int
//		numTicketOwnership int
//		numMasterNodesTop  int
//		numSign            int
//		numConnect         int
//		numDownload        int
//		numDownLoadNft int
//		numDone            int
//		numClose           int
//	}{
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: nil,
//				signErr:            nil,
//				masterNodesTopErr:  nil,
//				connectErr:         nil,
//				downloadErr:        nil,
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    2,
//			numTicketOwnership: 1,
//			numMasterNodesTop:  1,
//			numSign:            1,
//			numConnect:         3,
//			numDownload:        3,
//			numDownLoadNft: 3,
//			numDone:            0,
//			numClose:           3,
//		},
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: fmt.Errorf("failed to get ticket ownership"),
//				signErr:            nil,
//				masterNodesTopErr:  nil,
//				connectErr:         nil,
//				downloadErr:        nil,
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    1,
//			numTicketOwnership: 1,
//			numMasterNodesTop:  0,
//			numSign:            0,
//			numConnect:         0,
//			numDownload:        0,
//			numDownLoadNft: 0,
//			numDone:            0,
//			numClose:           0,
//		},
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: nil,
//				signErr:            fmt.Errorf("failed to sign data"),
//				masterNodesTopErr:  nil,
//				connectErr:         nil,
//				downloadErr:        nil,
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    1,
//			numTicketOwnership: 1,
//			numSign:            1,
//			numMasterNodesTop:  0,
//			numConnect:         0,
//			numDownload:        0,
//			numDownLoadNft: 0,
//			numDone:            0,
//			numClose:           0,
//		},
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: nil,
//				signErr:            nil,
//				masterNodesTopErr:  fmt.Errorf("failed to get top masternodes"),
//				connectErr:         nil,
//				downloadErr:        nil,
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    1,
//			numTicketOwnership: 1,
//			numSign:            1,
//			numMasterNodesTop:  1,
//			numConnect:         0,
//			numDownLoadNft: 0,
//			numDownload:        0,
//			numDone:            0,
//			numClose:           0,
//		},
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: nil,
//				signErr:            nil,
//				masterNodesTopErr:  nil,
//				connectErr:         fmt.Errorf("failed to dial supernoded address"),
//				downloadErr:        nil,
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    2,
//			numTicketOwnership: 1,
//			numSign:            1,
//			numMasterNodesTop:  1,
//			numConnect:         3,
//			numDownLoadNft: 0,
//			numDownload:        0,
//			numDone:            0,
//			numClose:           0,
//		},
//		{
//			fields: fields{Ticket: &NftDownloadingRequest{Txid: "txid", PastelID: "pastelid", PastelIDPassphrase: "passphrase"}},
//			args: args{
//				ctx:                context.Background(),
//				returnErr:          nil,
//				ticketOwnershipErr: nil,
//				signErr:            nil,
//				masterNodesTopErr:  nil,
//				connectErr:         nil,
//				downloadErr:        fmt.Errorf("failed to download"),
//				closeErr:           nil,
//				file:               []byte("testfile"),
//				ttxid:              "ttxid",
//				signature:          []byte("sign"),
//				returnMn: pastel.MasterNodes{
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
//					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
//				},
//				taskID: "downloadtask",
//			},
//			assertion:          assert.NoError,
//			numUpdateStatus:    2,
//			numTicketOwnership: 1,
//			numSign:            1,
//			numMasterNodesTop:  1,
//			numConnect:         3,
//			numDownLoadNft: 3,
//			numDownload:        3,
//			numDone:            0,
//			numClose:           0,
//		},
//	}
//
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			// t.Parallel()
//
//			nodeClient := test.NewMockClient(t)
//			if testCase.numConnect > 0 {
//				nodeClient.ListenOnConnect("", testCase.args.connectErr)
//			}
//			if testCase.numDownLoadNft > 0 {
//				nodeClient.ListenOnDownloadNft()
//			}
//			if testCase.numDownload > 0 {
//				nodeClient.ListenOnDownload(testCase.args.file, testCase.args.downloadErr)
//			}
//			if testCase.numDone > 0 {
//				nodeClient.ListenOnDone()
//			}
//			if testCase.numClose > 0 {
//				nodeClient.ListenOnClose(testCase.args.closeErr)
//			}
//
//			pastelClient := pastelMock.NewMockClient(t)
//			if testCase.numTicketOwnership > 0 {
//				pastelClient.ListenOnTicketOwnership(testCase.args.ttxid, testCase.args.ticketOwnershipErr)
//			}
//			if testCase.numMasterNodesTop > 0 {
//				pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.masterNodesTopErr)
//			}
//			if testCase.numSign > 0 {
//				pastelClient.ListenOnSign(testCase.args.signature, testCase.args.signErr)
//			}
//
//			service := &NftDownloadingService{
//				pastelClient: pastelClient.Client,
//				nodeClient:   nodeClient.Client,
//				config:       NewConfig(),
//			}
//
//			taskClient := taskMock.NewMockTask(t)
//			taskClient.
//				ListenOnID(testCase.args.taskID).
//				ListenOnUpdateStatus().
//				ListenOnSetStatusNotifyFunc()
//
//			task := &NftDownloadingTask{
//				WalletNodeTask: &common.WalletNodeTask{
//					Task:      taskClient.Task,
//					LogPrefix: logPrefix,
//				},
//				NftDownloadingService: service,
//				Request:            testCase.fields.Ticket,
//			}
//
//			//create context with timeout to automatically end process after 1 sec
//			ctx, cancel := context.WithTimeout(testCase.args.ctx, 2*time.Second)
//			defer cancel()
//			err := task.Run(ctx)
//
//			testCase.assertion(t, err)
//
//			// taskClient mock assertion
//			taskClient.AssertExpectations(t)
//			taskClient.AssertIDCall(1)
//			taskClient.AssertUpdateStatusCall(testCase.numUpdateStatus, mock.Anything)
//			taskClient.AssertSetStatusNotifyFuncCall(1, mock.Anything)
//
//			// pastelClient mock assertion
//			pastelClient.AssertExpectations(t)
//			pastelClient.AssertTicketOwnershipCall(testCase.numTicketOwnership, mock.Anything,
//				testCase.fields.Ticket.Txid,
//				testCase.fields.Ticket.PastelID,
//				testCase.fields.Ticket.PastelIDPassphrase,
//			)
//			pastelClient.AssertMasterNodesTopCall(testCase.numMasterNodesTop, mock.Anything)
//			pastelClient.AssertSignCall(testCase.numSign, mock.Anything, mock.Anything,
//				testCase.fields.Ticket.PastelID, testCase.fields.Ticket.PastelIDPassphrase, "ed448")
//
//			// nodeClient mock assertion
//			nodeClient.Connection.AssertExpectations(t)
//			nodeClient.Client.AssertExpectations(t)
//			nodeClient.DownloadNft.AssertExpectations(t)
//			nodeClient.AssertConnectCall(testCase.numConnect, mock.Anything, mock.Anything, mock.Anything)
//			nodeClient.AssertDownloadNftCall(testCase.numDownLoadNft)
//			nodeClient.AssertDownloadCall(testCase.numDownload, mock.Anything, testCase.fields.Ticket.Txid,
//				mock.Anything, string(testCase.args.signature), testCase.args.ttxid)
//			nodeClient.AssertDoneCall(testCase.numDone)
//		})
//	}
//}
//
//
//func TestNodesMatchFiles(t *testing.T) {
//	t.Parallel()
//
//	testCases := []struct {
//		nodes     List
//		assertion assert.ErrorAssertionFunc
//	}{
//		{
//			nodes: List{
//				&NftDownloadNodeClient{
//					address: "127.0.0.1",
//					file:    []byte("1234"),
//				},
//				&NftDownloadNodeClient{
//					address: "127.0.0.2",
//					file:    []byte("1234"),
//				},
//			},
//			assertion: assert.NoError,
//		},
//		{
//			nodes: List{
//				&NftDownloadNodeClient{
//					address: "127.0.0.1",
//					file:    []byte("1234"),
//				},
//				&NftDownloadNodeClient{
//					address: "127.0.0.2",
//					file:    []byte("1235"),
//				},
//			},
//			assertion: assert.Error,
//		},
//	}
//
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			t.Parallel()
//
//			err := testCase.nodes.MatchFiles()
//			testCase.assertion(t, err)
//		})
//	}
//}
//
//func TestNodesFiles(t *testing.T) {
//	t.Parallel()
//
//	testCases := []struct {
//		nodes List
//		want  []byte
//	}{
//		{
//			nodes: List{
//				&NftDownloadNodeClient{
//					address: "127.0.0.1",
//					file:    []byte("1234"),
//				},
//			},
//			want: []byte("1234"),
//		},
//	}
//
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			t.Parallel()
//
//			file := testCase.nodes.File()
//			assert.Equal(t, testCase.want, file)
//		})
//	}
//}
//
//func TestNodesDownload(t *testing.T) {
//	// FIXME: enable later
//	t.Skip()
//	t.Parallel()
//
//	type args struct {
//		ctx                               context.Context
//		txid, timestamp, signature, ttxid string
//		// downloadErrs error
//	}
//
//	type nodeAttribute struct {
//		address string
//	}
//
//	testCases := []struct {
//		nodes              []nodeAttribute
//		args               args
//		err                error
//		file               []byte
//		numberDownloadCall int
//	}{
//		{
//			nodes:              []nodeAttribute{{"127.0.0.1:4444"}, {"127.0.0.1:4445"}},
//			args:               args{context.Background(), "txid", "timestamp", "signature", "ttxid"},
//			err:                nil,
//			file:               []byte("test"),
//			numberDownloadCall: 1,
//		},
//		{
//			nodes:              []nodeAttribute{{"127.0.0.1:4444"}, {"127.0.0.1:4445"}},
//			args:               args{context.Background(), "txid", "timestamp", "signature", "ttxid"},
//			err:                nil,
//			file:               nil,
//			numberDownloadCall: 1,
//		},
//	}
//
//	for i, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
//			t.Parallel()
//
//			nodes := List{}
//			clients := []*test.Client{}
//
//			for _, a := range testCase.nodes {
//				//client mock
//				client := test.NewMockClient(t)
//				//listen on uploadImage call
//				client.ListenOnConnect("", nil)
//				client.ListenOnDownload(testCase.file, testCase.err)
//				clients = append(clients, client)
//
//				nodes.Add(&NftDownloadNodeClient{
//					address:              a.address,
//					DownloadNftInterface: client.DownloadNft,
//				})
//			}
//
//			_, err := nodes.Download(testCase.args.ctx, testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid, 5*time.Second, nil)
//			assert.Equal(t, testCase.err, err)
//
//			//mock assertion each client
//			for _, client := range clients {
//				client.DownloadNft.AssertExpectations(t)
//				client.AssertDownloadCall(testCase.numberDownloadCall, testCase.args.ctx,
//					testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid)
//			}
//		})
//	}
//}
