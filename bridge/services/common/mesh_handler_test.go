package common_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	nodeMock "github.com/pastelnetwork/gonode/bridge/node/test"
	"github.com/pastelnetwork/gonode/bridge/services/common"
	"github.com/pastelnetwork/gonode/bridge/services/download"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func TestConnectToNSuperNodes(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprint(i), ExtKey: "key"})
	}

	tests := map[string]struct {
		connections int
		nodeErr     error
		nodesRet    pastel.MasterNodes
		err         error
	}{
		"one":              {connections: 1, err: nil, nodesRet: nodes},
		"max":              {connections: 10, err: nil, nodesRet: nodes},
		"node-connect-err": {connections: 3, nodeErr: errors.New("test"), err: nil, nodesRet: nodes},
		"no-nodes-err":     {connections: 10, nodesRet: pastel.MasterNodes{}, err: errors.New("not enough masternodes")},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			nodeClientMock := nodeMock.NewMockClient(t).ListenOnDownloadData()
			if tc.nodeErr == nil {
				nodeClientMock.ListenOnConnect("", nil)
			} else {
				nodeClientMock.ListenOnConnect("0", nil)
				nodeClientMock.ListenOnConnect("1", nil)
				nodeClientMock.ListenOnConnect("2", nil)
				nodeClientMock.ListenOnConnect("3", tc.nodeErr)
			}

			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(tc.nodesRet, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
					MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
				}, nil)

			meshHandlerOpts := common.MeshHandlerOpts{
				NodeMaker:     &download.DownloadingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			meshHandler := common.NewMeshHandler(meshHandlerOpts)

			err := meshHandler.ConnectToNSuperNodes(context.Background(), tc.connections)
			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestSetupMeshOfNSupernodesNodes(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprint(i), ExtKey: "key"})
	}

	tests := map[string]struct {
		connections int
		nodeErr     error
		nodesRet    pastel.MasterNodes
		err         error
	}{
		"one": {connections: 1, err: nil, nodesRet: nodes},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			nodeClientMock := nodeMock.NewMockClient(t).ListenOnDownloadData().ListenOnSession(nil).
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnDone().ListenOnClose(nil).ListenOnSessID("abs").
				ListenOnMeshNodes(nil)
			if tc.nodeErr == nil {
				nodeClientMock.ListenOnConnect("", nil)
			} else {
				nodeClientMock.ListenOnConnect("0", nil)
				nodeClientMock.ListenOnConnect("1", nil)
				nodeClientMock.ListenOnConnect("2", nil)
				nodeClientMock.ListenOnConnect("3", tc.nodeErr)
			}

			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(tc.nodesRet, nil).
				ListenOnGetBlockCount(10, nil).ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{
				MerkleRoot: hex.EncodeToString([]byte("PrimaryID")),
			}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			meshHandlerOpts := common.MeshHandlerOpts{
				NodeMaker:     &download.DownloadingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			meshHandler := common.NewMeshHandler(meshHandlerOpts)

			err := meshHandler.ConnectToNSuperNodes(context.Background(), tc.connections)
			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			_, _, err = meshHandler.SetupMeshOfNSupernodesNodes(context.Background())
			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}
