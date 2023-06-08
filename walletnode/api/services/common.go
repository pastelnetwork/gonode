package services

import "context"

type Config struct {
	StaticFilesDir string `mapstructure:"static_files_dir" json:"static_files_dir,omitempty"`
}

// Common represents common service.
type Common struct {
	config *Config
}

// Run starts serving for operations that must be performed throughout the entire operation of the service.
func (service *Common) Run(_ context.Context) error {
	return nil
}

// NewCommon returns a new Common instance.
func NewCommon(config *Config) *Common {
	return &Common{
		config: config,
	}
}
