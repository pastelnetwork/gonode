package senseregister

import (
	"context"
	"image"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

const (
	testCreatorPastelID = "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW"
)

func newTestImageFile() (*files.File, error) {
	imageStorage := files.NewStorage(fs.NewFileStorage(os.TempDir()))
	imgFile := imageStorage.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create storage file: %w", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	png.Encode(f, img)
	imgFile.SetFormat(1)

	return imgFile, nil
}

func TestTaskRun(t *testing.T) {
	t.Parallel()

	type fields struct {
		Request *SenseRegistrationRequest
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
	}

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr error
	}{
		"success": {
			fields: fields{
				Request: &SenseRegistrationRequest{
					BurnTxID:              "txid",
					AppPastelID:           testCreatorPastelID,
					AppPastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
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
				Request: &SenseRegistrationRequest{
					BurnTxID:              "txid",
					AppPastelID:           testCreatorPastelID,
					AppPastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
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
	}

	for name, tc := range tests {
		testCase := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			senseFile, err := newTestImageFile()
			assert.NoError(t, err)

			// prepare task
			fg := pastel.Fingerprint{0.1, 0, 2}
			compressedFg, err := zstd.CompressLevel(nil, fg.Bytes(), 22)
			assert.Nil(t, err)
			testCase.args.fingerPrint = compressedFg

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
				ListenOnDone().
				ListenOnSendSignedTicket(1, nil).
				ListenOnClose(nil).ListenOnSendActionAct(nil)

			nodeClient.ConnectionInterface.On("RegisterSense").Return(nodeClient.RegisterSenseInterface)
			nodeClient.RegisterSenseInterface.On("MeshNodes", mock.Anything, mock.Anything).Return(nil)

			ddData := &pastel.DDAndFingerprints{
				InternetRareness:      &pastel.InternetRareness{},
				AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				FingerprintsStat:      &pastel.FingerprintsStat{},
				PerceptualImageHashes: &pastel.PerceptualImageHashes{},
				RarenessScores:        &pastel.RarenessScores{},
				Maxes:                 &pastel.Maxes{},
				Percentile:            &pastel.Percentile{},
			}
			compressed, err := pastel.ToCompressSignedDDAndFingerprints(ddData, []byte("signature"))
			assert.Nil(t, err)

			nodeClient.ListenOnProbeImage(compressed, true, testCase.args.returnErr)
			nodeClient.RegisterSenseInterface.On("SendRegMetadata", mock.Anything, mock.Anything).Return(nil)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
				ListenOnSign([]byte(testCase.args.signature), testCase.args.returnErr).
				ListenOnGetBlockCount(100, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Hash: "abc123", Height: 100}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{PqKey: ""}}, nil).
				ListenOnSendFromAddress("pre-burnt-txid", nil).
				ListenOnGetRawTransactionVerbose1(&pastel.GetRawTransactionVerbose1Result{Confirmations: 12}, nil).
				ListenOnVerify(true, nil).ListenOnGetBalance(10, nil).
				ListenOnActivateActionTicket("txid", nil).
				ListenOnGetActionFee(&pastel.GetActionFeesResult{CascadeFee: 10, SenseFee: 10}, nil)

			service := NewService(NewConfig(), pastelClientMock, nodeClient, nil, nil)
			service.config.WaitTxnValidInterval = 1
			go service.Run(testCase.args.ctx)

			Request := testCase.fields.Request
			Request.Image = senseFile
			task := NewSenseRegisterTask(service, Request)

			//create context with timeout to automatically end process after 5 sec
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 5*time.Second)
			defer cancel()

			err = task.Run(ctx)
			if testCase.wantErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
