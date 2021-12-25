package pastel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
			},
			sig: []byte("signature"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			reply, err := ToCompressSignedDDAndFingerprints(tc.dd, tc.sig)
			assert.Nil(t, err)

			gotDD, _, gotSig, err := ExtractCompressSignedDDAndFingerprints(reply)
			assert.Nil(t, err)
			if err == nil {
				assert.Equal(t, tc.sig, gotSig)
				assert.Equal(t, tc.dd.Block, gotDD.Block)
				assert.Equal(t, tc.dd.Principal, gotDD.Principal)
				assert.Equal(t, tc.dd.IsLikelyDupe, gotDD.IsLikelyDupe)
				assert.Equal(t, tc.dd.IsRareOnInternet, gotDD.IsRareOnInternet)
				assert.Equal(t, tc.dd.DupeDetectionSystemVersion, gotDD.DupeDetectionSystemVersion)
				assert.Equal(t, tc.dd.AlternativeNSFWScores.Porn, gotDD.AlternativeNSFWScores.Porn)
				assert.Equal(t, tc.dd.AlternativeNSFWScores.Hentai, gotDD.AlternativeNSFWScores.Hentai)
				assert.Equal(t, tc.dd.AlternativeNSFWScores.Sexy, gotDD.AlternativeNSFWScores.Sexy)
				assert.Equal(t, tc.dd.AlternativeNSFWScores.Drawings, gotDD.AlternativeNSFWScores.Drawings)
			} else {
				fmt.Println("err: ", err.Error())
			}
		})
	}
}

func TestGetIDFiles(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		dd      *DDAndFingerprints
		wantErr error
	}{
		"simple": {
			dd: &DDAndFingerprints{
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
			},
			wantErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		ddJSON, err := json.Marshal(tc.dd)
		assert.Nil(t, err)

		encoded := base64.StdEncoding.EncodeToString(ddJSON)
		data := encoded + ".Sig1.Sig2.Sig3"

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, _, err := GetIDFiles([]byte(data), 12, 50)
			if tc.wantErr == nil {
				assert.Nil(t, err)
			}
		})
	}
}
