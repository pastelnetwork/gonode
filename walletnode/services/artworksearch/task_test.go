package artworksearch

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/service/task"
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
	t.Skip()
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

			pastelClientMock.ListenOnActTickets(testCase.args.actTickets, testCase.args.actTicketsErr)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnRegisterArtwork().ListenOnClose(nil)

			if len(testCase.args.actTickets) != len(testCase.args.regTickets) {
				t.Fatalf("#act_tickets != # reg_tickets")
			}

			for i, ticket := range testCase.args.actTickets {
				pastelClientMock.ListenOnRegTicket(ticket.ActTicketData.RegTXID, testCase.args.regTickets[i], testCase.args.regTicketErr)
			}

			service := &Service{
				pastelClient: pastelClientMock.Client,
				nodeClient:   nodeClientMock,
				config:       NewConfig(),
			}

			task := NewTask(service, testCase.args.req)
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
				assert.Equal(t, testCase.want[i].RegTicketData.ArtTicketData.Author,
					result.RegTicketData.ArtTicketData.Author)

				assert.Equal(t, testCase.want[i].MatchIndex, result.MatchIndex)
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *Service
		req     *ArtSearchRequest
	}

	service := &Service{
		config: NewConfig(),
	}

	req := &ArtSearchRequest{}

	testCases := map[string]struct {
		args args
		want *Task
	}{
		"new-task": {
			args: args{
				service: service,
				req:     req,
			},
			want: &Task{
				Task:            task.New(StatusTaskStarted),
				Service:         service,
				request:         req,
				thumbnailHelper: thumbnail.New(service.pastelClient, service.nodeClient, service.config.ConnectTimeout),
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			task := NewTask(testCase.args.service, testCase.args.req)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.request, task.request)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
