// Package dupedetection provides functions to compute dupe detection fingerprints for specific image
package dupedetection

import (
	"image"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	pruntime "github.com/pastelnetwork/gonode/common/runtime"

	"github.com/disintegration/imaging"
	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
)

var models = make(map[string]*tg.Model)

type modelData struct {
	model string
	input string
}

var fingerprintSources = []modelData{
	{
		model: "EfficientNetB7.tf",
		input: "serving_default_input_1",
	},
	{
		model: "EfficientNetB6.tf",
		input: "serving_default_input_2",
	},
	{
		model: "InceptionResNetV2.tf",
		input: "serving_default_input_3",
	},
	{
		model: "DenseNet201.tf",
		input: "serving_default_input_4",
	},
	{
		model: "InceptionV3.tf",
		input: "serving_default_input_5",
	},
	{
		model: "NASNetLarge.tf",
		input: "serving_default_input_6",
	},
	{
		model: "ResNet152V2.tf",
		input: "serving_default_input_7",
	},
}

func tgModel(path string) *tg.Model {
	m, ok := models[path]
	if !ok {
		m = tg.LoadModel(path, []string{"serve"}, nil)
		models[path] = m
	}
	return m
}

func fromFloat32To64(input []float32) []float64 {
	output := make([]float64, len(input))
	for i, value := range input {
		output[i] = float64(value)
	}
	return output
}

func loadImage(imagePath string, width int, height int) (image.Image, error) {
	reader, err := os.Open(imagePath)
	if err != nil {
		return nil, errors.New(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.New(err)
	}

	img = imaging.Resize(img, width, height, imaging.Linear)
	return img, nil
}

// ComputeImageDeepLearningFeatures computes dupe detection fingerprints for image with imagePath
func ComputeImageDeepLearningFeatures(imagePath string) ([][]float64, error) {
	defer pruntime.PrintExecutionTime(time.Now())

	m, err := loadImage(imagePath, 224, 224)
	if err != nil {
		return nil, errors.New(err)
	}

	bounds := m.Bounds()

	var inputTensor [1][224][224][3]float32

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// height = y and width = x
			inputTensor[0][y][x][0] = float32(r >> 8)
			inputTensor[0][y][x][1] = float32(g >> 8)
			inputTensor[0][y][x][2] = float32(b >> 8)
		}
	}

	fingerprints := make([][]float64, len(fingerprintSources))
	for i, source := range fingerprintSources {
		model := tgModel(source.model)

		fakeInput, _ := tf.NewTensor(inputTensor)
		results := model.Exec([]tf.Output{
			model.Op("StatefulPartitionedCall", 0),
		}, map[tf.Output]*tf.Tensor{
			model.Op(source.input, 0): fakeInput,
		})

		predictions := results[0].Value().([][]float32)[0]
		fingerprints[i] = fromFloat32To64(predictions)
	}

	return fingerprints, nil
}
