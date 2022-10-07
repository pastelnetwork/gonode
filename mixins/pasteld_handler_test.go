package mixins

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/b85"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestVerifySignature(t *testing.T) {
	type args struct {
		validate bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				validate: true,
			},
			wantErr: nil,
		},
		"error": {
			args: args{
				validate: false,
			},
			wantErr: errors.New("test"),
		},

		"invalid": {
			args: args{
				validate: false,
			},
			wantErr: errors.New("signature verification failed"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(tc.args.validate, tc.wantErr)

			h := NewPastelHandler(pastelClientMock)
			isValid, err := h.VerifySignature(context.Background(), []byte{}, "", "", pastel.SignAlgorithmED448)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.args.validate, isValid)
			}

		})
	}
}

func TestGetEstimatedSenseFee(t *testing.T) {
	type args struct {
		size float64
	}

	testCases := map[string]struct {
		args    args
		res     *pastel.GetActionFeesResult
		wantErr error
	}{
		"success": {
			args: args{
				size: 1.0,
			},
			res:     &pastel.GetActionFeesResult{SenseFee: 10.0, CascadeFee: 15.0},
			wantErr: nil,
		},
		"error": {
			args: args{
				size: 2.0,
			},
			res:     &pastel.GetActionFeesResult{SenseFee: 10.0, CascadeFee: 15.0},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetActionFee(tc.res, tc.wantErr)

			h := NewPastelHandler(pastelClientMock)
			fee, err := h.GetEstimatedSenseFee(context.Background(), tc.args.size)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.res.SenseFee, fee)
			}

		})
	}
}

func TestGetEstimatedCascadeFee(t *testing.T) {
	type args struct {
		size float64
	}

	testCases := map[string]struct {
		args    args
		res     *pastel.GetActionFeesResult
		wantErr error
	}{
		"success": {
			args: args{
				size: 1.0,
			},
			res:     &pastel.GetActionFeesResult{SenseFee: 10.0, CascadeFee: 15.0},
			wantErr: nil,
		},
		"error": {
			args: args{
				size: 2.0,
			},
			res:     &pastel.GetActionFeesResult{SenseFee: 10.0, CascadeFee: 15.0},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetActionFee(tc.res, tc.wantErr)

			h := NewPastelHandler(pastelClientMock)
			fee, err := h.GetEstimatedCascadeFee(context.Background(), tc.args.size)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.res.CascadeFee, fee)
			}

		})
	}
}

func TestGetBlock(t *testing.T) {
	type args struct {
		size       float64
		block      int32
		blockErr   error
		verboseErr error
		res        *pastel.GetBlockVerbose1Result
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				size:  1.0,
				block: 10,
				res:   &pastel.GetBlockVerbose1Result{Hash: "hash"},
			},
			wantErr: nil,
		},
		"block-error": {
			args: args{
				size:     2.0,
				res:      &pastel.GetBlockVerbose1Result{Hash: "hash"},
				blockErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"verbose-error": {
			args: args{
				size:       2.0,
				res:        &pastel.GetBlockVerbose1Result{Hash: "hash"},
				verboseErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(tc.args.block, tc.args.blockErr).
				ListenOnGetBlockVerbose1(tc.args.res, tc.args.verboseErr)

			h := NewPastelHandler(pastelClientMock)
			block, hash, err := h.GetBlock(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, block, int(tc.args.block))
				assert.Equal(t, hash, tc.args.res.Hash)
			}

		})
	}
}

func TestCheckBalanceToPayRegistrationFee(t *testing.T) {
	type args struct {
		fee        float64
		maxfee     float64
		balanceErr error
		balance    float64
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				fee:     1.0,
				maxfee:  1.5,
				balance: 5.0,
			},
			wantErr: nil,
		},
		"max-free": {
			args: args{
				fee:     2.0,
				maxfee:  1.5,
				balance: 5.0,
			},
			wantErr: errors.New("registration fee is to expensive"),
		},
		"invalid-free": {
			args: args{
				fee:     0.0,
				maxfee:  1.5,
				balance: 5.0,
			},
			wantErr: errors.New("invalid fee amount to check"),
		},
		"insufficient-balance-free": {
			args: args{
				fee:     10.0,
				maxfee:  11.5,
				balance: 5.0,
			},
			wantErr: errors.New("not enough PSL"),
		},
		"error": {
			args: args{
				fee:        2.0,
				maxfee:     3.0,
				balanceErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBalance(tc.args.balance, tc.args.balanceErr)

			h := NewPastelHandler(pastelClientMock)
			err := h.CheckBalanceToPayRegistrationFee(context.Background(), "", tc.args.fee, tc.args.maxfee)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWaitTxIDValid(t *testing.T) {
	type args struct {
		txid string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				txid: "test-txid",
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnSign([]byte("signature"), nil).
				ListenOnGetBlockCount(100, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Hash: "abc123", Height: 100}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{PqKey: ""}}, nil).
				ListenOnSendFromAddress("pre-burnt-txid", nil).
				ListenOnGetRawTransactionVerbose1(&pastel.GetRawTransactionVerbose1Result{Confirmations: 12}, nil).
				ListenOnVerify(true, nil).ListenOnGetBalance(10, nil).
				ListenOnActivateActionTicket("txid", nil).
				ListenOnGetActionFee(&pastel.GetActionFeesResult{CascadeFee: 10, SenseFee: 10}, nil)

			h := NewPastelHandler(pastelClientMock)
			err := h.WaitTxidValid(context.Background(), "txid", 10, time.Microsecond)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestBurnSomeCoins(t *testing.T) {
	type args struct {
		amount  int64
		percent uint
		burnErr error
		txid    string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				txid:    "test-txid",
				amount:  10,
				percent: 5,
				burnErr: nil,
			},
			wantErr: nil,
		},
		"error": {
			args: args{
				txid:    "test-txid",
				amount:  10,
				percent: 5,
				burnErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSendFromAddress(tc.args.txid, tc.args.burnErr)

			h := NewPastelHandler(pastelClientMock)
			txid, err := h.BurnSomeCoins(context.Background(), "", tc.args.amount, tc.args.percent)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {

				assert.Nil(t, err)
				assert.Equal(t, tc.args.txid, txid)
			}
		})
	}
}

func TestRegTicket(t *testing.T) {
	type args struct {
		txid         string
		regTicket    pastel.RegTicket
		regTicketErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				txid:      "test-txid",
				regTicket: fakeRegiterTicket(),
			},
			wantErr: nil,
		},
		"error": {
			args: args{
				txid:         "test-txid",
				regTicket:    fakeRegiterTicket(),
				regTicketErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnRegTicket(tc.args.txid, tc.args.regTicket, tc.args.regTicketErr)

			h := NewPastelHandler(pastelClientMock)
			_ = h.GetBurnAddress()
			ticket, err := h.RegTicket(context.Background(), "test-txid")
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.args.regTicket.TXID, ticket.TXID)
			}
		})
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

	appTicketBytes, _ := json.Marshal(&appTicketData)

	appTicket := b85.Encode(appTicketBytes)

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
		TXID:          "test-txid",
		RegTicketData: regTicketData,
	}
	return regTicket
}
