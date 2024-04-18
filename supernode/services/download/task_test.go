package download

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/supernode/services/common"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/utils"

	// taskMock "github.com/pastelnetwork/gonode/common/service/task/test"

	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func fakeRQIDsData() ([]byte, []string) {
	rqfile := &rqnode.RawSymbolIDFile{
		ID:                "09f6c459-ec2a-4db1-a8fe-0648fd97b5cb",
		PastelID:          "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
		SymbolIdentifiers: []string{"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z"},
	}

	dataJSON, _ := json.Marshal(rqfile)
	encoded := utils.B64Encode(dataJSON)

	var buffer bytes.Buffer
	buffer.Write(encoded)
	buffer.WriteByte(46)
	buffer.WriteString("test-signature")
	buffer.WriteByte(46)
	buffer.WriteString(strconv.Itoa(55))

	compressedData, _ := utils.Compress(buffer.Bytes(), 4)

	return compressedData, []string{
		"24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z",
	}
}

func fakeRegiterTicket() pastel.RegTicket {
	appTicketData := pastel.AppTicket{
		CreatorName:                "creator_name",
		CreatorWebsite:             "artist_website",
		CreatorWrittenStatement:    "artist_written_statement",
		NFTTitle:                   "nft_title",
		NFTSeriesName:              "nft_series_name",
		NFTCreationVideoYoutubeURL: "nft_creation_video_youtube_url",
		NFTKeywordSet:              "nft_keyword_set",
		TotalCopies:                10,
		PreviewHash:                []byte("preview_hash"),
		Thumbnail1Hash:             []byte("thumbnail1_hash"),
		Thumbnail2Hash:             []byte("thumbnail2_hash"),
		DataHash:                   []byte("data_hash"),
		RQIDs:                      []string{"raptorq ID1", "raptorq ID2"},
		RQOti:                      []byte("rq_oti"),
	}
	//appTicket, _ := json.Marshal(&appTicketData)

	appTicketBytes, _ := json.Marshal(&appTicketData)

	appTicket := base64.RawStdEncoding.EncodeToString(appTicketBytes)

	nftTicketData := pastel.NFTTicket{
		Version:       1,
		Author:        "pastelID",
		BlockNum:      10,
		BlockHash:     "block_hash",
		Copies:        10,
		Royalty:       99,
		Green:         false,
		AppTicket:     appTicket,
		AppTicketData: appTicketData,
	}
	artTicket, _ := json.Marshal(&nftTicketData)
	ticketSignature := pastel.RegTicketSignatures{}
	regTicketData := pastel.RegTicketData{
		Type:          "type",
		CreatorHeight: 235,
		Signatures:    ticketSignature,
		Label:         "label",
		Green:         false,
		StorageFee:    1,
		TotalCopies:   10,
		Royalty:       99,
		Version:       1,
		NFTTicket:     artTicket,
		NFTTicketData: nftTicketData,
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
		service *NftDownloaderService
	}

	service := &NftDownloaderService{}

	testCases := []struct {
		args args
		want *NftDownloadingTask
	}{
		{
			args: args{
				service: service,
			},
			want: &NftDownloadingTask{
				SuperNodeTask:        &common.SuperNodeTask{Task: task.New(common.StatusTaskStarted)},
				NftDownloaderService: service,
			},
		},
	}
	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			task := NewNftDownloadingTask(testCase.args.service)
			assert.Equal(t, testCase.want.NftDownloaderService, task.NftDownloaderService)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}

func TestTaskGetSymbolIDs(t *testing.T) {
	type args struct {
		rqIDsData []byte
	}

	rqIDsData, haveRqIDs := fakeRQIDsData()

	testCases := []struct {
		args        args
		returnRqIDs []string
		assertion   assert.ErrorAssertionFunc
	}{
		{
			args: args{
				rqIDsData: rqIDsData,
			},
			returnRqIDs: haveRqIDs,
			assertion:   assert.NoError,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &NftDownloaderService{}
			task := NewNftDownloadingTask(service)
			rqIDs, err := task.getRQSymbolIDs(context.Background(), "id", testCase.args.rqIDsData)
			assert.Nil(t, err)
			if err == nil {
				assert.Equal(t, haveRqIDs, rqIDs)
			}
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
		nftTicketData pastel.NFTTicket
		appTicketData pastel.AppTicket
		assertion     assert.ErrorAssertionFunc
	}{
		{
			args: args{
				regTicket: &regTicket,
			},
			nftTicketData: regTicket.RegTicketData.NFTTicketData,
			appTicketData: regTicket.RegTicketData.NFTTicketData.AppTicketData,
			assertion:     assert.NoError,
		},
	}

	for i, testCase := range testCases {

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &NftDownloaderService{}
			task := NewNftDownloadingTask(service)
			err := task.decodeRegTicket(testCase.args.regTicket)
			testCase.assertion(t, err)
			//assert.Equal(t, testCase.nftTicketData, testCase.args.regTicket.RegTicketData.NFTTicketData)
			assert.Equal(t, testCase.appTicketData, testCase.args.regTicket.RegTicketData.NFTTicketData.AppTicketData)
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

		testCase := testCase
		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			service := &NftDownloaderService{}
			task := &NftDownloadingTask{
				SuperNodeTask:        &common.SuperNodeTask{Task: task.New(common.StatusTaskStarted)},
				NftDownloaderService: service,
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 6*time.Second)
			defer cancel()
			err := task.Run(ctx)
			testCase.assertion(t, err)
		})
	}
}

func TestTaskDownload(t *testing.T) {
	// FIXME : enable again
	t.Skip()
	type args struct {
		ctx                          context.Context
		txid                         string
		timestamp                    string
		signature                    string
		ttxid                        string
		regTicket                    pastel.RegTicket
		tradeTickets                 []pastel.TradeTicket
		isValid                      bool
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
			Height: 238,
			TXID:   "ttxid",
			Ticket: pastel.TradeTicketData{
				Type:     "trade",
				PastelID: "buyerPastelID", // PastelID of the buyer
				SellTXID: "sell_txid",     // txid with sale ticket
				BuyTXID:  "buy_txid",      // txid with buy ticket
				NFTTXID:  "NFT_txid",      // txid with either 1) NFT activation ticket or 2) trade ticket in it
			},
		},
	}
	file := []byte("test file")
	rqDecode := &rqnode.Decode{
		File: file,
	}
	rqIDsData, _ := fakeRQIDsData()

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
				timestamp:                    time.Now().UTC().Format(time.RFC3339),
				signature:                    "signature",
				ttxid:                        "ttxid",
				regTicket:                    regTicket,
				tradeTickets:                 tradeTickets,
				isValid:                      true,
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
			numberP2PRetrieve:               2,
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

			service := &NftDownloaderService{
				config: config,
				SuperNodeService: &common.SuperNodeService{
					PastelClient: pastelClient.Client,
					P2PClient:    p2pClient.Client,
					Worker:       task.NewWorker(),
				},
			}

			task := &NftDownloadingTask{
				SuperNodeTask:        &common.SuperNodeTask{Task: task.New(common.StatusTaskStarted)},
				NftDownloaderService: service,
			}

			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()

			go task.RunAction(ctx)

			file, _ := task.Download(ctx, testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid, "", false)

			assert.Equal(t, testCase.returnFile, file)

			// taskClient mock assertion
			// taskClient.AssertExpectations(t)
			// taskClient.AssertUpdateStatusCall(testCase.numUpdateStatus, mock.Anything)
			// taskClient.AssertRequiredStatusCall(testCase.numberRequiredStatus, mock.Anything)

			// pastelClient mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertRegTicketCall(testCase.numberRegTicket, mock.Anything, testCase.args.txid)
			pastelClient.AssertListAvailableTradeTicketsCall(testCase.numberListAvailableTradeTickets, mock.Anything)
			pastelID := testCase.args.regTicket.RegTicketData.NFTTicketData.Author
			for _, trade := range testCase.args.tradeTickets {
				if trade.TXID == testCase.args.ttxid {
					pastelID = trade.Ticket.PastelID
				}
			}
			pastelClient.AssertVerifyCall(testCase.numberVerify, mock.Anything, []byte(testCase.args.timestamp),
				testCase.args.signature, pastelID, "ed448")

			// p2pClient mock assertion
			p2pClient.AssertExpectations(t)
			//p2pClient.AssertRetrieveCall(testCase.numberP2PRetrieve, mock.Anything, mock.Anything)

			// raptorQClient mock assertion
			raptorQClient.Connection.AssertExpectations(t)
			raptorQClient.RaptorQ.AssertExpectations(t)
			raptorQClient.AssertConnectCall(testCase.numberRQConnect, mock.Anything, mock.Anything)
			raptorQClient.AssertDoneCall(testCase.numberRQDone)
			raptorQClient.AssertDecodeCall(testCase.numberRQDecode, mock.Anything, mock.Anything)
		})
	}
}
