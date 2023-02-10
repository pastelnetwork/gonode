package nftregister

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/supernode/services/common"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/tj/assert"
)

func TestDeterminePreviewQuality(t *testing.T) {
	testCases := map[string]struct {
		size int
		want float32
	}{
		"lower": {
			size: 1000,
			want: 30.0,
		},
		"upper": {
			size: 2000,
			want: 10.0,
		},
		"lower-end": {
			size: 999,
			want: 30.0,
		},
		"bw": {
			size: 1500,
			want: 10.0 + 20.0*float32(2000-1500)/1000.0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := determinePreviewQuality(tc.size)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDetermineMediumQuality(t *testing.T) {
	testCases := map[string]struct {
		size int
		want float32
	}{
		"a": {
			size: 100,
			want: 100,
		},
		"b": {
			size: 399,
			want: 50.0,
		},
		"c": {
			size: 2000,
			want: 10.0,
		},
		"d": {
			size: 1000,
			want: 10.0 + 20.0*float32(2000-1000)/1000.0,
		},
		"e": {
			size: 401,
			want: 30.0,
		},
		"f": {
			size: 999,
			want: 30.0,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := determineMediumQuality(tc.size)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMaxInt(t *testing.T) {
	testCases := map[string]struct {
		x, y, want int
	}{
		"x": {
			x:    1000,
			y:    30,
			want: 1000,
		},
		"y": {
			x:    312,
			y:    313,
			want: 313,
		},
		"equal": {
			x:    786,
			y:    786,
			want: 786,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()
			got := maxInt(tc.x, tc.y)
			assert.Equal(t, tc.want, got)
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
func TestCreateAndHashThumbnail(t *testing.T) {
	type args struct {
		task *NftRegistrationTask
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &NftRegistrationTask{
					NftRegistrationService: &NftRegistrationService{
						config:           &Config{},
						SuperNodeService: &common.SuperNodeService{},
					},
					Ticket: &pastel.NFTTicket{},
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			stg := files.NewStorage(fs.NewFileStorage(os.TempDir()))

			tc.args.task.Storage = stg
			file, err := newTestImageFile(stg)
			assert.Nil(t, err)
			tc.args.task.Nft = file

			srcImg, err := tc.args.task.Nft.LoadImage()
			if err != nil {
				fmt.Println("error: ", err.Error())
			}
			assert.Nil(t, err)

			originalW := srcImg.Bounds().Dx()
			originalH := srcImg.Bounds().Dy()
			// determine quality for preview thumbnail
			previewQuality := determinePreviewQuality(maxInt(originalW, originalH))

			_, err = tc.args.task.createAndHashThumbnail(srcImg, previewThumbnail, nil, previewQuality, fileSizeLimit)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateAndHashThumbnails(t *testing.T) {
	type args struct {
		task *NftRegistrationTask
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &NftRegistrationTask{
					NftRegistrationService: &NftRegistrationService{
						config:           &Config{},
						SuperNodeService: &common.SuperNodeService{},
					},
					Ticket: &pastel.NFTTicket{},
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			stg := files.NewStorage(fs.NewFileStorage(os.TempDir()))

			tc.args.task.Storage = stg
			file, err := newTestImageFile(stg)
			assert.Nil(t, err)
			tc.args.task.Nft = file

			coordinate := files.ThumbnailCoordinate{
				TopLeftX:     0,
				TopLeftY:     0,
				BottomRightX: 400,
				BottomRightY: 400,
			}
			_, _, _, err = tc.args.task.createAndHashThumbnails(coordinate)
			if tc.wantErr != nil {
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
