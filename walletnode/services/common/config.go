package common

import "time"

const (
	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay = 400 * time.Millisecond
	defaultAcceptNodesTimeout     = 180 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout   = time.Second * 20
)

// Config contains common configuration of the services.
type Config struct {
	NumberSuperNodes int

	ConnectToNodeTimeout   time.Duration
	ConnectToNextNodeDelay time.Duration
	AcceptNodesTimeout     time.Duration
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,

		ConnectToNodeTimeout:   defaultConnectToNodeTimeout,
		ConnectToNextNodeDelay: defaultConnectToNextNodeDelay,
		AcceptNodesTimeout:     defaultAcceptNodesTimeout,
	}
}
