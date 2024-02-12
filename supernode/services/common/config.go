package common

const (
	defaultNumberSuperNodes = 10
)

// Config contains common configuration of the services.
type Config struct {
	PastelID    string `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
	PassPhrase  string `mapstructure:"pass_phrase" json:"-"`
	NodeAddress string

	NumberSuperNodes int
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,
	}
}
