package api

const (
	defaultHostname = "localhost"
	defaultPort     = 8080
)

// Config contains settings of the Pastel client.
type Config struct {
	Hostname        string `mapstructure:"hostname" json:"hostname,omitempty"`
	Port            int    `mapstructure:"port" json:"port,omitempty"`
	Swagger         bool   `mapstructure:"swagger" json:"swagger,omitempty"`
	StaticFilesDir  string `mapstructure:"static_files_dir" json:"static_files_dir,omitempty"`
	CascadeFilesDir string `mapstructure:"cascade_files_dir" json:"cascade_files_dir,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
