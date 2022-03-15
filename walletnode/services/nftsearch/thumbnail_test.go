package nftsearch

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/mixins"

	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprint(i), ExtKey: "key"})
	}

	tests := map[string]struct {
		thumbnail   ThumbnailHandler
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
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nodeClientMock := nodeMock.NewMockClient(t)
			if tc.nodeErr == nil {
				nodeClientMock.ListenOnConnect("", nil)
			} else {
				nodeClientMock.ListenOnConnect("0", nil)
				nodeClientMock.ListenOnConnect("1", nil)
				nodeClientMock.ListenOnConnect("2", nil)
				nodeClientMock.ListenOnConnect("3", tc.nodeErr)
			}
			nodeClientMock.ListenOnClose(nil).ListenOnDownloadNft()
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)
			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(tc.nodesRet, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			meshHandlerOpts := common.MeshHandlerOpts{
				Task:          common.NewWalletNodeTask(logPrefix),
				NodeMaker:     &NftSearchingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			thumbnail := NewThumbnailHandler(common.NewMeshHandler(meshHandlerOpts))

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

func TestFetchOne(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtKey: "key", ExtAddress: "127.0.0.1:14445"})
	}

	tests := map[string]struct {
		want    map[int][]byte
		wantErr error
		//map[int][]byte{0: {2, 15}}
	}{
		"a": {want: map[int][]byte{0: {2, 78, 67, 33}}, wantErr: nil},
		"b": {want: map[int][]byte{0: {98, 22, 11}}, wantErr: nil},
		"c": {want: map[int][]byte{0: {11, 45, 19}}, wantErr: nil},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadNft().ListenOnDownloadThumbnail(tc.want, nil).ListenOnClose(nil)
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)
			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			service := &NftSearchingService{
				pastelHandler: mixins.NewPastelHandler(pastelClientMock.Client),
				nodeClient:    nodeClientMock,
				config:        NewConfig(),
			}

			nftGetSearchTask := NewNftGetSearchTask(service, "", "")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := nftGetSearchTask.thumbnail.Connect(ctx, 1, cancel)
			assert.Nil(t, err)

			data, err := nftGetSearchTask.thumbnail.FetchOne(ctx, "txid")
			assert.Nil(t, err)
			assert.Equal(t, data, tc.want[0])

		})
	}
}
