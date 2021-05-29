package tensor

import (
	"context"
	"image"
	"os"

	"github.com/pastelnetwork/gonode/common/log"
)

const logPrefix = "tensor"

// Tensor represents tensorflow models.
type Tensor struct {
	models Models
}

// Fingerpint implements probe.Probe.Fingerpint
func (tensor *Tensor) Fingerpint(ctx context.Context, img image.Image) ([][]float32, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

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

	return tensor.models.Exec(ctx, inputTensor)
}

// LoadModels loads models.
func (tensor *Tensor) LoadModels(ctx context.Context, baseDir string) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	return tensor.models.Load(ctx, baseDir)
}

// New returns a new Tensor instance.
func New() *Tensor {
	return &Tensor{
		models: NewModels(
			EfficientNetB7,
			EfficientNetB6,
			// InceptionResNetV2,
			// DenseNet201,
			// InceptionV3,
			// NASNetLarge,
			// ResNet152V2,
		),
	}
}

func init() {
	// disable tensorflow logger
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")
}
