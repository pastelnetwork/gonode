package pastel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCompressSignedDDAndFingerprints(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		dd  *DDAndFingerprints
		sig []byte
	}{
		"simple": {
			dd: &DDAndFingerprints{
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
			},
			sig: []byte("signature"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			reply, err := ToCompressSignedDDAndFingerprints(context.Background(), tc.dd, tc.sig)
			assert.Nil(t, err)

			gotDD, _, gotSig, err := ExtractCompressSignedDDAndFingerprints(reply)
			assert.Nil(t, err)

			assert.Equal(t, tc.sig, gotSig)
			assert.Equal(t, tc.dd.BlockHash, gotDD.BlockHash)
			assert.Equal(t, tc.dd.BlockHeight, gotDD.BlockHeight)
			assert.Equal(t, tc.dd.IsLikelyDupe, gotDD.IsLikelyDupe)
			assert.Equal(t, tc.dd.IsRareOnInternet, gotDD.IsRareOnInternet)
			assert.Equal(t, tc.dd.DupeDetectionSystemVersion, gotDD.DupeDetectionSystemVersion)
			assert.Equal(t, tc.dd.AlternativeNSFWScores.Porn, gotDD.AlternativeNSFWScores.Porn)
			assert.Equal(t, tc.dd.AlternativeNSFWScores.Hentai, gotDD.AlternativeNSFWScores.Hentai)
			assert.Equal(t, tc.dd.AlternativeNSFWScores.Sexy, gotDD.AlternativeNSFWScores.Sexy)
			assert.Equal(t, tc.dd.AlternativeNSFWScores.Drawings, gotDD.AlternativeNSFWScores.Drawings)

		})
	}
}
