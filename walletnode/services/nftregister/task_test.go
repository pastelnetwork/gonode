package nftregister

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	testCreatorPastelID = "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW"
	testPqID            = `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`
)

func float64SliceToBytes(floats []float64) []byte {
	buf := new(bytes.Buffer)

	for _, f := range floats {
		err := binary.Write(buf, binary.LittleEndian, f)
		if err != nil {
			fmt.Println("Error converting float64 to bytes:", err)
			return nil
		}
	}

	return buf.Bytes()
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
	type fields struct {
		Request *NftRegistrationRequest
	}

	type args struct {
		taskID            string
		ctx               context.Context
		networkFee        float64
		masterNodes       pastel.MasterNodes
		getTopMNsReply    *pb.GetTopMNsReply
		primarySessID     string
		pastelIDS         []string
		fingerPrint       []byte
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
				Request: &NftRegistrationRequest{
					MaximumFee:                0.5,
					CreatorPastelID:           testCreatorPastelID,
					CreatorPastelIDPassphrase: "passphrase",
				},
			},
			args: args{
				taskID:     "1",
				ctx:        context.Background(),
				networkFee: 0.4,
				masterNodes: pastel.MasterNodes{
					pastel.MasterNode{ExtAddress: "127.0.0.1:4446", ExtKey: "3"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4444", ExtKey: "1"},
					pastel.MasterNode{ExtAddress: "127.0.0.1:4445", ExtKey: "2"},
				},
				getTopMNsReply: &pb.GetTopMNsReply{MnTopList: []string{"127.0.0.1:4444", "127.0.0.1:4445", "127.0.0.1:4446", "127.0.0.1:4447", "127.0.0.1:4448",
					"127.0.0.1:4449", "127.0.0.0:4440", "127.0.0.0:4441", "127.0.0.0:4442", "127.0.0.1:4443"}},
				primarySessID: "sesid1",
				pastelIDS:     []string{"3", "1"},
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
					EncoderParam: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
				},
			},
		},*/

		"failure": {
			wantErr: errors.New("test"),
			fields: fields{
				Request: &NftRegistrationRequest{
					MaximumFee:                0.5,
					CreatorPastelID:           testCreatorPastelID,
					CreatorPastelIDPassphrase: "passphrase",
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
			t.Parallel()
			nftFile, err := newTestImageFile()
			assert.NoError(t, err)

			// prepare task
			fg := []float64{0.1, 0, 2}
			compressedFg, err := utils.Compress(float64SliceToBytes(fg), 4)
			assert.Nil(t, err)
			testCase.args.fingerPrint = compressedFg

			nodeClient := test.NewMockClient(t)
			nodeClient.
				ListenOnConnect("", testCase.args.returnErr).
				ListenOnRegisterNft().
				ListenOnSession(testCase.args.returnErr).
				ListenOnConnectTo(testCase.args.returnErr).
				ListenOnSessID(testCase.args.primarySessID).
				ListenOnAcceptedNodes(testCase.args.pastelIDS, testCase.args.returnErr).
				ListenOnRegisterGetDupeDetectionDBHash("", nil).
				ListenOnNFTGetDDServerStats(&pb.DDServerStatsReply{WaitingInQueue: 0, MaxConcurrent: 1}, nil).
				ListenOnDone().
				ListenOnUploadImageWithThumbnail([]byte("preview-hash"), []byte("medium-hash"), []byte("small-hash"), nil).
				ListenOnSendSignedTicket(1, nil).
				ListenOnClose(nil)

			nodeClient.ConnectionInterface.On("RegisterNft").Return(nodeClient.RegisterNftInterface)
			nodeClient.RegisterNftInterface.Mock.On(test.SendPreBurntFeeTxidMethod, mock.Anything, mock.AnythingOfType("string")).Return("100", nil)

			nodeClient.RegisterNftInterface.On("GetTopMNs", mock.Anything, mock.Anything).Return(testCase.args.getTopMNsReply, nil)
			nodeClient.RegisterNftInterface.On("MeshNodes", mock.Anything, mock.Anything).Return(nil)

			ddData := &pastel.DDAndFingerprints{
				BlockHash:   "BlockHash",
				BlockHeight: "BlockHeight",

				TimestampOfRequest: "Timestamp",
				SubmitterPastelID:  "PastelID",
				SN1PastelID:        "SN1PastelID",
				SN2PastelID:        "SN2PastelID",
				SN3PastelID:        "SN3PastelID",

				IsOpenAPIRequest: false,

				DupeDetectionSystemVersion: "v1.0",

				IsLikelyDupe:     true,
				IsRareOnInternet: true,

				OverallRarenessScore: 0.5,

				PctOfTop10MostSimilarWithDupeProbAbove25pct: 12.0,
				PctOfTop10MostSimilarWithDupeProbAbove33pct: 12.0,
				PctOfTop10MostSimilarWithDupeProbAbove50pct: 12.0,

				RarenessScoresTableJSONCompressedB64: "RarenessScoresTableJSONCompressedB64",

				InternetRareness: &pastel.InternetRareness{
					RareOnInternetSummaryTableAsJSONCompressedB64:    "RareOnInternetSummaryTableAsJSONCompressedB64",
					RareOnInternetGraphJSONCompressedB64:             "RareOnInternetGraphJSONCompressedB64",
					AlternativeRareOnInternetDictAsJSONCompressedB64: "AlternativeRareOnInternetDictAsJSONCompressedB64",
					MinNumberOfExactMatchesInPage:                    4,
					EarliestAvailableDateOfInternetResults:           "EarliestAvailableDateOfInternetResults",
				},

				OpenNSFWScore: 0.1,
				AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
					Drawings: 0.1,
					Hentai:   0.2,
					Neutral:  0.3,
					Porn:     0.4,
					Sexy:     0.5,
				},

				ImageFingerprintOfCandidateImageFile: []float64{1, 2, 3},

				HashOfCandidateImageFile: "HashOfCandidateImageFile",
			}
			compressed, err := pastel.ToCompressSignedDDAndFingerprints(context.Background(), ddData, []byte("signature"))
			assert.Nil(t, err)

			nodeClient.ListenOnProbeImage(compressed, true, false, testCase.args.returnErr)
			nodeClient.RegisterNftInterface.On("SendRegMetadata", mock.Anything, mock.Anything).Return(nil)
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.
				ListenOnStorageNetworkFee(testCase.args.networkFee, testCase.args.returnErr).
				ListenOnMasterNodesTop(testCase.args.masterNodes, testCase.args.returnErr).
				ListenOnSign([]byte(testCase.args.signature), testCase.args.returnErr).
				ListenOnGetBlockCount(100, nil).
				ListenOnGetBlockVerbose1(&pastel.GetBlockVerbose1Result{Hash: "abc123", Height: 100}, nil).
				ListenOnFindTicketByID(&pastel.IDTicket{IDTicketProp: pastel.IDTicketProp{PqKey: testPqID}}, nil).
				ListenOnSendFromAddress("pre-burnt-txid", nil).
				ListenOnGetRawTransactionVerbose1(&pastel.GetRawTransactionVerbose1Result{Confirmations: 12}, nil).
				ListenOnRegisterNFTTicket("nft-act-txid", nil).
				ListenOnVerify(true, nil).ListenOnGetBalance(100, nil).
				ListenOnRegisterActTicket("txid-act", nil).
				ListenOnGetRawMempoolMethod([]string{}, nil).
				ListenOnGetInactiveActionTickets(pastel.ActTickets{}, nil).
				ListenOnGetInactiveNFTTickets(pastel.RegTickets{}, nil)

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(testCase.args.encodeInfoReturns, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(testCase.args.connectErr)

			service := NewService(NewConfig(), pastelClientMock, nodeClient, nil, nil, nil, nil)
			service.rqClient = rqClientMock
			service.config.WaitTxnValidInterval = 1

			go service.Run(testCase.args.ctx)

			Request := testCase.fields.Request
			Request.Image = nftFile
			Request.MaximumFee = 50
			task := NewNFTRegistrationTask(service, Request)
			task.skipPrimaryNodeTxidVerify = true
			task.MaxRetries = 0
			//create context with timeout to automatically end process after 5 sec
			ctx, cancel := context.WithTimeout(testCase.args.ctx, 6*time.Second)
			defer cancel()

			err = task.Run(ctx)
			if testCase.wantErr != nil {
				assert.True(t, task.Status().IsFailure())
			} else {
				assert.True(t, task.Status().Is(common.StatusTaskCompleted))
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskCreateTicket(t *testing.T) {
	type args struct {
		task *NftRegistrationTask
	}

	issuedCopies := 10

	testCases := map[string]struct {
		args    args
		want    *pastel.NFTTicket
		wantErr error
	}{
		"fingerprint-error": {
			args: args{
				task: &NftRegistrationTask{
					dataHash: []byte{1, 2, 3},
					Request: &NftRegistrationRequest{
						CreatorPastelID: "test-id",
					},
					service: &NftRegistrationService{
						config: NewConfig(),
					},
					FingerprintsHandler: &mixins.FingerprintsHandler{},
					RqHandler: &mixins.RQHandler{
						RQIDs:          []string{"a"},
						RQIDsFile:      []byte{1, 2, 3},
						RQEncodeParams: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
					},
					ImageHandler: &mixins.NftImageHandler{
						PreviewHash:         []byte{1, 2, 3},
						MediumThumbnailHash: []byte{1, 2, 3},
						SmallThumbnailHash:  []byte{1, 2, 3},
					},
				},
			},
			wantErr: common.ErrEmptyFingerprints,
			want:    nil,
		},
		"data-hash-error": {
			args: args{
				task: &NftRegistrationTask{
					Request: &NftRegistrationRequest{
						CreatorPastelID: "test-id",
					},
					service: &NftRegistrationService{
						config: NewConfig(),
					},
					FingerprintsHandler: &mixins.FingerprintsHandler{
						DDAndFingerprintsIDs: []string{"a"},
						DDAndFpFile:          []byte{1, 2, 3},
						SNsSignatures:        [][]byte{[]byte{1, 2, 3}},
					},
					RqHandler: &mixins.RQHandler{
						RQIDs:          []string{"a"},
						RQIDsFile:      []byte{1, 2, 3},
						RQEncodeParams: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
					},
					ImageHandler: &mixins.NftImageHandler{
						PreviewHash:         []byte{1, 2, 3},
						MediumThumbnailHash: []byte{1, 2, 3},
						SmallThumbnailHash:  []byte{1, 2, 3},
					},
				},
			},
			wantErr: common.ErrEmptyDatahash,
			want:    nil,
		},
		"preview-hash-error": {
			args: args{
				task: &NftRegistrationTask{
					dataHash: []byte{1, 2},
					Request: &NftRegistrationRequest{
						CreatorPastelID: "test-id",
					},
					service: &NftRegistrationService{
						config: NewConfig(),
					},
					FingerprintsHandler: &mixins.FingerprintsHandler{
						DDAndFingerprintsIDs: []string{"a"},
						DDAndFpFile:          []byte{1, 2, 3},
						SNsSignatures:        [][]byte{[]byte{1, 2, 3}},
					},
					RqHandler: &mixins.RQHandler{
						RQIDs:          []string{"a"},
						RQIDsFile:      []byte{1, 2, 3},
						RQEncodeParams: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
					},
					ImageHandler: &mixins.NftImageHandler{},
				},
			},
			wantErr: common.ErrEmptyPreviewHash,
			want:    nil,
		},

		"raptorQ-symbols-error": {
			args: args{
				task: &NftRegistrationTask{
					dataHash: []byte{},
					Request: &NftRegistrationRequest{
						CreatorPastelID: "test-id",
					},
					service: &NftRegistrationService{
						config: NewConfig(),
					},
					FingerprintsHandler: &mixins.FingerprintsHandler{
						DDAndFingerprintsIDs: []string{"a"},
						DDAndFpFile:          []byte{1, 2, 3},
						SNsSignatures:        [][]byte{[]byte{1, 2, 3}},
					},
					RqHandler: &mixins.RQHandler{},
					ImageHandler: &mixins.NftImageHandler{
						PreviewHash:         []byte{1, 2, 3},
						MediumThumbnailHash: []byte{1, 2, 3},
						SmallThumbnailHash:  []byte{1, 2, 3},
					},
				},
			},
			wantErr: common.ErrEmptyRaptorQSymbols,
			want:    nil,
		},

		"success": {
			args: args{
				task: &NftRegistrationTask{
					Request: &NftRegistrationRequest{
						CreatorPastelID: "test-id",
						CreatorName:     "test-name",
						IssuedCopies:    &issuedCopies,
					},
					dataHash: []byte{1, 2},
					service: &NftRegistrationService{
						config: NewConfig(),
					},
					FingerprintsHandler: &mixins.FingerprintsHandler{
						DDAndFingerprintsIDs: []string{"a"},
						DDAndFpFile:          []byte{1, 2, 3},
						SNsSignatures:        [][]byte{[]byte{1, 2, 3}},
					},
					RqHandler: &mixins.RQHandler{
						RQIDs:          []string{"a"},
						RQIDsFile:      []byte{1, 2, 3},
						RQEncodeParams: rqnode.EncoderParameters{Oti: []byte{1, 2, 3}},
					},
					ImageHandler: &mixins.NftImageHandler{
						PreviewHash:         []byte{1, 2, 3},
						MediumThumbnailHash: []byte{1, 2, 3},
						SmallThumbnailHash:  []byte{1, 2, 3},
					},
					originalFileSizeInBytes: len([]byte{1, 2}),
					fileType:                "image/jpg",
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			err := tc.args.task.createNftTicket(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, tc.args.task.nftRegistrationTicket)
			}

		})
	}
}

func TestTaskSignTicket(t *testing.T) {
	type args struct {
		task        *NftRegistrationTask
		signErr     error
		signReturns []byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &NftRegistrationTask{
					Request: &NftRegistrationRequest{
						CreatorPastelID: "testid",
					},
					service: &NftRegistrationService{
						config: &Config{},
					},
					nftRegistrationTicket: &pastel.NFTTicket{},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &NftRegistrationTask{
					Request: &NftRegistrationRequest{
						CreatorPastelID: "testid",
					},
					service: &NftRegistrationService{
						config: &Config{},
					},
					nftRegistrationTicket: &pastel.NFTTicket{},
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
			tc.args.task.service.pastelHandler = mixins.NewPastelHandler(pastelClientMock)

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
