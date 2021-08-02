package dupedetection

import "path/filepath"

const (
	// DefaultInputDir that dd-service monitors for the new file to generate fingerprints
	DefaultInputDir = "input"
	// DefaultOutputDir that dd-service uses to put new fingerprints
	DefaultOutputDir = "output"
	// DefaultDataFile is path to the SQLite database used by dd-service
	DefaultDataFile = "dupe_detection_image_fingerprint_database.sqlite"
)

// Config contains settings of the dupe detection service
type Config struct {
	// input directory for monitoring the new file to generate fingerprints
	InputDir string `mapstructure:"input_dir" json:"input_dir,omitempty"`
	// output directory for using to put new fingerprints
	OutputDir string `mapstructure:"output_dir" json:"output_dir,omitempty"`
	// data directory for dupe detection
	DataFile string `mapstructure:"data_file" json:"data_file,omitempty"`
}

// SetWorkDir applies `workDir` to DataFile, OutputDir and InputDir if it was
// not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.InputDir) {
		config.InputDir = filepath.Join(workDir, config.InputDir)
	}

	if !filepath.IsAbs(config.OutputDir) {
		config.OutputDir = filepath.Join(workDir, config.OutputDir)
	}

	if !filepath.IsAbs(config.DataFile) {
		config.DataFile = filepath.Join(workDir, config.DataFile)
	}
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		InputDir:  DefaultInputDir,
		OutputDir: DefaultOutputDir,
		DataFile:  DefaultDataFile,
	}
}
