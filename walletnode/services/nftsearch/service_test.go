package nftsearch

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
)

func TestRegTicket(t *testing.T) {
	testIDA := "test-id-a"
	testIDB := "test-id-b"

	t.Parallel()

	regTicketA := pastel.RegTicket{
		TXID: testIDA,
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					CreatorName: "Alan Majchrowicz",
				},
			},
		},
	}

	regTicketB := pastel.RegTicket{
		TXID: testIDB,
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					CreatorName: "Andy",
					NFTTitle:    "alantic",
				},
			},
		},
	}

	assignBase64strs(t, &regTicketA)
	assignBase64strs(t, &regTicketB)

	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{})
	}

	type args struct {
		regTicketID  string
		regTicketErr error
	}

	testCases := map[string]struct {
		args args
		want pastel.RegTicket
		err  error
	}{
		"simple-a": {
			args: args{
				regTicketID: testIDA,
			},
			want: regTicketA,
			err:  nil,
		},
		"simple-b": {
			args: args{
				regTicketID: testIDB,
			},
			want: regTicketB,
			err:  nil,
		},
	}

	ctx := context.Background()
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase- %v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadNft().ListenOnDownloadThumbnail([]byte{}, nil).ListenOnClose(nil)

			pastelClientMock.ListenOnRegTicket(testCase.args.regTicketID, testCase.want, testCase.args.regTicketErr)

			service := NewNftSearchService(NewConfig(), pastelClientMock, nodeClientMock)

			result, err := service.RegTicket(ctx, testCase.args.regTicketID)
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.want.TXID, result.TXID)
			assert.Equal(t, testCase.want.RegTicketData.NFTTicketData.Author,
				result.RegTicketData.NFTTicketData.Author)
		})
	}
}

func TestGetThumbnail(t *testing.T) {
	t.Parallel()

	regTicketA := pastel.RegTicket{
		TXID: "test-id-a",
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					CreatorName: "Alan Majchrowicz",
				},
			},
		},
	}

	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtKey: "key", ExtAddress: "127.0.0.1:14445"})
	}

	type args struct {
		regTicket *pastel.RegTicket
	}

	testCases := map[string]struct {
		args args
		want []byte
	}{
		"simple-a": {
			args: args{
				regTicket: &regTicketA,
			},
			want: []byte{2, 15},
		},
		"simple-b": {
			args: args{
				regTicket: &regTicketA,
			},
			want: []byte{14, 44},
		},
	}

	ctx := context.Background()
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase- %v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadNft().ListenOnDownloadThumbnail(testCase.want, nil).ListenOnClose(nil)
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)
			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			service := &NftSearchingService{
				pastelHandler: mixins.NewPastelHandler(pastelClientMock.Client),
				nodeClient:    nodeClientMock,
				config:        NewConfig(),
			}

			result, err := service.GetThumbnail(ctx, testCase.args.regTicket, "", "")
			assert.Nil(t, err)
			assert.Equal(t, testCase.want, result)
		})
	}
}

func assignBase64strs(t *testing.T, ticket *pastel.RegTicket) {
	artTicketBytes, err := pastel.EncodeNFTTicket(&ticket.RegTicketData.NFTTicketData)
	assert.Nil(t, err)
	ticket.RegTicketData.NFTTicket = artTicketBytes
}
