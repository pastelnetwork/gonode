package thumbnail

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
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
	Connect(ctx context.Context, connections uint, secInfo *alts.SecInfo) error
	Fetch(ctx context.Context, key []byte) ([]byte, error)
	Close()
}

type thumbnailHelper struct {
	pastelClient pastel.Client
	nodeClient   nodeClient.ClientInterface

	reqCh      chan request
	timeOut    time.Duration
	isClosed   bool
	closeMutex sync.Mutex
	cancel     context.CancelFunc
}

// New returns a new instance of thumbnailHelper as Helper
func New(pastelClient pastel.Client, nodeClient nodeClient.ClientInterface, timeout time.Duration) Helper {
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

// Connect creates `connections` no. of connections with supernode & starts thumbnail request listeners
func (t *thumbnailHelper) Connect(ctx context.Context, connections uint, secInfo *alts.SecInfo) error {
	if connections > maxConnections {
		return errors.New(maxConnectionsErr)
	}

	nodes, err := t.pastelTopNodes(ctx)
	if err != nil {
		return fmt.Errorf("pastelTopNodes: %v", err)
	}

	if len(nodes) < int(connections) {
		return errors.New("not enough masternodes")
	}

	ctx, cancel := context.WithCancel(ctx)
	t.cancel = cancel
	connCnt := 0

	for _, node := range nodes {
		node := node
		if connCnt == int(connections) {
			break
		}

		if err := node.Connect(ctx, t.timeOut, secInfo); err != nil {
			log.WithContext(ctx).WithError(err).Error("Connect to master node failed")
			continue
		}

		connCnt = connCnt + 1
		go func() {
			t.listen(ctx, node)
			t.Close()
			if err := node.ConnectionInterface.Close(); err != nil {
				log.WithContext(ctx).WithError(err).Error("ThumbnailPuller.Start-Close node connection")
			}
		}()
	}

	if connCnt != int(connections) {
		t.Close()
		return errors.New("could not connect to enough nodes")
	}

	return nil
}

func (t *thumbnailHelper) listen(ctx context.Context, n *node.NftSearchNodeClient) {
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-t.reqCh:
			if !ok {
				return
			}

			log.WithContext(ctx).Debugf("thumb-key: %v-%v", req.key, n.PastelID())
			data, err := n.DownloadThumbnail(ctx, req.key)
			req.respCh <- &response{err: err, data: data}
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

func (t *thumbnailHelper) pastelTopNodes(ctx context.Context) (nodes node.List, err error) {
	mns, err := t.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
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
	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
	}
	close(t.reqCh)
}

// IsClosed returns whether helper is closed
func (t *thumbnailHelper) IsClosed() bool {
	t.closeMutex.Lock()
	defer t.closeMutex.Unlock()

	return t.isClosed
}
