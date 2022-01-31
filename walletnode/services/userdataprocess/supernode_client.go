package userdataprocess

import (
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// UserDataNodeMaker makes class ProcessUserdata for SuperNodeAPIInterface
type UserDataNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class ProcessUserdata for SuperNodeAPIInterface
func (maker UserDataNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &UserDataNode{ProcessUserdataInterface: conn.ProcessUserdata()}
}

// UserDataNode represent supernode connection.
type UserDataNode struct {
	node.ProcessUserdataInterface

	Result    *userdata.ProcessResult
	ResultGet *userdata.ProcessRequest
}
