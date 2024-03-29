package nftsearch

import (
	"context"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// ThumbnailHandler helps with fetching thumbnails
type ThumbnailHandler struct {
	meshHandler *common.MeshHandler

	fetchersChan chan request
	nodesDone    chan struct{}
	connMtx      *sync.RWMutex
}

// NewThumbnailHandler returns a new instance of ThumbnailHandler as Helper
func NewThumbnailHandler(meshHandler *common.MeshHandler) *ThumbnailHandler {
	return &ThumbnailHandler{
		meshHandler:  meshHandler,
		fetchersChan: make(chan request),
		connMtx:      &sync.RWMutex{},
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
func (h *ThumbnailHandler) Connect(ctx context.Context, num int, cancel context.CancelFunc) error {
	h.connMtx.Lock()
	defer h.connMtx.Unlock()

	if err := h.meshHandler.ConnectToNSuperNodes(ctx, num, []string{}); err != nil {
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

// FetchMultiple fetches multiple thumbnails from results list, sending them to the resultChan
func (h *ThumbnailHandler) FetchMultiple(ctx context.Context, searchResult []*RegTicketSearch, resultChan *chan *RegTicketSearch) error {
	err := h.fetchAll(ctx, searchResult, resultChan)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return errors.Errorf("fetch thumbnails: %w", err)
	}
	return nil
}

// FetchOne fetches single thumbnails by custom request
//
//	The key is base58(thumbnail_hash)
func (h *ThumbnailHandler) FetchOne(ctx context.Context, txid string) ([]byte, error) {
	data, err := h.fetch(ctx, txid, 1)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return nil, errors.Errorf("fetch thumbnails: %w", err)
	}
	return data[0], nil
}

// CloseAll disconnects from all SNs
func (h *ThumbnailHandler) CloseAll(ctx context.Context) error {
	h.connMtx.Lock()
	defer h.connMtx.Unlock()

	return h.meshHandler.CloseSNsConnections(ctx, h.nodesDone)
}

// set one fetcher for each connected SN
func (h *ThumbnailHandler) setFetchers(ctx context.Context) error {
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

func (h *ThumbnailHandler) fetcher(ctx context.Context, someNode *common.SuperNodeClient, nodeID string) error {
	nftSearchNode, ok := someNode.SuperNodeAPIInterface.(*NftSearchingNode)
	if !ok {
		return errors.Errorf("node %s is not NftSearchingNode", someNode.String())
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

// Use numnails > 1 to fetch both thumbnails for all search result tickets in searchResult
func (h *ThumbnailHandler) fetchAll(ctx context.Context, searchResult []*RegTicketSearch, resultChan *chan *RegTicketSearch) error {
	group, _ := errgroup.WithContext(ctx)

	for i, res := range searchResult {
		res := res
		res.MatchIndex = i

		group.Go(func() error {
			tgroup, tgctx := errgroup.WithContext(ctx)
			var thumbData map[int][]byte
			tgroup.Go(func() (err error) {
				thumbData, err = h.fetch(tgctx, res.RegTicket.TXID, 2)
				return err
			})

			if err := tgroup.Wait(); err != nil {
				log.WithContext(ctx).WithField("txid", res.TXID).WithError(err).Error("fetch Thumbnail")
				return fmt.Errorf("fetch thumbnail: txid: %s - err: %s", res.TXID, err)
			}

			res.Thumbnail = thumbData[0]
			res.ThumbnailSecondry = thumbData[1]
			// Post on result channel
			*resultChan <- res

			return nil
		})
	}
	return group.Wait()
}

// fetch gets the actual thumbnail data from the network as bytes to be wrapped by the calling function
func (h *ThumbnailHandler) fetch(ctx context.Context, txid string, numnails int) (data map[int][]byte, err error) {
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
