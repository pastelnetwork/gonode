package userdataprocess

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type UserDataNodeMaker struct {
	node.NodeMaker
}

func (maker UserDataNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &UserDataNode{ProcessUserdataInterface: conn.ProcessUserdata()}
}

// Node represent supernode connection.
type UserDataNode struct {
	node.ProcessUserdataInterface
}
