package senseregister

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	ddMock "github.com/pastelnetwork/gonode/dupedetection/ddclient/test"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	test "github.com/pastelnetwork/gonode/supernode/node/test/sense_register"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/tj/assert"
)

func makeConnected(task *SenseRegistrationTask, status common.Status) *SenseRegistrationTask {
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

func add2NodesAnd2TicketSignatures(task *SenseRegistrationTask) *SenseRegistrationTask {
	task.NetworkHandler.Accepted = common.SuperNodePeerList{
		&common.SuperNodePeer{ID: "A"},
		&common.SuperNodePeer{ID: "B"},
	}
	task.PeersTicketSignature = map[string][]byte{"A": []byte{1, 2, 3}, "B": []byte{1, 2, 3}}
	return task
}

func makeEmptySenseRegTask(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client, rqClient rqnode.ClientInterface, ddClient ddclient.DDServerClient) *SenseRegistrationTask {
	service := NewService(config, fileStorage, pastelClient, nodeClient, p2pClient, ddClient)
	task := NewSenseRegistrationTask(service)
	task.Ticket = &pastel.ActionTicket{}
	task.ActionTicketRegMetadata = &types.ActionRegMetadata{
		EstimatedFee: 100,
	}

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
			clientMock.ListenOnSendSenseTicketSignature(tc.args.sendArtErr).
				ListenOnConnect("", nil).ListenOnRegisterSense()

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, clientMock, nil, nil, nil)

			task.NetworkHandler.ConnectedTo = &common.SuperNodePeer{
				ClientInterface: clientMock,
				NodeMaker:       RegisterSenseNodeMaker{},
			}
			err := task.NetworkHandler.ConnectedTo.Connect(context.Background())
			assert.Nil(t, err)

			err = task.signAndSendSenseTicket(context.Background(), tc.args.primary)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskRegisterAction(t *testing.T) {
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
			pastelClientMock.ListenOnRegisterActionTicket(tc.args.regRetID, tc.args.regErr)

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			id, err := task.registerAction(context.Background())
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

			IsOpenAPIRequest:           false,
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

			task := makeEmptySenseRegTask(&Config{}, fsMock, pastelClientMock, nil, nil, nil, ddmock)
			task = add2NodesAnd2TicketSignatures(task)
			task.ActionTicketRegMetadata = &types.ActionRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			_, err := task.GenFingerprintsData(context.Background(), file, "testBlockHash", "testBlockHeight", "2022-03-31 16:55:28", "creatorPastelID",
				"jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis", "jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F",
				"jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ", false, "", "")
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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)

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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task = add2NodesAnd2TicketSignatures(task)

			err := task.VerifyPeersTicketSignature(context.Background(), task.Ticket)
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 19,
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 20.3,
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
				Vout: []pastel.VoutTxnResult{pastel.VoutTxnResult{Value: 18,
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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)

			err := <-task.WaitConfirmation(ctx, tc.args.txid,
				tc.args.minConfirmations, tc.args.interval, false, 0, 0)
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
			pastelClientMock.ListenOnSign([]byte("signature"), nil).
				ListenOnGetActionFee(&pastel.GetActionFeesResult{SenseFee: 10,
					CascadeFee: 10}, nil).ListenOnVerify(true, nil)

			var clientMock *test.Client
			if tc.wantErr == nil {
				clientMock = test.NewMockClient(t)
				clientMock.ListenOnSendSignedDDAndFingerprints(nil).
					ListenOnConnect("", nil).ListenOnRegisterSense()
			}

			task := makeEmptySenseRegTask(serviceCfg, fsMock, pastelClientMock, clientMock, nil, nil, ddmock)
			task = add2NodesAnd2TicketSignatures(task)
			task = makeConnected(task, tc.args.status)

			task.ActionTicketRegMetadata = &types.ActionRegMetadata{BlockHash: "testBlockHash", CreatorPastelID: "creatorPastelID"}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go task.RunAction(ctx)

			if tc.wantErr == nil {
				nodes := pastel.MasterNodes{}
				nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
				nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})

				pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnMasterNodesExtra(nodes, nil)

				peerDDAndFingerprints, _ := pastel.ToCompressSignedDDAndFingerprints(genfingerAndScoresFunc(), []byte("signature"))
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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
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

func TestTaskAddPeerSenseTicketSignature(t *testing.T) {
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
				status:         common.StatusImageProbed,
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
				status:         common.StatusImageProbed,
				masterNodesErr: nil,
				nodeID:         "A",
				acceptedNodeID: "B",
			},
			wantErr: errors.New("not in Accepted list"),
		},
		"success-close-sign-chn": {
			args: args{
				status:         common.StatusImageProbed,
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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)
			task.UpdateStatus(tc.args.status)

			go task.RunAction(ctx)

			task.NetworkHandler.Accepted = common.SuperNodePeerList{
				&common.SuperNodePeer{ID: tc.args.acceptedNodeID, Address: tc.args.acceptedNodeID},
			}
			task.PeersTicketSignature = map[string][]byte{tc.args.acceptedNodeID: []byte{}}

			err := task.AddPeerTicketSignature(tc.args.nodeID, []byte{}, common.StatusImageProbed)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTaskValidateDdFpIds(t *testing.T) {
	type args struct {
		status common.Status
		fg     *pastel.DDAndFingerprints
		ddSig  [][]byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				ddSig: [][]byte{[]byte("sig-1"), []byte("sig-2"), []byte("sig-3")},

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

			task := makeEmptySenseRegTask(&Config{}, nil, pastelClientMock, nil, nil, nil, nil)

			meshedNodes := []types.MeshedSuperNode{
				types.MeshedSuperNode{NodeID: "node-1"},
				types.MeshedSuperNode{NodeID: "node-2"},
				types.MeshedSuperNode{NodeID: "node-3"},
			}

			task.UpdateStatus(common.StatusConnected)
			task.NetworkHandler.MeshNodes(context.TODO(), meshedNodes)
			task.UpdateStatus(tc.args.status)

			task.Ticket = &pastel.ActionTicket{
				Caller:        "author-pastelid",
				APITicketData: pastel.APISenseTicket{},
				ActionType:    pastel.ActionTypeSense,
			}

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
			ticketData := &pastel.APISenseTicket{
				DataHash: []byte{},
			}
			task.Ticket.APITicketData = ticketData

			assert.Nil(t, err)

			err = task.validateDdFpIds(context.Background(), dd)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
