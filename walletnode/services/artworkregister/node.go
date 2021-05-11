package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Nodes represents multiple Node.
type Nodes []*Node

func (nodes *Nodes) add(node *Node) {
	*nodes = append(*nodes, node)
}

func (nodes *Nodes) activate() {
	for _, node := range *nodes {
		node.activated = true
	}
}

func (nodes *Nodes) disconnectInactive() {
	for _, node := range *nodes {
		if node.conn != nil && !node.activated {
			node.conn.Close()
		}
	}
}

func (nodes Nodes) findByPastelID(id string) *Node {
	for _, node := range nodes {
		if node.pastelID == id {
			return node
		}
	}
	return nil
}

// Node represent supernode connection.
type Node struct {
	node.RegisterArtowrk

	client node.Client
	conn   node.Connection

	activated bool

	address  string
	pastelID string
}

func (node *Node) connect(ctx context.Context, connID string, isPrimary bool) error {
	connCtx, connCancel := context.WithTimeout(ctx, connectToNodeTimeout)
	defer connCancel()

	if node.conn == nil {
		conn, err := node.client.Connect(connCtx, node.address)
		if err != nil {
			return err
		}
		node.conn = conn

		go func() {
			<-ctx.Done()
			conn.Close()
			node.conn = nil
		}()
	}

	stream, err := node.conn.RegisterArtowrk(ctx)
	if err != nil {
		return err
	}
	node.RegisterArtowrk = stream

	return stream.Handshake(ctx, connID, isPrimary)
}
