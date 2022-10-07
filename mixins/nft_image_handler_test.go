package mixins

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/base58"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/tj/assert"
)

func TestCreateCopyWithEncodedFingerprint(t *testing.T) {
	type args struct {
		findTicketReturns *pastel.IDTicket
		findTicketErr     error
		testPastelID      string
		fingerprints      Fingerprints
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				findTicketReturns: &pastel.IDTicket{
					IDTicketProp: pastel.IDTicketProp{
						PqKey: base58.Encode([]byte("pqkey")),
					},
				},
				testPastelID: "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
				fingerprints: Fingerprints{fingerprintAndScores: &pastel.DDAndFingerprints{}, signature: []byte("signature")},
			},
			wantErr: nil,
		},
		"error": {
			args: args{
				testPastelID:  "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
				findTicketErr: errors.New("test"),
			},
			wantErr: errors.New("test"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil).
				ListenOnFindTicketByID(tc.args.findTicketReturns, tc.args.findTicketErr)

			stg := files.NewStorage(fs.NewFileStorage(os.TempDir()))
			file, err := newTestImageFile(stg)
			assert.Nil(t, err)
			ph := NewPastelHandler(pastelClientMock)

			h := NewImageHandler(ph)
			err = h.CreateCopyWithEncodedFingerprint(context.Background(), tc.args.testPastelID, "passphrase", tc.args.fingerprints, file)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}

			h.ImageEncodedWithFingerprints = file
			_, err = h.GetHash()
			assert.Nil(t, err)
		})
	}
}

func TestMatchHashes(t *testing.T) {
	type args struct {
		preview []byte
		small   []byte
		medium  []byte
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				preview: []byte("preview-thumbnail"),
				small:   []byte("small-thumbnail"),
				medium:  []byte("medium-thumbnail"),
			},
			wantErr: nil,
		},
		"preview-error": {
			args: args{
				preview: []byte("preview-thumbnail"),
				small:   []byte("small-thumbnail"),
				medium:  []byte("medium-thumbnail"),
			},
			wantErr: errors.New("preview"),
		},
		"medium-error": {
			args: args{
				preview: []byte("preview-thumbnail"),
				small:   []byte("small-thumbnail"),
				medium:  []byte("medium-thumbnail"),
			},
			wantErr: errors.New("medium"),
		},
		"small-error": {
			args: args{
				preview: []byte("preview-thumbnail"),
				small:   []byte("small-thumbnail"),
				medium:  []byte("medium-thumbnail"),
			},
			wantErr: errors.New("small"),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			pastelClientMock := pastelMock.NewMockClient(t)
			pastelClientMock.ListenOnSign([]byte("signature"), nil)

			ph := NewPastelHandler(pastelClientMock)

			h := NewImageHandler(ph)
			h.AddHashes(&hashes{})
			h.ClearHashes()

			err := h.MatchThumbnailHashes()
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "wrong number"))

			h.AddNewHashes(tc.args.preview, tc.args.medium, tc.args.small, "pastel-id-a")
			h.AddNewHashes(tc.args.preview, tc.args.medium, tc.args.small, "pastel-id-b")
			if tc.wantErr != nil {
				if tc.wantErr.Error() == "preview" {
					h.AddNewHashes([]byte("invalid-preview"), tc.args.medium, tc.args.small, "pastel-id-c")
				} else if tc.wantErr.Error() == "medium" {
					h.AddNewHashes(tc.args.preview, []byte("invalid-medium"), tc.args.small, "pastel-id-c")
				} else if tc.wantErr.Error() == "small" {
					h.AddNewHashes(tc.args.preview, tc.args.medium, []byte("invalid-small"), "pastel-id-c")
				}
			} else {
				h.AddNewHashes(tc.args.preview, tc.args.medium, tc.args.small, "pastel-id-c")
			}

			err = h.MatchThumbnailHashes()
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.False(t, h.IsEmpty())
				assert.Nil(t, err)
			}
		})
	}
}

func newTestImageFile(stg *files.Storage) (*files.File, error) {
	imgFile := stg.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("create storage file: %w", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	if err := png.Encode(f, img); err != nil {
		return nil, err
	}

	return imgFile, nil
}
