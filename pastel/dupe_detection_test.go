package pastel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProbeImageReply(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		dd  *DDAndFingerprints
		sig []byte
	}{
		"simple": {
			dd: &DDAndFingerprints{
				DupeDetectionSystemVersion: "1",
				Block:                      "block-hash",
				PastelRarenessScore:        0.55,
				IsLikelyDupe:               true,
				IsRareOnInternet:           true,
				InternetRarenessScore:      0.111,
				AlternateNSFWScores: &AlternativeNSFWScore{
					Sexy:    0.234,
					Hentai:  1.0,
					Drawing: 0.131,
					Porn:    0.9999,
				},
				Fingerprints: []float32{1.0, 2, 4, 3.3},
				Score: &DDScores{
					CombinedRarenessScore:         3.1,
					XgboostPredictedRarenessScore: 0.4,
				},
			},
			sig: []byte("signature"),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			reply, err := GetProbeImageReply(tc.dd, tc.sig)
			assert.Nil(t, err)

			gotDD, gotSig, err := GetDDFingerprintsAndSigFromProbeImageReply(reply)
			assert.Nil(t, err)
			assert.Equal(t, tc.sig, gotSig)
			assert.Equal(t, tc.dd.Block, gotDD.Block)
			assert.Equal(t, tc.dd.Score.XgboostPredictedRarenessScore, gotDD.Score.XgboostPredictedRarenessScore)
			assert.Equal(t, tc.dd.Score.CombinedRarenessScore, gotDD.Score.CombinedRarenessScore)
			assert.Equal(t, tc.dd.InternetRarenessScore, gotDD.InternetRarenessScore)
			assert.Equal(t, tc.dd.IsLikelyDupe, gotDD.IsLikelyDupe)
			assert.Equal(t, tc.dd.IsRareOnInternet, gotDD.IsRareOnInternet)
			assert.Equal(t, tc.dd.DupeDetectionSystemVersion, gotDD.DupeDetectionSystemVersion)
			assert.Equal(t, tc.dd.AlternateNSFWScores.Porn, gotDD.AlternateNSFWScores.Porn)
			assert.Equal(t, tc.dd.AlternateNSFWScores.Hentai, gotDD.AlternateNSFWScores.Hentai)
			assert.Equal(t, tc.dd.AlternateNSFWScores.Sexy, gotDD.AlternateNSFWScores.Sexy)
			assert.Equal(t, tc.dd.AlternateNSFWScores.Drawing, gotDD.AlternateNSFWScores.Drawing)

		})
	}
}
