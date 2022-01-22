package test

import (
	"context"
	"fmt"
	taskMock "github.com/pastelnetwork/gonode/common/service/task/test"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/sense_register"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	service "github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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

	testCases := []struct {
		nodes                     []nodeAttribute
		args                      args
		err                       error
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
			args:                      args{context.Background(), &files.File{}},
			err:                       nil,
			validBurnTxID:             true,
			compressedFingersAndScore: testCompressedFingerAndScores,
			numberProbeImageCall:      1,
		},
		{
			nodes: []nodeAttribute{
				{"127.0.0.1:4444", nil},
				{"127.0.0.1:4445", fmt.Errorf("failed to open stream")},
			},
			args:                      args{context.Background(), &files.File{}},
			err:                       fmt.Errorf("failed to open stream"),
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
			clients := []*test.Client{}

			meshHandler := mixins.NewMeshHandlerSimple(nil, service.RegisterSenseNodeMaker{})

			for _, a := range testCase.nodes {
				//client mock
				client := test.NewMockClient(t)
				//listen on uploadImage call
				client.ListenOnProbeImage(testCase.compressedFingersAndScore, testCase.validBurnTxID, testCase.err)
				clients = append(clients, client)

				someNode := common.NewSuperNode(client, a.address, "", service.RegisterSenseNodeMaker{})
				someNode.SuperNodeAPIInterface = &service.SenseRegisterNode{client.RegisterSense}

				//maker := service.RegisterSenseNodeMaker{}
				//someNode.SuperNodeAPIInterface = maker.MakeNode(client.Connection)
				nodes.Add(someNode)

			}
			meshHandler.Nodes = nodes

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(true, nil)
			pslHandler := mixins.NewPastelHandler(pastelClientMock)

			taskClient := taskMock.NewMockTask(t)
			taskClient.ListenOnUpdateStatus()
			nodeTask := &common.WalletNodeTask{Task: taskClient}

			fpHandler := mixins.NewFingerprintsHandler(nodeTask, pslHandler)
			srvHandler := service.NewSenseRegisterHandler(meshHandler, fpHandler)

			err := srvHandler.ProbeImage(testCase.args.ctx, testCase.args.file, "")
			if testCase.err != nil {
				assert.True(t, strings.Contains(err.Error(), testCase.err.Error()))
			} else {
				assert.Nil(t, err)
			}

			//mock assertion each client
			for _, client := range clients {
				client.RegisterSense.AssertExpectations(t)
				client.AssertProbeImageCall(testCase.numberProbeImageCall, testCase.args.ctx, testCase.args.file)
			}
		})
	}
}
