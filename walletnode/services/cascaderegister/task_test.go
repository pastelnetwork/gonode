package cascaderegister

import (
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/storage/ticketstore"

	"github.com/pastelnetwork/gonode/walletnode/services/download"

	rqnode "github.com/pastelnetwork/gonode/raptorq/node"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
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
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".pastel")

	// Make sure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		t.Fatalf("Failed to create base path: %v", err)
	}

	tempDir, err := createTempDirInPath(basePath)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up after test

	type fields struct {
		Request *common.ActionRegistrationRequest
	}

	type args struct {
		taskID            string
		ctx               context.Context
		networkFee        float64
		masterNodes       pastel.MasterNodes
		primarySessID     string
		pastelIDS         []string
		signature         []byte
		returnErr         error
		connectErr        error
		encodeInfoReturns *rqnode.EncodeInfo
	}

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr error
	}{
		/*"success": {
			fields: fields{
				Request: &common.ActionRegistrationRequest{
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
				signature:     []byte("sign"),
				returnErr:     nil,
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
						"test-file": {
							ID:                uuid.New().String(),
							SymbolIdentifiers: []string{"test-s1, test-s2"},
							BlockHash:         "test-block-hash",
							PastelID:          "test-pastel-id",
						},
					},
					EncoderParam: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
				},
			},
		},*/

		"failure": {
			wantErr: errors.New("test"),
			fields: fields{
				Request: &common.ActionRegistrationRequest{
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
				signature:     []byte("sign"),
				returnErr:     errors.New("test"),
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
						"test-file": {
							ID:                uuid.New().String(),
							SymbolIdentifiers: []string{"test-s1, test-s2"},
							BlockHash:         "test-block-hash",
							PastelID:          "test-pastel-id",
						},
					},
					EncoderParam: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
				},
			},
		},
	}

	for name, tc := range tests {
		testCase := tc

		t.Run(name, func(t *testing.T) {
			cascadeFile, err := newTestImageFile()
			assert.NoError(t, err)

			// prepare task
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
				ListenOnDone().
				//ListenOnSendSignedTicket(100, nil).
				ListenOnSendActionAct(nil).
				ListenOnClose(nil)
			nodeClient.RegisterCascadeInterface.
				On("SendSignedTicket", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("", nil).Times(1)
			nodeClient.RegisterCascadeInterface.
				On("SendSignedTicket", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("100", nil).Times(1)
			nodeClient.RegisterCascadeInterface.
				On("SendSignedTicket", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return("", nil).Times(1)

			nodeClient.ConnectionInterface.On("RegisterCascade").Return(nodeClient.RegisterCascadeInterface)
			nodeClient.RegisterCascadeInterface.On("MeshNodes", mock.Anything, mock.Anything).Return(nil)
			nodeClient.RegisterCascadeInterface.On("SendRegMetadata", mock.Anything, mock.Anything).Return(nil)

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
				ListenOnGetActionFee(&pastel.GetActionFeesResult{CascadeFee: 10, SenseFee: 10}, nil).On("FindActionRegTicketsByLabel", mock.Anything, mock.Anything, mock.Anything).
				Return(pastel.ActionTicketDatas{}, nil).On("FindNFTRegTicketsByLabel", mock.Anything, mock.Anything).Return(pastel.RegTickets{}, nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(testCase.args.encodeInfoReturns, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(testCase.args.connectErr)

			ticketDB, err := ticketstore.OpenTicketingDb(configurer.DefaultPath())
			assert.NoError(t, err)

			downloadService := download.NewNftDownloadService(download.NewConfig(), pastelClientMock, nodeClient, nil)
			service := NewService(NewConfig(), pastelClientMock, nodeClient, nil, nil, *downloadService, nil, ticketDB)
			service.rqClient = rqClientMock
			service.config.WaitTxnValidInterval = 1

			go service.Run(testCase.args.ctx)

			Request := testCase.fields.Request
			Request.Image = cascadeFile
			task := NewCascadeRegisterTask(service, Request)

			//create context with timeout to automatically end process after 5 sec
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 5*time.Second)
			defer cancel()

			err = task.Run(ctx)
			if testCase.wantErr != nil {
				assert.True(t, task.Status().IsFailure())
			} else {
				task.Status().Is(common.StatusTaskCompleted)
				assert.Nil(t, err)
			}
		})
	}
}

func createTempDirInPath(basePath string) (string, error) {
	dir, err := os.MkdirTemp(basePath, ".pastel")
	if err != nil {
		return "", err
	}
	return dir, nil
}
