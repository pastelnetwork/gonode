package tfmodel

import (
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/common/sys"
)

// AllConfigs contains configs of all models.
var AllConfigs = []Config{
	EfficientNetB7,
	EfficientNetB6,
	InceptionResNetV2,
	DenseNet201,
	InceptionV3,
	NASNetLarge,
	ResNet152V2,
}

// Configs of the available models.
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

func init() {
	// NOTE: Reduce loaded models, where `0` means not to load at all. Used for integration testing and testing at the developer stage.
	if number := sys.GetIntEnv("LOAD_TFMODELS", -1); number >= 0 {
		AllConfigs = AllConfigs[:number]
	}
}
