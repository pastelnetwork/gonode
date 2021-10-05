package ddscan

import "path/filepath"

const (
	// DefaultInputDir that dd-service monitors for the new file to generate fingerprints
	DefaultInputDir = "input_files"
	// DefaultOutputDir that dd-service uses to put new fingerprints
	DefaultOutputDir = "output_files"
	// DefaultDataFile that dd-service uses to store fingerprints
	DefaultDataFile = "dupe_detection_image_fingerprint_database.sqlite"
	// DefaultSupportDir that dd-service uses for support files
	// Should be ~/pastel_dupe_detection_service/support_files/dupe_detection_image_fingerprint_database.sqlite
	DefaultSupportDir = "support_files"
)

// Config contains settings of the dupe detection service
type Config struct {
	// InputDir is directory to contain input image for dupe detection service processing
	InputDir string `mapstructure:"input_dir" json:"input_dir,omitempty"`

	// OutputDir is directory to contain output of dupe detection service
	OutputDir string `mapstructure:"output_dir" json:"output_dir,omitempty"`

	// SupportDir is directory to contain support files of dupe detection service
	SupportDir string `mapstructure:"support_dir" json:"support_dir,omitempty"`

	// DataName is sqlite database name used by dd-service
	DataFile string `mapstructure:"data_file" json:"data_file"`
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

	if !filepath.IsAbs(config.DataFile) {
		config.DataFile = filepath.Join(workDir, config.SupportDir, config.DataFile)
	}
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		InputDir:   DefaultInputDir,
		OutputDir:  DefaultOutputDir,
		SupportDir: DefaultSupportDir,
		DataFile:   DefaultDataFile,
	}
}
