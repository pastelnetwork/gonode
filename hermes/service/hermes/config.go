package hermes

import (
	"path/filepath"
	"time"
)

const (
	// DefaultInputDir that dd-service monitors for the new file to generate fingerprint
	DefaultInputDir = "input_files"
	// DefaultOutputDir that dd-service uses to put new fingerprint
	DefaultOutputDir = "output_files"
	// DefaultDataFile that dd-service uses to store fingerprint
	DefaultDataFile = "registered_image_fingerprints_db.sqlite"
	// DefaultSupportDir that dd-service uses for support files
	// Should be ~/pastel_dupe_detection_service/support_files/registered_image_fingerprints_db.sqlite
	DefaultSupportDir         = "support_files"
	defaultActivationDeadline = 300

	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay = 400 * time.Millisecond
	defaultAcceptNodesTimeout     = 600 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout   = time.Second * 20
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

	// ActivationDeadline is number of blocks allowed to pass for a ticket to be activated
	ActivationDeadline uint32 `mapstructure:"activation_deadline" json:"activation_deadline,omitempty"`
	SNHost             string `mapstructure:"sn_host" json:"sn_host"`
	SNPort             int    `mapstructure:"sn_port" json:"sn_port"`

	NumberSuperNodes int

	ConnectToNodeTimeout      time.Duration
	ConnectToNextNodeDelay    time.Duration
	AcceptNodesTimeout        time.Duration
	CreatorPastelID           string
	CreatorPastelIDPassphrase string
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
		InputDir:           DefaultInputDir,
		OutputDir:          DefaultOutputDir,
		SupportDir:         DefaultSupportDir,
		DataFile:           DefaultDataFile,
		ActivationDeadline: defaultActivationDeadline,
		NumberSuperNodes:   defaultNumberSuperNodes,

		ConnectToNodeTimeout:   defaultConnectToNodeTimeout,
		ConnectToNextNodeDelay: defaultConnectToNextNodeDelay,
		AcceptNodesTimeout:     defaultAcceptNodesTimeout,
	}
}
