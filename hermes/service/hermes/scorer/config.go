package scorer

import "time"

// Config contains settings of the dupe detection service
type Config struct {
	NumberSuperNodes int

	ConnectToNodeTimeout      time.Duration
	ConnectToNextNodeDelay    time.Duration
	AcceptNodesTimeout        time.Duration
	CreatorPastelID           string
	CreatorPastelIDPassphrase string
}
