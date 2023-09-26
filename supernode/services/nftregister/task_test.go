package nftregister

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/ddclient/test"
	"github.com/pastelnetwork/gonode/p2p"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqMock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/supernode/node"
	test "github.com/pastelnetwork/gonode/supernode/node/test/nft_register"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

func makeConnected(task *NftRegistrationTask, status common.Status) *NftRegistrationTask {
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

	task.UpdateStatus(common.StatusConnected)
	task.NetworkHandler.MeshNodes(context.TODO(), meshedNodes)
	task.UpdateStatus(status)

	return task
}

func add2NodesAnd2TicketSignatures(task *NftRegistrationTask) *NftRegistrationTask {
	task.NetworkHandler.Accepted = common.SuperNodePeerList{
		&common.SuperNodePeer{ID: "A"},
		&common.SuperNodePeer{ID: "B"},
	}
	task.PeersTicketSignature = map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}}
	return task
}

func makeEmptyNftRegTask(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client,
	ddClient ddclient.DDServerClient, rqClient rqnode.ClientInterface) *NftRegistrationTask {
	service := NewService(config, fileStorage, pastelClient, nodeClient, p2pClient, ddClient, nil)
	task := NewNftRegistrationTask(service)
	task.storage.RqClient = rqClient
	task.Ticket = &pastel.NFTTicket{}
	task.ActionTicketRegMetadata = &types.ActionRegMetadata{EstimatedFee: 100}

	return task
}

func TestTaskSignAndSendArtTicket(t *testing.T) {
	type args struct {
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
				primary: true,
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				signErr: errors.New("test"),
				primary: true,
			},
			wantErr: errors.New("test"),
		},
		"primary-err": {
			args: args{
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

			clientMock := test.NewMockClient(t)
			clientMock.ListenOnSendNftTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterNft()

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, clientMock, nil, nil, nil)

			task.NetworkHandler.ConnectedTo = &common.SuperNodePeer{
				ClientInterface: clientMock,
				NodeMaker:       RegisterNftNodeMaker{},
			}
			err := task.NetworkHandler.ConnectedTo.Connect(context.Background())
			assert.Nil(t, err)

			err = task.signAndSendNftTicket(context.Background(), tc.args.primary)
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
		regErr   error
		regRetID string
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args:    args{},
			wantErr: nil,
		},
		"err": {
			args: args{
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
			pastelClientMock.ListenOnGetRawMempoolMethod([]string{}, nil).
				ListenOnGetInactiveActionTickets(pastel.ActTickets{}, nil).
				ListenOnGetInactiveNFTTickets(pastel.RegTickets{}, nil).
				ListenOnRegisterNFTTicket(tc.args.regRetID, tc.args.regErr)

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			id, err := task.registerNft(context.Background())
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
	}

	type args struct {
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
			},
			wantErr: nil,
		},
		"file-err": {
			args: args{
				genErr:  nil,
				fileErr: errors.New("test"),
				genResp: genfingerAndScoresFunc(),
			},
			wantErr: errors.New("test"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
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

			ddmock := ddMock.NewMockClient(t)
			ddmock.ListenOnImageRarenessScore(tc.args.genResp, tc.args.genErr)

			task := makeEmptyNftRegTask(&Config{}, fsMock, pastelClientMock, nil, nil, ddmock, nil)
			task = add2NodesAnd2TicketSignatures(task)
			task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID", BlockHeight: "testBlockHeight", Timestamp: "2022-03-31 16:55:28"}

			_, err := task.GenFingerprintsData(context.Background(), file, "testBlockHash", "testBlockHeight", "2022-03-31 16:55:28", "creatorPastelID", "jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis", "jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F", "jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ", false, "", "")
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
		nodeID         string
		masterNodesErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"err": {
			args: args{
				masterNodesErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
		"node-err": {
			args: args{
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
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)

			_, err := task.NetworkHandler.PastelNodeByExtKey(context.Background(), tc.args.nodeID)
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
				connectErr: errors.New("test"),
				fileErr:    nil,
				addIDsFile: true,
			},
			wantErr: errors.New("test"),
		},
		"file-err": {
			args: args{
				connectErr: nil,
				fileErr:    errors.New("test"),
				addIDsFile: true,
			},
			wantErr: errors.New("read image"),
		},
		"rqids-len-err": {
			args: args{
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

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			task := makeEmptyNftRegTask(&Config{}, fsMock, nil, nil, nil, nil, rqClientMock)

			storage := files.NewStorage(fsMock)
			task.Nft = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			if tc.args.addIDsFile {
				task.rawRqFile = []byte{'a'}
			}

			err := task.compareRQSymbolID(context.Background())
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
				encodeErr:  errors.New("test"),
				connectErr: nil,
				fileErr:    nil,
				storeErr:   nil,
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
			p2pClient.ListenOnStore("", tc.args.storeErr).On("StoreBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.args.storeErr)

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			task := makeEmptyNftRegTask(&Config{}, fsMock, nil, nil, p2pClient, nil, rqClientMock)

			storage := files.NewStorage(fsMock)
			task.Nft = files.NewFile(storage, "test")
			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err = task.storeRaptorQSymbols(context.Background())
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
		storeErr error
		fileErr  error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				fileErr:  nil,
				storeErr: nil,
			},
			wantErr: nil,
		},

		"file-err": {
			args: args{
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
			p2pClient.ListenOnStore("", tc.args.storeErr).On("StoreBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.args.storeErr)

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF)

			task := makeEmptyNftRegTask(&Config{}, fsMock, nil, nil, p2pClient, nil, nil)

			storage := files.NewStorage(fsMock)
			task.SmallThumbnail = files.NewFile(storage, "test-small")
			task.MediumThumbnail = files.NewFile(storage, "test-medium")
			task.PreviewThumbnail = files.NewFile(storage, "test-preview")

			fsMock.ListenOnOpen(fileMock, tc.args.fileErr)

			err := task.storeThumbnails(context.Background())
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
		verifyErr error
		verifyRet bool
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				verifyRet: true,
				verifyErr: nil,
			},
			wantErr: nil,
		},
		"verify-err": {
			args: args{
				verifyRet: true,
				verifyErr: errors.New("test"),
			},
			wantErr: errors.New("verify signature"),
		},
		"verify-failure": {
			args: args{
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

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			err := task.verifyPeersSignature(context.Background())
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
				minConfirmations: 2,
				interval:         100 * time.Millisecond,
				ctxTimeout:       20 * time.Second,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 20,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
			},
			retErr:  nil,
			wantErr: errors.New("timeout"),
		},
		"success": {
			args: args{
				minConfirmations: 1,
				interval:         50 * time.Millisecond,
				ctxTimeout:       500 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 10,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
			},
			wantErr: nil,
		},
		"ctx-done-err": {
			args: args{
				minConfirmations: 1,
				interval:         500 * time.Millisecond,
				ctxTimeout:       10 * time.Millisecond,
			},
			retRes: &pastel.GetRawTransactionVerbose1Result{
				Confirmations: 1,
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 20,
					ScriptPubKey: pastel.ScriptPubKey{
						Addresses: []string{"tPpasteLBurnAddressXXXXXXXXXXX3wy7u"}}}},
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

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task.registrationFee = 100

			err := <-task.WaitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.interval, true, float64(task.registrationFee), 10)
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
	}

	type args struct {
		fileErr error
		genErr  error
		genResp *pastel.DDAndFingerprints
		status  common.Status
	}

	serviceCfg := NewConfig()
	serviceCfg.PastelID = "PrimaryID"
	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				status:  common.StatusConnected,
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				genErr:  nil,
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				status:  common.StatusImageProbed,
			},
			wantErr: errors.New("required status"),
		},
		"gen-err": {
			args: args{
				genErr:  errors.New("test"),
				fileErr: nil,
				genResp: genfingerAndScoresFunc(),
				status:  common.StatusConnected,
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

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil).ListenOnVerify(true, nil)

			var clientMock *test.Client
			if tc.wantErr == nil {
				clientMock = test.NewMockClient(t)
				clientMock.ListenOnSendSignedDDAndFingerprints(nil).
					ListenOnConnect("", nil).ListenOnRegisterNft()
			}

			task := makeEmptyNftRegTask(serviceCfg, fsMock, pastelClientMock, clientMock, nil, ddmock, nil)
			task = add2NodesAnd2TicketSignatures(task)
			task = makeConnected(task, tc.args.status)

			task.nftRegMetadata = &types.NftRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID", BlockHeight: "testBlockHeight", Timestamp: "2022-03-31 16:55:28"}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go task.RunAction(ctx)

			if tc.wantErr == nil {
				nodes := pastel.MasterNodes{}
				nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})

				pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnMasterNodesExtra(nodes, nil)

				peerDDAndFingerprints, _ := pastel.ToCompressSignedDDAndFingerprints(context.Background(), genfingerAndScoresFunc(), []byte("signature"))
				go task.AddSignedDDAndFingerprints("A", peerDDAndFingerprints)
				go task.AddSignedDDAndFingerprints("B", peerDDAndFingerprints)
			}
			_, err := task.ProbeImage(context.Background(), file)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func makeTask1(pastelClient pastel.Client, status common.Status) *NftRegistrationTask {
	task := makeEmptyNftRegTask(&Config{}, nil, pastelClient, nil, nil, nil, nil)
	task.UpdateStatus(status)
	task.Ticket = &pastel.NFTTicket{
		Author: "author-id-b",
		AppTicketData: pastel.AppTicket{
			CreatorName: "Andy",
			NFTTitle:    "alantic",
		},
	}
	return task
}

func TestTaskGetRegistrationFee(t *testing.T) {
	type args struct {
		retFee int64
		retErr error
		status common.Status
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				status: common.StatusImageAndThumbnailCoordinateUploaded,
			},
			wantErr: nil,
		},

		"status-err": {
			args: args{
				status: common.StatusConnected,
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
					Request: &pastel.NFTTicket{
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

			task := makeTask1(pastelClientMock, tc.args.status)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go task.RunAction(ctx)
			artTicketBytes, err := pastel.EncodeNFTTicket(task.Ticket)
			assert.Nil(t, err)

			_, err = task.GetNftRegistrationFee(context.Background(), artTicketBytes,
				[]byte{}, "", []byte{}, []byte{}, []byte{})
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
		nodeID         string
		masterNodesErr error
		status         common.Status
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				status:         common.StatusPrimaryMode,
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				status:         common.StatusConnected,
				masterNodesErr: nil,
				nodeID:         "A",
			},
			wantErr: errors.New("status"),
		},
		"pastel-err": {
			args: args{
				status:         common.StatusPrimaryMode,
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

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			err := task.NetworkHandler.SessionNode(context.Background(), tc.args.nodeID)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskAddPeerNftTicketSignature(t *testing.T) {
	type args struct {
		status         common.Status
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
				status:         common.StatusRegistrationFeeCalculated,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: nil,
		},
		"status-err": {
			args: args{
				status:         common.StatusConnected,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "A",
			},
			wantErr: errors.New("status"),
		},
		"no-node-err": {
			args: args{
				status:         common.StatusRegistrationFeeCalculated,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("not in Accepted list"),
		},
		"success-close-sign-chn": {
			args: args{
				status:         common.StatusRegistrationFeeCalculated,
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

			nodes := pastel.MasterNodes{}
			for i := 0; i < 10; i++ {
				nodes = append(nodes, pastel.MasterNode{})
			}
			nodes = append(nodes, pastel.MasterNode{ExtKey: tc.args.nodeID})
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnMasterNodesTop(nodes, tc.args.masterNodesErr).ListenOnMasterNodesExtra(nodes, tc.args.masterNodesErr)

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			task.NetworkHandler.Accepted = common.SuperNodePeerList{
				&common.SuperNodePeer{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID},
			}
			task.PeersTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := task.AddPeerTicketSignature(tc.args.nodeID, []byte{}, common.StatusRegistrationFeeCalculated)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskUploadImageWithThumbnail(t *testing.T) {
	type args struct {
		status common.Status
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				status: common.StatusImageProbed,
			},
			wantErr: nil,
		},
		"failure-status": {
			args: args{
				status: common.StatusConnected,
			},
			wantErr: errors.New("status"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			stg := files.NewStorage(fs.NewFileStorage(os.TempDir()))

			task := makeEmptyNftRegTask(&Config{}, stg, nil, nil, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			file, err := newTestImageFile(stg)
			assert.Nil(t, err)
			task.Nft = file

			coordinate := files.ThumbnailCoordinate{
				TopLeftX:     0,
				TopLeftY:     0,
				BottomRightX: 400,
				BottomRightY: 400,
			}
			_, _, _, err = task.UploadImageWithThumbnail(ctx, file, coordinate)
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
		status common.Status
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
				},
				status: common.StatusImageProbed,
			},

			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(true, nil)

			task := makeEmptyNftRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)

			meshedNodes := []types.MeshedSuperNode{
				types.MeshedSuperNode{NodeID: "node-1"},
				types.MeshedSuperNode{NodeID: "node-2"},
				types.MeshedSuperNode{NodeID: "node-3"},
			}

			task.UpdateStatus(common.StatusConnected)
			task.NetworkHandler.MeshNodes(context.TODO(), meshedNodes)
			task.UpdateStatus(tc.args.status)

			task.Ticket = &pastel.NFTTicket{
				Author:        "author-pastelid",
				AppTicketData: pastel.AppTicket{},
			}

			var rq, dd []byte
			ddJSON, err := json.Marshal(tc.args.fg)
			assert.Nil(t, err)

			ddStr := base64.StdEncoding.EncodeToString(ddJSON)
			ddStr = ddStr + "." + base64.StdEncoding.EncodeToString(tc.args.ddSig[0]) + "." +
				base64.StdEncoding.EncodeToString(tc.args.ddSig[1]) + "." +
				base64.StdEncoding.EncodeToString(tc.args.ddSig[2])

			compressedDd, err := utils.Compress([]byte(ddStr), 4)
			assert.Nil(t, err)
			dd = utils.B64Encode(compressedDd)

			rqJSON, err := json.Marshal(tc.args.rqFile)
			assert.Nil(t, err)
			rqStr := base64.StdEncoding.EncodeToString(rqJSON)
			rqStr = rqStr + "." + base64.StdEncoding.EncodeToString(tc.args.rqSig)
			compressedRq, err := utils.Compress([]byte(rqStr), 4)
			assert.Nil(t, err)
			rq = utils.B64Encode(compressedRq)

			task.Ticket.AppTicketData.DDAndFingerprintsIc = rand.Uint32()
			task.Ticket.AppTicketData.DDAndFingerprintsMax = 50
			task.Ticket.AppTicketData.RQIc = rand.Uint32()
			task.Ticket.AppTicketData.RQMax = 50

			task.Ticket.AppTicketData.DDAndFingerprintsIDs, _, err = pastel.GetIDFiles(context.Background(), []byte(ddStr),
				task.Ticket.AppTicketData.DDAndFingerprintsIc,
				task.Ticket.AppTicketData.DDAndFingerprintsMax)

			assert.Nil(t, err)

			task.Ticket.AppTicketData.RQIDs, _, err = pastel.GetIDFiles(context.Background(), []byte(rqStr),
				task.Ticket.AppTicketData.RQIc,
				task.Ticket.AppTicketData.RQMax)

			assert.Nil(t, err)

			err = task.validateRqIDsAndDdFpIds(context.Background(), rq, dd)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
