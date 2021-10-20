package externaldupedetetion

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	stateMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/externaldupedetection/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNode(address, pastelID string) *node.Node {
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

func newTestImageFile() (*artwork.File, error) {
	imageStorage := artwork.NewStorage(fs.NewFileStorage(os.TempDir()))
	imgFile := imageStorage.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create storage file: %w", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	png.Encode(f, img)

	return imgFile, nil
}

func TestTaskRun(t *testing.T) {
	t.Parallel()
	t.Skip()

	type fields struct {
		Request *Request
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
			fields: fields{&Request{MaximumFee: 0.5, PastelID: "1", PastelIDPassphrase: "2"}},
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
			assert.NoError(t, err)
			testCase.args.fingerPrint = compressedFg

			t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
				var task *Task

				nodeClient := test.NewMockClient(t)
				nodeClient.
					ListenOnConnect("", testCase.args.returnErr).
					ListenOnExternalDupeDetection().
					ListenOnExternalDupeDetectionSession(testCase.args.returnErr).
					ListenOnExternalDupeDetectionConnectTo(testCase.args.returnErr).
					ListenOnExternalDupeDetectionSessID(testCase.args.primarySessID).
					ListenOnExternalDupeDetectionAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
					ListenOnDone().
					ListenOnExternalDupeDetectionUploadImage(nil).
					ListenOnExternalDupeDetectionSendSignedEDDTicket(1, nil).
					ListenOnClose(nil)

				counter := &struct {
					sync.Mutex
					val int
				}{}
				preburnCustomHandler := func() (string, error) {
					counter.Lock()
					defer counter.Unlock()
					counter.val++
					if counter.val == 1 {
						return "reg-art-txid", nil
					}
					return "", nil
				}
				nodeClient.ExternalDupeDetection.Mock.On(test.SendPreBurnedFeeTxIDMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.ExternalDupeDetection.Mock.On(test.SendPreBurnedFeeTxIDMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.ExternalDupeDetection.Mock.On(test.SendPreBurnedFeeTxIDMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.ExternalDupeDetection.Mock.On(test.SendPreBurnedFeeTxIDMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())

				//need to remove generate thumbnail file
				customProbeImageFunc := func(ctx context.Context, file *artwork.File) *pastel.FingerAndScores {
					file.Remove()
					return &pastel.FingerAndScores{ZstdCompressedFingerprint: testCase.args.fingerPrint}
				}
				nodeClient.ListenOnExternalDupeDetectionProbeImage(customProbeImageFunc, testCase.args.returnErr)

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
					ListenOnRegisterExDDTicket("art-act-txid", nil)

				rqClientMock := rqMock.NewMockClient(t)
				rqClientMock.ListenOnEncodeInfo(testCase.args.encodeInfoReturns, nil)
				rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
				rqClientMock.ListenOnConnect(testCase.args.connectErr)

				service := &Service{
					pastelClient: pastelClientMock.Client,
					nodeClient:   nodeClient.Client,
					rqClient:     rqClientMock,
					config:       NewConfig(),
				}

				taskClient := stateMock.NewMockTask(t)
				taskClient.
					ListenOnID(testCase.args.taskID).
					ListenOnUpdateStatus().
					ListenOnSetStatusNotifyFunc()

				Request := testCase.fields.Request
				Request.Image = artworkFile
				Request.MaximumFee = 100
				task = &Task{
					Task:    taskClient.Task,
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
					Request.PastelID,
					Request.PastelIDPassphrase,
					mock.Anything,
				)

				// //nodeClient mock assertion
				nodeClient.AssertExternalDupeDetectionAcceptedNodesCall(1, mock.Anything)
				nodeClient.AssertExternalDupeDetectionSessIDCall(testCase.numSessIDCall)
				nodeClient.AssertExternalDupeDetectionSessionCall(testCase.numSessionCall, mock.Anything, false)
				nodeClient.AssertExternalDupeDetectionConnectToCall(testCase.numConnectToCall, mock.Anything, mock.Anything, testCase.args.primarySessID)
				nodeClient.AssertExternalDupeDetectionProbeImageCall(testCase.numProbeImageCall, mock.Anything, mock.IsType(&artwork.File{}))
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
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionSession(testCase.args.returnErr).
				ListenOnExternalDupeDetectionConnectTo(testCase.args.returnErr).
				ListenOnExternalDupeDetectionSessID(testCase.args.primarySessID).
				ListenOnExternalDupeDetectionAcceptedNodes(testCase.args.pastelIDS, testCase.args.acceptNodeErr)

			nodes := node.List{}
			for _, n := range testCase.args.nodes {
				nodes.Add(node.NewNode(nodeClient.Client, n.address, n.pastelID))
			}

			service := &Service{
				config: NewConfig(),
			}

			task := &Task{Service: service, Request: &Request{}}
			got, err := task.meshNodes(testCase.args.ctx, nodes, testCase.args.primaryIndex)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, pullPastelAddressIDNodes(got))

			nodeClient.AssertExternalDupeDetectionAcceptedNodesCall(1, mock.Anything)
			nodeClient.AssertExternalDupeDetectionSessIDCall(testCase.numSessIDCall)
			nodeClient.AssertExternalDupeDetectionSessionCall(testCase.numSessionCall, mock.Anything, false)
			nodeClient.AssertExternalDupeDetectionConnectToCall(testCase.numConnectToCall, mock.Anything, testCase.args.primaryPastelID, testCase.args.primarySessID)
			nodeClient.Client.AssertExpectations(t)
			nodeClient.Connection.AssertExpectations(t)
		})
	}
}

func TestTaskIsSuitableStorageNetworkFee(t *testing.T) {
	t.Parallel()

	type fields struct {
		Request *Request
	}

	type args struct {
		ctx        context.Context
		networkFee float64
		returnErr  error
	}

	testCases := []struct {
		fields    fields
		args      args
		want      bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				Request: &Request{MaximumFee: 0.5},
			},
			args: args{
				ctx:        context.Background(),
				networkFee: 0.49,
			},
			want:      true,
			assertion: assert.NoError,
		},
		{
			fields: fields{
				Request: &Request{MaximumFee: 0.5},
			},
			args: args{
				ctx:        context.Background(),
				networkFee: 0.51,
			},
			want:      false,
			assertion: assert.NoError,
		},
		{
			args: args{
				ctx:       context.Background(),
				returnErr: fmt.Errorf("connection timeout"),
			},
			want:      false,
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create new mock service
			pastelClient := pastelMock.NewMockClient(t)
			pastelClient.ListenOnStorageNetworkFee(testCase.args.networkFee, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClient.Client,
			}

			task := &Task{
				Service: service,
				Request: testCase.fields.Request,
			}

			got, err := task.isSuitableStorageFee(testCase.args.ctx)
			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, got)

			//pastelClient mock assertion
			pastelClient.AssertExpectations(t)
			pastelClient.AssertStorageNetworkFeeCall(1, testCase.args.ctx)
		})
	}
}

func TestTaskPastelTopNodes(t *testing.T) {
	t.Parallel()

	type fields struct {
		Task    task.Task
		Request *Request
	}

	type args struct {
		ctx       context.Context
		returnMn  pastel.MasterNodes
		returnErr error
	}

	testCases := []struct {
		fields    fields
		args      args
		want      node.List
		assertion assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				Request: &Request{
					MaximumFee: 0.30,
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.1, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4444", "1"),
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			fields: fields{
				Request: &Request{
					MaximumFee: 0.3,
				},
			},
			args: args{
				ctx: context.Background(),
				returnMn: pastel.MasterNodes{
					pastel.MasterNode{Fee: 0.5, ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{Fee: 0.2, ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				returnErr: nil,
			},
			want: node.List{
				newTestNode("127.0.0.1:4445", "2"),
			},
			assertion: assert.NoError,
		}, {
			args: args{
				ctx:       context.Background(),
				returnMn:  nil,
				returnErr: fmt.Errorf("connection timeout"),
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
			pastelClient.ListenOnMasterNodesTop(testCase.args.returnMn, testCase.args.returnErr)
			service := &Service{
				pastelClient: pastelClient.Client,
			}

			task := &Task{
				Task:    testCase.fields.Task,
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
		Request *Request
	}

	service := &Service{}
	Request := &Request{}

	testCases := []struct {
		args args
		want *Task
	}{
		{
			args: args{
				service: service,
				Request: Request,
			},
			want: &Task{
				Task:    task.New(StatusTaskStarted),
				Service: service,
				Request: Request,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			task := NewTask(testCase.args.service, testCase.args.Request)
			assert.Equal(t, testCase.want.Service, task.Service)
			assert.Equal(t, testCase.want.Request, task.Request)
			assert.Equal(t, testCase.want.Status().SubStatus, task.Status().SubStatus)
		})
	}
}

func TestTaskCreateTicket(t *testing.T) {
	t.Parallel()

	type args struct {
		task *Task
	}

	testCases := map[string]struct {
		args    args
		want    *pastel.ExDDTicket
		wantErr error
	}{
		"fingerprint-error": {
			args: args{
				task: &Task{
					fingerprintSignature: []byte{},
					rqids:                []string{},
					datahash:             []byte{},
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyFingerprints,
			want:    nil,
		},
		"fingerprint-hash-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintSignature: []byte{},
					rqids:                []string{},
					datahash:             []byte{},
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyFingerprintsHash,
			want:    nil,
		},
		"fingerprint-signature-error": {
			args: args{
				task: &Task{
					fingerprint:      []byte{},
					fingerprintsHash: []byte{},
					rqids:            []string{},
					datahash:         []byte{},
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyFingerprintSignature,
			want:    nil,
		},
		"data-hash-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					rqids:                []string{},
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyDatahash,
			want:    nil,
		},
		"raptorQ-symbols-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					datahash:             []byte{},
					Request: &Request{
						PastelID: "test-id",
					},
					Service:              &Service{},
					fingerprintAndScores: &pastel.FingerAndScores{},
				},
			},
			wantErr: errEmptyRaptorQSymbols,
			want:    nil,
		},
		"success": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					datahash:             []byte{},
					fingerprintAndScores: &pastel.FingerAndScores{},
					rqids:                []string{},
					Request: &Request{
						PastelID:           "test-id",
						PastelIDPassphrase: "test-passphrase",
						ImageHash:          "test-image-hash",
						SendingAddress:     "test-sending-addeess",
						MaximumFee:         10,
						RequestIdentifier:  "test-request-identifier",
					},
					Service: &Service{},
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

			tc.want = &pastel.ExDDTicket{
				Version:             1,
				BlockHash:           tc.args.task.creatorBlockHash,
				BlockNum:            tc.args.task.creatorBlockHeight,
				Identifier:          tc.args.task.Request.RequestIdentifier,
				ImageHash:           tc.args.task.Request.ImageHash,
				MaximumFee:          tc.args.task.Request.MaximumFee,
				SendingAddress:      tc.args.task.Request.SendingAddress,
				EffectiveTotalFee:   float64(tc.args.task.detectionFee),
				SupernodesSignature: []map[string]string{},
			}
			err := tc.args.task.createEDDTicket(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tc.args.task.ticket)
				assert.Equal(t, *tc.want, *tc.args.task.ticket)
			}

		})
	}
}

func TestTaskGetBlock(t *testing.T) {
	t.Parallel()

	type args struct {
		task            *Task
		blockCountErr   error
		blockVerboseErr error
		blockNum        int32
		blockInfo       *pastel.GetBlockVerbose1Result
	}

	testCases := map[string]struct {
		args            args
		wantblockHash   string
		wantBlockHeight int
		wantErr         error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
			},
			wantBlockHeight: 10,
			wantErr:         nil,
		},
		"block-count-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
				blockCountErr: errors.New("block-count-err"),
			},
			wantBlockHeight: 10,
			wantErr:         errors.New("block-count-err"),
		},
		"block-verbose-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "test-id",
					},
					Service: &Service{},
				},
				blockNum: int32(10),
				blockInfo: &pastel.GetBlockVerbose1Result{
					Hash: "000000007465737468617368",
				},
				blockVerboseErr: errors.New("verbose-err"),
			},
			wantBlockHeight: 10,
			wantErr:         errors.New("verbose-err"),
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

			tc.wantblockHash = tc.args.blockInfo.Hash

			err := tc.args.task.getBlock(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantblockHash, tc.args.task.creatorBlockHash)
				assert.Equal(t, tc.wantBlockHeight, tc.args.task.creatorBlockHeight)
			}
		})
	}
}

func TestTaskEncodeFingerprint(t *testing.T) {
	type args struct {
		task                *Task
		fingerprint         []byte
		img                 *artwork.File
		signReturns         []byte
		findTicketIDReturns *pastel.IDTicket
		signErr             error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{},
				},
				img:         &artwork.File{},
				signReturns: []byte("test-signature"),
				fingerprint: []byte("test-fingerprint"),
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
			},
			wantErr: errors.New("decode image"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.args.signReturns, tc.args.signErr)
			pastelClientMock.ListenOnFindTicketByID(tc.args.findTicketIDReturns, nil)

			tc.args.task.Service.pastelClient = pastelClientMock
			fileStorageMock := storageMock.NewMockFileStorage()
			storage := artwork.NewStorage(fileStorageMock)

			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF).ListenOnName("test")

			file := artwork.NewFile(storage, "test-file")
			fileStorageMock.ListenOnOpen(fileMock, nil)

			err := tc.args.task.encodeFingerprint(context.Background(), tc.args.fingerprint, file)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestTaskSignTicket(t *testing.T) {
	type args struct {
		task        *Task
		signErr     error
		signReturns []byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ExDDTicket{},
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
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.args.signReturns, tc.args.task.creatorSignature)
			}
		})
	}
}

func TestWaitTxnValid(t *testing.T) {
	type args struct {
		task                            *Task
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
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
			pastelClientMock.ListenOnGetRawTransactionVerbose1(tc.args.getRawTransactionVerbose1Ret,
				tc.args.getRawTransactionVerbose1RetErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			ctx, cancel := context.WithCancel(context.Background())
			if tc.args.ctxDone {
				cancel()
			}

			err := tc.args.task.waitTxidValid(ctx, "txid", 5, time.Second, 500*time.Millisecond)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			cancel()
		})
	}
}

func TestTaskPreburnedDetectionFee(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}

	type args struct {
		task                  *Task
		sendFromAddressRetErr error
		burnTxnIDRet          string
		preBurnedFeeRetTxID   string
		preBurnedFeeRetErr    error
		nodes                 []nodeArg
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"reg-fee-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.ExDDTicket{},
				},
			},
			wantErr: errors.New("registration fee"),
		},
		"send-from-address-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:       &pastel.ExDDTicket{},
					detectionFee: 40,
				},
				sendFromAddressRetErr: errors.New("test"),
			},
			wantErr: errors.New("failed to burn"),
		},
		"send-preburnt-fee-err": {
			args: args{
				preBurnedFeeRetErr: errors.New("test"),
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:       &pastel.ExDDTicket{},
					detectionFee: 40,
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				sendFromAddressRetErr: nil,
				burnTxnIDRet:          "test-id",
			},
			wantErr: errors.New("failed to send pre-burn-txid"),
		},
		"regEDDTxId-empty-err": {
			args: args{
				preBurnedFeeRetTxID: "test-id",
				preBurnedFeeRetErr:  nil,
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:       &pastel.ExDDTicket{},
					detectionFee: 40,
				},
				sendFromAddressRetErr: nil,
				burnTxnIDRet:          "test-id",
			},
			wantErr: errors.New("regEDDTxID is empty"),
		},
		"success": {
			args: args{
				preBurnedFeeRetTxID: "test-id",
				preBurnedFeeRetErr:  nil,
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:       &pastel.ExDDTicket{},
					detectionFee: 40,
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				sendFromAddressRetErr: nil,
				burnTxnIDRet:          "test-id",
			},
			wantErr: nil,
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSendFromAddress(tc.args.burnTxnIDRet, tc.args.sendFromAddressRetErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionSession(nil).
				ListenOnExternalDupeDetectionConnectTo(nil).
				ListenOnExternalDupeDetectionSessID("").
				ListenOnExternalDupeDetectionAcceptedNodes([]string{}, nil).
				ListenOnExternalDupeDetectionSendPreBurnedFeeEDDTxID(tc.args.preBurnedFeeRetTxID, tc.args.preBurnedFeeRetErr)

			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err := tc.args.task.preburnedDetectionFee(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskUploadImage(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}
	type args struct {
		task                *Task
		nodes               []nodeArg
		findTicketIDReturns *pastel.IDTicket
		uploadImageErr      error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
			},
			wantErr: nil,
		},
		"upload-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				uploadImageErr: errors.New("test"),
			},
			wantErr: errors.New("upload encoded image failed "),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			artworkFile, err := newTestImageFile()
			assert.NoError(t, err)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("test-signature"), nil)
			pastelClientMock.ListenOnFindTicketByID(tc.args.findTicketIDReturns, nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.args.task.Task = task.New(&state.Status{})
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionUploadImage(tc.args.uploadImageErr)

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.uploadImage(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskProbeImage(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}
	type args struct {
		task        *Task
		nodes       []nodeArg
		probeImgErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{
							thumbnailSize: 224,
						},
					},
					Request: &Request{
						PastelID: "testid",
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: nil,
		},
		"probe-img-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{
							thumbnailSize: 224,
						},
					},
					Request: &Request{
						PastelID: "testid",
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				probeImgErr: errors.New("test"),
			},
			wantErr: errors.New("failed to probe image"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			artworkFile, err := newTestImageFile()
			assert.NoError(t, err)

			//need to remove generate thumbnail file
			customProbeImageFunc := func(ctx context.Context, file *artwork.File) *pastel.FingerAndScores {
				file.Remove()
				finger := pastel.Fingerprint{0.1}
				compressedFinger, _ := zstd.Compress(nil, finger.Bytes())
				fingerAndScore := pastel.FingerAndScores{
					ZstdCompressedFingerprint: compressedFinger,
				}
				return &fingerAndScore
			}
			tc.args.task.Task = task.New(&state.Status{})
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionProbeImage(customProbeImageFunc, tc.args.probeImgErr)

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second, &alts.SecInfo{}))
				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.probeImage(context.Background())
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				fmt.Printf("%v", err)
				assert.NoError(t, err)
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
		task                   *Task
		sendSignedTicketRet    int64
		sendSignedTicketRetErr error
		preBurnedFeeRetTxID    string
		preBurnedFeeRetErr     error
		nodes                  []nodeArg
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"upload-err": {
			args: args{
				sendSignedTicketRetErr: errors.New("test"),
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: errors.New("upload signed ticket"),
		},
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: nil,
		},
		"reg-fee-err": {
			args: args{
				task: &Task{
					Request: &Request{
						PastelID:   "testid",
						MaximumFee: -1,
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: errors.New("detection fee"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionSession(nil).
				ListenOnExternalDupeDetectionConnectTo(nil).
				ListenOnExternalDupeDetectionSessID("").
				ListenOnExternalDupeDetectionAcceptedNodes([]string{}, nil).
				ListenOnExternalDupeDetectionSendSignedEDDTicket(tc.args.sendSignedTicketRet, tc.args.sendSignedTicketRetErr).
				ListenOnExternalDupeDetectionSendPreBurnedFeeEDDTxID(tc.args.preBurnedFeeRetTxID, tc.args.preBurnedFeeRetErr)

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
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
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
		task              *Task
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
				task: &Task{
					Request: &Request{
						PastelID:   "testid",
						MaximumFee: 1,
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
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
				task: &Task{
					Request: &Request{
						PastelID: "testid",
					},
					Service: &Service{
						config: &Config{NumberSuperNodes: 1},
					},
					ticket:          &pastel.ExDDTicket{},
					rqSymbolIDFiles: map[string][]byte{},
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
				ListenOnMasterNodesTop(tc.args.masterNodes, tc.args.masterNodesTopErr)

			tc.args.task.Service.pastelClient = pastelClientMock

			//need to remove generate thumbnail file
			tc.args.task.Task = task.New(&state.Status{})
			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnExternalDupeDetection().
				ListenOnExternalDupeDetectionSession(nil).
				ListenOnExternalDupeDetectionConnectTo(nil).
				ListenOnExternalDupeDetectionSessID("").
				ListenOnExternalDupeDetectionAcceptedNodes([]string{}, nil).
				ListenOnClose(tc.wantErr)
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
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
