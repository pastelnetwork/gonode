package qrsignature

import (
	"fmt"
	"image"
)

const (
	dataQRWidth    = 100
	dataQRCapacity = QRCodeCapacityBinary
)

type DataEncoder interface {
	Encode(data []byte) ([]image.Image, error)
}

type Data struct {
	raw      []byte
	dataType DataType

	encoder DataEncoder
	images  []image.Image
}

func (data *Data) Encode() error {
	imgs, err := data.encoder.Encode(data.raw)
	if err != nil {
		return fmt.Errorf("%q: %w", data.dataType, err)
	}
	data.images = imgs
	return nil
}

func NewData(raw []byte, dataType DataType) *Data {
	return &Data{
		raw:      raw,
		dataType: dataType,

		encoder: NewQRCodeWriter(dataQRCapacity, dataQRWidth),
	}
}
