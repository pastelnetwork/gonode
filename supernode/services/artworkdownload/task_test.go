package artworkdownload

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"
	// taskMock "github.com/pastelnetwork/gonode/common/service/task/test"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func fakeRQIDsData(isValid bool) ([]byte, []string) {
	rqIDs := []string{
		"raptorQ ID1",
		"raptorQ ID2",
	}
	data := "fileID"
	data += "\n" + "blockhash"
	data += "\n" + "pastelID"
	if isValid {
		for _, id := range rqIDs {
			data += "\n" + id
		}
	}
	rqIDsData := []byte(data)
	return rqIDsData, rqIDs
}

func fakeRegiterTicket() pastel.RegTicket {
	appTicketData := pastel.AppTicket{
		AuthorPastelID:                 "pastelID",
		BlockTxID:                      "block_tx_id",
		BlockNum:                       10,
		ArtistName:                     "artist_name",
		ArtistWebsite:                  "artist_website",
		ArtistWrittenStatement:         "artist_written_statement",
		ArtworkTitle:                   "artwork_title",
		ArtworkSeriesName:              "artwork_series_name",
		ArtworkCreationVideoYoutubeURL: "artwork_creation_video_youtube_url",
		ArtworkKeywordSet:              "artwork_keyword_set",
		TotalCopies:                    10,
		PreviewHash:                    []byte("preview_hash"),
		Thumbnail1Hash:                 []byte("thumbnail1_hash"),
		Thumbnail2Hash:                 []byte("thumbnail2_hash"),
		DataHash:                       []byte("data_hash"),
		FingerprintsHash:               []byte("fingerprints_hash"),
		FingerprintsSignature:          []byte("fingerprints_signature"),
		PastelRarenessScore:            0.8,
		OpenNSFWScore:                  0.08,
		InternetRarenessScore:          0.4,
		RQIDs:                          []string{"raptorq ID1", "raptorq ID2"},
		RQOti:                          []byte("rq_oti"),
	}
	appTicket, _ := json.Marshal(&appTicketData)
	// appTicket := base64.StdEncoding.EncodeToString(data)
	// fmt.Println(string(appTicket))
	artTicketData := pastel.ArtTicket{
		Version:       1,
		Author:        []byte("pastelID"),
		BlockNum:      10,
		BlockHash:     []byte("block_hash"),
		Copies:        10,
		Royalty:       99,
		Green:         "green_address",
		AppTicket:     appTicket,
		AppTicketData: appTicketData,
	}
	artTicket, _ := json.Marshal(&artTicketData)
	// fmt.Println(string(artTicket))
	ticketSignature := pastel.TicketSignatures{}
	regTicketData := pastel.RegTicketData{
		Type:          "type",
		ArtistHeight:  235,
		Signatures:    ticketSignature,
		Key1:          "key1",
		Key2:          "key2",
		IsGreen:       true,
		StorageFee:    1,
		TotalCopies:   10,
		Royalty:       99,
		Version:       1,
		ArtTicket:     artTicket,
		ArtTicketData: artTicketData,
	}
	regTicket := pastel.RegTicket{
		Height:        235,
		TXID:          "txid",
		RegTicketData: regTicketData,
	}
	return regTicket
}

func TestNewTask(t *testing.T) {

	type args struct {
		service *Service
	}

	service := &Service{}

	testCases := []struct {
		args args
		want *Task
	}{
		{
			args: args{
				service: service,
			},
			want: &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
			},
		},
	}
	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			task := NewTask(testCase.args.service)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}

func TestTaskGetSymbolIDs(t *testing.T) {
	type args struct {
		rqIDsData []byte
	}

	rqIDsData, rqIDs := fakeRQIDsData(true)
	invalidRQIDsData, _ := fakeRQIDsData(false)

	testCases := []struct {
		args        args
		returnRqIDS []string
		assertion   assert.ErrorAssertionFunc
	}{
		{
			args: args{
				rqIDsData: rqIDsData,
			},
			returnRqIDS: rqIDs,
			assertion:   assert.NoError,
		},
		{
			args: args{
				rqIDsData: invalidRQIDsData,
			},
			returnRqIDS: nil,
			assertion:   assert.Error,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &Service{}
			task := NewTask(service)
			rqIDs, err := task.getRQSymbolIDs(testCase.args.rqIDsData)
			assert.Equal(t, testCase.returnRqIDS, rqIDs)
			testCase.assertion(t, err)
		})
	}
}

func TestDecodeRegTicket(t *testing.T) {

	type args struct {
		regTicket *pastel.RegTicket
	}
	regTicket := fakeRegiterTicket()
	testCases := []struct {
		args          args
		artTicketData pastel.ArtTicket
		appTicketData pastel.AppTicket
		assertion     assert.ErrorAssertionFunc
	}{
		{
			args: args{
				regTicket: &regTicket,
			},
			artTicketData: regTicket.RegTicketData.ArtTicketData,
			appTicketData: regTicket.RegTicketData.ArtTicketData.AppTicketData,
			assertion:     assert.NoError,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &Service{}
			task := NewTask(service)
			err := task.decodeRegTicket(testCase.args.regTicket)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.artTicketData, testCase.args.regTicket.RegTicketData.ArtTicketData)
			assert.Equal(t, testCase.appTicketData, testCase.args.regTicket.RegTicketData.ArtTicketData.AppTicketData)
		})
	}
}

func TestTaskRun(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			args: args{
				ctx: context.Background(),
			},
			assertion: assert.NoError,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &Service{}
			task := &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()
			err := task.Run(ctx)
			testCase.assertion(t, err)
		})
	}
}

func TestTaskDownload(t *testing.T) {

	type args struct {
		ctx                          context.Context
		txid                         string
		timestamp                    string
		signature                    string
		ttxid                        string
		regTicket                    pastel.RegTicket
		tradeTickets                 []pastel.TradeTicket
		isValid                      bool
		newActionChan                chan struct{}
		requiredStatusErr            error
		regTicketErr                 error
		listAvailableTradeTicketsErr error
		verifyErr                    error
		rqConnectErr                 error
		rqDecode                     *rqnode.Decode
		rqDecodeErr                  error
		p2pRetrieveData              []byte
		p2pRetrieveErr               error
	}

	regTicket := fakeRegiterTicket()
	tradeTickets := []pastel.TradeTicket{
		{
			Type:     "trade",
			PastelID: "buyerPastelID", // PastelID of the buyer
			SellTXID: "sell_txid",     // txid with sale ticket
			BueTXID:  "buy_txid",      // txid with buy ticket
			ArtTXID:  "ttxid",         // txid with either 1) art activation ticket or 2) trade ticket in it
		},
	}
	file := []byte("test file")
	rqDecode := &rqnode.Decode{
		File: file,
	}
	rqIDsData, _ := fakeRQIDsData(true)

	testCases := []struct {
		args                            args
		returnFile                      []byte
		numUpdateStatus                 int
		numberRequiredStatus            int
		numberNewAction                 int
		numberRegTicket                 int
		numberListAvailableTradeTickets int
		numberVerify                    int
		numberRQConnect                 int
		numberRQDone                    int
		numberRQRaptorQ                 int
		numberRQDecode                  int
		numberP2PRetrieve               int
		assertion                       assert.ErrorAssertionFunc
	}{
		{
			args: args{
				ctx:                          context.Background(),
				txid:                         "txid",
				timestamp:                    time.Now().Format(time.RFC3339),
				signature:                    "signature",
				ttxid:                        "ttxid",
				regTicket:                    regTicket,
				tradeTickets:                 tradeTickets,
				isValid:                      true,
				newActionChan:                make(chan struct{}),
				requiredStatusErr:            nil,
				regTicketErr:                 nil,
				listAvailableTradeTicketsErr: nil,
				verifyErr:                    nil,
				rqConnectErr:                 nil,
				rqDecode:                     rqDecode,
				rqDecodeErr:                  nil,
				p2pRetrieveData:              rqIDsData,
				p2pRetrieveErr:               nil,
			},
			returnFile:                      nil,
			numUpdateStatus:                 1,
			numberRequiredStatus:            1,
			numberNewAction:                 1,
			numberRegTicket:                 1,
			numberListAvailableTradeTickets: 1,
			numberVerify:                    1,
			numberRQConnect:                 1,
			numberRQDone:                    1,
			numberRQRaptorQ:                 1,
			numberRQDecode:                  0,
			numberP2PRetrieve:               4,
			assertion:                       assert.Error,
		},
	}
	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			config := &Config{}

			pastelClient := pastelMock.NewMockClient(t)
			if testCase.numberRegTicket > 0 {
				pastelClient.ListenOnRegTicket(testCase.args.txid, testCase.args.regTicket, testCase.args.regTicketErr)
			}
			if testCase.numberListAvailableTradeTickets > 0 {
				pastelClient.ListenOnListAvailableTradeTickets(testCase.args.tradeTickets, testCase.args.listAvailableTradeTicketsErr)
			}
			if testCase.numberVerify > 0 {
				pastelClient.ListenOnVerify(testCase.args.isValid, testCase.args.verifyErr)
			}

			p2pClient := p2pMock.NewMockClient(t)
			if testCase.numberP2PRetrieve > 0 {
				p2pClient.ListenOnRetrieve(testCase.args.p2pRetrieveData, testCase.args.p2pRetrieveErr)
			}

			raptorQClient := rqmock.NewMockClient(t)
			if testCase.numberRQConnect > 0 {
				raptorQClient.ListenOnConnect(testCase.args.rqConnectErr)
			}
			if testCase.numberRQDone > 0 {
				raptorQClient.ListenOnDone()
			}
			if testCase.numberRQRaptorQ > 0 {
				raptorQClient.ListenOnRaptorQ()
			}
			if testCase.numberRQDecode > 0 {
				raptorQClient.ListenOnDecode(testCase.args.rqDecode, testCase.args.rqDecodeErr)
			}

			// taskClient := taskMock.NewMockTask(t)
			// if testCase.numUpdateStatus > 0 {
			// 	taskClient.ListenOnUpdateStatus()
			// }
			// if testCase.numberRequiredStatus > 0 {
			// 	taskClient.ListenOnRequiredStatus(testCase.args.requiredStatusErr)
			// }
			// if testCase.numberNewAction > 0 {
			// 	taskClient.ListenOnNewAction(testCase.args.newActionChan)
			// }

			service := &Service{
				config:        config,
				pastelClient:  pastelClient.Client,
				p2pClient:     p2pClient.Client,
				raptorQClient: raptorQClient.Client,
				Worker:        task.NewWorker(),
			}

			task := &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
			}

			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()

			go task.RunAction(ctx)

			file, err := task.Download(ctx, testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.returnFile, file)

			// taskClient mock assertion
			// taskClient.AssertExpectations(t)
			// taskClient.AssertUpdateStatusCall(testCase.numUpdateStatus, mock.Anything)
			// taskClient.AssertRequiredStatusCall(testCase.numberRequiredStatus, mock.Anything)

			// pastelClient mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertRegTicketCall(testCase.numberRegTicket, mock.Anything, testCase.args.txid)
			pastelClient.AssertListAvailableTradeTicketsCall(testCase.numberListAvailableTradeTickets, mock.Anything)
			pastelID := testCase.args.regTicket.RegTicketData.ArtTicketData.AppTicketData.AuthorPastelID
			for _, trade := range testCase.args.tradeTickets {
				if trade.ArtTXID == testCase.args.ttxid {
					pastelID = trade.PastelID
				}
			}
			pastelClient.AssertVerifyCall(testCase.numberVerify, mock.Anything, []byte(testCase.args.timestamp),
				testCase.args.signature, pastelID)

			// p2pClient mock assertion
			p2pClient.AssertExpectations(t)
			p2pClient.AssertRetrieveCall(testCase.numberP2PRetrieve, mock.Anything, mock.Anything)

			// raptorQClient mock assertion
			raptorQClient.Connection.AssertExpectations(t)
			raptorQClient.RaptorQ.AssertExpectations(t)
			raptorQClient.AssertConnectCall(testCase.numberRQConnect, mock.Anything, mock.Anything)
			raptorQClient.AssertDoneCall(testCase.numberRQDone)
			raptorQClient.AssertDecodeCall(testCase.numberRQDecode, mock.Anything, mock.Anything)
		})
	}
}
