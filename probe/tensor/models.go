package tensor

import (
	"context"
)

// Models represents multiple Model.
type Models []*Model

// Load loads models data from the baseDir.
func (models *Models) Load(ctx context.Context, baseDir string) error {
	for _, model := range *models {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := model.Load(ctx, baseDir); err != nil {
				return err
			}
		}
	}
	return nil
}

// Exec executes the nodes/tensors that must be present in the loaded models.
func (models *Models) Exec(ctx context.Context, value interface{}) ([][]float32, error) {
	var fingerprints [][]float32

	for i, model := range *models {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			predictions, err := model.Exec(ctx, value)
			if err != nil {
				return nil, err
			}
			fingerprints[i] = predictions
		}
	}
	return fingerprints, nil
}

// NewModels returns a new Models instance.
func NewModels(infos ...ModelInfo) Models {
	var models Models
	for _, info := range infos {
		models = append(models, NewModel(info))
	}
	return models
}
