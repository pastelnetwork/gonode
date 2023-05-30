package download

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/assert"
)

func TestTaskRun(t *testing.T) {

	regTicketA := pastel.RegTicket{
		TXID: "txid",
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					CreatorName:            "Alan Majchrowicz",
					MakePubliclyAccessible: true,
				},
			},
		},
	}

	regTicketB := pastel.RegTicket{
		TXID: "txid",
		RegTicketData: pastel.RegTicketData{
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					CreatorName:            "Andy",
					NFTTitle:               "alantic",
					MakePubliclyAccessible: false,
				},
			},
		},
	}

	assignBase64strs(t, &regTicketA)
	assignBase64strs(t, &regTicketB)

	t.Parallel()
	type fields struct {
		Request *NftDownloadingRequest
	}

	type args struct {
		taskID        string
		ctx           context.Context
		networkFee    float64
		masterNodes   pastel.MasterNodes
		primarySessID string
		pastelIDS     []string
		fingerPrint   []byte
		signature     []byte
		returnErr     error
		ownerErr      error
		downloadErr   error
		wantRegTicket pastel.RegTicket
	}

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr error
	}{
		"success": {
			wantErr: nil,
			fields: fields{
				Request: &NftDownloadingRequest{
					Txid:               "txid",
					PastelID:           "abc",
					PastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				wantRegTicket: regTicketA,
				ownerErr:      nil,
				downloadErr:   nil,
				taskID:        "1",
				ctx:           context.Background(),
				networkFee:    0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4448", ExtKey: "4"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4449", ExtKey: "5"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4450", ExtKey: "6"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4451", ExtKey: "7"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4452", ExtKey: "8"},
				},

				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3", "4", "5", "6", "7"},
				fingerPrint:   []byte("match"),
				signature:     []byte("sign"),
				returnErr:     nil,
			},
		},

		"failure": {
			wantErr: errors.New("test"),
			fields: fields{
				Request: &NftDownloadingRequest{
					Txid:               "txid",
					PastelID:           "abc",
					PastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				downloadErr:   nil,
				ownerErr:      nil,
				taskID:        "1",
				ctx:           context.Background(),
				wantRegTicket: regTicketB,
				networkFee:    0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3"},
				fingerPrint:   []byte("match"),
				signature:     []byte("sign"),
				returnErr:     errors.New("test"),
			},
		},
		"non-owner": {
			wantErr: errors.New("test"),
			fields: fields{
				Request: &NftDownloadingRequest{
					Txid:               "txid",
					PastelID:           "abc",
					PastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				downloadErr:   nil,
				taskID:        "1",
				ctx:           context.Background(),
				wantRegTicket: regTicketB,
				networkFee:    0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3"},
				fingerPrint:   []byte("match"),
				signature:     []byte("sign"),
				returnErr:     nil,
				ownerErr:      errors.New("test"),
			},
		},
		"non-owner-but-skip-verify-check": {
			wantErr: nil,
			fields: fields{
				Request: &NftDownloadingRequest{
					Txid:               "txid",
					PastelID:           "abc",
					PastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				downloadErr:   nil,
				taskID:        "1",
				ctx:           context.Background(),
				wantRegTicket: regTicketA,
				networkFee:    0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4448", ExtKey: "4"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4449", ExtKey: "5"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4450", ExtKey: "6"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4451", ExtKey: "7"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4452", ExtKey: "8"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3"},
				fingerPrint:   []byte("match"),
				signature:     []byte("sign"),
				returnErr:     nil,
				ownerErr:      errors.New("test"),
			},
		},
	}

	for name, tc := range tests {
		testCase := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
				ListenOnDone().ListenOnClose(nil).ListenOnDownload([]byte{1, 2, 3}, testCase.args.downloadErr)

			nodeClient.ConnectionInterface.On("DownloadNft").Return(nodeClient.DownloadNftInterface)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnTicketOwnership("txid", testCase.args.ownerErr).
				ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
				ListenOnSign([]byte(testCase.args.signature), testCase.args.returnErr).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{}}, nil).On("RegTicket", mock.Anything, mock.Anything).Return(testCase.args.wantRegTicket, nil)

			service := NewNftDownloadService(NewConfig(), pastelClientMock, nodeClient)

			Request := testCase.fields.Request
			task := NewNftDownloadTask(service, Request)

			//create context with timeout to automatically end process after 5 sec
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 5*time.Second)
			defer cancel()

			err := task.Run(ctx)
			if testCase.wantErr != nil {

				assert.True(t, task.Status().IsFailure())
			} else {
				assert.Nil(t, err)
				assert.Nil(t, task.Error())
				assert.True(t, task.Status().Is(common.StatusTaskCompleted))
			}
		})
	}
}

func TestMatchFiles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		files   []downFile
		wantErr error
		wantN   int
	}{
		{
			files: []downFile{
				{file: []byte{1, 2, 3}},
				{file: []byte{1, 2, 3}},
				{file: []byte{1, 2, 3}},
			},

			wantErr: nil,
		},
		{
			files: []downFile{
				{file: []byte{1, 2, 3}},
				{file: []byte{2, 2, 3}},
				{file: []byte{1, 3, 3}},
				{file: []byte{1, 5, 3}},
				{file: []byte{1, 9, 3}},
				{file: []byte{7, 2, 3}},
				{file: []byte{6, 6, 3}},
				{file: []byte{7, 2, 3}},
				{file: []byte{7, 2, 3}},
			},

			wantErr: nil,
			wantN:   5,
		},
		{
			files: []downFile{
				{file: []byte{1, 2, 3}},
				{file: []byte{2, 2, 3}},
				{file: []byte{1, 3, 3}},
				{file: []byte{1, 5, 3}},
				{file: []byte{1, 9, 3}},
				{file: []byte{1, 2, 3}},
				{file: []byte{6, 6, 3}},
				{file: []byte{7, 2, 3}},
				{file: []byte{1, 2, 3}},
			},

			wantErr: nil,
			wantN:   0,
		},
		{
			files: []downFile{
				{file: []byte{2, 2, 3}},
				{file: []byte{2, 2, 3}},
				{file: []byte{1, 2, 3}},
			},

			wantErr: errors.New("test"),
		},
		{
			files: []downFile{
				{file: []byte{1, 2, 3}},
				{file: []byte{1, 2, 3}},
				{file: []byte{}}},

			wantErr: errors.New("test"),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			service := NewNftDownloadService(NewConfig(), nil, nil)
			task := NewNftDownloadTask(service, &NftDownloadingRequest{})
			task.files = testCase.files

			n, _, err := task.MatchFiles()
			if testCase.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, testCase.wantN, n)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func assignBase64strs(t *testing.T, ticket *pastel.RegTicket) {
	artTicketBytes, err := pastel.EncodeNFTTicket(&ticket.RegTicketData.NFTTicketData)
	assert.Nil(t, err)
	ticket.RegTicketData.NFTTicket = artTicketBytes
	ticket.RegTicketData.NFTTicketData = pastel.NFTTicket{}
}
