package artworksearch

import (
	"context"
	"fmt"
	"testing"

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
			ArtTicketData: pastel.ArtTicket{
				AppTicketData: pastel.AppTicket{
					AuthorPastelID: "author-id",
					ArtistName:     "Alan Majchrowicz",
				},
			},
		},
	}

	regTicketB := pastel.RegTicket{
		TXID: testIDB,
		RegTicketData: pastel.RegTicketData{
			ArtTicketData: pastel.ArtTicket{
				AppTicketData: pastel.AppTicket{
					AuthorPastelID: "author-id-b",
					ArtistName:     "Andy",
					ArtworkTitle:   "alantic",
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
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadArtwork().ListenOnDownloadThumbnail([]byte{}, nil).ListenOnClose(nil)

			pastelClientMock.ListenOnRegTicket(testCase.args.regTicketID, testCase.want, testCase.args.regTicketErr)

			service := NewService(NewConfig(), pastelClientMock, nil, nodeClientMock)

			result, err := service.RegTicket(ctx, testCase.args.regTicketID)
			assert.Equal(t, testCase.err, err)
			assert.Equal(t, testCase.want.TXID, result.TXID)
			assert.Equal(t, testCase.want.RegTicketData.ArtTicketData.Author,
				result.RegTicketData.ArtTicketData.Author)
		})
	}
}

func TestGetThumbnail(t *testing.T) {
	t.Parallel()

	regTicketA := pastel.RegTicket{
		TXID: "test-id-a",
		RegTicketData: pastel.RegTicketData{
			ArtTicketData: pastel.ArtTicket{
				AppTicketData: pastel.AppTicket{
					AuthorPastelID: "author-id",
					ArtistName:     "Alan Majchrowicz",
				},
			},
		},
	}

	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{})
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
			want: []byte{},
		},
		"simple-b": {
			args: args{
				regTicket: &regTicketA,
			},
			want: []byte{},
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
			nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadArtwork().ListenOnDownloadThumbnail([]byte{}, nil).ListenOnClose(nil)

			service := &Service{
				pastelClient: pastelClientMock.Client,
				nodeClient:   nodeClientMock,
				config:       NewConfig(),
			}

			result, err := service.GetThumbnail(ctx, testCase.args.regTicket)
			assert.Nil(t, err)
			assert.Equal(t, testCase.want, result)
		})
	}
}

func assignBase64strs(t *testing.T, ticket *pastel.RegTicket) {
	artTicketBytes, err := pastel.EncodeArtTicket(&ticket.RegTicketData.ArtTicketData)
	assert.Nil(t, err)
	ticket.RegTicketData.ArtTicket = artTicketBytes
}
