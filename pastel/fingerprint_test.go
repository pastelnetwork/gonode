package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func genfingerAndScoresFunc() *DDAndFingerprints {
	return &DDAndFingerprints{
		Block:                      "Block",
		Principal:                  "Principal",
		DupeDetectionSystemVersion: "v1.0",

		IsLikelyDupe:     true,
		IsRareOnInternet: true,

		RarenessScores: &RarenessScores{
			CombinedRarenessScore:         0,
			XgboostPredictedRarenessScore: 0,
			NnPredictedRarenessScore:      0,
			OverallAverageRarenessScore:   0,
		},
		InternetRareness: &InternetRareness{
			MatchesFoundOnFirstPage: 0,
			NumberOfPagesOfResults:  0,
			URLOfFirstMatchInPage:   "",
		},

		OpenNSFWScore: 0.1,
		AlternativeNSFWScores: &AlternativeNSFWScores{
			Drawings: 0.1,
			Hentai:   0.2,
			Neutral:  0.3,
			Porn:     0.4,
			Sexy:     0.5,
		},

		ImageFingerprintOfCandidateImageFile: []float32{1, 2, 3},
		FingerprintsStat: &FingerprintsStat{
			NumberOfFingerprintsRequiringFurtherTesting1: 1,
			NumberOfFingerprintsRequiringFurtherTesting2: 2,
			NumberOfFingerprintsRequiringFurtherTesting3: 3,
			NumberOfFingerprintsRequiringFurtherTesting4: 4,
			NumberOfFingerprintsRequiringFurtherTesting5: 5,
			NumberOfFingerprintsRequiringFurtherTesting6: 6,
			NumberOfFingerprintsOfSuspectedDupes:         7,
		},

		HashOfCandidateImageFile: "HashOfCandidateImageFile",
		PerceptualImageHashes: &PerceptualImageHashes{
			PDQHash:        "PdqHash",
			PerceptualHash: "PerceptualHash",
			AverageHash:    "AverageHash",
			DifferenceHash: "DifferenceHash",
			NeuralHash:     "NeuralhashHash",
		},
		PerceptualHashOverlapCount: 1,

		Maxes: &Maxes{
			PearsonMax:           1.0,
			SpearmanMax:          2.0,
			KendallMax:           3.0,
			HoeffdingMax:         4.0,
			MutualInformationMax: 5.0,
			HsicMax:              6.0,
			XgbimportanceMax:     7.0,
		},
		Percentile: &Percentile{
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

// TestCompareFingerPrintAndScore test CompareFingerPrintAndScore()
func TestCompareFingerPrintAndScore(t *testing.T) {

	lhs := genfingerAndScoresFunc()
	rhs1 := genfingerAndScoresFunc()
	assert.Nil(t, CompareFingerPrintAndScore(lhs, rhs1))

	rhs2 := genfingerAndScoresFunc()
	rhs2.Block = "newBlock"
	assert.NotNil(t, CompareFingerPrintAndScore(lhs, rhs2))

	rhs3 := genfingerAndScoresFunc()
	rhs3.OpenNSFWScore = 0.101
	assert.NotNil(t, CompareFingerPrintAndScore(lhs, rhs3))

	rhs4 := genfingerAndScoresFunc()
	rhs4.OpenNSFWScore = 0.10000001
	assert.Nil(t, CompareFingerPrintAndScore(lhs, rhs4))
}

// TestCombineFingerPrintAndScores test CombineFingerPrintAndScores()
func TestCombineFingerPrintAndScores(t *testing.T) {
	a := genfingerAndScoresFunc()
	b := genfingerAndScoresFunc()
	c := genfingerAndScoresFunc()
	_, err := CombineFingerPrintAndScores(a, b, c)
	assert.Nil(t, err)

	a = genfingerAndScoresFunc()
	b = genfingerAndScoresFunc()
	c = genfingerAndScoresFunc()

	b.Block = "newBlock"
	_, err = CombineFingerPrintAndScores(a, b, c)
	assert.NotNil(t, err)

	a = genfingerAndScoresFunc()
	b = genfingerAndScoresFunc()
	c = genfingerAndScoresFunc()

	b.OpenNSFWScore = 0.101
	_, err = CombineFingerPrintAndScores(a, b, c)
	assert.NotNil(t, err)

	a = genfingerAndScoresFunc()
	b = genfingerAndScoresFunc()
	c = genfingerAndScoresFunc()

	// test IsRareOnInternet
	b.IsRareOnInternet = false
	final, err := CombineFingerPrintAndScores(a, b, c)
	assert.Nil(t, err)
	assert.Equal(t, true, final.IsRareOnInternet)

	a.IsRareOnInternet = false
	final, err = CombineFingerPrintAndScores(a, b, c)
	assert.Nil(t, err)
	assert.Equal(t, false, final.IsRareOnInternet)

	b.IsRareOnInternet = true
	final, err = CombineFingerPrintAndScores(a, b, c)
	assert.Nil(t, err)
	assert.Equal(t, true, final.IsRareOnInternet)
}
