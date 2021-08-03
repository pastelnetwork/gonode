package dupedetection

import "path/filepath"

const (
	// DefaultInputDir that dd-service monitors for the new file to generate fingerprints
	DefaultInputDir = "dupe_detection_input_files"
	// DefaultOutputDir that dd-service uses to put new fingerprints
	DefaultOutputDir = "dupe_detection_output_files"
	// DefaultOutputDir that dd-service uses to put new fingerprints(5 minutes)
	DefaultTimeout = 300
)

// Config contains settings of the dupe detection service
type Config struct {
	// InputDir is directory to contain input image for dupe detection service processing
	InputDir string `mapstructure:"input_dir" json:"input_dir,omitempty"`

	// OutputDir is directory to contain output of dupe detection service
	OutputDir string `mapstructure:"output_dir" json:"output_dir,omitempty"`

	// WaitForOutputTimeout is timeout when waiting for output from dupe detection service in second
	WaitForOutputTimeout int64 `mapstructure:"timeout" json:"timeout,omitempty"`
}

// SetWorkDir applies `workDir` to OutputDir and InputDir if they were
// not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.InputDir) {
		config.InputDir = filepath.Join(workDir, config.InputDir)
	}

	if !filepath.IsAbs(config.OutputDir) {
		config.OutputDir = filepath.Join(workDir, config.OutputDir)
	}
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		InputDir:             DefaultInputDir,
		OutputDir:            DefaultOutputDir,
		WaitForOutputTimeout: DefaultTimeout,
	}
}
