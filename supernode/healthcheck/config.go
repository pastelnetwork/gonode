package healthcheck

import "time"

const (
	defaultPort           = 4445
	defaultUpdateDuration = 150 * time.Second // 2.5 minute
)

// Config contains settings of the supernode server.
type Config struct {
	Enable         bool          `mapstructure:"enable" json:"enable,omitempty"`
	LocalOnly      bool          `mapstructure:"local-only" json:"local-only,omitempty"`
	Port           int           `mapstructure:"port" json:"port,omitempty"`
	UpdateInterval time.Duration `mapstructure:"update-interval" json:"update-interval,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Enable:         false,
		LocalOnly:      true,
		Port:           defaultPort,
		UpdateInterval: defaultUpdateDuration,
	}
}
