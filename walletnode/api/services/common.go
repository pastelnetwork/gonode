package services

import (
	"context"
	"fmt"
	"time"
	"math/rand"
)

var randIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))

// Config represents a config for the common service.
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

func randIDFunc() string {

	// Generate a random 6-digit ID
	uniqueID := randIDGen.Intn(1000000)
	return fmt.Sprintf("%06d", uniqueID)
}