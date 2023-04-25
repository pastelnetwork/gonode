package collectionregister

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/tj/assert"
)

func add2NodesAnd2TicketSignatures(task *CollectionRegistrationTask) *CollectionRegistrationTask {
	task.NetworkHandler.Accepted = common.SuperNodePeerList{
		&common.SuperNodePeer{ID: "A"},
		&common.SuperNodePeer{ID: "B"},
	}
	task.PeersTicketSignature = map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}}
	return task
}

func makeEmptyCollectionRegTask(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client, rqClient rqnode.ClientInterface) *CollectionRegistrationTask {
	service := NewService(config, fileStorage, pastelClient, nodeClient, p2pClient, rqClient)
	task := NewCollectionRegistrationTask(service)
	task.Ticket = &pastel.CollectionTicket{}
	task.ActionTicketRegMetadata = &types.ActionRegMetadata{
		EstimatedFee: 100,
	}

	return task
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
			pastelClientMock.ListenOnRegisterCollectionTicket(tc.args.regRetID, tc.args.regErr)

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
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
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)

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

func TestTaskVerifyCollectionPeersSignature(t *testing.T) {
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
			pastelClientMock.ListenOnVerifyCollectionTicket(tc.args.verifyRet, tc.args.verifyErr)

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			err := task.VerifyPeersCollectionTicketSignature(context.Background(), task.Ticket)
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 19,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 20.3,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 18,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
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

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)

			err := <-task.WaitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.interval, false, 0, 0)
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
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
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

func TestTaskAddPeerCollectionTicketSignature(t *testing.T) {
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
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
		"no-node-err": {
			args: args{
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("not in Accepted list"),
		},
		"success-close-sign-chn": {
			args: args{
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
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyCollectionRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			task.NetworkHandler.Accepted = common.SuperNodePeerList{
				&common.SuperNodePeer{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID},
			}
			task.PeersTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := task.AddPeerCollectionTicketSignature(tc.args.nodeID, []byte{})
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
