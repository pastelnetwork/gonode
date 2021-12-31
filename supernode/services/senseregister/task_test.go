package senseregister

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/types"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/ddclient/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/tj/assert"
)

func TestTaskSignAndSendArtTicket(t *testing.T) {
	type args struct {
		task        *Task
		signErr     error
		sendArtErr  error
		signReturns []byte
		primary     bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ActionTicket{},
				},
				primary: true,
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ActionTicket{},
				},
				signErr: errors.New("test"),
				primary: true,
			},
			wantErr: errors.New("test"),
		},
		"primary-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.ActionTicket{},
				},
				sendArtErr: errors.New("test"),
				signErr:    nil,
				primary:    false,
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

			clientMock := test.NewMockClient(t)
			clientMock.ListenOnSendArtTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterArtwork()

			tc.args.task.nodeClient = clientMock
			tc.args.task.connectedTo = &Node{client: clientMock}
			err := tc.args.task.connectedTo.connect(context.Background())
			assert.Nil(t, err)

			err = tc.args.task.signAndSendArtTicket(context.Background(), tc.args.primary)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskRegisterArt(t *testing.T) {
	type args struct {
		task     *Task
		regErr   error
		regRetID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
				regErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnRegisterArtTicket(tc.args.regRetID, tc.args.regErr).
				ListenOnRegisterNFTTicket(tc.args.regRetID, tc.args.regErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			id, err := tc.args.task.registerArt(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.args.regRetID, id)
			}
		})
	}
}

func TestTaskGenFingerprintsData(t *testing.T) {
	genfingerAndScoresFunc := func() *pastel.DDAndFingerprints {
		return &pastel.DDAndFingerprints{
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
	}

	type args struct {
		task    *Task
		fileErr error
		genErr  error
		genResp *pastel.DDAndFingerprints
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				genErr:  nil,
				fileErr: errors.New("test"),
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			file := artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.args.task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock
			_, err := tc.args.task.genFingerprintsData(context.Background(), file)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskPastelNodesByExtKey(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				masterNodesErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"node-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				nodeID:         "B",
				masterNodesErr: nil,
			},
			wantErr: errors.New("not found"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			_, err := tc.args.task.pastelNodeByExtKey(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskVerifyPeersSignature(t *testing.T) {
	type args struct {
		task      *Task
		verifyErr error
		verifyRet bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: true,
				verifyErr: nil,
			},
			wantErr: nil,
		},
		"verify-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: true,
				verifyErr: errors.New("test"),
			},
			wantErr: errors.New("verify signature"),
		},
		"verify-failure": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					peersArtTicketSignature: map[string][]byte{"A": []byte("test")},
				},
				verifyRet: false,
				verifyErr: nil,
			},
			wantErr: errors.New("mistmatch"),
		},
	}
	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(tc.args.verifyRet, tc.args.verifyErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			err := tc.args.task.verifyPeersSingature(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskWaitConfirmation(t *testing.T) {
	type args struct {
		task             *Task
		txid             string
		interval         time.Duration
		minConfirmations int64
		ctxTimeout       time.Duration
	}

	testCases := map[string]struct {
		args    args
		wantErr error
		retRes  *pastel.GetRawTransactionVerbose1Result
		retErr  error
	}{
		"min-confirmations-timeout": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				minConfirmations: 2,
				interval:         100 * time.Millisecond,
				ctxTimeout:       20 * time.Second,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			retErr:  nil,
			wantErr: errors.New("timeout"),
		},
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				minConfirmations: 1,
				interval:         50 * time.Millisecond,
				ctxTimeout:       500 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: nil,
		},
		"ctx-done-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				minConfirmations: 1,
				interval:         500 * time.Millisecond,
				ctxTimeout:       10 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
			},
			wantErr: errors.New("context"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), tc.args.ctxTimeout)
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetBlockCount(1, nil)
			pastelClientMock.ListenOnGetRawTransactionVerbose1(tc.retRes, tc.retErr)
			tc.args.task.Service.pastelClient = pastelClientMock

			err := <-tc.args.task.waitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.interval)
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
	// TODO: update later
	t.Skip()
	genfingerAndScoresFunc := func() *pastel.DDAndFingerprints {
		return &pastel.DDAndFingerprints{
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
	}

	type args struct {
		task    *Task
		fileErr error
		genErr  error
		genResp *pastel.DDAndFingerprints
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:                    task.New(StatusConnected),
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:                    task.New(StatusImageProbed),
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("required status"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task: &Task{
					Task: task.New(StatusConnected),
					Service: &Service{
						config: &Config{},
					},
					Ticket:                  &pastel.NFTTicket{},
					accepted:                Nodes{&Node{ID: "A"}, &Node{ID: "B"}},
					peersArtTicketSignature: map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}},
				},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := artwork.NewStorage(fsMock)
			file := artwork.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			//fingerprints := []float32{12.3, 34.4}
			//tc.args.genResp.Fingerprints = fingerprints
			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.args.task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			_, err := tc.args.task.ProbeImage(context.Background(), file)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskGetRegistrationFee(t *testing.T) {
	type args struct {
		task   *Task
		retFee int64
		retErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task: task.New(StatusImageAndThumbnailCoordinateUploaded),
					Ticket: &pastel.NFTTicket{
						Author: "author-id-b",
						AppTicketData: pastel.AppTicket{
							CreatorName: "Andy",
							NFTTitle:    "alantic",
						},
					},
				},
			},
			wantErr: nil,
		},

		"status-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task: task.New(StatusConnected),
					Ticket: &pastel.NFTTicket{
						Author: "author-id-b",
						AppTicketData: pastel.AppTicket{
							CreatorName: "Andy",
							NFTTitle:    "alantic",
						},
					},
				},
			},
			wantErr: errors.New("require status"),
		},
		/*"fee-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task: task.New(StatusImageAndThumbnailCoordinateUploaded),
					Ticket: &pastel.NFTTicket{
						Author: "author-id-b",
						AppTicketData: pastel.AppTicket{
							CreatorName: "Andy",
							NFTTitle:    "alantic",
						},
					},
				},
				retErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},*/
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnGetRegisterNFTFee(tc.args.retFee, tc.args.retErr)
			pastelClientMock.ListenOnVerify(true, nil)
			tc.args.task.Service.pastelClient = pastelClientMock

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)
			artTicketBytes, err := pastel.EncodeNFTTicket(tc.args.task.Ticket)
			assert.Nil(t, err)

			_, err = tc.args.task.GetRegistrationFee(context.Background(), artTicketBytes,
				[]byte{}, "", "", []byte{}, []byte{}, []byte{})
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskSessionNode(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusPrimaryMode),
					Ticket: &pastel.NFTTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusConnected),
					Ticket: &pastel.NFTTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: errors.New("status"),
		},
		"pastel-err": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Task:   task.New(StatusPrimaryMode),
					Ticket: &pastel.NFTTicket{},
				},
				masterNodesErr: errors.New("test"),
				nodeID:         "A",
			},
			wantErr: errors.New("get node"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			err := tc.args.task.SessionNode(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskAddPeerArticketSignature(t *testing.T) {
	type args struct {
		task           *Task
		nodeID         string
		masterNodesErr error
		acceptedNodeID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					peersArtTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                     task.New(StatusRegistrationFeeCalculated),
					Ticket:                   &pastel.NFTTicket{},
					allSignaturesReceivedChn: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				task: &Task{
					peersArtTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                     task.New(StatusConnected),
					Ticket:                   &pastel.NFTTicket{},
					allSignaturesReceivedChn: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: errors.New("status"),
		},
		"no-node-err": {
			args: args{
				task: &Task{
					peersArtTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					Task:                     task.New(StatusRegistrationFeeCalculated),
					Ticket:                   &pastel.NFTTicket{},
					allSignaturesReceivedChn: make(chan struct{}),
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("accepted"),
		},
		"success-close-sign-chn": {
			args: args{
				task: &Task{
					peersArtTicketSignatureMtx: &sync.Mutex{},
					Service: &Service{
						config: &Config{},
					},
					allSignaturesReceivedChn: make(chan struct{}),
					Task:                     task.New(StatusRegistrationFeeCalculated),
					Ticket:                   &pastel.NFTTicket{},
				},
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: tc.args.nodeID})
			tc.args.task.accepted = Nodes{&Node{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID}}

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
			tc.args.task.pastelClient = pastelClientMock
			tc.args.task.Service.pastelClient = pastelClientMock

			tc.args.task.peersArtTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := tc.args.task.AddPeerArticketSignature(tc.args.nodeID, []byte{})
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskValidateRqIDsAndDdFpIds(t *testing.T) {
	type args struct {
		task   *Task
		fg     *pastel.DDAndFingerprints
		rqFile *rqnode.RawSymbolIDFile
		ddSig  [][]byte
		rqSig  []byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				ddSig: [][]byte{[]byte("sig-1"), []byte("sig-2"), []byte("sig-3")},
				rqSig: []byte("rq-sig"),
				rqFile: &rqnode.RawSymbolIDFile{
					ID:                "id",
					SymbolIdentifiers: []string{"symbol-1", "symbol-2", "symbol-3", "symbol-4", "symbol-5"},
					BlockHash:         "block-hash",
					PastelID:          "author-pastelid",
				},

				fg: &pastel.DDAndFingerprints{
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
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					meshedNodes: []types.MeshedSuperNode{
						types.MeshedSuperNode{NodeID: "node-1"},
						types.MeshedSuperNode{NodeID: "node-2"},
						types.MeshedSuperNode{NodeID: "node-3"},
					},
					Task: task.New(StatusImageProbed),
					Ticket: &pastel.NFTTicket{
						Author:        "author-pastelid",
						AppTicketData: pastel.AppTicket{},
					},
				},
			},

			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(true, nil)
			tc.args.task.Service.pastelClient = pastelClientMock
			var dd []byte

			ddJSON, err := json.Marshal(tc.args.fg)
			assert.Nil(t, err)

			ddStr := base64.StdEncoding.EncodeToString(ddJSON)
			ddStr = ddStr + "." + base64.StdEncoding.EncodeToString(tc.args.ddSig[0]) + "." +
				base64.StdEncoding.EncodeToString(tc.args.ddSig[1]) + "." +
				base64.StdEncoding.EncodeToString(tc.args.ddSig[2])

			compressedDd, err := zstd.CompressLevel(nil, []byte(ddStr), 22)
			assert.Nil(t, err)
			dd = utils.B64Encode(compressedDd)

			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIc = rand.Uint32()
			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsMax = 50
			tc.args.task.Ticket.AppTicketData.RQIc = rand.Uint32()
			tc.args.task.Ticket.AppTicketData.RQMax = 50

			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIDs, _, err = pastel.GetIDFiles([]byte(ddStr),
				tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIc,
				tc.args.task.Ticket.AppTicketData.DDAndFingerprintsMax)

			assert.Nil(t, err)

			err = tc.args.task.validateRqIDsAndDdFpIds(context.Background(), dd)
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
