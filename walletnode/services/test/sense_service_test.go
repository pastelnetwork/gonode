package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	taskMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/common/storage/files"
	service2 "github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/sense_register"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	service "github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"github.com/stretchr/testify/assert"
)

func TestNodesSendImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		file *files.File
	}

	type nodeAttribute struct {
		address   string
		returnErr error
	}

	fingerAndScores := &pastel.DDAndFingerprints{
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

	testCompressedFingerAndScores, genErr := pastel.ToCompressSignedDDAndFingerprints(context.Background(), fingerAndScores, []byte("testSignature"))
	assert.Nil(t, genErr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		nodes                     []nodeAttribute
		args                      args
		err                       error
		errString                 string
		validBurnTxID             bool
		compressedFingersAndScore []byte
		numberProbeImageCall      int
	}{
		{
			nodes: []nodeAttribute{
				{"127.0.0.1:4444", nil},
				{"127.0.0.1:4445", nil},
				{"127.0.0.1:4446", nil},
			},
			args:                      args{ctx, &files.File{}},
			err:                       nil,
			errString:                 "",
			validBurnTxID:             true,
			compressedFingersAndScore: testCompressedFingerAndScores,
			numberProbeImageCall:      1,
		},
		{
			nodes: []nodeAttribute{
				{"127.0.0.1:4444", nil},
				{"127.0.0.1:4445", fmt.Errorf("failed to open stream")},
			},
			args:                      args{ctx, &files.File{}},
			err:                       fmt.Errorf("failed to open stream"),
			errString:                 "Unable to probe",
			validBurnTxID:             false,
			compressedFingersAndScore: testCompressedFingerAndScores,
			numberProbeImageCall:      1,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			nodes := common.SuperNodeList{}
			meshHandler := common.NewMeshHandlerSimple(nil, service.RegisterSenseNodeMaker{}, nil)

			for _, a := range testCase.nodes {
				//client mock
				client := test.NewMockClient(t)
				//listen on uploadImage call
				client.ListenOnProbeImage(testCase.compressedFingersAndScore, testCase.validBurnTxID, testCase.errString, false, testCase.err)

				someNode := common.NewSuperNode(client, a.address, "", service.RegisterSenseNodeMaker{})
				someNode.SuperNodeAPIInterface = &service.SenseRegistrationNode{RegisterSenseInterface: client.RegisterSenseInterface}

				//maker := service.RegisterSenseNodeMaker{}
				//someNode.SuperNodeAPIInterface = maker.MakeNode(client.Connection)
				nodes.Add(someNode)

			}
			meshHandler.Nodes = nodes

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(true, nil)
			pslHandler := service2.NewPastelHandler(pastelClientMock)

			taskClient := taskMock.NewMockTask(t)
			taskClient.ListenOnUpdateStatus()
			nodeTask := &common.WalletNodeTask{Task: taskClient}

			fpHandler := service2.NewFingerprintsHandler(pslHandler)
			srvTask := service.SenseRegistrationTask{WalletNodeTask: nodeTask}
			srvTask.MeshHandler = meshHandler
			srvTask.FingerprintsHandler = fpHandler

			err := srvTask.ProbeImage(testCase.args.ctx, testCase.args.file, "")
			if testCase.err != nil {
				assert.True(t, strings.Contains(err.Error(), testCase.err.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
