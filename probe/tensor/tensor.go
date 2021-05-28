package tensor

import (
	"context"
	"image"
)

type Tensor struct {
	Models
}

// Fingerpints computes fingerprints for the given image.
func (tensor *Tensor) Fingerpints(ctx context.Context, img image.Image) ([][]float32, error) {
	bounds := img.Bounds()

	var inputTensor [1][224][224][3]float32

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// height = y and width = x
			inputTensor[0][y][x][0] = float32(r >> 8)
			inputTensor[0][y][x][1] = float32(g >> 8)
			inputTensor[0][y][x][2] = float32(b >> 8)
		}
	}

	return tensor.Models.Exec(ctx, inputTensor)
}

func New() *Tensor {
	return &Tensor{
		Models: NewModels(
			EfficientNetB7,
			EfficientNetB6,
			InceptionResNetV2,
			DenseNet201,
			InceptionV3,
			NASNetLarge,
			ResNet152V2,
		),
	}
}
