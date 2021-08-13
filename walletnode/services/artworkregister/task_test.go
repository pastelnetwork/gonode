package artworkregister

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
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
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func appendStr(b []byte, s string) []byte {
	b = append(b, []byte(s)...)
	b = append(b, '\n')
	return b
}

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
			fields: fields{&Request{MaximumFee: 0.5, ArtistPastelID: "1", ArtistPastelIDPassphrase: "2"}},
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
						"test-file": rqnode.RawSymbolIDFile{
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
				var task *Task

				nodeClient := test.NewMockClient(t)
				nodeClient.
					ListenOnConnect("", testCase.args.returnErr).
					ListenOnRegisterArtwork().
					ListenOnSession(testCase.args.returnErr).
					ListenOnConnectTo(testCase.args.returnErr).
					ListenOnSessID(testCase.args.primarySessID).
					ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
					ListenOnDone().
					ListenOnUploadImageWithThumbnail([]byte("preview-hash"), []byte("medium-hash"), []byte("small-hash"), nil).
					ListenOnSendSignedTicket(1, nil).
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
				nodeClient.RegisterArtwork.Mock.On(test.SendPreBurntFeeTxidMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.RegisterArtwork.Mock.On(test.SendPreBurntFeeTxidMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.RegisterArtwork.Mock.On(test.SendPreBurntFeeTxidMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())
				nodeClient.RegisterArtwork.Mock.On(test.SendPreBurntFeeTxidMethod, mock.Anything, mock.AnythingOfType("string")).Once().Return(preburnCustomHandler())

				//need to remove generate thumbnail file
				customProbeImageFunc := func(ctx context.Context, file *artwork.File) *pastel.FingerAndScores {
					file.Remove()
					return &pastel.FingerAndScores{ZstdCompressedFingerprint: testCase.args.fingerPrint}
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
					Request.ArtistPastelID,
					Request.ArtistPastelIDPassphrase,
					mock.Anything,
				)

				// //nodeClient mock assertion
				nodeClient.AssertAcceptedNodesCall(1, mock.Anything)
				nodeClient.AssertSessIDCall(testCase.numSessIDCall)
				nodeClient.AssertSessionCall(testCase.numSessionCall, mock.Anything, false)
				nodeClient.AssertConnectToCall(testCase.numConnectToCall, mock.Anything, mock.Anything, testCase.args.primarySessID)
				nodeClient.AssertProbeImageCall(testCase.numProbeImageCall, mock.Anything, mock.IsType(&artwork.File{}))
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
				ListenOnRegisterArtwork().
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

			task := &Task{Service: service}
			got, err := task.meshNodes(testCase.args.ctx, nodes, testCase.args.primaryIndex)

			testCase.assertion(t, err)
			assert.Equal(t, testCase.want, pullPastelAddressIDNodes(got))

			nodeClient.AssertAcceptedNodesCall(1, mock.Anything)
			nodeClient.AssertSessIDCall(testCase.numSessIDCall)
			nodeClient.AssertSessionCall(testCase.numSessionCall, mock.Anything, false)
			nodeClient.AssertConnectToCall(testCase.numConnectToCall, mock.Anything, testCase.args.primaryPastelID, testCase.args.primarySessID)
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
		want    *pastel.NFTTicket
		wantErr error
	}{
		"fingerprint-error": {
			args: args{
				task: &Task{
					fingerprintSignature: []byte{},
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					previewHash:          []byte{},
					rqids:                []string{},
					datahash:             []byte{},
					Request: &Request{
						ArtistPastelID: "test-id",
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
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					previewHash:          []byte{},
					rqids:                []string{},
					datahash:             []byte{},
					Request: &Request{
						ArtistPastelID: "test-id",
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
					fingerprint:         []byte{},
					fingerprintsHash:    []byte{},
					previewHash:         []byte{},
					smallThumbnailHash:  []byte{},
					mediumThumbnailHash: []byte{},
					rqids:               []string{},
					datahash:            []byte{},
					Request: &Request{
						ArtistPastelID: "test-id",
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
					previewHash:          []byte{},
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					rqids:                []string{},
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyDatahash,
			want:    nil,
		},
		"preview-hash-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					datahash:             []byte{},
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					rqids:                []string{},
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyPreviewHash,
			want:    nil,
		},
		"medium-thumbnail-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					previewHash:          []byte{},
					datahash:             []byte{},
					smallThumbnailHash:   []byte{},
					rqids:                []string{},
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptyMediumThumbnailHash,
			want:    nil,
		},
		"small-thumbnail-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					mediumThumbnailHash:  []byte{},
					previewHash:          []byte{},
					datahash:             []byte{},
					rqids:                []string{},
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
			},
			wantErr: errEmptySmallThumbnailHash,
			want:    nil,
		},
		"raptorQ-symbols-error": {
			args: args{
				task: &Task{
					fingerprint:          []byte{},
					fingerprintsHash:     []byte{},
					fingerprintSignature: []byte{},
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					previewHash:          []byte{},
					datahash:             []byte{},
					Request: &Request{
						ArtistPastelID: "test-id",
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
					smallThumbnailHash:   []byte{},
					mediumThumbnailHash:  []byte{},
					previewHash:          []byte{},
					datahash:             []byte{},
					fingerprintAndScores: &pastel.FingerAndScores{},
					rqids:                []string{},
					Request: &Request{
						ArtistPastelID: "test-id",
						ArtistName:     "test-name",
						IssuedCopies:   10,
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

			pastelID := base58.Decode(tc.args.task.Request.ArtistPastelID)

			tc.want = &pastel.NFTTicket{
				Version:  1,
				Author:   string(pastelID),
				BlockNum: tc.args.task.creatorBlockHeight,
				Copies:   tc.args.task.Request.IssuedCopies,
				AppTicketData: pastel.AppTicket{
					AuthorPastelID:             tc.args.task.Request.ArtistPastelID,
					BlockTxID:                  tc.args.task.blockTxID,
					BlockNum:                   0,
					CreatorName:                tc.args.task.Request.ArtistName,
					CreatorWebsite:             safeString(tc.args.task.Request.ArtistWebsiteURL),
					CreatorWrittenStatement:    safeString(tc.args.task.Request.Description),
					NFTCreationVideoYoutubeURL: safeString(tc.args.task.Request.YoutubeURL),
					NFTKeywordSet:              safeString(tc.args.task.Request.Keywords),
					TotalCopies:                tc.args.task.Request.IssuedCopies,
					PreviewHash:                tc.args.task.previewHash,
					Thumbnail1Hash:             tc.args.task.mediumThumbnailHash,
					Thumbnail2Hash:             tc.args.task.smallThumbnailHash,
					DataHash:                   tc.args.task.datahash,
					FingerprintsHash:           tc.args.task.fingerprintsHash,
					FingerprintsSignature:      tc.args.task.fingerprintSignature,
					RQIDs:                      tc.args.task.rqids,
					RQOti:                      tc.args.task.rqEncodeParams.Oti,
				},
			}

			err := tc.args.task.createArtTicket(context.Background())
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
		task            *Task
		blockCountErr   error
		blockVerboseErr error
		blockNum        int32
		blockInfo       *pastel.GetBlockVerbose1Result
	}

	testCases := map[string]struct {
		args                  args
		wantArtistblockHash   []byte
		wantArtistBlockHeight int
		wantErr               error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "test-id",
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
				task: &Task{
					Request: &Request{
						ArtistPastelID: "test-id",
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
				task: &Task{
					Request: &Request{
						ArtistPastelID: "test-id",
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

			blockhash, err := hex.DecodeString(tc.args.blockInfo.Hash)
			assert.Nil(t, err)
			tc.wantArtistblockHash = blockhash

			err = tc.args.task.getBlock(context.Background())
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

func TestTaskConvertToSymbolIdFile(t *testing.T) {
	t.Parallel()

	testErrStr := "test-err"
	type args struct {
		task    *Task
		signErr error
		inFile  rqnode.RawSymbolIDFile
	}

	testCases := map[string]struct {
		args        args
		wantContent []byte
		wantErr     error
		wantSign    []byte
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
				inFile: rqnode.RawSymbolIDFile{
					ID:                uuid.New().String(),
					SymbolIdentifiers: []string{"test-s1", "test-s2"},
					BlockHash:         "test-block-hash",
					PastelID:          "test-pastel-id",
				},
			},
			wantSign: []byte("test-signature"),
			wantErr:  nil,
		},
		"error": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "test-id",
					},
					Service: &Service{},
				},
				signErr: errors.New(testErrStr),
			},
			wantErr: errors.New(testErrStr),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			var rawContent []byte
			rawContent = appendStr(rawContent, tc.args.inFile.ID)
			rawContent = appendStr(rawContent, tc.args.inFile.BlockHash)
			rawContent = appendStr(rawContent, tc.args.inFile.PastelID)
			for _, id := range tc.args.inFile.SymbolIdentifiers {
				rawContent = appendStr(rawContent, id)
			}
			rawContent = append(rawContent, tc.wantSign...)
			fmt.Printf("%s\n", string(rawContent))

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign(tc.wantSign, tc.args.signErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			_, content, err := tc.args.task.convertToSymbolIDFile(context.Background(), tc.args.inFile)
			fmt.Printf("%x\n", content)

			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), testErrStr))
			} else {
				assert.Nil(t, err)
				compressContent, zstdErr := zstd.CompressLevel(nil, rawContent, 22)
				fmt.Printf("%x\n", compressContent)
				assert.Nil(t, zstdErr)
				assert.True(t, bytes.Equal(content, compressContent))
			}
		})
	}
}

func TestTaskGenRQIdentifiersFiles(t *testing.T) {
	artworkFile, err := newTestImageFile()
	assert.NoError(t, err)

	type args struct {
		task                *Task
		connectErr          error
		encodeInfoReturns   *rqnode.EncodeInfo
		findTicketIDReturns *pastel.IDTicket
		encodeInfoErr       error
		readErr             error
		signErr             error
	}
	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Task: task.New(StatusTaskStarted),
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					imageEncodedWithFingerprints: artworkFile,
				},
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
						"test-file": rqnode.RawSymbolIDFile{
							ID:                uuid.New().String(),
							SymbolIdentifiers: []string{"test-s1, test-s2"},
							BlockHash:         "test-block-hash",
							PastelID:          "test-pastel-id",
						},
					},
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				readErr: io.EOF,
			},
			wantErr: nil,
		},
		"connect-error": {
			args: args{
				task: &Task{
					Task: task.New(StatusTaskStarted),
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					imageEncodedWithFingerprints: artworkFile,
				},
				readErr: io.EOF,
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: make(map[string]rqnode.RawSymbolIDFile),
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				connectErr: errors.New("test-err"),
			},
			wantErr: errors.New("test-err"),
		},

		"encode-info-error": {
			args: args{
				task: &Task{
					Task: task.New(StatusTaskStarted),
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					imageEncodedWithFingerprints: artworkFile,
				},
				readErr: io.EOF,
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: make(map[string]rqnode.RawSymbolIDFile),
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				encodeInfoErr: errors.New("test-err"),
			},
			wantErr: errors.New("test-err"),
		},
		"read-error": {
			args: args{
				task: &Task{
					Task: task.New(StatusTaskStarted),
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					imageEncodedWithFingerprints: artworkFile,
				},
				readErr: errors.New("test-err"),
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: make(map[string]rqnode.RawSymbolIDFile),
				},
			},
			wantErr: errors.New("read image"),
		},
		"sign-err": {
			args: args{
				task: &Task{
					Task: task.New(StatusTaskStarted),
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					imageEncodedWithFingerprints: artworkFile,
				},
				encodeInfoReturns: &rqnode.EncodeInfo{
					SymbolIDFiles: map[string]rqnode.RawSymbolIDFile{
						"test-file": rqnode.RawSymbolIDFile{
							ID:                uuid.New().String(),
							SymbolIdentifiers: []string{"test-s1, test-s2"},
							BlockHash:         "test-block-hash",
							PastelID:          "test-pastel-id",
						},
					},
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				readErr: io.EOF,
				signErr: errors.New("tes-err"),
			},
			wantErr: errors.New("rqids file"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("test-signature"), tc.args.signErr)
			pastelClientMock.ListenOnFindTicketByID(tc.args.findTicketIDReturns, nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(tc.args.encodeInfoReturns, tc.args.encodeInfoErr)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr)

			tc.args.task.Service.rqClient = rqClientMock

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, tc.args.readErr)

			storage := artwork.NewStorage(fsMock)
			tc.args.task.imageEncodedWithFingerprints = artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, nil)

			tc.args.task.Request.Image = artwork.NewFile(storage, "test")

			err := tc.args.task.genRQIdentifiersFiles(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				// TODO: fix this when the return error is define correctly
				// assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				fmt.Printf("%s", err)
				assert.Nil(t, err)
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
						ArtistPastelID: "testid",
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
			wantErr: errors.New("failed to decode image"),
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.NFTTicket{},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
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
						ArtistPastelID: "testid",
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
						ArtistPastelID: "testid",
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
						ArtistPastelID: "testid",
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
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
			cancel()
		})
	}

}

func TestTaskPreburntRegistrationFee(t *testing.T) {
	type nodeArg struct {
		address  string
		pastelID string
	}

	type args struct {
		task                  *Task
		sendFromAddressRetErr error
		burnTxnIDRet          string
		preBurntFeeRetTxID    string
		preBurntFeeRetErr     error
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket: &pastel.NFTTicket{},
				},
			},
			wantErr: errors.New("registration fee"),
		},
		"send-from-address-err": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					registrationFee: 40,
				},
				sendFromAddressRetErr: errors.New("test"),
			},
			wantErr: errors.New("failed to burn"),
		},
		"send-preburnt-fee-err": {
			args: args{
				preBurntFeeRetErr: errors.New("test"),
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					registrationFee: 40,
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
		"regArtTxId-empty-err": {
			args: args{
				preBurntFeeRetTxID: "test-id",
				preBurntFeeRetErr:  nil,
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					registrationFee: 40,
				},
				sendFromAddressRetErr: nil,
				burnTxnIDRet:          "test-id",
			},
			wantErr: errors.New("regNFTTxid is empty"),
		},
		"success": {
			args: args{
				preBurntFeeRetTxID: "test-id",
				preBurntFeeRetErr:  nil,
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					registrationFee: 40,
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
				ListenOnRegisterArtwork().
				ListenOnSession(nil).
				ListenOnConnectTo(nil).
				ListenOnSessID("").
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnSendPreBurntFeeTxID(tc.args.preBurntFeeRetTxID, tc.args.preBurntFeeRetErr)

			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err := tc.args.task.preburntRegistrationFee(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
		task                    *Task
		nodes                   []nodeArg
		findTicketIDReturns     *pastel.IDTicket
		previewHash             []byte
		mediumHash              []byte
		smallHash               []byte
		uploadImageThumbnailErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Request: &Request{
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
				findTicketIDReturns: &pastel.IDTicket{
					TXID: "test-txid",
				},
				uploadImageThumbnailErr: errors.New("test"),
			},
			wantErr: errors.New("upload encoded image "),
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
				ListenOnRegisterArtwork().
				ListenOnUploadImageWithThumbnail(tc.args.previewHash, tc.args.mediumHash,
					tc.args.smallHash, tc.args.uploadImageThumbnailErr)

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.uploadImage(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
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
						ArtistPastelID: "testid",
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
					},
					ticket:          &pastel.NFTTicket{},
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
				ListenOnRegisterArtwork().
				ListenOnProbeImage(customProbeImageFunc, tc.args.probeImgErr)

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second))
				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err = tc.args.task.probeImage(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				fmt.Printf("%v", err)
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
		task                   *Task
		sendSignedTicketRet    int64
		sendSignedTicketRetErr error
		preBurntFeeRetTxID     string
		preBurntFeeRetErr      error
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
						MaximumFee:     -1,
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
					rqSymbolIDFiles: map[string][]byte{},
				},
				nodes: []nodeArg{
					{"127.0.0.1", "1"},
					{"127.0.0.2", "2"},
				},
			},
			wantErr: errors.New("registration fee"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", nil).
				ListenOnRegisterArtwork().
				ListenOnSession(nil).
				ListenOnConnectTo(nil).
				ListenOnSessID("").
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnSendSignedTicket(tc.args.sendSignedTicketRet, tc.args.sendSignedTicketRetErr).
				ListenOnSendPreBurntFeeTxID(tc.args.preBurntFeeRetTxID, tc.args.preBurntFeeRetErr)

			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second))

				nodes.Add(newNode)
			}
			tc.args.task.nodes = nodes

			err := tc.args.task.sendSignedTicket(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
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
						ArtistPastelID: "testid",
						MaximumFee:     1,
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{},
					},
					ticket:          &pastel.NFTTicket{},
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
						ArtistPastelID: "testid",
					},
					Service: &Service{
						config: &Config{NumberSuperNodes: 1},
					},
					ticket:          &pastel.NFTTicket{},
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
				ListenOnRegisterArtwork().
				ListenOnSession(nil).
				ListenOnConnectTo(nil).
				ListenOnSessID("").
				ListenOnAcceptedNodes([]string{}, nil).
				ListenOnClose(tc.wantErr)
			tc.args.task.Service.nodeClient = nodeClient

			tc.args.task.Request.Image = artworkFile
			nodes := node.List{}
			for _, n := range tc.args.nodes {
				newNode := node.NewNode(nodeClient.Client, n.address, n.pastelID)
				newNode.SetPrimary(true)
				assert.Nil(t, newNode.Connect(context.Background(), 1*time.Second))

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
