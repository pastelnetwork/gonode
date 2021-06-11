package qrsignature

import (
	"fmt"
	"image"

	"github.com/pastelnetwork/gonode/common/service/artwork"
)

const (
// positionImageSize is size in pixels for the QR-code of the position vector.
// positionImageSize    = 185
// positionDataCapacity = QRCodeCapacityAlphanumeric
)

type Encoder struct {
	file *artwork.File
}

func (enc Encoder) Encode(signatures Signatures) (image.Image, error) {
	for _, signature := range signatures {
		if err := signature.Encode(); err != nil {
			return nil, err
		}
	}

	img, err := enc.file.LoadImage()
	if err != nil {
		return nil, err
	}

	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	newW, newH := signatures.minArtImageDimension(srcW, srcH)

	fmt.Println(newW, "-------", newH)
	// if signatures.imageOverflow(enc.img) {
	// 	return nil, errors.New("image is overflowing with signatures")
	// }

	return nil, nil
}

func NewEncoder(file *artwork.File) *Encoder {
	return &Encoder{
		file: file,
	}
}
