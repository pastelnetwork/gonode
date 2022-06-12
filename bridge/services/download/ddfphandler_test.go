package download

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/mixins"

	nodeMock "github.com/pastelnetwork/gonode/bridge/node/test"
	"github.com/pastelnetwork/gonode/bridge/services/common"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func TestDDConnect(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprint(i), ExtKey: "key"})
	}

	tests := map[string]struct {
		dd          ddFpHandler
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
			pastelClientMock.ListenOnMasterNodesTop(tc.nodesRet, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			meshHandlerOpts := common.MeshHandlerOpts{
				NodeMaker:     &DownloadingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			thumbnail := newThumbnailHandler(common.NewMeshHandler(meshHandlerOpts))

			_, cancel := context.WithCancel(context.Background())
			err := thumbnail.Connect(context.Background(), tc.connections, cancel)
			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestDDFetch(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtKey: "key", ExtAddress: "127.0.0.1:14445"})
	}

	tests := map[string]struct {
		want    []byte
		wantErr error
	}{
		"a": {want: []byte{2, 78, 67, 33}, wantErr: nil},
		"b": {want: []byte{98, 22, 11}, wantErr: nil},
		"c": {want: []byte{11, 45, 19}, wantErr: nil},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadDDAndFP(tc.want, nil).ListenOnClose(nil)
			nodeClientMock.ConnectionInterface.On("DownloadData").Return(nodeClientMock.DownloadDataInterface)
			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			meshHandlerOpts := common.MeshHandlerOpts{
				NodeMaker:     &DownloadingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			dd := newDDFPHandler(common.NewMeshHandler(meshHandlerOpts))
			ctx, cancel := context.WithCancel(context.Background())
			err := dd.Connect(ctx, 1, cancel)
			assert.Nil(t, err)

			data, err := dd.Fetch(ctx, "txid")
			assert.Nil(t, err)
			assert.Equal(t, data, tc.want)
		})
	}
}
