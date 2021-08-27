package artworkregister

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	storageMock "github.com/pastelnetwork/gonode/common/storage/test"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/tj/assert"
)

func newTestImageFile() (*artwork.File, error) {
	imageStorage := artwork.NewStorage(fs.NewFileStorage(os.TempDir()))
	imgFile := imageStorage.NewFile()

	f, err := imgFile.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create storage file: %w", err)
	}
	defer f.Close()
	f.Write([]byte("hello"))
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	png.Encode(f, img)

	_, err = imgFile.File
	if err != nil {
		return nil, err
	}

	return imgFile, nil
}
func TestCreateAndHashThumbnail(t *testing.T) {
	type args struct {
		task     *Task
		storeErr error
		fileErr  error
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				task: &Task{
					Service: &Service{
						config: &Config{},
					},
					Ticket: &pastel.NFTTicket{},
				},
				fileErr:  nil,
				storeErr: nil,
			},
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			fsMock := storageMock.NewMockFileStorage()
			fileMock := storageMock.NewMockFile()
			fileMock.ListenOnClose(nil).ListenOnRead(0, io.EOF).ListenOnName("test")

			storage := artwork.NewStorage(fsMock)

			tc.args.task.Storage = storage
			file, err := newTestImageFile()
			assert.Nil(t, err)
			tc.args.task.Artwork = file

			fsMock.ListenOnOpen(file, tc.args.fileErr).ListenOnCreate(file, nil)

			srcImg, err := tc.args.task.Artwork.LoadImage()
			if err != nil {
				fmt.Println("errrr: ", err.Error())
			}
			assert.Nil(t, err)

			originalW := srcImg.Bounds().Dx()
			originalH := srcImg.Bounds().Dy()
			// determine quality for preview thumbnail
			previewQuality := determinePreviewQuality(maxInt(originalW, originalH))

			_, err = tc.args.task.createAndHashThumbnail(srcImg, previewThumbnail, nil, previewQuality, fileSizeLimit)
			if tc.wantErr != nil {
				if err != nil {
					fmt.Println("errrr: ", err.Error())
				}
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), tc.wantErr.Error()))
			} else {
				if err != nil {
					fmt.Println("errrrr: ", err.Error())
				}
				assert.Nil(t, err)
			}
		})
	}
}
