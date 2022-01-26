package nftsearch

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"sync"
)

type ThumbnailHandler struct {
	meshHandler *mixins.MeshHandler

	fetchersChan       chan request
	closeFetchersMutex sync.Mutex
	nodesDone          chan struct{}
}

// New returns a new instance of ThumbnailHandler as Helper
func NewThumbnailHandler(meshHandler *mixins.MeshHandler) *ThumbnailHandler {
	return &ThumbnailHandler{
		meshHandler:  meshHandler,
		fetchersChan: make(chan request),
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
func (h *ThumbnailHandler) Connect(ctx context.Context, num int, cancel context.CancelFunc) error {

	if err := h.meshHandler.ConnectToNSuperNodes(ctx, num); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	h.nodesDone = h.meshHandler.ConnectionsSupervisor(ctx, cancel)

	fetchersErrs, err := h.setFetchers(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not setup thumbnail fetcher")
		return errors.Errorf("setup thumbnail fetchers: %w (%w)", err, fetchersErrs)
	}
	return nil
}

// FetchMultiple fetches multiple thumbnails from results list
func (h *ThumbnailHandler) FetchMultiple(ctx context.Context, searchResult []*RegTicketSearch, resultChan *chan *RegTicketSearch) error {
	err := h.fetchAll(ctx, searchResult, resultChan)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return errors.Errorf("fetch thumbnails: %w", err)
	}
	return nil
}

// FetchOne fetches single thumbnails by custom request
func (h *ThumbnailHandler) FetchOne(ctx context.Context, key []byte) ([]byte, error) {
	data, err := h.fetch(ctx, key)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return nil, errors.Errorf("fetch thumbnails: %w", err)
	}
	return data, nil
}

// CloseAll disconnects from all SNs
func (h *ThumbnailHandler) CloseAll(ctx context.Context) error {
	return h.meshHandler.CloseSNsConnections(ctx, h.nodesDone)
}

// set one fetcher for each connected SN
func (h *ThumbnailHandler) setFetchers(ctx context.Context) ([]error, error) {
	group, _ := errgroup.WithContext(ctx)
	errChan := make(chan error, len(h.meshHandler.Nodes))

	for _, someNode := range h.meshHandler.Nodes {
		group.Go(func() error {
			return h.fetcher(ctx, someNode, someNode.PastelID())
		})
	}
	err := group.Wait()

	close(errChan)

	downloadErrors := []error{}
	for subErr := range errChan {
		downloadErrors = append(downloadErrors, subErr)
	}

	return downloadErrors, err
}

func (h *ThumbnailHandler) fetcher(ctx context.Context, someNode *common.SuperNodeClient, nodeId string) error {
	nftSearchNode, ok := someNode.SuperNodeAPIInterface.(*NftSearchNode)
	if !ok {
		//TODO: use assert here?
		return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case req, ok := <-h.fetchersChan:
			if !ok {
				return nil
			}

			log.WithContext(ctx).Debugf("thumb-key: %v-%v", req.key, nodeId)
			data, err := nftSearchNode.DownloadThumbnail(ctx, req.key)
			req.respCh <- &response{err: err, data: data}
		}
	}
}

func (h *ThumbnailHandler) fetchAll(ctx context.Context, searchResult []*RegTicketSearch, resultChan *chan *RegTicketSearch) error {
	group, _ := errgroup.WithContext(ctx)

	for i, res := range searchResult {
		res := res
		res.MatchIndex = i

		group.Go(func() error {
			tgroup, tgctx := errgroup.WithContext(ctx)
			var t1Data, t2Data []byte
			tgroup.Go(func() (err error) {
				t1Data, err = h.fetch(tgctx, res.RegTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash)
				return err
			})

			tgroup.Go(func() (err error) {
				t2Data, err = h.fetch(tgctx, res.RegTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail2Hash)
				return err
			})

			if err := tgroup.Wait(); err != nil {
				log.WithContext(ctx).WithField("txid", res.TXID).WithError(err).Error("fetch Thumbnail")
				return fmt.Errorf("fetch thumbnail: txid: %s - err: %s", res.TXID, err)
			}

			res.Thumbnail = t1Data
			res.ThumbnailSecondry = t2Data
			// Post on result channel
			*resultChan <- res

			log.WithContext(ctx).WithField("search_result", res).Debug("Posted search result")

			return nil
		})
	}
	return group.Wait()
}

// Fetch gets thumbnail
func (h *ThumbnailHandler) fetch(ctx context.Context, key []byte) (data []byte, err error) {
	respCh := make(chan *response)
	req := request{key: key, respCh: respCh}

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
