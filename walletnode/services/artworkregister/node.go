package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/artwork"
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
		if node.PastelID == id {
			return node
		}
	}
	return nil
}

func (nodes *Nodes) sendImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			return node.UploadImage(ctx, file)
		})
	}
	return group.Wait()
}

// Node represent supernode connection.
type Node struct {
	node.RegisterArtwork
	client node.Client
	conn   node.Connection

	activated bool

	Address  string
	PastelID string
}

func (node *Node) connect(ctx context.Context) error {
	if node.conn != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, connectToNodeTimeout)
	defer connCancel()

	conn, err := node.client.Connect(connCtx, node.Address)
	if err != nil {
		return err
	}
	node.conn = conn
	node.RegisterArtwork = conn.RegisterArtwork()
	return nil
}
