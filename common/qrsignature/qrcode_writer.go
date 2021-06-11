package qrsignature

import (
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pastelnetwork/gonode/common/collection"
	"github.com/pastelnetwork/gonode/common/errors"
)

// QRCodeCapacity represents maximum character storage dataCapacity that can be encoded with a single QRCodeWriter-code.
// https://stackoverflow.com/questions/12764334/qr-code-max-char-length
type QRCodeCapacity int

const (
	// QRCodeCapacityAlphanumeric represents size for Alphanumeric
	// 0–9, A–Z (upper-case only), space, $, %, *, +, -, ., /, :
	QRCodeCapacityAlphanumeric QRCodeCapacity = 4296

	// QRCodeCapacityBinary represents size for Binary
	// ISO 8859-1
	QRCodeCapacityBinary QRCodeCapacity = 2953
)

// QRCodeWriter represents QR code encoder.
type QRCodeWriter struct {
	dataCapacity QRCodeCapacity
	imageWidth   int
}

// Encode encodes given data into multiple QR codes.
func (qr *QRCodeWriter) Encode(data []byte) ([]image.Image, error) {
	imgsData := collection.SplitBytesBySize(data, int(qr.dataCapacity))
	imgs := make([]image.Image, len(imgsData))

	for i, data := range imgsData {
		qrEnc := qrcode.NewQRCodeWriter()
		img, err := qrEnc.Encode(string(data), gozxing.BarcodeFormat_QR_CODE, qr.imageWidth, qr.imageWidth, nil)
		if err != nil {
			return nil, errors.Errorf("failed to encode data: %w", err)
		}
		imgs[i] = img
	}
	return imgs, nil
}

// NewQRCodeWriter returns a new QRCodeWriter instance.
func NewQRCodeWriter(dataCapacity QRCodeCapacity, imageWidth int) *QRCodeWriter {
	return &QRCodeWriter{
		dataCapacity: dataCapacity,
		imageWidth:   imageWidth,
	}
}
