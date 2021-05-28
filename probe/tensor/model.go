package tensor

import (
	"path/filepath"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	"github.com/pastelnetwork/gonode/common/errors"
)

type Model struct {
	ModelName
	data *tf.SavedModel
}

// String implements fmt.Stringer.String
func (model *Model) String() string {
	return model.path
}

// Load loads model data from the baseDir.
// The graph loaded is identified by the set of tags specified when exporting it.
// This operation creates a session with specified `options`.
func (model *Model) Load(baseDir string) error {
	modelPath := filepath.Join(baseDir, model.path)

	data, err := tf.LoadSavedModel(modelPath, []string{"serve"}, nil)
	if err != nil {
		return errors.Errorf("failed to load tensor model %q: %w", modelPath, err)
	}
	model.data = data

	return nil
}

// Exec executes the nodes/tensors that must be present in the loaded model.
// feedDict values to feed to placeholders (that must have been saved in the model definition).
func (model *Model) Exec(value interface{}) ([]float32, error) {
	fetcheOutput, err := model.operation("StatefulPartitionedCall", 0)
	if err != nil {
		return nil, err
	}
	fetches := []tf.Output{fetcheOutput}

	tensor, err := tf.NewTensor(value)
	if err != nil {
		return nil, errors.Errorf("failed to create new tensor: %w", err)
	}
	feedOutput, err := model.operation(model.input, 0)
	if err != nil {
		return nil, err
	}
	feeds := map[tf.Output]*tf.Tensor{
		feedOutput: tensor,
	}

	results, err := model.data.Session.Run(feeds, fetches, nil)
	if err != nil {
		return nil, errors.Errorf("failed to run model graph %q: %w", model.path, err)
	}
	return results[0].Value().([][]float32)[0], nil
}

// Operation extracts the output in position idx of the tensor with the specified name from the model graph
func (model *Model) operation(name string, idx int) (tf.Output, error) {
	op := model.data.Graph.Operation(name)
	if op == nil {
		return tf.Output{}, errors.Errorf("op %q not found", name)
	}
	nout := op.NumOutputs()
	if nout <= idx {
		return tf.Output{}, errors.Errorf("op %q has %d outputs. Requested output number %d", name, nout, idx)
	}
	return op.Output(idx), nil
}

func NewModel(name ModelName) *Model {
	return &Model{
		ModelName: name,
	}
}
