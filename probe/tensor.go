//go:generate mockery --name=Tensor

package probe

import (
	"context"
	"image"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/probe/tfmodel"
)

const logTensorPrefix = "tensor"

// Tensor represents image analysis based on machine learning methods.
type Tensor interface {
	// Fingerprints computes and returns fingerprints for the given image by models.
	Fingerprints(ctx context.Context, img image.Image) (Fingerprints, error)

	// LoadModels loads all models.
	LoadModels(ctx context.Context) error
}

// tensor represents tensorflow models.
type tensor struct {
	tfModels tfmodel.List
	baseDir  string
}

// Fingerprints implements tensor.Tensor.Fingerprints
func (tensor *tensor) Fingerprints(ctx context.Context, img image.Image) (Fingerprints, error) {
	ctx = log.ContextWithPrefix(ctx, logTensorPrefix)

	var inputTensor [1][224][224][3]float32

	bounds := img.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// height = y and width = x
			inputTensor[0][y][x][0] = float32(r >> 8)
			inputTensor[0][y][x][1] = float32(g >> 8)
			inputTensor[0][y][x][2] = float32(b >> 8)
		}
	}

	fg, err := tensor.tfModels.Exec(ctx, inputTensor)
	if err != nil {
		return nil, err
	}

	fingerprints := make(Fingerprints, len(fg))
	for i, f := range fg {
		fingerprints[i] = f
	}
	return fingerprints, nil

}

// LoadModels implements tensor.Tensor.LoadModels
func (tensor *tensor) LoadModels(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logTensorPrefix)

	log.WithContext(ctx).Info("Loading models...")
	defer log.WithContext(ctx).WithDuration(time.Now()).Info("All models loaded")

	return tensor.tfModels.Load(ctx, tensor.baseDir)
}

// NewTensor returns a new Tensor interface implementation.
func NewTensor(baseDir string, tfModelConfigs []tfmodel.Config) Tensor {
	return &tensor{
		baseDir:  baseDir,
		tfModels: tfmodel.NewList(tfModelConfigs),
	}
}

func init() {
	// disable tensorflow logger
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")
}
