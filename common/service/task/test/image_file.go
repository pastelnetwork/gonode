package test

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/storage/fs"
)

//CreateBlankImage create blank png image for testing purpose only
func CreateBlankImage(filePath string, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	png.Encode(file, img)
	return nil
}

// NewTestImageFile create new image artwork.File instance for testing purpose
func NewTestImageFile(tmpDir, fileName string) (*artwork.File, error) {
	err := CreateBlankImage(fmt.Sprintf("%s/%s", tmpDir, fileName), 400, 400)
	if err != nil {
		return nil, err
	}

	imageStorage := artwork.NewStorage(fs.NewFileStorage(tmpDir))
	imageFile := artwork.NewFile(imageStorage, fileName)
	return imageFile, nil
}
