package cascaderegister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/files"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	test "github.com/pastelnetwork/gonode/supernode/node/test/cascade_register"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

func add2NodesAnd2TicketSignatures(task *CascadeRegistrationTask) *CascadeRegistrationTask {
	task.NetworkHandler.Accepted = common.SuperNodePeerList{
		&common.SuperNodePeer{ID: "A"},
		&common.SuperNodePeer{ID: "B"},
	}
	task.PeersTicketSignature = map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}}
	return task
}

func makeEmptyCascadeRegTask(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client, rqClient rqnode.ClientInterface) *CascadeRegistrationTask {
	service := NewService(config, fileStorage, pastelClient, nodeClient, p2pClient, rqClient)
	task := NewCascadeRegistrationTask(service)
	task.Ticket = &pastel.ActionTicket{}

	return task
}

func TestTaskSignAndSendArtTicket(t *testing.T) {
	type args struct {
		signErr     error
		sendArtErr  error
		signReturns []byte
		primary     bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				primary: true,
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				signErr: errors.New("test"),
				primary: true,
			},
			wantErr: errors.New("test"),
		},
		"primary-err": {
			args: args{
				sendArtErr: errors.New("test"),
				signErr:    nil,
				primary:    false,
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.args.signReturns, tc.args.signErr)

			clientMock := test.NewMockClient(t)
			clientMock.ListenOnSendCascadeTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterCascade()

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, clientMock, nil, nil)

			task.NetworkHandler.ConnectedTo = &common.SuperNodePeer{
				ClientInterface: clientMock,
				NodeMaker:       RegisterCascadeNodeMaker{},
			}
			err := task.NetworkHandler.ConnectedTo.Connect(context.Background())
			assert.Nil(t, err)

			err = task.signAndSendCascadeTicket(context.Background(), tc.args.primary)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskRegisterAction(t *testing.T) {
	type args struct {
		regErr   error
		regRetID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args:    args{},
			wantErr: nil,
		},
		"err": {
			args: args{
				regErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnRegisterActionTicket(tc.args.regRetID, tc.args.regErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			id, err := task.registerAction(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.args.regRetID, id)
			}
		})
	}
}

func TestTaskPastelNodesByExtKey(t *testing.T) {
	type args struct {
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				masterNodesErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"node-err": {
			args: args{
				nodeID:         "B",
				masterNodesErr: nil,
			},
			wantErr: errors.New("not found"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)

			_, err := task.NetworkHandler.PastelNodeByExtKey(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskVerifyPeersSignature(t *testing.T) {
	type args struct {
		verifyErr error
		verifyRet bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				verifyRet: true,
				verifyErr: nil,
			},
			wantErr: nil,
		},
		"verify-err": {
			args: args{
				verifyRet: true,
				verifyErr: errors.New("test"),
			},
			wantErr: errors.New("verify signature"),
		},
		"verify-failure": {
			args: args{
				verifyRet: false,
				verifyErr: nil,
			},
			wantErr: errors.New("mistmatch"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(tc.args.verifyRet, tc.args.verifyErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			err := task.VerifyPeersTicketSignature(context.Background(), task.Ticket)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskWaitConfirmation(t *testing.T) {
	type args struct {
		txid             string
		interval         time.Duration
		minConfirmations int64
		ctxTimeout       time.Duration
	}

	testCases := map[string]struct {
		args    args
		wantErr error
		retRes  *pastel.GetRawTransactionVerbose1Result
		retErr  error
	}{
		"min-confirmations-timeout": {
			args: args{
				minConfirmations: 2,
				interval:         100 * time.Millisecond,
				ctxTimeout:       20 * time.Second,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			retErr:  nil,
			wantErr: errors.New("timeout"),
		},
		"success": {
			args: args{
				minConfirmations: 1,
				interval:         50 * time.Millisecond,
				ctxTimeout:       500 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: nil,
		},
		"ctx-done-err": {
			args: args{
				minConfirmations: 1,
				interval:         500 * time.Millisecond,
				ctxTimeout:       10 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: errors.New("context"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), tc.args.ctxTimeout)
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(1, nil)
			pastelClientMock.ListenOnGetRawTransactionVerbose1(tc.retRes, tc.retErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)

			err := <-task.WaitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.interval)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}

			cancel()
		})
	}

}

func TestTaskSessionNode(t *testing.T) {
	type args struct {
		nodeID         string
		masterNodesErr error
		status         common.Status
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				status:         common.StatusPrimaryMode,
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				status:         common.StatusConnected,
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: errors.New("status"),
		},
		"pastel-err": {
			args: args{
				status:         common.StatusPrimaryMode,
				masterNodesErr: errors.New("test"),
				nodeID:         "A",
			},
			wantErr: errors.New("get node"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			err := task.NetworkHandler.SessionNode(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskAddPeerCascadeTicketSignature(t *testing.T) {
	type args struct {
		status         common.Status
		nodeID         string
		masterNodesErr error
		acceptedNodeID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				status:         common.StatusImageProbed,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				status:         common.StatusConnected,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: errors.New("status"),
		},
		"no-node-err": {
			args: args{
				status:         common.StatusImageProbed,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("not in Accepted list"),
		},
		"success-close-sign-chn": {
			args: args{
				status:         common.StatusImageProbed,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: tc.args.nodeID})
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)

			task := makeEmptyCascadeRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			task.NetworkHandler.Accepted = common.SuperNodePeerList{
				&common.SuperNodePeer{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID},
			}
			task.PeersTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := task.AddPeerTicketSignature(tc.args.nodeID, []byte{}, common.StatusImageProbed)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskCompareRQSymbolID(t *testing.T) {
	type args struct {
		connectErr error
		fileErr    error
		addIDsFile bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"conn-err": {
			args: args{
				connectErr: errors.New("test"),
				fileErr:    nil,
				addIDsFile: true,
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				connectErr: nil,
				fileErr:    errors.New("test"),
				addIDsFile: true,
			},
			wantErr: errors.New("read image"),
		},
		"rqids-len-err": {
			args: args{
				connectErr: nil,
				fileErr:    nil,
				addIDsFile: false,
			},
			wantErr: errors.New("no symbols identifiers"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr)

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			task := makeEmptyCascadeRegTask(&Config{}, fsMock, nil, nil, nil, rqClientMock)

			storage := files.NewStorage(fsMock)
			task.Asset = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			if tc.args.addIDsFile {
				task.rawRqFile = []byte{'a'}
			}

			err := task.validateRQSymbolID(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskStoreRaptorQSymbols(t *testing.T) {
	type args struct {
		encodeErr  error
		connectErr error
		fileErr    error
		storeErr   error
		encodeResp *rqnode.Encode
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    errors.New("test"),
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"conn-err": {
			args: args{
				encodeErr:  nil,
				connectErr: errors.New("test"),
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"encode-err": {
			args: args{
				encodeErr:  errors.New("test"),
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"store-err": {
			args: args{
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   errors.New("test"),
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqFile := rq.SymbolIDFile{ID: "A"}
			bytes, err := json.Marshal(rqFile)
			assert.Nil(t, err)

			tc.args.encodeResp.Symbols = map[string][]byte{"A": bytes}

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr).
				ListenOnEncode(tc.args.encodeResp, tc.args.encodeErr)

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnStore("", tc.args.storeErr)

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			task := makeEmptyCascadeRegTask(&Config{}, fsMock, nil, nil, p2pClient, rqClientMock)

			storage := files.NewStorage(fsMock)
			task.Asset = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err = task.storeRaptorQSymbols(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
