package thumbnail

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	nodeClient "github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch/node"
)

const (
	maxConnections    = 10
	maxConnectionsErr = "cannot connect to more than 10 nodes"
)

// Helper is interface contract for Thumbnail Getter
type Helper interface {
	Connect(ctx context.Context, connections uint) error
	Fetch(ctx context.Context, key []byte) ([]byte, error)
	Close()
}

type thumbnailHelper struct {
	pastelClient pastel.Client
	nodeClient   nodeClient.Client

	reqCh      chan request
	timeOut    time.Duration
	isClosed   bool
	closeMutex sync.Mutex
}

// New returns a new instance of thumbnailHelper as Helper
func New(pastelClient pastel.Client, nodeClient nodeClient.Client, timeout time.Duration) Helper {
	return &thumbnailHelper{
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		timeOut:      timeout,
		reqCh:        make(chan request),
	}
}

type response struct {
	data []byte
	err  error
}

type request struct {
	key    []byte
	respCh chan *response
}

// Connect creates `conncetions` no. of connections with supernode & starts thumbnail request listeners
func (t *thumbnailHelper) Connect(ctx context.Context, connections uint) error {
	if connections > maxConnections {
		return errors.New(maxConnectionsErr)
	}

	nodes, err := t.pastelTopNodes(ctx, connections)
	if err != nil {
		return fmt.Errorf("pastelTopNodes: %v", err)
	}

	for _, node := range nodes {
		node := node
		if err := node.Connect(ctx, t.timeOut); err != nil {
			t.Close()

			return err
		}

		go func() {
			t.listen(ctx, node)
			t.Close()
			if err := node.Connection.Close(); err != nil {
				log.WithContext(ctx).WithError(err).Error("ThumbnailPuller.Start-Close node connection")
			}
		}()
	}

	return nil
}

func (t *thumbnailHelper) listen(ctx context.Context, n *node.Node) {
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-t.reqCh:
			if !ok {
				return
			}

			fmt.Println("thumb-key:", req.key, n.PastelID())
			data, err := n.DownloadThumbnail(ctx, req.key)
			var resp *response
			if err != nil {
				resp = &response{err: err, data: nil}
			} else {
				resp = &response{err: nil, data: data}
			}
			req.respCh <- resp
		}
	}
}

// Fetch gets thumbnail
func (t *thumbnailHelper) Fetch(ctx context.Context, key []byte) (data []byte, err error) {
	// FIXME : verify hash of returned data
	if t.IsClosed() {
		return data, errors.New("connection is closed")
	}

	respCh := make(chan *response)
	req := request{key: key, respCh: respCh}

	go func() {
		t.reqCh <- req
	}()

	select {
	case <-ctx.Done():
		return data, nil
	case resp := <-respCh:
		return resp.data, resp.err
	}
}

func (t *thumbnailHelper) pastelTopNodes(ctx context.Context, connections uint) (nodes node.List, err error) {
	mns, err := t.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for i, mn := range mns {
		if uint(i) == connections {
			break
		}

		nodes = append(nodes, node.NewNode(t.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// Close closes request channel. Once Close is called, the helper can no longer fetch thumbnails
func (t *thumbnailHelper) Close() {
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()

	if t.isClosed {
		return
	}

	t.isClosed = true
	close(t.reqCh)
}

// IsClosed returns whether helper is closed
func (t *thumbnailHelper) IsClosed() bool {
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()

	return t.isClosed
}
