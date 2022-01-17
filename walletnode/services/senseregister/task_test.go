package senseregister

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"image"
	"image/png"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	stateMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/register_sense"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNode(address, pastelID string) *node.SenseRegisterNode {
	return node.NewNode(nil, address, pastelID)
}

func pullPastelAddressIDNodes(nodes node.List) []string {
	var v []string
	for _, n := range nodes {
		v = append(v, fmt.Sprintf("%s:%s", n.PastelID(), n.String()))
	}

	sort.Strings(v)
	return v
}

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
	t.Skip()

	type fields struct {
		Request *SenseRegisterRequest
	}

	type args struct {
		taskID            string
		ctx               context.Context
		networkFee        float64
		masterNodes       pastel.MasterNodes
		primarySessID     string
		pastelIDS         []string
		fingerPrint       []byte
		signature         []byte
		returnErr         error
		connectErr        error
		encodeInfoReturns *rqnode.EncodeInfo
	}

	testCases := []struct {
		fields            fields
		args              args
		assertion         assert.ErrorAssertionFunc
		numSessIDCall     int
		numUpdateStatus   int
		numSignCall       int
		numSessionCall    int
		numConnectToCall  int
		numProbeImageCall int
	}{
		{
			fields: fields{&SenseRegisterRequest{BurnTxID: "1", AppPastelID: "1", AppPastelIDPassphrase: "2"}},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4446", ExtKey: "2"},
					pastel.MasterNode{Fee: 0.3, ExtAddress: "127.0.0.1:4447", ExtKey: "3"},
					pastel.MasterNode{Fee: 0.4, ExtAddress: "127.0.0.1:4448", ExtKey: "4"},
				},
				primarySessID: "sesid1",
				pastelIDS:     []string{"2", "3", "4"},
				fingerPrint:   []byte("match"),
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
				},
			},
			assertion:         assert.NoError,
			numSessIDCall:     3,
			numUpdateStatus:   4,
			numSignCall:       4,
			numSessionCall:    4,
			numConnectToCall:  3,
			numProbeImageCall: 4,
		},
	}

	t.Run("group", func(t *testing.T) {
		artworkFile, err := newTestImageFile()
		assert.NoError(t, err)

		for i, testCase := range testCases {
			testCase := testCase

			// prepare task
			fg := pastel.Fingerprint{0.1, 0, 2}
			compressedFg, err := zstd.CompressLevel(nil, fg.Bytes(), 22)
			assert.Nil(t, err)
			testCase.args.fingerPrint = compressedFg

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				var task *SenseRegisterTask

				nodeClient := test.NewMockClient(t)
				nodeClient.
					ListenOnConnect("", testCase.args.returnErr).
					ListenOnRegisterSense().
					ListenOnSession(testCase.args.returnErr).
					ListenOnConnectTo(testCase.args.returnErr).
					ListenOnSessID(testCase.args.primarySessID).
					ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
					ListenOnDone().
					ListenOnSendSignedTicket("1", nil).
					ListenOnClose(nil)

				//need to remove generate thumbnail file
				customProbeImageFunc := func(ctx context.Context, file *files.File) *pastel.DDAndFingerprints {
					file.Remove()
					return &pastel.DDAndFingerprints{ZstdCompressedFingerprint: testCase.args.fingerPrint}
				}
				nodeClient.ListenOnProbeImage(customProbeImageFunc, testCase.args.returnErr)

				pastelClientMock := pastelMock.NewMockClient(t)
				pastelClientMock.
					ListenOnStorageNetworkFee(testCase.args.networkFee, testCase.args.returnErr).
					ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
					ListenOnSign([]byte(testCase.args.signature), testCase.args.returnErr).
					ListenOnGetBlockCount(100, nil).
					ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{}, nil).
					ListenOnFindTicketByID(&pastel.IDTicket{}, nil).
					ListenOnSendFromAddress("pre-burnt-txid", nil).
					ListenOnGetRawTransactionVerbose1(&pastel.GetRawTransactionVerbose1Result{Confirmations: 10}, nil).
					ListenOnRegisterNFTTicket("art-act-txid", nil)

				rqClientMock := rqMock.NewMockClient(t)
				rqClientMock.ListenOnEncodeInfo(testCase.args.encodeInfoReturns, nil)
				rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
				rqClientMock.ListenOnConnect(testCase.args.connectErr)

				service := &Service{
					pastelClient: pastelClientMock.Client,
					nodeClient:   nodeClient.Client,
					config:       NewConfig(),
				}

				taskClient := stateMock.NewMockTask(t)
				taskClient.
					ListenOnID(testCase.args.taskID).
					ListenOnUpdateStatus().
					ListenOnSetStatusNotifyFunc()

				Request := testCase.fields.Request
				Request.Image = artworkFile
				Request.BurnTxID = "1"
				task = &SenseRegisterTask{
					WalletNodeTask: &common.WalletNodeTask{
						Task:      taskClient.Task,
						LogPrefix: logPrefix,
					},
					Service: service,
					Request: Request,
				}

				//create context with timeout to automatically end process after 1 sec
				ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
				defer cancel()
				testCase.assertion(t, task.Run(ctx))

				taskClient.AssertExpectations(t)
				taskClient.AssertIDCall(1)
				taskClient.AssertUpdateStatusCall(testCase.numUpdateStatus, mock.Anything)
				taskClient.AssertSetStatusNotifyFuncCall(1, mock.Anything)

				// //pastelClient mock assertion
				pastelClientMock.AssertStorageNetworkFeeCall(1, mock.Anything)
				pastelClientMock.AssertSignCall(testCase.numSignCall,
					mock.Anything,
					mock.Anything,
					Request.AppPastelID,
					Request.AppPastelIDPassphrase,
					mock.Anything,
				)

				// //nodeClient mock assertion
				nodeClient.AssertAcceptedNodesCall(1, mock.Anything)
				nodeClient.AssertSessIDCall(testCase.numSessIDCall)
				nodeClient.AssertSessionCall(testCase.numSessionCall, mock.Anything, false)
				nodeClient.AssertConnectToCall(testCase.numConnectToCall, mock.Anything, mock.Anything, testCase.args.primarySessID)
				nodeClient.AssertProbeImageCall(testCase.numProbeImageCall, mock.Anything, mock.IsType(&files.File{}))
			})
		}
	})

}

func TestTaskMeshNodes(t *testing.T) {
	t.Parallel()

	type nodeArg struct {
		address  string
		pastelID string
	}

	type args struct {
		ctx             context.Context
		nodes           []nodeArg
		primaryIndex    int
		primaryPastelID string
		primarySessID   string
		pastelIDS       []string
		returnErr       error
		acceptNodeErr   error
	}

	testCases := []struct {
		args             args
		want             []string
		assertion        assert.ErrorAssertionFunc
		numSessionCall   int
		numSessIDCall    int
		numConnectToCall int
	}{
		{
			args: args{
				ctx:          context.Background(),
				primaryIndex: 1,
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
					{"127.0.0.3", "3"},
					{"127.0.0.4", "4"},
					{"127.0.0.5", "5"},
					{"127.0.0.6", "6"},
					{"127.0.0.7", "7"},
				},
				primaryPastelID: "2",
				primarySessID:   "xdcfjc",
				pastelIDS:       []string{"1", "4", "7"},
				returnErr:       nil,
				acceptNodeErr:   nil,
			},
			assertion:        assert.NoError,
			numSessionCall:   7,
			numSessIDCall:    6,
			numConnectToCall: 6,
			want:             []string{"1:127.0.0.1", "2:127.0.0.2", "4:127.0.0.4", "7:127.0.0.7"},
		}, {
			args: args{
				ctx:          context.Background(),
				primaryIndex: 0,
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.1", "2"},
					{"127.0.0.1", "3"},
				},
				primaryPastelID: "1",
				primarySessID:   "xdxdf",
				pastelIDS:       []string{"2", "3"},
				returnErr:       nil,
				acceptNodeErr:   fmt.Errorf("primary node not accepted"),
			},
			assertion:        assert.Error,
			numSessionCall:   3,
			numSessIDCall:    2,
			numConnectToCall: 2,
			want:             nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			//create new client mock
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnRegisterSense().
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.acceptNodeErr)

			nodes := node.List{}
			for _, n := range testCase.args.nodes {
				nodes.Add(node.NewNode(nodeClient.Client, n.address, n.pastelID))
			}

			service := &Service{
				config: NewConfig(),
			}

			task := &SenseRegisterTask{Service: service, Request: &SenseRegisterRequest{}}
			got, err := task.meshNodes(testCase.args.ctx, nodes, testCase.args.primaryIndex)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, pullPastelAddressIDNodes(got))

			nodeClient.AssertAcceptedNodesCall(1, mock.Anything)
			nodeClient.AssertSessIDCall(testCase.numSessIDCall)
			nodeClient.AssertSessionCall(testCase.numSessionCall, mock.Anything, false)
			nodeClient.AssertConnectToCall(testCase.numConnectToCall, mock.Anything, types.MeshedSuperNode{NodeID: testCase.args.primaryPastelID, SessID: testCase.args.primarySessID})
			nodeClient.Client.AssertExpectations(t)
			nodeClient.Connection.AssertExpectations(t)
		})

	}
}

func TestTaskPastelTopNodes(t *testing.T) {
	t.Parallel()

	type fields struct {
		Task    task.Task
		Request *SenseRegisterRequest
	}

	type args struct {
		ctx             context.Context
		returnMn        pastel.MasterNodes
		returnMnTopErr  error
		returnFindIDErr error
	}

	testCases := []struct {
		fields    fields
		args      args
		want      node.List
		assertion assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				Request: &SenseRegisterRequest{
					BurnTxID: "xdxdf",
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnMnTopErr:  nil,
				returnFindIDErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4444", "1"),
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			fields: fields{
				Request: &SenseRegisterRequest{
					BurnTxID: "xdxdf",
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.5, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnMnTopErr:  nil,
				returnFindIDErr: fmt.Errorf("connection timeout"),
			},
			want:      nil,
			assertion: assert.NoError,
		}, {
			args: args{
				ctx:             context.Background(),
				returnMn:        nil,
				returnMnTopErr:  fmt.Errorf("connection timeout"),
				returnFindIDErr: nil,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClient := pastelMock.NewMockClient(t)
			pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnMnTopErr)
			if testCase.args.returnMnTopErr == nil {
				pastelClient.ListenOnFindTicketByID(&pastel.IDTicket{}, testCase.args.returnFindIDErr)
			}

			service := &Service{
				pastelClient: pastelClient.Client,
			}

			task := &SenseRegisterTask{
				WalletNodeTask: &common.WalletNodeTask{
					Task:      testCase.fields.Task,
					LogPrefix: logPrefix,
				},
				Service: service,
				Request: testCase.fields.Request,
			}
			got, err := task.pastelTopNodes(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertMasterNodesTopCall(1, mock.Anything)
		})
	}

}

func TestNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		service *Service
		Request *SenseRegisterRequest
	}

	service := &Service{}
	Request := &SenseRegisterRequest{}

	testCases := []struct {
		args args
		want *SenseRegisterTask
	}{
		{
			args: args{
				service: service,
				Request: Request,
			},
			want: &SenseRegisterTask{
				WalletNodeTask: common.NewWalletNodeTask(logPrefix),
				Service:        service,
				Request:        Request,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			task := NewSenseRegisterTask(testCase.args.service, testCase.args.Request)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Request, task.Request)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}

func TestTaskCreateTicket(t *testing.T) {
	t.Parallel()

	type args struct {
		task *SenseRegisterTask
	}

	testCases := map[string]struct {
		args    args
		want    *pastel.ActionTicket
		wantErr error
	}{
		"data-hash-error": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "test-id",
					},
					Service: &Service{
						config: NewConfig(),
					},
				},
			},
			wantErr: errEmptyDatahash,
			want:    nil,
		},
		"success": {
			args: args{
				task: &SenseRegisterTask{
					fingerprint:          []byte{},
					datahash:             []byte{},
					fingerprintAndScores: &pastel.DDAndFingerprints{},
					Request: &SenseRegisterRequest{
						AppPastelID:           "test-id",
						AppPastelIDPassphrase: "test-name",
					},
					Service: &Service{
						config: NewConfig(),
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			blockNum := 1
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(int32(blockNum), nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{}, nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.want = &pastel.ActionTicket{
				Version:    1,
				Caller:     tc.args.task.Request.AppPastelID,
				BlockNum:   tc.args.task.creatorBlockHeight,
				BlockHash:  tc.args.task.creatorBlockHash,
				ActionType: pastel.ActionTypeSense,
				APITicketData: &pastel.APISenseTicket{
					DataHash:             tc.args.task.datahash,
					DDAndFingerprintsIc:  tc.args.task.ddAndFingerprintsIc,
					DDAndFingerprintsMax: tc.args.task.config.DDAndFingerprintsMax,
					DDAndFingerprintsIDs: tc.args.task.ddAndFingerprintsIDs,
				},
			}

			err := tc.args.task.createSenseTicket(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, tc.args.task.ticket)
				assert.Equal(t, *tc.want, *tc.args.task.ticket)
			}
		})
	}
}

func TestTaskGetBlock(t *testing.T) {
	t.Parallel()

	type args struct {
		task            *SenseRegisterTask
		blockCountErr   error
		blockVerboseErr error
		blockNum        int32
		blockInfo       *pastel.GetBlockVerbose1Result
	}

	testCases := map[string]struct {
		args                  args
		wantArtistblockHash   string
		wantArtistBlockHeight int
		wantErr               error
	}{
		"success": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
			},
			wantArtistBlockHeight: 10,
			wantErr:               nil,
		},
		"block-count-err": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
				blockCountErr: errors.New("block-count-err"),
			},
			wantArtistBlockHeight: 10,
			wantErr:               errors.New("block-count-err"),
		},
		"block-verbose-err": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
				blockVerboseErr: errors.New("verbose-err"),
			},
			wantArtistBlockHeight: 10,
			wantErr:               errors.New("verbose-err"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(tc.args.blockNum, tc.args.blockCountErr).
				ListenOnGetBlockVerbose1(tc.args.blockInfo, tc.args.blockVerboseErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.wantArtistblockHash = tc.args.blockInfo.Hash

			err := tc.args.task.getBlock(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantArtistblockHash, tc.args.task.creatorBlockHash)
				assert.Equal(t, tc.wantArtistBlockHeight, tc.args.task.creatorBlockHeight)
			}
		})
	}
}

func TestTaskSignTicket(t *testing.T) {
	type args struct {
		task        *SenseRegisterTask
		signErr     error
		signReturns []byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ActionTicket{},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ActionTicket{},
				},
				signErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.args.signReturns, tc.args.signErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			err := tc.args.task.signTicket(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.args.signReturns, tc.args.task.creatorSignature)
			}
		})
	}
}

func TestWaitTxnValid(t *testing.T) {
	type args struct {
		task                            *SenseRegisterTask
		getRawTransactionVerbose1RetErr error
		getRawTransactionVerbose1Ret    *pastel.GetRawTransactionVerbose1Result
		ctxDone                         bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"ctx-done-err": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
				},
				ctxDone: true,
			},
			wantErr: errors.New("context done"),
		},
		"get-raw-transaction-err": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
				},
				getRawTransactionVerbose1Ret:    &pastel.GetRawTransactionVerbose1Result{},
				getRawTransactionVerbose1RetErr: errors.New("test"),
			},
			wantErr: errors.New("timeout"),
		},
		"insufficient-confirmations": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
				},
				getRawTransactionVerbose1Ret:    &pastel.GetRawTransactionVerbose1Result{Confirmations: 4},
				getRawTransactionVerbose1RetErr: nil,
			},
			wantErr: errors.New("timeout"),
		},
		"success": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
				},
				getRawTransactionVerbose1Ret:    &pastel.GetRawTransactionVerbose1Result{Confirmations: 5},
				getRawTransactionVerbose1RetErr: nil,
			},
			wantErr: nil,
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(1, nil)
			pastelClientMock.ListenOnGetRawTransactionVerbose1(tc.args.getRawTransactionVerbose1Ret,
				tc.args.getRawTransactionVerbose1RetErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			ctx, cancel := context.WithCancel(context.Background())
			if tc.args.ctxDone {
				cancel()
			}

			err := tc.args.task.waitTxidValid(ctx, "txid", 5, 500*time.Millisecond)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
			cancel()
		})
	}

}

func TestTaskProbeImage(t *testing.T) {
	fingerAndScores := &pastel.DDAndFingerprints{
		Block:                      "Block",
		Principal:                  "Principal",
		DupeDetectionSystemVersion: "v1.0",

		IsLikelyDupe:     true,
		IsRareOnInternet: true,

		RarenessScores: &pastel.RarenessScores{
			CombinedRarenessScore:         0,
			XgboostPredictedRarenessScore: 0,
			NnPredictedRarenessScore:      0,
			OverallAverageRarenessScore:   0,
		},
		InternetRareness: &pastel.InternetRareness{
			MatchesFoundOnFirstPage: 0,
			NumberOfPagesOfResults:  0,
			URLOfFirstMatchInPage:   "",
		},

		OpenNSFWScore: 0.1,
		AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
			Drawings: 0.1,
			Hentai:   0.2,
			Neutral:  0.3,
			Porn:     0.4,
			Sexy:     0.5,
		},

		ImageFingerprintOfCandidateImageFile: []float32{1, 2, 3},
		FingerprintsStat: &pastel.FingerprintsStat{
			NumberOfFingerprintsRequiringFurtherTesting1: 1,
			NumberOfFingerprintsRequiringFurtherTesting2: 2,
			NumberOfFingerprintsRequiringFurtherTesting3: 3,
			NumberOfFingerprintsRequiringFurtherTesting4: 4,
			NumberOfFingerprintsRequiringFurtherTesting5: 5,
			NumberOfFingerprintsRequiringFurtherTesting6: 6,
			NumberOfFingerprintsOfSuspectedDupes:         7,
		},

		HashOfCandidateImageFile: "HashOfCandidateImageFile",
		PerceptualImageHashes: &pastel.PerceptualImageHashes{
			PDQHash:        "PdqHash",
			PerceptualHash: "PerceptualHash",
			AverageHash:    "AverageHash",
			DifferenceHash: "DifferenceHash",
			NeuralHash:     "NeuralhashHash",
		},
		PerceptualHashOverlapCount: 1,

		Maxes: &pastel.Maxes{
			PearsonMax:           1.0,
			SpearmanMax:          2.0,
			KendallMax:           3.0,
			HoeffdingMax:         4.0,
			MutualInformationMax: 5.0,
			HsicMax:              6.0,
			XgbimportanceMax:     7.0,
		},
		Percentile: &pastel.Percentile{
			PearsonTop1BpsPercentile:             1.0,
			SpearmanTop1BpsPercentile:            2.0,
			KendallTop1BpsPercentile:             3.0,
			HoeffdingTop10BpsPercentile:          4.0,
			MutualInformationTop100BpsPercentile: 5.0,
			HsicTop100BpsPercentile:              6.0,
			XgbimportanceTop100BpsPercentile:     7.0,
		},
	}

	testCompressedFingerAndScores, genErr := pastel.ToCompressSignedDDAndFingerprints(fingerAndScores, []byte("testSignature"))
	assert.Nil(t, genErr)

	type nodeArg struct {
		address  string
		pastelID string
	}
	type args struct {
		task            *SenseRegisterTask
		nodes           []nodeArg
		isValidBurnTxID bool
		probeImgErr     error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				isValidBurnTxID: true,
				task: &SenseRegisterTask{
					Service: &Service{
						config: &Config{
							thumbnailSize: 224,
						},
					},
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				probeImgErr: nil,
			},
			wantErr: nil,
		},
		"probe-img-err": {
			args: args{
				isValidBurnTxID: true,
				task: &SenseRegisterTask{
					Service: &Service{
						config: &Config{
							thumbnailSize: 224,
						},
					},
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				probeImgErr: errors.New("test"),
			},
			wantErr: errors.New("send image"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			artworkFile, err := newTestImageFile()
			assert.NoError(t, err)

			//need to remove generate thumbnail file
			customProbeImageFunc := func(ctx context.Context, file *files.File) []byte {
				file.Remove()
				return testCompressedFingerAndScores
			}
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(true, nil)

			tc.args.task.Service.pastelClient = pastelClientMock

			newT := &common.WalletNodeTask{
				Task:      task.New(&state.Status{}),
				LogPrefix: logPrefix,
			}
			tc.args.task.WalletNodeTask = newT
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnRegisterSense().
				ListenOnProbeImage(customProbeImageFunc, tc.args.isValidBurnTxID, tc.args.probeImgErr)

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				newNode.SetValidBurnTxID(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))
				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.probeImage(context.Background())

			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				if err != nil {
					fmt.Println(err)
				}
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskSendSignedTicket(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}
	type args struct {
		task                   *SenseRegisterTask
		sendSignedTicketRet    string
		sendSignedTicketRetErr error
		nodes                  []nodeArg
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{

		"success": {
			args: args{
				sendSignedTicketRet: "test",
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID:           "testid",
						AppPastelIDPassphrase: "testpassphrase",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnRegisterSense().
				ListenOnSession(nil).
				ListenOnConnectTo(nil).
				ListenOnSessID("").
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnSendSignedTicket(tc.args.sendSignedTicketRet, tc.args.sendSignedTicketRetErr)

			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err := tc.args.task.sendSignedTicket(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				if err != nil {
					fmt.Println("err: ", err.Error())
				}
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskConnectToTopRankNodes(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}
	type args struct {
		task              *SenseRegisterTask
		nodes             []nodeArg
		masterNodesTopErr error
		masterNodes       pastel.MasterNodes
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
					pastel.MasterNode{Fee: 0.3, ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
					pastel.MasterNode{Fee: 0.4, ExtAddress: "127.0.0.1:4447", ExtKey: "4"},
				},
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID:           "testid",
						AppPastelIDPassphrase: "testpassphrase",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: nil,
		},
		"master-nodes-err": {
			args: args{
				masterNodesTopErr: errors.New("test"),
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: errors.New("masternode top"),
		},
		"insufficient-sns-err": {
			args: args{
				masterNodesTopErr: nil,
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "testid",
					},
					Service: &Service{
						config: &Config{NumberSuperNodes: 1},
					},
					ticket: &pastel.ActionTicket{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: errors.New("unable to find enough Supernodes"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			artworkFile, err := newTestImageFile()
			assert.NoError(t, err)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("test-signature"), nil).
				ListenOnGetBlockCount(100, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{}, nil).
				ListenOnMasterNodesTop(tc.args.masterNodes, tc.args.masterNodesTopErr)

			tc.args.task.Service.pastelClient = pastelClientMock

			//need to remove generate thumbnail file
			newT := &common.WalletNodeTask{
				Task:      task.New(&state.Status{}),
				LogPrefix: logPrefix,
			}
			tc.args.task.WalletNodeTask = newT
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnRegisterSense().
				ListenOnSession(nil).
				ListenOnConnectTo(nil).
				ListenOnSessID("").
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnClose(tc.wantErr)

			if tc.wantErr == nil {
				nodeClient.ListenOnMeshNodes(nil)
			}
			tc.args.task.Service.nodeClient = nodeClient

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.connectToTopRankNodes(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskGenerateDDAndFingerprintsIDs(t *testing.T) {
	t.Parallel()

	type args struct {
		task *SenseRegisterTask
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &SenseRegisterTask{
					Request: &SenseRegisterRequest{
						AppPastelID: "test-id",
					},
					Service: &Service{
						config: NewConfig(),
					},
					signatures: make([][]byte, 3),
					fingerprintAndScores: &pastel.DDAndFingerprints{
						Block:                      "Block",
						Principal:                  "Principal",
						DupeDetectionSystemVersion: "v1.0",

						IsLikelyDupe:     true,
						IsRareOnInternet: true,

						RarenessScores: &pastel.RarenessScores{
							CombinedRarenessScore:         0,
							XgboostPredictedRarenessScore: 0,
							NnPredictedRarenessScore:      0,
							OverallAverageRarenessScore:   0,
						},
						InternetRareness: &pastel.InternetRareness{
							MatchesFoundOnFirstPage: 0,
							NumberOfPagesOfResults:  0,
							URLOfFirstMatchInPage:   "",
						},

						OpenNSFWScore: 0.1,
						AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
							Drawings: 0.1,
							Hentai:   0.2,
							Neutral:  0.3,
							Porn:     0.4,
							Sexy:     0.5,
						},

						ImageFingerprintOfCandidateImageFile: []float32{1, 2, 3},
						FingerprintsStat: &pastel.FingerprintsStat{
							NumberOfFingerprintsRequiringFurtherTesting1: 1,
							NumberOfFingerprintsRequiringFurtherTesting2: 2,
							NumberOfFingerprintsRequiringFurtherTesting3: 3,
							NumberOfFingerprintsRequiringFurtherTesting4: 4,
							NumberOfFingerprintsRequiringFurtherTesting5: 5,
							NumberOfFingerprintsRequiringFurtherTesting6: 6,
							NumberOfFingerprintsOfSuspectedDupes:         7,
						},

						HashOfCandidateImageFile: "HashOfCandidateImageFile",
						PerceptualImageHashes: &pastel.PerceptualImageHashes{
							PDQHash:        "PdqHash",
							PerceptualHash: "PerceptualHash",
							AverageHash:    "AverageHash",
							DifferenceHash: "DifferenceHash",
							NeuralHash:     "NeuralhashHash",
						},
						PerceptualHashOverlapCount: 1,

						Maxes: &pastel.Maxes{
							PearsonMax:           1.0,
							SpearmanMax:          2.0,
							KendallMax:           3.0,
							HoeffdingMax:         4.0,
							MutualInformationMax: 5.0,
							HsicMax:              6.0,
							XgbimportanceMax:     7.0,
						},
						Percentile: &pastel.Percentile{
							PearsonTop1BpsPercentile:             1.0,
							SpearmanTop1BpsPercentile:            2.0,
							KendallTop1BpsPercentile:             3.0,
							HoeffdingTop10BpsPercentile:          4.0,
							MutualInformationTop100BpsPercentile: 5.0,
							HsicTop100BpsPercentile:              6.0,
							XgbimportanceTop100BpsPercentile:     7.0,
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			err := tc.args.task.generateDDAndFingerprintsIDs()
			if tc.wantErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
