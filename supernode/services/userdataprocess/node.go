package userdataprocess

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/supernode/node"
)

const (
	defaultConnectToNodeTimeout = time.Second * 5
)

// Nodes represents muptiple Nodes
type Nodes []*Node

// Add adds a new node to the list
func (nodes *Nodes) Add(node *Node) {
	*nodes = append(*nodes, node)
}

// ByID returns a node from the list by the given id.
func (nodes Nodes) ByID(id string) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

// Remove removes a node from the list by the given id.
func (nodes *Nodes) Remove(id string) {
	for i, node := range *nodes {
		if node.ID == id {
			*nodes = append((*nodes)[:i], (*nodes)[:i+1]...)
			break
		}
	}
}

// Node represents a single supernode
type Node struct {
	node.ProcessUserdata
	client node.Client
	conn   node.Connection

	ID      string
	Address string
}

func (node *Node) connect(ctx context.Context) error {
	connCtx, connCancel := context.WithTimeout(ctx, defaultConnectToNodeTimeout)
	defer connCancel()

	conn, err := node.client.Connect(connCtx, node.Address)
	if err != nil {
		return err
	}

	node.conn = conn
	node.ProcessUserdata = conn.ProcessUserdata()
	return nil
}
