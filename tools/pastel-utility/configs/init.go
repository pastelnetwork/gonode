package configs

import "encoding/json"

// Init contains config of the Init command
type Init struct {
	WorkingDir string `json:"workdir,omitempty"`
	Network    string `json:"network,omitempty"`
	Force      bool   `json:"force,omitempty"`
	Peers      string `json:"peers"`
}

func (i *Init) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

// NewInit returns a new Init instance.
func NewInit() *Init {
	return &Init{
		Force: false,
	}
}
