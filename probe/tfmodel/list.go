package tfmodel

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
)

// List represents multiple TFModel.
type List []*TFModel

// Load loads models data from the baseDir.
func (models *List) Load(ctx context.Context, baseDir string) error {
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
func (models *List) Exec(ctx context.Context, value interface{}) ([][]float32, error) {
	fingerprints := make([][]float32, len(*models))

	if len(*models) == 0 {
		log.WithContext(ctx).Warn("Empty model list")
		return fingerprints, nil
	}

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

// NewList returns a new List instance.
func NewList(configs []Config) List {
	models := make(List, len(configs))

	for i, config := range configs {
		models[i] = NewModel(config)
	}
	return models
}
