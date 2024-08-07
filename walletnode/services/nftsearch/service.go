package nftsearch

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/queries"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "walletnode-nft-search"
)

// NftSearchingService represents a service for the NFT search.
type NftSearchingService struct {
	*task.Worker

	config        *Config
	nodeClient    node.ClientInterface
	pastelHandler *mixins.PastelHandler
	historyDB     queries.LocalStoreInterface
}

// Run starts worker.
func (service *NftSearchingService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *NftSearchingService) Tasks() []*NftSearchingTask {
	var tasks []*NftSearchingTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftSearchingTask))
	}
	return tasks
}

// GetTask returns the task of the NFT search by the given id.
func (service *NftSearchingService) GetTask(id string) *NftSearchingTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftSearchingTask)
	}
	return nil
}

// AddTask runs a new task of the NFT search and returns its taskID.
func (service *NftSearchingService) AddTask(p *nft.NftSearchPayload) string {

	request := FromNftSearchRequest(p)
	task := NewNftSearchTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// GetThumbnail gets thumbnail
func (service *NftSearchingService) GetThumbnail(ctx context.Context, regTicket *pastel.RegTicket, pastelID string, passphrase string) (data []byte, err error) {
	nftGetSearchTask := NewNftGetSearchTask(service, pastelID, passphrase)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := nftGetSearchTask.thumbnail.Connect(ctx, 1, cancel); err != nil {
		return nil, errors.Errorf("connect and setup fetchers: %w", err)
	}
	data, err = nftGetSearchTask.thumbnail.FetchOne(ctx, regTicket.TXID)
	if err != nil {
		return nil, errors.Errorf("nftsearch get thumbnail fetchone error, there may be multiple thumbnails: %w", err)
	}

	return data, nftGetSearchTask.thumbnail.CloseAll(ctx)

}

// GetDDAndFP gets dupe detection and fingerprint file
func (service *NftSearchingService) GetDDAndFP(ctx context.Context, regTicket *pastel.RegTicket, pastelID string, passphrase string) (data []byte, err error) {
	nftGetSearchTask := NewNftGetSearchTask(service, pastelID, passphrase)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := nftGetSearchTask.ddAndFP.Connect(ctx, 1, cancel); err != nil {
		return nil, errors.Errorf("connect and setup fetchers: %w", err)
	}

	defer func() {
		if err := nftGetSearchTask.ddAndFP.CloseAll(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Error("close dd&fp fetchers")
		}
	}()

	data, err = nftGetSearchTask.ddAndFP.Fetch(ctx, regTicket.TXID)
	if err != nil {
		return nil, errors.Errorf("nftsearch get dd and fp fetchone error: %w", err)
	}

	return data, nil
}

// RegTicket pull NFT registration ticket from cNode & decodes base64 encoded fields
func (service *NftSearchingService) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := service.pastelHandler.PastelClient.RegTicket(ctx, RegTXID)
	if err != nil {
		return nil, errors.Errorf("fetch: %w", err)
	}

	nftTicketData, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		return nil, errors.Errorf("convert NFT ticket: %w", err)
	}

	regTicket.RegTicketData.NFTTicketData = *nftTicketData

	return &regTicket, nil
}

// NewNftSearchService returns a new NFT Search Service instance.
//
//	NB: Because NewNftApiHandler calls AddTask, an NftSearchTask will actually
//		be instantiated instead of a generic Task.
func NewNftSearchService(config *Config, pastelClient pastel.Client,
	nodeClient node.ClientInterface, historyDB queries.LocalStoreInterface) *NftSearchingService {

	return &NftSearchingService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		historyDB:     historyDB,
	}
}
