package ddscan

const (
	// DefaultDataFile that dd-service uses to store fingerprints
	// Should be ~/pastel_dupe_detection_service/dupe_detection_support_files/dupe_detection_image_fingerprint_database.sqlite
	DefaultDataFile = "~/pastel_dupe_detection_service/dupe_detection_support_files/dupe_detection_image_fingerprint_database.sqlite"
)

// Config contains settings of the dupe detection service
type Config struct {
	// DataName is sqlite database name used by dd-service
	DataFile string `mapstructure:"data_file" json:"data_file"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		DataFile: DefaultDataFile,
	}
}
