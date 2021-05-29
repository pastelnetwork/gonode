package probe

import (
	"context"
	"image"
	"path/filepath"

	"github.com/pastelnetwork/gonode/probe/tensor"
)

const (
	tensorDir = "tensormodels"
)

// Probe represents image analysis techniques.
type Probe interface {
	// Fingerpint computes fingerprints for the given image.
	Fingerpint(ctx context.Context, img image.Image) ([][]float32, error)

	// Run bootstraps all services.
	Run(ctx context.Context) error
}

type probe struct {
	*tensor.Tensor
	workDir string
}

// Run implements probe.Probe.Run
func (service *probe) Run(ctx context.Context) error {
	tensorDir := filepath.Join(service.workDir, tensorDir)

	return service.Tensor.LoadModels(ctx, tensorDir)
}

// New returns a implemented Probe interface.
func New(workDir string) Probe {
	return &probe{
		workDir: workDir,
		Tensor:  tensor.New(),
	}
}
