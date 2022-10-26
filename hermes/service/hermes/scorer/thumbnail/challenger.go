package thumbnail

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pastelnetwork/gonode/common/random"

	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer/network"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	numberOfThumbnailsToChallenge = 10
)

// Run runs the challenger once
func (h *ThumbnailsChallenger) Run(ctx context.Context) error {
	if err := h.run(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("run thumbnail challenger failed")
	} else {
		log.WithContext(ctx).Info("thumbnail challenger run success")
	}

	return nil
}

func (h *ThumbnailsChallenger) getChallengeRegTickets(ctx context.Context) (pastel.RegTickets, error) {
	regTickets, err := h.pastelClient.RegTickets(ctx)
	if err != nil {
		return nil, errors.Errorf("get reg tickets: %w", err)
	}

	var ret pastel.RegTickets
	for i := 0; i < numberOfThumbnailsToChallenge; i++ {
		n := random.Number(int64(len(regTickets) - 1))
		ret = append(ret, regTickets[n])
	}

	return ret, nil
}

func (h *ThumbnailsChallenger) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := h.connect(ctx, -1, cancel); err != nil {
		return fmt.Errorf("error connecting to nodes: %w", err)
	}
	log.WithContext(ctx).Info("thumbnail challenger connected")

	regTickets, err := h.getChallengeRegTickets(ctx)
	if err != nil {
		return errors.Errorf("get registered ticket: %w", err)
	}

	for _, t := range regTickets {
		if err := h.fetchThmbnials(ctx, t.TXID); err != nil {
			log.WithContext(ctx).WithField("txid", t.TXID).WithError(err).Error("fetch thumbnails failed")
			continue
		}

		badNodes := h.matchResults()
		if len(badNodes) > 0 {
			log.WithContext(ctx).WithField("count", len(badNodes)).Info("bad nodes count")
		}
		for _, node := range badNodes {
			if err := h.pastelClient.IncrementPoseBanScore(ctx, strings.Split(node.txid, "-")[0], node.idx); err != nil {
				log.WithContext(ctx).WithField("txid", t.TXID).WithError(err).Error("failed to increment pose-ban score")
			}
		}
	}
	log.WithContext(ctx).Info("thumbnail challenge completed")

	return h.closeAll(ctx)
}

// ThumbnailsChallenger helps with thumbnail challenges
type ThumbnailsChallenger struct {
	meshHandler  *network.MeshHandler
	pastelClient pastel.Client

	results   []response
	nodesDone chan struct{}
	resMtx    *sync.Mutex
}

type response struct {
	th1      []byte
	th2      []byte
	pastelID string
	txid     string
	idx      int
	err      error
}

// NewThumbnailsChallenger returns a new instance of thumbnailsChallenger
func NewThumbnailsChallenger(meshHandler *network.MeshHandler, pastelClient pastel.Client) *ThumbnailsChallenger {
	return &ThumbnailsChallenger{
		meshHandler:  meshHandler,
		resMtx:       &sync.Mutex{},
		pastelClient: pastelClient,
	}
}

// Connect creates `connections` no. of connections with supernode & starts thumbnail request listeners
func (h *ThumbnailsChallenger) connect(ctx context.Context, num int, cancel context.CancelFunc) error {
	if err := h.meshHandler.ConnectToNSuperNodes(ctx, num); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	h.nodesDone = h.meshHandler.ConnectionsSupervisor(ctx, cancel)

	return nil
}

// CloseAll disconnects from all SNs
func (h *ThumbnailsChallenger) closeAll(ctx context.Context) error {
	return h.meshHandler.CloseSNsConnections(ctx, h.nodesDone)
}

func (h *ThumbnailsChallenger) addResult(res response) {
	h.resMtx.Lock()
	defer h.resMtx.Unlock()

	h.results = append(h.results, res)
}

func (h *ThumbnailsChallenger) matchResults() (badNodes []response) {
	h.resMtx.Lock()
	defer h.resMtx.Unlock()

	maj := len(h.results) / 2
	idx := -1
	for a, resp := range h.results {
		matches := 0

		for b, resp2 := range h.results {
			if a == b {
				continue
			}

			if bytes.Equal(resp.th1, resp2.th1) {
				matches++
			}
		}

		if matches > maj {
			idx = a
			break
		}
	}

	if idx == -1 {
		return
	}

	for _, resp := range h.results {
		if !bytes.Equal(resp.th1, h.results[idx].th1) {
			badNodes = append(badNodes, resp)
		} else if !bytes.Equal(resp.th2, h.results[idx].th2) {
			badNodes = append(badNodes, resp)
		}
	}

	return badNodes
}

func (h *ThumbnailsChallenger) fetchThmbnials(ctx context.Context, txid string) error {
	if len(h.meshHandler.Nodes) == 0 {
		return fmt.Errorf("no nodes to listen")
	}

	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range h.meshHandler.Nodes {
		nftSearchNode, ok := someNode.SuperNodeAPIInterface.(*network.NftSearchingNode)
		if !ok {
			return errors.Errorf("node %s is not NftSearchingNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			data, err := nftSearchNode.DownloadThumbnail(ctx, txid, 2)
			resp := response{
				err:      err,
				pastelID: someNode.PastelID(),
				txid:     someNode.TxID(),
				idx:      someNode.Idx(),
			}

			if err != nil {
				resp.th1 = data[0]
				resp.th2 = data[1]
			}

			h.addResult(resp)

			return nil
		})
	}

	return group.Wait()
}
