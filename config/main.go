package config

const (
	defaultLogLevel = "info"
)

type Main struct {
	LogLevel string `mapstructure:"log-level" json:"log-level"`
	LogFile  string `mapstructure:"log-file" json:"log-file"`
	Quiet    bool   `mapstructure:"quiet" json:"quiet"`
}

func NewMain() *Main {
	return &Main{
		LogLevel: defaultLogLevel,
	}
}
