package thumbnail

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

type ThumbnailService struct {
	h                   *thumbnailHandler
	pastelConnections   int
	connRefreshInterval time.Duration
	connMtx             sync.Mutex
	err                 error
}

// FetchMultiple fetches multiple thumbnails from results list, sending them to the resultChan
func (s *ThumbnailService) FetchMultiple(ctx context.Context, searchResult []*common.RegTicketSearch, resultChan *chan *common.RegTicketSearch) error {
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	err := s.h.fetchAll(ctx, searchResult, resultChan)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return errors.Errorf("fetch thumbnails: %w", err)
	}
	return nil
}

// FetchOne fetches single thumbnails by custom request
//  The key is base58(thumbnail_hash)
func (s *ThumbnailService) FetchOne(ctx context.Context, txid string) ([]byte, error) {
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	data, err := s.h.fetch(ctx, txid, 1)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not fetch thumbnails")
		return nil, errors.Errorf("fetch thumbnails: %w", err)
	}

	return data[0], nil
}

// NewService returns a new instance of ThumbnailService as Helper
func NewService(meshHandler *common.MeshHandler, connections int,
	connRefreshInterval time.Duration) *ThumbnailService {

	return &ThumbnailService{
		h:                   newThumbnailHandler(meshHandler),
		pastelConnections:   connections,
		connRefreshInterval: connRefreshInterval,
	}
}

// Run starts the service
func (s *ThumbnailService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return s.run(ctx)
	})

	if err := group.Wait(); err != nil {
		s.err = err
		return err
	}

	return nil
}

func (s *ThumbnailService) connect(ctx context.Context) error {
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	newCtx, cancel := context.WithCancel(ctx)
	if err := s.h.connect(newCtx, s.pastelConnections, cancel); err != nil {
		return errors.Errorf("connect and setup fetchers: %w", err)
	}

	return nil
}

// CloseAll disconnects from all SNs
func (s *ThumbnailService) CloseAll(ctx context.Context) error {
	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	return s.h.meshHandler.CloseSNsConnections(ctx, s.h.nodesDone)
}

func (s *ThumbnailService) run(ctx context.Context) error {
	fmt.Println("connecting thumbnail")
	if err := s.connect(ctx); err != nil {
		return err
	}
	fmt.Println("connected thumbnail")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.connRefreshInterval):
			if err := s.CloseAll(ctx); err != nil {
				return errors.Errorf("error while closing connections: %w", err)
			}

			if err := s.connect(ctx); err != nil {
				return err
			}
		}
	}
}
