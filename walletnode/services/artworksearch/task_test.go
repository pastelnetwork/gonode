package artworksearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"sort"
	"testing"

	"github.com/pastelnetwork/gonode/pastel"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
	thumbnail "github.com/pastelnetwork/gonode/walletnode/services/artworksearch/thumbnail"
	"github.com/stretchr/testify/assert"
)

const (
	testIDA = "test-id-a"
	testIDB = "test-id-b"
)

func TestRunTask(t *testing.T) {

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

	type args struct {
		actTickets    pastel.ActTickets
		regTickets    pastel.RegTickets
		req           *ArtSearchRequest
		actTicketsErr error
		regTicketErr  error
	}

	testCases := map[string]struct {
		args args
		want []RegTicketSearch
		fail bool
	}{
		"match": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &ArtSearchRequest{
					Query:      "alan",
					ArtistName: true,
					Limit:      10,
				},
			},
			want: []RegTicketSearch{RegTicketSearch{RegTicket: &regTicketA}},
			fail: false,
		},

		"multiple-match": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}},
					pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDB}}},
				regTickets: pastel.RegTickets{regTicketA, regTicketB},
				req: &ArtSearchRequest{
					Query:      "Alan",
					ArtistName: true,
					ArtTitle:   true,
					Limit:      10,
				},
			},
			want: []RegTicketSearch{RegTicketSearch{RegTicket: &regTicketB, MatchIndex: 0},
				RegTicketSearch{RegTicket: &regTicketA, MatchIndex: 1},
			},
			fail: false,
		},

		"no-match": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &ArtSearchRequest{
					Query:      "nowhere-to-be-found",
					ArtTitle:   true,
					ArtistName: true,
					Series:     true,
				},
			},
			want: []RegTicketSearch{},
			fail: false,
		},

		"act-tickets-err": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &ArtSearchRequest{
					Query:      "alan",
					ArtTitle:   true,
					ArtistName: true,
				},
				actTicketsErr: errors.New("test-err"),
			},

			want: []RegTicketSearch{},
			fail: true,
		},

		"reg-ticket-err": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDB}}},
				regTickets: pastel.RegTickets{regTicketB},
				req: &ArtSearchRequest{
					Query:      "alan",
					ArtTitle:   true,
					ArtistName: true,
				},
				regTicketErr: errors.New("test-err"),
			},

			want: []RegTicketSearch{},
			fail: true,
		},
	}

	ctx := context.Background()
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase- %v", name), func(t *testing.T) {
			pastelClientMock := pastelMock.NewMockClient(t)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}

			pastelClientMock.ListenOnActTickets(testCase.args.actTickets, testCase.args.actTicketsErr)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnRegisterArtwork().
				ListenOnClose(nil).ListenOnDownloadArtwork().ListenOnDownloadThumbnail([]byte{}, nil)

			if len(testCase.args.actTickets) != len(testCase.args.regTickets) {
				t.Fatalf("#act_tickets != # reg_tickets")
			}

			for i, ticket := range testCase.args.actTickets {
				pastelClientMock.ListenOnRegTicket(ticket.ActTicketData.RegTXID, testCase.args.regTickets[i], testCase.args.regTicketErr)
			}

			service := &NftSearchService{
				pastelClient: pastelClientMock.Client,
				nodeClient:   nodeClientMock,
				config:       NewConfig(),
			}

			task := NewNftSearchTask(service, testCase.args.req)
			resultChan := task.SubscribeSearchResult()

			go task.Run(ctx)

			results := []*RegTicketSearch{}
		loop:
			for {
				select {
				case <-ctx.Done():
					t.Fail()
					break loop
				case srch, ok := <-resultChan:
					if !ok {
						assert.Equal(t, testCase.fail, task.Status().IsFailure())
						break loop
					}

					results = append(results, srch)
				}
			}

			assert.Equal(t, len(testCase.want), len(results))

			sort.Slice(results, func(i, j int) bool {
				return results[i].MatchIndex < results[j].MatchIndex
			})

			for i, result := range results {
				assert.Equal(t, testCase.want[i].TXID, result.TXID)
				assert.Equal(t, testCase.want[i].RegTicketData.NFTTicketData.Author,
					result.RegTicketData.NFTTicketData.Author)

				assert.Equal(t, testCase.want[i].MatchIndex, result.MatchIndex)
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *NftSearchService
		req     *ArtSearchRequest
	}

	service := &NftSearchService{
		config: NewConfig(),
	}

	req := &ArtSearchRequest{}

	testCases := map[string]struct {
		args args
		want *NftSearchTask
	}{
		"new-task": {
			args: args{
				service: service,
				req:     req,
			},
			want: &NftSearchTask{
				WalletNodeTask:   common.NewWalletNodeTask(logPrefix),
				NftSearchService: service,
				request:          req,
				thumbnailHelper:  thumbnail.New(service.pastelClient, service.nodeClient, service.config.ConnectToNodeTimeout),
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			task := NewNftSearchTask(testCase.args.service, testCase.args.req)
			assert.Equal(t, testCase.want.NftSearchService, task.NftSearchService)
			assert.Equal(t, testCase.want.request, task.request)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
