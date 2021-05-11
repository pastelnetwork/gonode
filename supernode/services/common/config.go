package common

// Config contains common configuration of the servcies.
type Config struct {
	PastelID string `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{}
}
