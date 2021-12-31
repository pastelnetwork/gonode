package node

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/pastel"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/register_sense"
	"github.com/stretchr/testify/assert"
)

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *Node
	}
	testCases := []struct {
		nodes List
		args  args
		want  List
	}{
		{
			nodes: List{},
			args:  args{node: &Node{address: "127.0.0.1"}},
			want: List{
				&Node{address: "127.0.0.1"},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(testCase.args.node)
			assert.Equal(t, testCase.want, testCase.nodes)
		})
	}
}

func TestNodesActivate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
	}{
		{
			nodes: List{
				&Node{address: "127.0.0.1"},
				&Node{address: "127.0.0.2"},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Activate()
			for _, n := range testCase.nodes {
				assert.True(t, n.activated)
			}
		})
	}
}

func TestNodesDisconnectInactive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
		conn  []struct {
			client       *test.Client
			activated    bool
			numberOfCall int
		}
	}{
		{
			nodes: List{},
			conn: []struct {
				client       *test.Client
				activated    bool
				numberOfCall int
			}{
				{
					client:       test.NewMockClient(t),
					numberOfCall: 1,
					activated:    false,
				},
				{
					client:       test.NewMockClient(t),
					numberOfCall: 0,
					activated:    true,
				},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			for _, c := range testCase.conn {
				c.client.ListenOnClose(nil)

				node := &Node{
					Connection: c.client.Connection,
					activated:  c.activated,
					mtx:        &sync.RWMutex{},
				}

				testCase.nodes = append(testCase.nodes, node)
			}

			testCase.nodes.DisconnectInactive()

			for j, c := range testCase.conn {
				c := c

				t.Run(fmt.Sprintf("close-called-%d", j), func(t *testing.T) {
					c.client.AssertCloseCall(c.numberOfCall)
				})

			}

		})
	}

}

func TestNodesFindByPastelID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}
	testCases := []struct {
		nodes List
		args  args
		want  *Node
	}{
		{
			nodes: List{
				&Node{pastelID: "1"},
				&Node{pastelID: "2"},
			},
			args: args{"2"},
			want: &Node{pastelID: "2"},
		}, {
			nodes: List{
				&Node{pastelID: "1"},
				&Node{pastelID: "2"},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, testCase.nodes.FindByPastelID(testCase.args.id))
		})
	}
}

func TestNodesSendImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		file *artwork.File
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
		compressedFingersAndScore []byte
		numberProbeImageCall      int
	}{
		{
			nodes:                     []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", nil}},
			args:                      args{context.Background(), &artwork.File{}},
			err:                       nil,
			compressedFingersAndScore: testCompressedFingerAndScores,
			numberProbeImageCall:      1,
		},
		{
			nodes:                     []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", fmt.Errorf("failed to open stream")}},
			args:                      args{context.Background(), &artwork.File{}},
			err:                       fmt.Errorf("failed to open stream"),
			compressedFingersAndScore: testCompressedFingerAndScores,
			numberProbeImageCall:      1,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			nodes := List{}
			clients := []*test.Client{}

			for _, a := range testCase.nodes {
				//client mock
				client := test.NewMockClient(t)
				//listen on uploadImage call
				client.ListenOnProbeImage(testCase.compressedFingersAndScore, testCase.err)
				clients = append(clients, client)

				nodes.Add(&Node{
					address:       a.address,
					RegisterSense: client.RegisterSense,
				})
			}

			err := nodes.ProbeImage(testCase.args.ctx, testCase.args.file)
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
