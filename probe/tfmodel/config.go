package tfmodel

import (
	"path/filepath"
	"strings"
)

// List of the available models.
var (
	EfficientNetB7    = NewConfig("EfficientNetB7.tf", "serving_default_input_1")
	EfficientNetB6    = NewConfig("EfficientNetB6.tf", "serving_default_input_2")
	InceptionResNetV2 = NewConfig("InceptionResNetV2.tf", "serving_default_input_3")
	DenseNet201       = NewConfig("DenseNet201.tf", "serving_default_input_4")
	InceptionV3       = NewConfig("InceptionV3.tf", "serving_default_input_5")
	NASNetLarge       = NewConfig("NASNetLarge.tf", "serving_default_input_6")
	ResNet152V2       = NewConfig("ResNet152V2.tf", "serving_default_input_7")
)

// Config represents information about the exported model.
type Config struct {
	name  string
	path  string
	input string
}

// NewConfig returns a new Config instantce.
func NewConfig(path, input string) Config {
	return Config{
		name:  strings.TrimSuffix(path, filepath.Ext(path)),
		path:  path,
		input: input,
	}
}
