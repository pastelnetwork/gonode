package download

import "time"

const (
	defaultConnectToNextNodeDelay    = 400 * time.Millisecond
	defaultAcceptNodesTimeout        = 600 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout      = time.Second * 20
	defaultConnectionRefreshInterval = 2
	defaultNumberOfConnections       = 7
)

// Config contains settings of the registering nft.
type Config struct {
	PastelID                  string `mapstructure:"pastel_id" json:"pastel_id"`
	Passphrase                string `mapstructure:"passphrase" json:"passphrase"`
	ConnectionRefreshInterval int    `mapstructure:"connection_refresh_interval" json:"connection_refresh_interval"`
	Connections               int    `mapstructure:"connections" json:"connections"`

	ConnectToNodeTimeout   time.Duration
	ConnectToNextNodeDelay time.Duration
	AcceptNodesTimeout     time.Duration
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		Connections:               defaultNumberOfConnections,
		ConnectionRefreshInterval: defaultConnectionRefreshInterval,
		ConnectToNodeTimeout:      defaultConnectToNodeTimeout,
		ConnectToNextNodeDelay:    defaultConnectToNextNodeDelay,
		AcceptNodesTimeout:        defaultAcceptNodesTimeout,
	}
}
