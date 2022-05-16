package hermes

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "hermes"
)

// Service represents hermes service
type Service struct {
	pastelClient pastel.Client
	p2pClient    p2p.Client
	config       *Config
	currentBlock uint64
}

// Run starts Hermes service
func (service *Service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.run(ctx)
	})

	return group.Wait()
}

func (service *Service) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.WithContext(ctx).Error("exit hermes service because context done")
			return nil
		case <-time.After(30 * time.Second):
			if err := service.cleanupInactiveTickets(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("cleanupInactiveTickets failure, retrying...")
			}
		}
	}
}

func (service *Service) cleanupInactiveTickets(ctx context.Context) (err error) {
	err = service.cleanupRegTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("cleanupRegTickets failure")
	}

	err = service.cleanupActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("cleanupActionTickets failure")
	}

	return err
}

func (service *Service) cleanupActionTickets(ctx context.Context) error {
	actionRegTickets, err := service.pastelClient.ActionTicketsFromBlockHeight(ctx, pastel.TicketTypeInactive, service.currentBlock)
	if err != nil {
		return errors.Errorf("get ActionTickets: %w", err)
	}

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(actionRegTickets); i++ {
		if uint64(actionRegTickets[i].ActionTicketData.CalledAt) <= service.currentBlock {
			continue
		}

		decTicket, err := pastel.DecodeActionTicket(actionRegTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionRegTickets[i].ActionTicketData.ActionTicketData = *decTicket

		actionTicket := actionRegTickets[i].ActionTicketData.ActionTicketData
		service.cleanupActionTicketData(ctx, actionTicket)
	}

	return nil
}

func (service *Service) removeRegTicketData(ctx context.Context, ticketData pastel.AppTicket) error {
	// Remove RQ File IDs
	for _, rqFileID := range ticketData.RQIDs {
		if err := service.p2pClient.Delete(ctx, rqFileID); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
		}
	}

	// Remove Thumbnails
	if err := service.p2pClient.Delete(ctx, string(ticketData.Thumbnail1Hash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail1 hash")
	}
	if err := service.p2pClient.Delete(ctx, string(ticketData.Thumbnail2Hash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting thumbnail2 hash")
	}

	if err := service.p2pClient.Delete(ctx, string(ticketData.PreviewHash)); err != nil {
		log.WithContext(ctx).WithError(err).Error("error deleting preview thumbnail hash")
	}

	// Remove ddAndFpFileIDs
	for _, ddFpFileID := range ticketData.DDAndFingerprintsIDs {
		if err := service.p2pClient.Delete(ctx, ddFpFileID); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
		}
	}

	return nil
}

func (service *Service) cleanupRegTickets(ctx context.Context) error {
	nftRegTickets, err := service.pastelClient.RegTicketsFromBlockHeight(ctx, pastel.TicketTypeInactive, service.currentBlock)
	if err != nil {
		return errors.Errorf("get RegTickets: %w", err)
	}

	//loop through nft tickets and store newly found nft reg tickets
	for i := 0; i < len(nftRegTickets); i++ {
		if uint64(nftRegTickets[i].Height) <= service.currentBlock {
			continue
		}

		decTicket, err := pastel.DecodeNFTTicket(nftRegTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		nftRegTickets[i].RegTicketData.NFTTicketData = *decTicket
		ticketData := nftRegTickets[i].RegTicketData.NFTTicketData.AppTicketData

		service.removeRegTicketData(ctx, ticketData)
	}

	return nil
}

func (service *Service) cleanupActionTicketData(ctx context.Context, actionTicket pastel.ActionTicket) error {
	switch actionTicket.ActionType {
	case pastel.ActionTypeCascade:
		cascadeTicket, err := actionTicket.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		// Remove RQ File IDs
		for _, rqFileID := range cascadeTicket.RQIDs {
			if err := service.p2pClient.Delete(ctx, rqFileID); err != nil {
				log.WithContext(ctx).WithError(err).Error("error deleting rqFileID")
			}
		}

		if err := service.p2pClient.Delete(ctx, string(cascadeTicket.DataHash)); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
		}

	case pastel.ActionTypeSense:
		senseTicket, err := actionTicket.APISenseTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket).
				Warnf("Could not get sense ticket for action ticket data")
		}

		for _, ddFpFileID := range senseTicket.DDAndFingerprintsIDs {
			if err := service.p2pClient.Delete(ctx, ddFpFileID); err != nil {
				log.WithContext(ctx).WithError(err).Error("error deleting ddFpFileID")
			}
		}

		if err := service.p2pClient.Delete(ctx, string(senseTicket.DataHash)); err != nil {
			log.WithContext(ctx).WithError(err).Error("error deleting thumbnail hash")
		}
	}

	return nil
}

// NewService returns a new Hermes Service instance.
func NewService(config *Config, pastelClient pastel.Client, p2pClient p2p.Client) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		p2pClient:    p2pClient,
		currentBlock: 1,
	}
}
