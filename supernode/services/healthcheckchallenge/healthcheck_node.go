package healthcheckchallenge

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// HCNodeMaker makes concrete instance of HealthCheckChallengeNode
type HCNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of HealthCheckChallengeNode
func (maker HCNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &HCNode{HealthCheckChallengeInterface: conn.HealthCheckChallenge()}
}

// HCNode (for HealthCheck Challenge Node) represent supernode connection.
type HCNode struct {
	node.HealthCheckChallengeInterface
}
