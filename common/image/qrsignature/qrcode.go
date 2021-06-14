package qrsignature

import (
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	// QRCodeCapacityBinary represents size for Binary
	// ISO 8859-1
	QRCodeCapacityBinary QRCodeCapacity = 2953
)

// QRCodeCapacity represents maximum character storage dataCapacity that can be encoded with a single QRCodeWriter-code.
// https://stackoverflow.com/questions/12764334/qr-code-max-char-length
type QRCodeCapacity int

// Validate validates a length of the data.
func (cap QRCodeCapacity) Validate(data string) error {
	if int(cap) < len(data) {
		return errors.Errorf("exceeds the maximum allowed amount of data by %d characters", len(data)-int(cap))
	}
	return nil
}

// QRCode represents QR code image and its coordinates on the canvas.
type QRCode struct {
	image.Image
	imageSize int
	X         int
	Y         int
	writer    *qrcode.QRCodeWriter
	reader    gozxing.Reader
}

// Decode reads data from the QR code.
func (qr *QRCode) Decode() (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(qr.Image)
	if err != nil {
		return "", errors.Errorf("failed to get bitmap from the QR code image :%w", err)
	}
	result, err := qr.reader.Decode(bmp, nil)
	if err != nil {
		return "", errors.Errorf("failed to decode a QR code :%w", err)
	}

	return result.GetText(), nil
}

// Encode generates a new QR code image by the given data.
func (qr *QRCode) Encode(data string) error {
	img, err := qr.writer.Encode(data, gozxing.BarcodeFormat_QR_CODE, qr.imageSize, qr.imageSize, nil)
	if err != nil {
		return errors.Errorf("failed to encode data: %w", err)
	}
	qr.Image = img
	return nil
}

// NewQRCode returns a new QRCode instance.
func NewQRCode(imageSize int) *QRCode {
	return &QRCode{
		imageSize: imageSize,
		X:         -1,
		Y:         -1,
		writer:    qrcode.NewQRCodeWriter(),
		reader:    qrcode.NewQRCodeReader(),
		Image: &image.NRGBA{
			Rect: image.Rectangle{
				Max: image.Point{
					X: imageSize,
					Y: imageSize,
				},
			},
		},
	}
}
