package nftregister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/ddclient/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	test "github.com/pastelnetwork/gonode/supernode/node/test/nft_register"
	"github.com/tj/assert"
)

func makeNftRegTaskWith2NodesAnd2TicketSignaturesConnected(status common.Status) *NftRegistrationTask {
	task := makeNftRegTaskWith2NodesAnd2TicketSignatures()
	task.UpdateStatus(status)

	meshedNodes := []types.MeshedSuperNode{
		types.MeshedSuperNode{
			NodeID: "PrimaryID",
		},
		types.MeshedSuperNode{
			NodeID: "A",
		},
		types.MeshedSuperNode{
			NodeID: "B",
		},
	}

	task.NetworkHandler.MeshNodes(nil, meshedNodes)
	task.RegTaskHelper.ServerPastelID = "PrimaryID"

	return task
}

func makeNftRegTaskWith2NodesAnd2TicketSignatures() *NftRegistrationTask {
	service := &NftRegistrationService{
		config:           &Config{},
		SuperNodeService: &common.SuperNodeService{},
	}

	task := NewNftRegistrationTask(service)
	task.Ticket = &pastel.NFTTicket{}
	task.NetworkHandler.Accepted = common.SuperNodePeerList{
		&common.SuperNodePeer{ID: "A"},
		&common.SuperNodePeer{ID: "B"},
	}
	task.PeersTicketSignature = map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}}

	return task
}

func makeEmptyNftRegTask() *NftRegistrationTask {
	service := &NftRegistrationService{
		config:           &Config{},
		SuperNodeService: &common.SuperNodeService{},
	}

	task := NewNftRegistrationTask(service)
	task.Ticket = &pastel.NFTTicket{}

	return task
}

func TestTaskSignAndSendArtTicket(t *testing.T) {
	type args struct {
		task        *NftRegistrationTask
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
				task:    makeEmptyNftRegTask(),
				primary: true,
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task:    makeEmptyNftRegTask(),
				signErr: errors.New("test"),
				primary: true,
			},
			wantErr: errors.New("test"),
		},
		"primary-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
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
			tc.args.task.PastelClient = pastelClientMock

			clientMock := test.NewMockClient(t)
			clientMock.ListenOnSendNftTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterNft()

			tc.args.task.nodeClient = clientMock
			tc.args.task.NetworkHandler.ConnectedTo = &common.SuperNodePeer{
				ClientInterface: clientMock,
				NodeMaker:       RegisterNftNodeMaker{},
			}
			err := tc.args.task.NetworkHandler.ConnectedTo.Connect(context.Background())
			assert.Nil(t, err)

			err = tc.args.task.signAndSendNftTicket(context.Background(), tc.args.primary)
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
		task     *NftRegistrationTask
		regErr   error
		regRetID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: makeNftRegTaskWith2NodesAnd2TicketSignatures(),
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task:   makeNftRegTaskWith2NodesAnd2TicketSignatures(),
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
			pastelClientMock.ListenOnRegisterNFTTicket(tc.args.regRetID, tc.args.regErr)
			tc.args.task.PastelClient = pastelClientMock

			id, err := tc.args.task.registerNft(context.Background())
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
		task    *NftRegistrationTask
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
				task:    makeNftRegTaskWith2NodesAnd2TicketSignatures(),
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				genErr:  nil,
				fileErr: errors.New("test"),
				genResp: genfingerAndScoresFunc(),
				task:    makeNftRegTaskWith2NodesAnd2TicketSignatures(),
			},
			wantErr: errors.New("test"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task:    makeNftRegTaskWith2NodesAnd2TicketSignatures(),
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

			storage := files.NewStorage(fsMock)
			file := files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil)
			tc.args.task.PastelClient = pastelClientMock
			tc.args.task.DupeDetectionHandler.RegTaskHelper.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}

			tc.args.task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.DupeDetectionHandler.DdClient = ddmock
			_, err := tc.args.task.GenFingerprintsData(context.Background(), file, "testBlockHash", "creatorPastelID")
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
		task           *NftRegistrationTask
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task:           makeEmptyNftRegTask(),
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				task:           makeEmptyNftRegTask(),
				masterNodesErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"node-err": {
			args: args{
				task:           makeEmptyNftRegTask(),
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
			tc.args.task.PastelClient = pastelClientMock
			tc.args.task.NetworkHandler.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}

			_, err := tc.args.task.NetworkHandler.PastelNodeByExtKey(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskCompareRQSymbolID(t *testing.T) {
	type args struct {
		task       *NftRegistrationTask
		connectErr error
		fileErr    error
		addIDsFile bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"conn-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				connectErr: errors.New("test"),
				fileErr:    nil,
				addIDsFile: true,
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				connectErr: nil,
				fileErr:    errors.New("test"),
				addIDsFile: true,
			},
			wantErr: errors.New("read image"),
		},
		"rqids-len-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				connectErr: nil,
				fileErr:    nil,
				addIDsFile: false,
			},
			wantErr: errors.New("no symbols identifiers"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr)

			tc.args.task.RQClient = rqClientMock
			tc.args.task.storage.RqClient = rqClientMock

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := files.NewStorage(fsMock)
			tc.args.task.Nft = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			if tc.args.addIDsFile {
				tc.args.task.rawRqFile = []byte{'a'}
			}

			err := tc.args.task.compareRQSymbolID(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskStoreRaptorQSymbols(t *testing.T) {
	type args struct {
		task       *NftRegistrationTask
		encodeErr  error
		connectErr error
		fileErr    error
		storeErr   error
		encodeResp *rqnode.Encode
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task:       makeEmptyNftRegTask(),
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    errors.New("test"),
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"conn-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				encodeErr:  nil,
				connectErr: errors.New("test"),
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"encode-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				encodeErr:  errors.New("test"),
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
		"store-err": {
			args: args{
				task:       makeEmptyNftRegTask(),
				encodeErr:  nil,
				connectErr: nil,
				fileErr:    nil,
				storeErr:   errors.New("test"),
				encodeResp: &rqnode.Encode{},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			rqFile := rq.SymbolIDFile{ID: "A"}
			bytes, err := json.Marshal(rqFile)
			assert.Nil(t, err)

			tc.args.encodeResp.Symbols = map[string][]byte{"A": bytes}

			rqClientMock := rqMock.NewMockClient(t)
			rqClientMock.ListenOnEncodeInfo(&rqnode.EncodeInfo{}, nil)
			rqClientMock.ListenOnRaptorQ().ListenOnClose(nil)
			rqClientMock.ListenOnConnect(tc.args.connectErr).
				ListenOnEncode(tc.args.encodeResp, tc.args.encodeErr)

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnStore("", tc.args.storeErr)
			tc.args.task.P2PClient = p2pClient
			tc.args.task.storage.P2PClient = p2pClient

			tc.args.task.RQClient = rqClientMock
			tc.args.task.storage.RqClient = rqClientMock

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := files.NewStorage(fsMock)
			tc.args.task.Nft = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err = tc.args.task.storeRaptorQSymbols(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskStoreThumbnails(t *testing.T) {
	type args struct {
		task     *NftRegistrationTask
		storeErr error
		fileErr  error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task:     makeEmptyNftRegTask(),
				fileErr:  nil,
				storeErr: nil,
			},
			wantErr: nil,
		},

		"store-err": {
			args: args{
				task:     makeEmptyNftRegTask(),
				fileErr:  nil,
				storeErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				task:     makeEmptyNftRegTask(),
				fileErr:  errors.New("test"),
				storeErr: nil,
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			p2pClient := p2pMock.NewMockClient(t)
			p2pClient.ListenOnStore("", tc.args.storeErr)
			tc.args.task.P2PClient = p2pClient
			tc.args.task.storage.P2PClient = p2pClient

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			storage := files.NewStorage(fsMock)

			tc.args.task.SmallThumbnail = files.NewFile(storage, "test-small")
			tc.args.task.MediumThumbnail = files.NewFile(storage, "test-medium")
			tc.args.task.PreviewThumbnail = files.NewFile(storage, "test-preview")

			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err := tc.args.task.storeThumbnails(context.Background())
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
		task      *NftRegistrationTask
		verifyErr error
		verifyRet bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task:      makeNftRegTaskWith2NodesAnd2TicketSignatures(),
				verifyRet: true,
				verifyErr: nil,
			},
			wantErr: nil,
		},
		"verify-err": {
			args: args{
				task:      makeNftRegTaskWith2NodesAnd2TicketSignatures(),
				verifyRet: true,
				verifyErr: errors.New("test"),
			},
			wantErr: errors.New("verify signature"),
		},
		"verify-failure": {
			args: args{
				task:      makeNftRegTaskWith2NodesAnd2TicketSignatures(),
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
			tc.args.task.PastelClient = pastelClientMock
			tc.args.task.DupeDetectionHandler.RegTaskHelper.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}

			err := tc.args.task.verifyPeersSignature(context.Background())
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
		task             *NftRegistrationTask
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
				task:             makeEmptyNftRegTask(),
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
				task:             makeEmptyNftRegTask(),
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
				task:             makeEmptyNftRegTask(),
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
			tc.args.task.PastelClient = pastelClientMock
			tc.args.task.DupeDetectionHandler.RegTaskHelper.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}

			err := <-tc.args.task.WaitConfirmation(ctx, tc.args.txid,
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
		task    *NftRegistrationTask
		fileErr error
		genErr  error
		genResp *pastel.DDAndFingerprints
	}

	//serviceCfg := NewConfig()
	//serviceCfg.PastelID = "PrimaryID"
	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task:    makeNftRegTaskWith2NodesAnd2TicketSignaturesConnected(common.StatusConnected),
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task:    makeNftRegTaskWith2NodesAnd2TicketSignaturesConnected(common.StatusImageProbed),
			},
			wantErr: errors.New("required status"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				task:    makeNftRegTaskWith2NodesAnd2TicketSignaturesConnected(common.StatusConnected),
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

			storage := files.NewStorage(fsMock)
			file := files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)
			tc.args.task.ddClient = ddmock
			tc.args.task.DupeDetectionHandler.DdClient = ddmock

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil).ListenOnVerify(true, nil)
			tc.args.task.PastelClient = pastelClientMock
			tc.args.task.DupeDetectionHandler.RegTaskHelper.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}
			tc.args.task.NetworkHandler.PastelHandler = &mixins.PastelHandler{PastelClient: pastelClientMock}

			tc.args.task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go tc.args.task.RunAction(ctx)

			if tc.wantErr == nil {
				clientMock := test.NewMockClient(t)
				clientMock.ListenOnSendSignedDDAndFingerprints(nil).
					ListenOnConnect("", nil).ListenOnRegisterNft()

				tc.args.task.nodeClient = clientMock
				tc.args.task.NetworkHandler.NodeClient = clientMock

				nodes := pastel.MasterNodes{}
				nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})

				pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

				peerDDAndFingerprints, _ := pastel.ToCompressSignedDDAndFingerprints(genfingerAndScoresFunc(), []byte("signature"))
				go tc.args.task.AddSignedDDAndFingerprints("A", peerDDAndFingerprints)
				go tc.args.task.AddSignedDDAndFingerprints("B", peerDDAndFingerprints)
			}
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

//func TestTaskGetRegistrationFee(t *testing.T) {
//	type args struct {
//		task   *NftRegistrationTask
//		retFee int64
//		retErr error
//	}
//
//	testCases := map[string]struct {
//		args    args
//		wantErr error
//	}{
//		"success": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task: task.New(StatusImageAndThumbnailCoordinateUploaded),
//					Ticket: &pastel.NFTTicket{
//						Author: "author-id-b",
//						AppTicketData: pastel.AppTicket{
//							CreatorName: "Andy",
//							NFTTitle:    "alantic",
//						},
//					},
//				},
//			},
//			wantErr: nil,
//		},
//
//		"status-err": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task: task.New(StatusConnected),
//					Ticket: &pastel.NFTTicket{
//						Author: "author-id-b",
//						AppTicketData: pastel.AppTicket{
//							CreatorName: "Andy",
//							NFTTitle:    "alantic",
//						},
//					},
//				},
//			},
//			wantErr: errors.New("require status"),
//		},
//		/*"fee-err": {
//			args: args{
//				task: &Task{
//					Service: &Service{
//						config: &Config{},
//					},
//					Task: task.New(StatusImageAndThumbnailCoordinateUploaded),
//					Request: &pastel.NFTTicket{
//						Author: "author-id-b",
//						AppTicketData: pastel.AppTicket{
//							CreatorName: "Andy",
//							NFTTitle:    "alantic",
//						},
//					},
//				},
//				retErr: errors.New("test"),
//			},
//			wantErr: errors.New("test"),
//		},*/
//	}
//
//	for name, tc := range testCases {
//		tc := tc
//
//		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
//			t.Parallel()
//
//			pastelClientMock := pastelMock.NewMockClient(t)
//			pastelClientMock.ListenOnGetRegisterNFTFee(tc.args.retFee, tc.args.retErr)
//			pastelClientMock.ListenOnVerify(true, nil)
//			tc.args.task.NftRegistrationService.pastelClient = pastelClientMock
//
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//
//			go tc.args.task.RunAction(ctx)
//			artTicketBytes, err := pastel.EncodeNFTTicket(tc.args.task.Ticket)
//			assert.Nil(t, err)
//
//			_, err = tc.args.task.GetNftRegistrationFee(context.Background(), artTicketBytes,
//				[]byte{}, "", "", []byte{}, []byte{}, []byte{})
//			if tc.wantErr != nil {
//				assert.NotNil(t, err)
//				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
//			} else {
//				assert.Nil(t, err)
//			}
//		})
//	}
//}
//
//func TestTaskSessionNode(t *testing.T) {
//	type args struct {
//		task           *NftRegistrationTask
//		nodeID         string
//		masterNodesErr error
//	}
//
//	testCases := map[string]struct {
//		args    args
//		wantErr error
//	}{
//		"success": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:   task.New(StatusPrimaryMode),
//					Ticket: &pastel.NFTTicket{},
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//			},
//			wantErr: nil,
//		},
//		"status-err": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:   task.New(StatusConnected),
//					Ticket: &pastel.NFTTicket{},
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//			},
//			wantErr: errors.New("status"),
//		},
//		"pastel-err": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:   task.New(StatusPrimaryMode),
//					Ticket: &pastel.NFTTicket{},
//				},
//				masterNodesErr: errors.New("test"),
//				nodeID:         "A",
//			},
//			wantErr: errors.New("get node"),
//		},
//	}
//
//	for name, tc := range testCases {
//		tc := tc
//
//		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
//			t.Parallel()
//
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//
//			go tc.args.task.RunAction(ctx)
//
//			nodes := pastel.MasterNodes{}
//			for i := 0; i < 10; i++ {
//				nodes = append(nodes, pastel.MasterNode{})
//			}
//			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
//
//			pastelClientMock := pastelMock.NewMockClient(t)
//			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
//			tc.args.task.pastelClient = pastelClientMock
//			tc.args.task.NftRegistrationService.pastelClient = pastelClientMock
//
//			err := tc.args.task.SessionNode(context.Background(), tc.args.nodeID)
//			if tc.wantErr != nil {
//				assert.NotNil(t, err)
//				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
//			} else {
//				assert.Nil(t, err)
//			}
//		})
//	}
//}
//
//func TestTaskAddPeerArticketSignature(t *testing.T) {
//	type args struct {
//		task           *NftRegistrationTask
//		nodeID         string
//		masterNodesErr error
//		acceptedNodeID string
//	}
//
//	testCases := map[string]struct {
//		args    args
//		wantErr error
//	}{
//		"success": {
//			args: args{
//				task: &NftRegistrationTask{
//					peersArtTicketSignatureMtx: &sync.Mutex{},
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:                     task.New(StatusRegistrationFeeCalculated),
//					Ticket:                   &pastel.NFTTicket{},
//					allSignaturesReceivedChn: make(chan struct{}),
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//				acceptedNodeID: "A",
//			},
//			wantErr: nil,
//		},
//		"status-err": {
//			args: args{
//				task: &NftRegistrationTask{
//					peersArtTicketSignatureMtx: &sync.Mutex{},
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:                     task.New(StatusConnected),
//					Ticket:                   &pastel.NFTTicket{},
//					allSignaturesReceivedChn: make(chan struct{}),
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//				acceptedNodeID: "A",
//			},
//			wantErr: errors.New("status"),
//		},
//		"no-node-err": {
//			args: args{
//				task: &NftRegistrationTask{
//					peersArtTicketSignatureMtx: &sync.Mutex{},
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:                     task.New(StatusRegistrationFeeCalculated),
//					Ticket:                   &pastel.NFTTicket{},
//					allSignaturesReceivedChn: make(chan struct{}),
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//				acceptedNodeID: "B",
//			},
//			wantErr: errors.New("accepted"),
//		},
//		"success-close-sign-chn": {
//			args: args{
//				task: &NftRegistrationTask{
//					peersArtTicketSignatureMtx: &sync.Mutex{},
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					allSignaturesReceivedChn: make(chan struct{}),
//					Task:                     task.New(StatusRegistrationFeeCalculated),
//					Ticket:                   &pastel.NFTTicket{},
//				},
//				masterNodesErr: nil,
//				nodeID:         "A",
//				acceptedNodeID: "A",
//			},
//			wantErr: nil,
//		},
//	}
//
//	for name, tc := range testCases {
//		tc := tc
//
//		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
//			t.Parallel()
//
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//
//			go tc.args.task.RunAction(ctx)
//
//			nodes := pastel.MasterNodes{}
//			for i := 0; i < 10; i++ {
//				nodes = append(nodes, pastel.MasterNode{})
//			}
//			nodes = append(nodes, pastel.MasterNode{ExtKey: tc.args.nodeID})
//			tc.args.task.accepted = common.NftRegistrationNodes{&NftRegistrationNode{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID}}
//
//			pastelClientMock := pastelMock.NewMockClient(t)
//			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr)
//			tc.args.task.pastelClient = pastelClientMock
//			tc.args.task.NftRegistrationService.pastelClient = pastelClientMock
//
//			tc.args.task.peersArtTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}
//
//			err := tc.args.task.AddPeerTicketSignature(tc.args.nodeID, []byte{})
//			if tc.wantErr != nil {
//				assert.NotNil(t, err)
//				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
//			} else {
//				assert.Nil(t, err)
//			}
//		})
//	}
//}
//
//func TestTaskUploadImageWithThumbnail(t *testing.T) {
//	type args struct {
//		task *NftRegistrationTask
//	}
//
//	testCases := map[string]struct {
//		args    args
//		wantErr error
//	}{
//		"success": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:   task.New(StatusImageProbed),
//					Ticket: &pastel.NFTTicket{},
//				},
//			},
//			wantErr: nil,
//		},
//		"failure-status": {
//			args: args{
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					Task:   task.New(StatusConnected),
//					Ticket: &pastel.NFTTicket{},
//				},
//			},
//			wantErr: errors.New("status"),
//		},
//	}
//
//	for name, tc := range testCases {
//		tc := tc
//
//		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
//			t.Parallel()
//
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//
//			go tc.args.task.RunAction(ctx)
//
//			stg := files.NewStorage(fs.NewFileStorage(os.TempDir()))
//
//			tc.args.task.Storage = stg
//			file, err := newTestImageFile(stg)
//			assert.Nil(t, err)
//			tc.args.task.Nft = file
//
//			coordinate := files.ThumbnailCoordinate{
//				TopLeftX:     0,
//				TopLeftY:     0,
//				BottomRightX: 400,
//				BottomRightY: 400,
//			}
//			_, _, _, err = tc.args.task.UploadImageWithThumbnail(ctx, file, coordinate)
//			if tc.wantErr != nil {
//				assert.NotNil(t, err)
//				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
//			} else {
//				assert.Nil(t, err)
//			}
//		})
//	}
//}
//
//func TestTaskValidateRqIDsAndDdFpIds(t *testing.T) {
//	type args struct {
//		task   *NftRegistrationTask
//		fg     *pastel.DDAndFingerprints
//		rqFile *rqnode.RawSymbolIDFile
//		ddSig  [][]byte
//		rqSig  []byte
//	}
//
//	testCases := map[string]struct {
//		args    args
//		wantErr error
//	}{
//		"success": {
//			args: args{
//				ddSig: [][]byte{[]byte("sig-1"), []byte("sig-2"), []byte("sig-3")},
//				rqSig: []byte("rq-sig"),
//				rqFile: &rqnode.RawSymbolIDFile{
//					ID:                "id",
//					SymbolIdentifiers: []string{"symbol-1", "symbol-2", "symbol-3", "symbol-4", "symbol-5"},
//					BlockHash:         "block-hash",
//					PastelID:          "author-pastelid",
//				},
//
//				fg: &pastel.DDAndFingerprints{
//					Block:                      "Block",
//					Principal:                  "Principal",
//					DupeDetectionSystemVersion: "v1.0",
//
//					IsLikelyDupe:     true,
//					IsRareOnInternet: true,
//
//					RarenessScores: &pastel.RarenessScores{
//						CombinedRarenessScore:         0,
//						XgboostPredictedRarenessScore: 0,
//						NnPredictedRarenessScore:      0,
//						OverallAverageRarenessScore:   0,
//					},
//					InternetRareness: &pastel.InternetRareness{
//						MatchesFoundOnFirstPage: 0,
//						NumberOfPagesOfResults:  0,
//						URLOfFirstMatchInPage:   "",
//					},
//
//					OpenNSFWScore: 0.1,
//					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
//						Drawings: 0.1,
//						Hentai:   0.2,
//						Neutral:  0.3,
//						Porn:     0.4,
//						Sexy:     0.5,
//					},
//
//					ImageFingerprintOfCandidateImageFile: []float32{1, 2, 3},
//					FingerprintsStat: &pastel.FingerprintsStat{
//						NumberOfFingerprintsRequiringFurtherTesting1: 1,
//						NumberOfFingerprintsRequiringFurtherTesting2: 2,
//						NumberOfFingerprintsRequiringFurtherTesting3: 3,
//						NumberOfFingerprintsRequiringFurtherTesting4: 4,
//						NumberOfFingerprintsRequiringFurtherTesting5: 5,
//						NumberOfFingerprintsRequiringFurtherTesting6: 6,
//						NumberOfFingerprintsOfSuspectedDupes:         7,
//					},
//
//					HashOfCandidateImageFile: "HashOfCandidateImageFile",
//					PerceptualImageHashes: &pastel.PerceptualImageHashes{
//						PDQHash:        "PdqHash",
//						PerceptualHash: "PerceptualHash",
//						AverageHash:    "AverageHash",
//						DifferenceHash: "DifferenceHash",
//						NeuralHash:     "NeuralhashHash",
//					},
//					PerceptualHashOverlapCount: 1,
//
//					Maxes: &pastel.Maxes{
//						PearsonMax:           1.0,
//						SpearmanMax:          2.0,
//						KendallMax:           3.0,
//						HoeffdingMax:         4.0,
//						MutualInformationMax: 5.0,
//						HsicMax:              6.0,
//						XgbimportanceMax:     7.0,
//					},
//					Percentile: &pastel.Percentile{
//						PearsonTop1BpsPercentile:             1.0,
//						SpearmanTop1BpsPercentile:            2.0,
//						KendallTop1BpsPercentile:             3.0,
//						HoeffdingTop10BpsPercentile:          4.0,
//						MutualInformationTop100BpsPercentile: 5.0,
//						HsicTop100BpsPercentile:              6.0,
//						XgbimportanceTop100BpsPercentile:     7.0,
//					},
//				},
//				task: &NftRegistrationTask{
//					NftRegistrationService: &NftRegistrationService{
//						config: &Config{},
//					},
//					meshedNodes: []types.MeshedSuperNode{
//						types.MeshedSuperNode{NodeID: "node-1"},
//						types.MeshedSuperNode{NodeID: "node-2"},
//						types.MeshedSuperNode{NodeID: "node-3"},
//					},
//					Task: task.New(StatusImageProbed),
//					Ticket: &pastel.NFTTicket{
//						Author:        "author-pastelid",
//						AppTicketData: pastel.AppTicket{},
//					},
//				},
//			},
//
//			wantErr: nil,
//		},
//	}
//
//	for name, tc := range testCases {
//		tc := tc
//
//		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
//			pastelClientMock := pastelMock.NewMockClient(t)
//			pastelClientMock.ListenOnVerify(true, nil)
//			tc.args.task.NftRegistrationService.pastelClient = pastelClientMock
//			var rq, dd []byte
//
//			ddJSON, err := json.Marshal(tc.args.fg)
//			assert.Nil(t, err)
//
//			ddStr := base64.StdEncoding.EncodeToString(ddJSON)
//			ddStr = ddStr + "." + base64.StdEncoding.EncodeToString(tc.args.ddSig[0]) + "." +
//				base64.StdEncoding.EncodeToString(tc.args.ddSig[1]) + "." +
//				base64.StdEncoding.EncodeToString(tc.args.ddSig[2])
//
//			compressedDd, err := zstd.CompressLevel(nil, []byte(ddStr), 22)
//			assert.Nil(t, err)
//			dd = utils.B64Encode(compressedDd)
//
//			rqJSON, err := json.Marshal(tc.args.rqFile)
//			assert.Nil(t, err)
//			rqStr := base64.StdEncoding.EncodeToString(rqJSON)
//			rqStr = rqStr + "." + base64.StdEncoding.EncodeToString(tc.args.rqSig)
//			compressedRq, err := zstd.CompressLevel(nil, []byte(rqStr), 22)
//			assert.Nil(t, err)
//			rq = utils.B64Encode(compressedRq)
//
//			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIc = rand.Uint32()
//			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsMax = 50
//			tc.args.task.Ticket.AppTicketData.RQIc = rand.Uint32()
//			tc.args.task.Ticket.AppTicketData.RQMax = 50
//
//			tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIDs, _, err = pastel.GetIDFiles([]byte(ddStr),
//				tc.args.task.Ticket.AppTicketData.DDAndFingerprintsIc,
//				tc.args.task.Ticket.AppTicketData.DDAndFingerprintsMax)
//
//			assert.Nil(t, err)
//
//			tc.args.task.Ticket.AppTicketData.RQIDs, _, err = pastel.GetIDFiles([]byte(rqStr),
//				tc.args.task.Ticket.AppTicketData.RQIc,
//				tc.args.task.Ticket.AppTicketData.RQMax)
//
//			assert.Nil(t, err)
//
//			err = tc.args.task.validateRqIDsAndDdFpIds(context.Background(), rq, dd)
//			if tc.wantErr != nil {
//				assert.NotNil(t, err)
//				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
//			} else {
//				if err != nil {
//					fmt.Println("err: ", err.Error())
//				}
//				assert.Nil(t, err)
//			}
//		})
//	}
//}
