package download

import (
	"context"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/bridge/services/common"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

// ThumbnailHandler helps with fetching thumbnails
type thumbnailHandler struct {
	meshHandler *common.MeshHandler

	fetchersChan chan request
	nodesDone    chan struct{}
	connMtx      sync.RWMutex
}

// newThumbnailHandler returns a new instance of ThumbnailHandler as Helper
func newThumbnailHandler(meshHandler *common.MeshHandler) *thumbnailHandler {
	return &thumbnailHandler{
		meshHandler:  meshHandler,
		fetchersChan: make(chan request),
	}
}

type response struct {
	data map[int][]byte
	err  error
}

type request struct {
	txid     string
	numnails int
	respCh   chan *response
}

// Connect creates `connections` no. of connections with supernode & starts thumbnail request listeners
func (h *thumbnailHandler) Connect(ctx context.Context, num int, cancel context.CancelFunc) error {
	h.connMtx.Lock()
	defer h.connMtx.Unlock()

	if err := h.meshHandler.ConnectToNSuperNodes(ctx, num); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	h.nodesDone = h.meshHandler.ConnectionsSupervisor(ctx, cancel)

	if err := h.setFetchers(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not setup thumbnail fetcher")
		return errors.Errorf("setup thumbnail fetchers: %w", err)
	}
	return nil
}

// CloseAll disconnects from all SNs
func (h *thumbnailHandler) CloseAll(ctx context.Context) error {
	h.connMtx.Lock()
	defer h.connMtx.Unlock()

	return h.meshHandler.CloseSNsConnections(ctx, h.nodesDone)
}

// set one fetcher for each connected SN
func (h *thumbnailHandler) setFetchers(ctx context.Context) error {
	if len(h.meshHandler.Nodes) == 0 {
		return fmt.Errorf("no nodes to listen")
	}

	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range h.meshHandler.Nodes {

		someNode := someNode
		group.Go(func() error {
			return h.fetcher(ctx, someNode, someNode.PastelID())
		})
	}

	return group.Wait()
}

func (h *thumbnailHandler) fetcher(ctx context.Context, someNode *common.SuperNodeClient, nodeID string) error {
	nftSearchNode, ok := someNode.SuperNodeAPIInterface.(*DownloadNode)
	if !ok {
		return errors.Errorf("node %s is not DownloadNode", someNode.String())
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req, ok := <-h.fetchersChan:
				if !ok {
					return
				}

				log.WithContext(ctx).Debugf("thumb-txid: %v-%v", req.txid, nodeID)
				data, err := nftSearchNode.DownloadThumbnail(ctx, req.txid, req.numnails)
				req.respCh <- &response{err: err, data: data}
			}
		}
	}()

	return nil
}

// fetch gets the actual thumbnail data from the network as bytes to be wrapped by the calling function
func (h *thumbnailHandler) Fetch(ctx context.Context, txid string, numnails int) (data map[int][]byte, err error) {
	h.connMtx.RLock()
	defer h.connMtx.RUnlock()

	respCh := make(chan *response)
	req := request{txid: txid, respCh: respCh, numnails: numnails}

	go func() {
		h.fetchersChan <- req
	}()

	select {
	case <-ctx.Done():
		return data, nil
	case resp := <-respCh:
		return resp.data, resp.err
	}
}
