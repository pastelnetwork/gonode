package nftsearch

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/mixins"

	"context"
	"sort"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"

	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
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
		req           *NftSearchingRequest
		actTicketsErr error
		regTicketErr  error
		ddAndFp       pastel.DDAndFingerprints
	}

	testCases := map[string]struct {
		args args
		want []RegTicketSearch
		fail bool
	}{
		"weird-ddfpvals": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &NftSearchingRequest{
					Query:        "alan",
					ArtistName:   true,
					Limit:        10,
					IsLikelyDupe: false,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         false,
					OverallRarenessScore: 0.5,
				},
			},
			want: []RegTicketSearch{RegTicketSearch{RegTicket: &regTicketA}},
			fail: false,
		},
		"match": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &NftSearchingRequest{
					Query:        "alan",
					ArtistName:   true,
					Limit:        10,
					IsLikelyDupe: false,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         false,
					OverallRarenessScore: 0.5,
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
				req: &NftSearchingRequest{
					Query:      "Alan",
					ArtistName: true,
					ArtTitle:   true,
					Limit:      10,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         false,
					OverallRarenessScore: 0.5,
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
				req: &NftSearchingRequest{
					Query:      "nowhere-to-be-found",
					ArtTitle:   true,
					ArtistName: true,
					Series:     true,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         false,
					OverallRarenessScore: 0.5,
				},
			},
			want: []RegTicketSearch{},
			fail: false,
		},

		"no-match-is-likely-dupe": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}},
					pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDB}}},
				regTickets: pastel.RegTickets{regTicketA, regTicketB},
				req: &NftSearchingRequest{
					Query:        "Alan",
					ArtistName:   true,
					ArtTitle:     true,
					Limit:        10,
					IsLikelyDupe: false,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         true,
					OverallRarenessScore: 0.5,
				},
			},
			want: []RegTicketSearch{},
			fail: false,
		},
		"no-match-is-opennsfw-range": {
			args: args{
				actTickets: pastel.ActTickets{pastel.ActTicket{ActTicketData: pastel.ActTicketData{RegTXID: testIDA}}},
				regTickets: pastel.RegTickets{regTicketA},
				req: &NftSearchingRequest{
					Query:        "alan",
					ArtistName:   true,
					Limit:        10,
					IsLikelyDupe: false,
					MaxNsfwScore: 0.2,
				},
				ddAndFp: pastel.DDAndFingerprints{
					OpenNSFWScore:        float32(0.5),
					IsLikelyDupe:         false,
					OverallRarenessScore: 0.5,
				},
			},
			want: []RegTicketSearch{},
			fail: false,
		},
	}

	ctx := context.Background()
	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase- %v", name), func(t *testing.T) {
			ddAndFpMarshalled, _ := json.Marshal(testCase.args.ddAndFp)

			pastelClientMock := pastelMock.NewMockClient(t)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{
					ExtKey:     "key",
					ExtAddress: "127.0.0.1:14445",
				})
			}

			pastelClientMock.ListenOnActTickets(testCase.args.actTickets, testCase.args.actTicketsErr)
			pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil)

			nodeClientMock := nodeMock.NewMockClient(t)
			nodeClientMock.ListenOnConnect("", nil).ListenOnRegisterNft().
				ListenOnClose(nil).ListenOnDownloadNft().ListenOnDownloadThumbnail(map[int][]byte{}, nil).ListenOnDownloadDDAndFP(ddAndFpMarshalled, nil)
			nodeClientMock.ConnectionInterface.On("DownloadNft").Return(nodeClientMock.DownloadNftInterface)

			doneCh := make(<-chan struct{})
			nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

			if len(testCase.args.actTickets) != len(testCase.args.regTickets) {
				t.Fatalf("#act_tickets != # reg_tickets")
			}

			for i, ticket := range testCase.args.actTickets {
				pastelClientMock.ListenOnRegTicket(ticket.ActTicketData.RegTXID, testCase.args.regTickets[i], testCase.args.regTicketErr)
			}

			service := &NftSearchingService{
				pastelHandler: mixins.NewPastelHandler(pastelClientMock.Client),
				nodeClient:    nodeClientMock,
				config:        NewConfig(),
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

func TestNewNftGetSearchTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *NftSearchingService
		req     *NftSearchingRequest
	}

	service := &NftSearchingService{
		config: NewConfig(),
	}

	req := &NftSearchingRequest{}

	testCases := map[string]struct {
		args args
		want *NftSearchingTask
	}{
		"new-task": {
			args: args{
				service: service,
				req:     req,
			},
			want: &NftSearchingTask{
				WalletNodeTask: common.NewWalletNodeTask(logPrefix),
				thumbnail:      NewThumbnailHandler(common.NewMeshHandlerSimple(nil, nil, nil)),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			task := NewNftSearchTask(testCase.args.service, testCase.args.req)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}
