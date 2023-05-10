package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func genfingerAndScoresFunc() *DDAndFingerprints {
	return &DDAndFingerprints{
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

		InternetRareness: &InternetRareness{
			RareOnInternetSummaryTableAsJSONCompressedB64:    "RareOnInternetSummaryTableAsJSONCompressedB64",
			RareOnInternetGraphJSONCompressedB64:             "RareOnInternetGraphJSONCompressedB64",
			AlternativeRareOnInternetDictAsJSONCompressedB64: "AlternativeRareOnInternetDictAsJSONCompressedB64",
			MinNumberOfExactMatchesInPage:                    4,
			EarliestAvailableDateOfInternetResults:           "EarliestAvailableDateOfInternetResults",
		},

		OpenNSFWScore: 0.1,
		AlternativeNSFWScores: &AlternativeNSFWScores{
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

// TestCompareFingerPrintAndScore test CompareFingerPrintAndScore()
func TestCompareFingerPrintAndScore(t *testing.T) {

	lhs := genfingerAndScoresFunc()
	rhs1 := genfingerAndScoresFunc()
	assert.Nil(t, CompareFingerPrintAndScore(lhs, rhs1))

	rhs2 := genfingerAndScoresFunc()
	rhs2.BlockHash = "newBlock"
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

	b.BlockHeight = "newBlock"
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
