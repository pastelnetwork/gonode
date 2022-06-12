package configs

import (
	"github.com/pastelnetwork/gonode/bridge/services/download"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	DownloadData download.Config `mapstructure:",squash" json:"download_data,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		DownloadData: *download.NewConfig(),
	}
}
