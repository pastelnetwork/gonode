package tensor

var (
	EfficientNetB7    = NewModelName("EfficientNetB7.tf", "serving_default_input_1")
	EfficientNetB6    = NewModelName("EfficientNetB6.tf", "serving_default_input_2")
	InceptionResNetV2 = NewModelName("InceptionResNetV2.tf", "serving_default_input_3")
	DenseNet201       = NewModelName("DenseNet201.tf", "serving_default_input_4")
	InceptionV3       = NewModelName("InceptionV3.tf", "serving_default_input_5")
	NASNetLarge       = NewModelName("NASNetLarge.tf", "serving_default_input_6")
	ResNet152V2       = NewModelName("ResNet152V2.tf", "serving_default_input_7")
)

type ModelName struct {
	path  string
	input string
}

func NewModelName(path, input string) ModelName {
	return ModelName{
		path:  path,
		input: input,
	}
}
