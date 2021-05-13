package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/walletnode/node"
	"golang.org/x/sync/errgroup"
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

func (nodes *Nodes) uploadImage(ctx context.Context, filename string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			return node.UploadImage(ctx, filename)
		})
	}
	return group.Wait()
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

func (node *Node) openStream(ctx context.Context, connID string, isPrimary bool) error {
	if err := node.connect(ctx); err != nil {
		return err
	}

	stream, err := node.conn.RegisterArtowrk(ctx)
	if err != nil {
		return err
	}
	node.RegisterArtowrk = stream

	return stream.Handshake(ctx, connID, isPrimary)
}

func (node *Node) connect(ctx context.Context) error {
	if node.conn != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, connectToNodeTimeout)
	defer connCancel()

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
	return nil
}
