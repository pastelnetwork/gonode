package tensor

import "context"

type Models []*Model

func (models *Models) Load(ctx context.Context, baseDir string) error {
	for _, model := range *models {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := model.Load(baseDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func (models *Models) Exec(ctx context.Context, value interface{}) ([][]float32, error) {
	var fingerprints [][]float32

	for i, model := range *models {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			predictions, err := model.Exec(value)
			if err != nil {
				return nil, err
			}
			fingerprints[i] = predictions
		}
	}
	return fingerprints, nil
}

func NewModels(names ...ModelName) Models {
	var models Models
	for _, name := range names {
		models = append(models, NewModel(name))
	}
	return models
}
