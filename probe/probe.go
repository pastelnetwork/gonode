package probe

import (
	"context"
	"image"
	"path/filepath"

	"github.com/pastelnetwork/gonode/probe/tensor"
)

const (
	tensorDir = "tensor"
)

type Probe interface {
	Fingerpints(ctx context.Context, img image.Image) ([][]float32, error)
	Run(ctx context.Context) error
}

type probe struct {
	*tensor.Tensor
	workDir string
}

func (service *probe) Run(ctx context.Context) error {
	tensorDir := filepath.Join(service.workDir, tensorDir)

	return service.Tensor.Models.Load(ctx, tensorDir)
}

func New(workDir string) Probe {
	return &probe{
		workDir: workDir,
		Tensor:  tensor.New(),
	}
}
