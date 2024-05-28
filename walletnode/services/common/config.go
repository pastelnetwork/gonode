package common

import "time"

const (
	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay       = 600 * time.Millisecond
	defaultAcceptNodesTimeout           = 600 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout         = time.Second * 25
	defaultDownloadConnectToNodeTimeout = time.Second * 10
	defaultHashCheckMaxRetries          = 2
)

// Config contains common configuration of the services.
type Config struct {
	NumberSuperNodes int

	ConnectToNodeTimeout         time.Duration
	ConnectToNextNodeDelay       time.Duration
	AcceptNodesTimeout           time.Duration
	HashCheckMaxRetries          int
	DownloadConnectToNodeTimeout time.Duration
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes:             defaultNumberSuperNodes,
		DownloadConnectToNodeTimeout: defaultDownloadConnectToNodeTimeout,

		ConnectToNodeTimeout:   defaultConnectToNodeTimeout,
		ConnectToNextNodeDelay: defaultConnectToNextNodeDelay,
		AcceptNodesTimeout:     defaultAcceptNodesTimeout,
		HashCheckMaxRetries:    defaultHashCheckMaxRetries,
	}
}
