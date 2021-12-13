package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareFingerPrintAndScore(t *testing.T) {
	genfingerAndScoresFunc := func() *DDAndFingerprints {
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
				UrlOfFirstMatchInPage:   "",
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
	lhs := genfingerAndScoresFunc()
	rhs1 := genfingerAndScoresFunc()
	assert.Nil(t, CompareFingerPrintAndScore(lhs, rhs1))

	rhs2 := genfingerAndScoresFunc()
	rhs2.Block = "newBlock"
	assert.NotNil(t, CompareFingerPrintAndScore(lhs, rhs2))

	rhs3 := genfingerAndScoresFunc()
	rhs3.OpenNSFWScore = 0.1000001
	assert.NotNil(t, CompareFingerPrintAndScore(lhs, rhs3))

	rhs4 := genfingerAndScoresFunc()
	rhs4.OpenNSFWScore = 0.10000001
	assert.Nil(t, CompareFingerPrintAndScore(lhs, rhs4))
}
