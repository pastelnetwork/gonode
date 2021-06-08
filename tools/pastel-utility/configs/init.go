package configs

// Init contains config of the Init command
type Init struct {
	WorkingDir string `json:"workdir,omitempty"`
	Network    string `json:"network,omitempty"`
	Force      bool   `json:"force,omitempty"`
	Peers      string `json:"peers"`
}

// NewInit returns a new Init instance.
func NewInit() *Init {
	return &Init{
		Force: false,
	}
}
