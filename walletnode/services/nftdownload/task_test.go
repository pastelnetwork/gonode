package nftdownload

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/assert"
)

func TestTaskRun(t *testing.T) {
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
				ownerErr:    nil,
				downloadErr: nil,
				taskID:      "1",
				ctx:         context.Background(),
				networkFee:  0.4,
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
				downloadErr: nil,
				ownerErr:    nil,
				taskID:      "1",
				ctx:         context.Background(),
				networkFee:  0.4,
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
				downloadErr: nil,
				taskID:      "1",
				ctx:         context.Background(),
				networkFee:  0.4,
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
				ListenOnDone().ListenOnClose(nil).ListenOnDownload(nil, testCase.args.downloadErr)

			nodeClient.ConnectionInterface.On("DownloadNft").Return(nodeClient.DownloadNftInterface)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnTicketOwnership("txid", testCase.args.ownerErr).
				ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
				ListenOnSign([]byte(testCase.args.signature), testCase.args.returnErr).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{}}, nil)

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

			err := task.MatchFiles()
			if testCase.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
