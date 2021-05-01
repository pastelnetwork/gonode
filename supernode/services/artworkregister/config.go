package artworkregister

const (
	defaultNumberSecondaryNodes = 2
)

// Config contains settings of the registering artwork.
type Config struct {
	NumberSecondaryNodes int `mapstructure:"number_secondary_nodes" json:"number_secondary_nodes,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSecondaryNodes: defaultNumberSecondaryNodes,
	}
}
