package nftregister

import (
	"image"
	"io"

	"github.com/pastelnetwork/gonode/common/storage/files"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"

	"github/pastelnetwork/gonode/go-webp/encoder"
	"github/pastelnetwork/gonode/go-webp/webp"

	"github.com/disintegration/imaging"
	"golang.org/x/crypto/sha3"
)

type thumbnailType int

const (
	// previewThumbnail is best effort webp image from the original Nft
	previewThumbnail thumbnailType = iota

	// mediumThumbnail is wepb imnage from cropped image
	mediumThumbnail

	// smallhumbnail is snaller webp imnage from cropped image
	smallThumbnail

	fileSizeLimit int = 400000
)

func maxInt(x, y int) int {
	if x > y {
		return x
	}

	return y
}

// determine quality for preview image - with same size of orignal image
// the quality will vary from 10% -> 30%
func determinePreviewQuality(size int) float32 {
	if size <= 1000 {
		return 30.0
	}

	if size >= 2000 {
		return 10.0
	}

	return 10.0 + 20.0*float32(2000-size)/1000.0
}

// determine quality for medium image - with cropped size

func determineMediumQuality(size int) float32 {
	if size <= 100 {
		return 100.0
	}

	if size <= 400 {
		return 50.0
	}

	if size >= 2000 {
		return 10.0
	}

	if size >= 1000 {
		return 10.0 + 20.0*float32(2000-size)/1000.0
	}

	// other with size in (400 -> 100)
	return 30.0
}

// is a half of medium quality
func determineSmallQuality(size int) float32 {
	return determineMediumQuality(size) / 2
}

func (task *NftRegistrationTask) createAndHashThumbnails(coordinate files.ThumbnailCoordinate) (preview []byte, medium []byte, small []byte, err error) {
	srcImg, err := task.Nft.LoadImage()
	if err != nil {
		err = errors.Errorf("load image from Nft %s %w", task.Nft.Name(), err)
		return
	}

	originalW := srcImg.Bounds().Dx()
	originalH := srcImg.Bounds().Dy()

	// determine quality for preview thumbnail
	previewQuality := determinePreviewQuality(maxInt(originalW, originalH))
	preview, err = task.createAndHashThumbnail(srcImg, previewThumbnail, nil, previewQuality, fileSizeLimit)
	if err != nil {
		err = errors.Errorf("generate thumbnail %w", err)
		return
	}

	// determine user-selected coordinate
	rect := image.Rect(int(coordinate.TopLeftX), int(coordinate.TopLeftY), int(coordinate.BottomRightX), int(coordinate.BottomRightY))

	// determine quality for medium thumbnail
	cropW := int(coordinate.BottomRightX - coordinate.TopLeftX)
	cropH := int(coordinate.BottomRightY - coordinate.TopLeftY)

	mediumQuality := determineMediumQuality(maxInt(cropW, cropH))
	medium, err = task.createAndHashThumbnail(srcImg, mediumThumbnail, &rect, mediumQuality, fileSizeLimit)
	if err != nil {
		err = errors.Errorf("hash medium thumbnail failed %w", err)
		return
	}

	smallQuality := determineSmallQuality(maxInt(cropW, cropH))
	small, err = task.createAndHashThumbnail(srcImg, smallThumbnail, &rect, smallQuality, fileSizeLimit)
	if err != nil {
		err = errors.Errorf("hash small thumbnail failed %w", err)
		return
	}

	return
}

func (task *NftRegistrationTask) createAndHashThumbnail(srcImg image.Image, thumbnail thumbnailType, rect *image.Rectangle, quality float32, targetFileSize int) ([]byte, error) {
	f := task.Storage.NewFile()
	if f == nil {
		return nil, errors.Errorf("create thumbnail file failed")
	}

	previewFile, err := f.Create()
	if err != nil {
		return nil, errors.Errorf("create file %s: %w", f.Name(), err)
	}
	defer previewFile.Close()

	var thumbnailImg image.Image
	if rect != nil {
		thumbnailImg = imaging.Crop(srcImg, *rect)
	} else {
		thumbnailImg = srcImg
	}

	log.Debugf("Encode with target size %d and quality %f", targetFileSize, quality)
	encoderOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return nil, errors.Errorf("create lossless encoder option %w", err)
	}
	encoderOptions.TargetSize = targetFileSize

	if err := webp.Encode(previewFile, thumbnailImg, encoderOptions); err != nil {
		return nil, errors.Errorf("encode to webp format: %w", err)
	}
	log.Debugf("preview thumbnail %s", f.Name())

	previewFile.Seek(0, io.SeekStart)
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, previewFile); err != nil {
		return nil, errors.Errorf("hash failed: %w", err)
	}

	switch thumbnail {
	case previewThumbnail:
		task.PreviewThumbnail = f
	case mediumThumbnail:
		task.MediumThumbnail = f
	case smallThumbnail:
		task.SmallThumbnail = f
	}

	return hasher.Sum(nil), nil
}
