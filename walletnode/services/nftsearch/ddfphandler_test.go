package nftsearch

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/pastelnetwork/gonode/mixins"

	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
)

func TestDDFPConnect(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	TopMNs := &pb.GetTopMNsReply{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprintf("%s:%d", fmt.Sprint(i), 14445), ExtKey: "key"})
		TopMNs.MnTopList = append(TopMNs.MnTopList, fmt.Sprintf("%s:%d", fmt.Sprint(i), 14445))
	}

	tests := map[string]struct {
		ddFP        DDFPHandler
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
				nodeClientMock.ListenOnConnect("0:14445", nil)
				nodeClientMock.ListenOnConnect("1:14445", nil)
				nodeClientMock.ListenOnConnect("2:14445", nil)
				nodeClientMock.ListenOnConnect("3:14445", tc.nodeErr)
				nodeClientMock.ListenOnConnect("4:14445", nil)
				nodeClientMock.ListenOnConnect("5:14445", nil)
				nodeClientMock.ListenOnConnect("6:14445", nil)
				nodeClientMock.ListenOnConnect("7:14445", nil)
				nodeClientMock.ListenOnConnect("8:14445", nil)
				nodeClientMock.ListenOnConnect("9:14445", nil)
				nodeClientMock.ListenOnConnect("10:14445", nil)
			}
			nodeClientMock.ListenOnClose(nil).ListenOnDownloadNft()
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)
			nodeClientMock.DownloadNftInterface.On("GetTopMNs", mock.Anything, mock.Anything).Return(TopMNs, nil)
			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(tc.nodesRet, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			meshHandlerOpts := common.MeshHandlerOpts{
				Task:          common.NewWalletNodeTask(logPrefix, nil),
				NodeMaker:     &NftSearchingNodeMaker{},
				PastelHandler: mixins.NewPastelHandler(pastelClientMock),
				NodeClient:    nodeClientMock,
				Configs:       &common.MeshHandlerConfig{},
			}

			ddFP := NewDDFPHandler(common.NewMeshHandler(meshHandlerOpts))

			_, cancel := context.WithCancel(context.Background())
			err := ddFP.Connect(context.Background(), tc.connections, cancel)
			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	TopMNs := &pb.GetTopMNsReply{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtKey: "key", ExtAddress: "127.0.0.1:14445"})
		TopMNs.MnTopList = append(TopMNs.MnTopList, "127.0.0.1:14445")
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
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadNft().ListenOnDownloadDDAndFP(tc.want, nil).ListenOnClose(nil)
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)
			nodeClientMock.DownloadNftInterface.On("GetTopMNs", mock.Anything, mock.Anything).Return(TopMNs, nil)
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

			err := nftGetSearchTask.ddAndFP.Connect(ctx, 1, cancel)
			assert.Nil(t, err)

			data, err := nftGetSearchTask.ddAndFP.Fetch(ctx, "txid")
			assert.Nil(t, err)
			assert.Equal(t, data, tc.want)

		})
	}
}
