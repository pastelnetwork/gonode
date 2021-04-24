package register

const (
	defaultNumberSuperNodes = 3
)

// Config contains settings of the registering artwork.
type Config struct {
	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,
	}
}
