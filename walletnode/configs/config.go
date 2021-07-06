package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/raptorq"
)

const (
	defaultLogLevel = "info"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	LogLevel string `mapstructure:"log-level" json:"log-level,omitempty"`
	LogFile  string `mapstructure:"log-file" json:"log-file,omitempty"`
	Quiet    bool   `mapstructure:"quiet" json:"quiet"`
	TempDir  string `mapstructure:"temp-dir" json:"temp-dir"`
	WorkDir  string `mapstructure:"work-dir" json:"work-dir"`

	Node    `mapstructure:"node" json:"node,omitempty"`
	Pastel  *pastel.Config  `mapstructure:"pastel-api" json:"pastel-api,omitempty"`
	P2P     *p2p.Config     `mapstructure:"p2p" json:"p2p,omitempty"`
	RaptorQ *raptorq.Config `mapstructure:"raptorq" json:"raptorq,omitempty"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		LogLevel: defaultLogLevel,

		Node:    NewNode(),
		Pastel:  pastel.NewConfig(),
		P2P:     p2p.NewConfig(),
		RaptorQ: raptorq.NewConfig(),
	}
}
