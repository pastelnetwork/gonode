package healthcheck

import "time"

const (
	defaultUpdateDuration = 150 * time.Second // 2.5 minute
)

// Config contains settings of the supernode server.
type Config struct {
	UpdateInterval time.Duration `mapstructure:"update-interval" json:"update-interval,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		UpdateInterval: defaultUpdateDuration,
	}
}
