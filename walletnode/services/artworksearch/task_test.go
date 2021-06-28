package artworksearch

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/service/task"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

const (
	testIDA = "test-id-a"
	testIDB = "test-id-b"
)

func TestRunTask(t *testing.T) {
	t.Parallel()

	toBase64 := func(from interface{}) ([]byte, error) {
		testBytes, err := json.Marshal(from)
		return []byte(b64.StdEncoding.EncodeToString([]byte(testBytes))), err
	}

	assignBase64strs := func(ticket *pastel.RegTicket) {
		base64Bytes, err := toBase64(ticket.RegTicketData.ArtTicketData.AppTicketData)
		assert.Nil(t, err)
		ticket.RegTicketData.ArtTicketData.AppTicket = base64Bytes

		base64Bytes, err = toBase64(ticket.RegTicketData.ArtTicketData)
		assert.Nil(t, err)
		ticket.RegTicketData.ArtTicket = base64Bytes
	}

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

	assignBase64strs(&regTicketA)
	assignBase64strs(&regTicketB)

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

	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase- %v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			p2pClientMock := p2pMock.NewMockClient(t)

			pastelClientMock.ListenOnActTickets(testCase.args.actTickets, testCase.args.actTicketsErr)

			if len(testCase.args.actTickets) != len(testCase.args.regTickets) {
				t.Fatalf("#act_tickets != # reg_tickets")
			}

			for i, ticket := range testCase.args.actTickets {
				pastelClientMock.ListenOnRegTicket(ticket.ActTicketData.RegTXID, testCase.args.regTickets[i], testCase.args.regTicketErr)
			}

			p2pClientMock.ListenOnRetrieve([]byte{}, nil)
			service := &Service{
				pastelClient: pastelClientMock.Client,
				p2pClient:    p2pClientMock,
			}

			task := NewTask(service, testCase.args.req)
			resultChan := task.SubscribeSearchResult()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

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

	service := &Service{}
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
				Task:    task.New(StatusTaskStarted),
				Service: service,
				request: req,
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
