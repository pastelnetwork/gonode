package hermes

const (
	defaultActivationDeadline = 300
)

// Config contains settings of the hermes server
type Config struct {
	// ActivationDeadline is number of blocks allowed to pass for a ticket to be activated
	ActivationDeadline uint32 `mapstructure:"activation_deadline" json:"activation_deadline,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ActivationDeadline: defaultActivationDeadline,
	}
}
