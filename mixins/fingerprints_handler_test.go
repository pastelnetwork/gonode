package mixins

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestMatch(t *testing.T) {
	type args struct {
		fg1 Fingerprints
		fg2 Fingerprints
		fg3 Fingerprints

		verify    bool
		verifyErr error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				verify: true,
				fg1: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg2: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg3: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
			},
			wantErr: nil,
		},

		"match-failure": {
			args: args{
				verify: true,
				fg1: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness: &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
						Drawings: 0.9,
						Neutral:  0.1,
						Hentai:   0.6,
					},
				}},
				fg2: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness: &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
						Drawings: 0.9,
						Neutral:  0.1,
						Hentai:   0.6,
					},
				}},
				fg3: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness: &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
						Drawings: 0.1,
						Neutral:  0.9,
						Hentai:   0.2,
					},
				}},
			},
			wantErr: errors.New("not matched"),
		},

		"verification-failure": {
			args: args{
				verify: false,
				fg1: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg2: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg3: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
			},
			wantErr: errors.New("verification failed"),
		},

		"verify-err": {
			args: args{
				verify:    false,
				verifyErr: errors.New("test"),
				fg1: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg2: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
				fg3: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{
					InternetRareness:      &pastel.InternetRareness{},
					AlternativeNSFWScores: &pastel.AlternativeNSFWScores{},
				}},
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnVerify(tc.args.verify, tc.args.verifyErr)

			ph := NewPastelHandler(pastelClientMock)
			h := NewFingerprintsHandler(ph)

			h.Add(&Fingerprints{})
			h.Clear()

			bytes1, err := json.Marshal(tc.args.fg1.fingerprintAndScores)
			assert.Nil(t, err)
			h.AddNew(tc.args.fg1.fingerprintAndScores, bytes1, []byte("signature"), "pastelid-a")

			err = h.Match(context.Background())
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "wrong number"))

			bytes2, err := json.Marshal(tc.args.fg2.fingerprintAndScores)
			assert.Nil(t, err)
			h.AddNew(tc.args.fg2.fingerprintAndScores, bytes2, []byte("signature"), "pastelid-a")

			bytes3, err := json.Marshal(tc.args.fg3.fingerprintAndScores)
			assert.Nil(t, err)
			h.AddNew(tc.args.fg3.fingerprintAndScores, bytes3, []byte("signature"), "pastelid-a")

			err = h.Match(context.Background())
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				fmt.Println(err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)

				err = h.GenerateDDAndFingerprintsIDs(context.Background(), 5)
				assert.Nil(t, err)
				assert.False(t, h.IsEmpty())
			}
		})
	}
}
