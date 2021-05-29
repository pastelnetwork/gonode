package tensor

import (
	"path/filepath"
	"strings"
)

// List of the available models.
var (
	EfficientNetB7    = NewModelInfo("EfficientNetB7.tf", "serving_default_input_1")
	EfficientNetB6    = NewModelInfo("EfficientNetB6.tf", "serving_default_input_2")
	InceptionResNetV2 = NewModelInfo("InceptionResNetV2.tf", "serving_default_input_3")
	DenseNet201       = NewModelInfo("DenseNet201.tf", "serving_default_input_4")
	InceptionV3       = NewModelInfo("InceptionV3.tf", "serving_default_input_5")
	NASNetLarge       = NewModelInfo("NASNetLarge.tf", "serving_default_input_6")
	ResNet152V2       = NewModelInfo("ResNet152V2.tf", "serving_default_input_7")
)

// ModelInfo represents information about the exported model.
type ModelInfo struct {
	name  string
	path  string
	input string
}

// NewModelInfo returns a new ModelInfo instantce.
func NewModelInfo(path, input string) ModelInfo {
	return ModelInfo{
		name:  strings.TrimSuffix(path, filepath.Ext(path)),
		path:  path,
		input: input,
	}
}
