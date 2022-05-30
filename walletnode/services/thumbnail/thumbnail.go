package thumbnail

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// thumbnailHandler helps with fetching thumbnails
type thumbnailHandler struct {
	meshHandler *common.MeshHandler

	fetchersChan chan request
	nodesDone    chan struct{}
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
func (h *thumbnailHandler) connect(ctx context.Context, num int, cancel context.CancelFunc) error {

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
	nftSearchNode, ok := someNode.SuperNodeAPIInterface.(*common.NftSearchingNode)
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

//Use numnails > 1 to fetch both thumbnails for all search result tickets in searchResult
func (h *thumbnailHandler) fetchAll(ctx context.Context, searchResult []*common.RegTicketSearch, resultChan *chan *common.RegTicketSearch) error {
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

			log.WithContext(ctx).WithField("search_result", res).Debug("Posted search result")

			return nil
		})
	}
	return group.Wait()
}

// fetch gets the actual thumbnail data from the network as bytes to be wrapped by the calling function
func (h *thumbnailHandler) fetch(ctx context.Context, txid string, numnails int) (data map[int][]byte, err error) {
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
