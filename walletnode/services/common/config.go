package common

import "time"

const (
	connectTimeout = time.Second * 2
)

// Config contains common configuration of the servcies.
type Config struct {
	ConnectTimeout time.Duration
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		ConnectTimeout: connectTimeout,
	}
}
