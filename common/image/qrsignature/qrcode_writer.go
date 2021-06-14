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
func (cap QRCodeCapacity) Validate(data []byte) error {
	if int(cap) < len(data) {
		return errors.Errorf("exceeds the maximum allowed amount of data by %d characters", len(data)-int(cap))
	}
	return nil
}

// QRCodeWriter represents QR code encoder.
type QRCodeWriter struct {
	writer    *qrcode.QRCodeWriter
	imageSize int
}

// Encode encodes given data into multiple QR codes.
func (qr *QRCodeWriter) Encode(data []byte) (image.Image, error) {
	img, err := qr.writer.Encode(string(data), gozxing.BarcodeFormat_QR_CODE, qr.imageSize, qr.imageSize, nil)
	if err != nil {
		return nil, errors.Errorf("failed to encode data: %w", err)
	}
	return img, nil
}

// NewQRCodeWriter returns a new QRCodeWriter instance.
func NewQRCodeWriter(imageSize int) *QRCodeWriter {
	return &QRCodeWriter{
		writer:    qrcode.NewQRCodeWriter(),
		imageSize: imageSize,
	}
}
