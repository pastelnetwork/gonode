package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/probe/pkg/dupedetection"
)

type Config struct {
	DupeDetection *dupedetection.Config `mapstructure:"dupe-detection" json:"dupe-detection,omitempty"`
}

func (config *Config) String() string {
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		DupeDetection: dupedetection.NewConfig(),
	}
}
