package nftsearch

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"sync"

	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// NftSearchTask is the task of searching for nft.
type NftSearchTask struct {
	*common.WalletNodeTask

	thumbnail *ThumbnailHandler

	service *NftSearchService
	Request *NftSearchRequest

	searchResult   []*RegTicketSearch
	resultChan     chan *RegTicketSearch
	searchResMutex sync.Mutex

	err error
}

// Run starts the task
func (task *NftSearchTask) Run(ctx context.Context) error {
	defer close(task.resultChan)
	task.err = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
}

func (task *NftSearchTask) run(ctx context.Context) error {
	if err := task.search(ctx); err != nil {
		return errors.Errorf("search tickets: %w", err)
	}

	pastelConnections := task.service.config.NumberSuperNodes
	if len(task.searchResult) < pastelConnections {
		pastelConnections = len(task.searchResult)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := task.thumbnail.Connect(ctx, pastelConnections, cancel); err != nil {
		return errors.Errorf("connect and setup fetchers: %w", err)
	}
	if err := task.thumbnail.FetchMultiple(ctx, task.searchResult, &task.resultChan); err != nil {
		return errors.Errorf("fetch multiple thumbnails: %w", err)
	}

	return task.thumbnail.CloseAll(ctx)
}

func (task *NftSearchTask) search(ctx context.Context) error {
	actTickets, err := task.service.pastelHandler.PastelClient.ActTickets(ctx, pastel.ActTicketAll, task.Request.MinBlock)
	if err != nil {
		return fmt.Errorf("act ticket: %s", err)
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, ticket := range actTickets {
		ticket := ticket

		if !common.InIntRange(ticket.Height, nil, task.Request.MaxBlock) {
			continue
		}

		if task.Request.Artist != nil && *task.Request.Artist != ticket.ActTicketData.PastelID {
			continue
		}

		group.Go(func() error {
			regTicket, err := task.service.pastelHandler.PastelClient.RegTicket(gctx, ticket.ActTicketData.RegTXID)
			if err != nil {
				log.WithContext(ctx).WithField("txid", ticket.TXID).WithError(err).Error("Reg Request")

				return fmt.Errorf("reg ticket - txid: %s - err: %s", ticket.TXID, err)
			}

			if srch, isMatched := task.filterRegTicket(&regTicket); isMatched {
				task.addMatchedResult(srch)
			}

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return fmt.Errorf("reg ticket: %s", err)
	}

	sort.Slice(task.searchResult, func(i, j int) bool {
		return task.searchResult[i].MaxScore > task.searchResult[j].MaxScore
	})

	if len(task.searchResult) > task.Request.Limit {
		task.searchResult = task.searchResult[:task.Request.Limit]
	}

	if len(task.searchResult) == 0 {
		log.WithContext(ctx).WithField("request", task.Request).Debug("No matching results")
	}
	return nil
}

// filterRegTicket filters ticket against Request params & checks if its a match
func (task *NftSearchTask) filterRegTicket(regTicket *pastel.RegTicket) (srch *RegTicketSearch, matched bool) {
	// --------- WIP: PSL-142------------------
	/* if !inFloatRange(float64(regTicket.RegTicketData.NFTTicketData.AppTicketData.PastelRarenessScore),
		task.Request.MinRarenessScore, task.Request.MaxRarenessScore) {
		return srch, false
	}

		if !inFloatRange(float64(regTicket.RegTicketData.NFTTicketData.AppTicketData.OpenNSFWScore),
			task.Request.MinNsfwScore, task.Request.MaxNsfwScore) {
			return srch, false
		}

	if !inFloatRange(float64(regTicket.RegTicketData.NFTTicketData.AppTicketData.InternetRarenessScore),
		task.Request.MinInternetRarenessScore, task.Request.MaxInternetRarenessScore) {
		return srch, false
	}*/

	if !common.InIntRange(regTicket.RegTicketData.NFTTicketData.AppTicketData.TotalCopies,
		task.Request.MinCopies, task.Request.MaxCopies) {
		return srch, false
	}

	regSearch := &RegTicketSearch{
		RegTicket: regTicket,
	}

	return regSearch.Search(task.Request)
}

// addMatchedResult adds to search result
func (task *NftSearchTask) addMatchedResult(res *RegTicketSearch) {
	task.searchResMutex.Lock()
	defer task.searchResMutex.Unlock()

	task.searchResult = append(task.searchResult, res)
}

// SubscribeSearchResult returns a new search result of the state.
func (task *NftSearchTask) SubscribeSearchResult() <-chan *RegTicketSearch {
	return task.resultChan
}

// Error returns task err
func (task *NftSearchTask) Error() error {
	return task.err
}

func (task *NftSearchTask) removeArtifacts() {
}

// NewNftSearchTask returns a new NftSearchTask instance.
func NewNftSearchTask(service *NftSearchService, request *NftSearchRequest) *NftSearchTask {
	task := &NftSearchTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
		resultChan:     make(chan *RegTicketSearch),
	}

	meshHandler := mixins.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &NftSearchNodeMaker{},
		service.pastelHandler,
		request.UserPastelID, request.UserPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	task.thumbnail = NewThumbnailHandler(meshHandler)

	return task
}

type NftGetSearchTask struct {
	*common.WalletNodeTask
	thumbnail *ThumbnailHandler
}

// NewNftSearchTask returns a new NftSearchTask instance.
func NewNftGetSearchTask(service *NftSearchService, pastelID string, passphrase string) *NftGetSearchTask {
	task := &NftGetSearchTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
	}

	meshHandler := mixins.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &NftSearchNodeMaker{},
		service.pastelHandler,
		pastelID, passphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	task.thumbnail = NewThumbnailHandler(meshHandler)

	return task
}
