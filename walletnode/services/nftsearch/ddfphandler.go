package nftsearch

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// DDFPHandler helps with fetching DD and FP Files
type DDFPHandler struct {
	meshHandler *common.MeshHandler

	fetchersChan chan ddfpRequest
	nodesDone    chan struct{}
}

// NewDDFPHandler returns a new instance of DDFPHandler as Helper
func NewDDFPHandler(meshHandler *common.MeshHandler) *DDFPHandler {
	return &DDFPHandler{
		meshHandler:  meshHandler,
		fetchersChan: make(chan ddfpRequest),
	}
}

type ddfpResponse struct {
	data []byte
	err  error
}

type ddfpRequest struct {
	txid   string
	respCh chan *ddfpResponse
}

// Connect creates `connections` no. of connections with supernode & starts dd and fp ddfpRequest listeners
func (h *DDFPHandler) Connect(ctx context.Context, num int, cancel context.CancelFunc) error {

	if err := h.meshHandler.ConnectToNSuperNodes(ctx, num); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	h.nodesDone = h.meshHandler.ConnectionsSupervisor(ctx, cancel)

	if err := h.setFetchers(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not setup DD and FP fetcher")
		return errors.Errorf("setup dd and fp fetchers: %w", err)
	}
	return nil
}

// CloseAll disconnects from all SNs
func (h *DDFPHandler) CloseAll(ctx context.Context) error {
	return h.meshHandler.CloseSNsConnections(ctx, h.nodesDone)
}

// set one fetcher for each connected SN
func (h *DDFPHandler) setFetchers(ctx context.Context) error {
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

func (h *DDFPHandler) fetcher(ctx context.Context, someNode *common.SuperNodeClient, nodeID string) error {
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

				log.WithContext(ctx).Debugf("dd-and-fp-txid: %v-%v", req.txid, nodeID)
				data, err := nftSearchNode.DownloadDDAndFingerprints(ctx, req.txid)
				req.respCh <- &ddfpResponse{err: err, data: data}
			}
		}
	}()

	return nil
}

// fetch gets the actual thumbnail data from the network as bytes to be wrapped by the calling function
func (h *DDFPHandler) Fetch(ctx context.Context, txid string) (data []byte, err error) {
	respCh := make(chan *ddfpResponse)
	req := ddfpRequest{txid: txid, respCh: respCh}

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
