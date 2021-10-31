package common

import "time"

const (
	defaultConnectToNodeTimeout = time.Second * 10
)

// Config contains common configuration of the servcies.
type Config struct {
	ConnectToNodeTimeout time.Duration
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		ConnectToNodeTimeout: defaultConnectToNodeTimeout,
	}
}
