package nftsearch

import (
	"context"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// NftSearchingTask is the task of searching for nft.
type NftSearchingTask struct {
	*common.WalletNodeTask

	thumbnail *ThumbnailHandler

	service *NftSearchingService
	// Request is search request from API call
	Request *NftSearchingRequest

	searchResult   []*RegTicketSearch
	resultChan     chan *RegTicketSearch
	searchResMutex sync.Mutex

	err error
}

// Run starts the task
func (task *NftSearchingTask) Run(ctx context.Context) error {
	defer close(task.resultChan)
	task.err = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
}

func (task *NftSearchingTask) run(ctx context.Context) error {
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

func (task *NftSearchingTask) search(ctx context.Context) error {
	actTickets, err := task.service.pastelHandler.PastelClient.ActTickets(ctx, pastel.ActTicketAll, task.Request.MinBlock)
	if err != nil {
		return fmt.Errorf("act ticket: %s", err)
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, ticket := range actTickets {
		ticket := ticket
		//filter list of activation tickets by blocknum if provided
		if !common.InIntRange(ticket.Height, nil, task.Request.MaxBlock) {
			continue
		}
		//filter list of activation tickets by artist pastelid if artist is provided
		if task.Request.Artist != nil && *task.Request.Artist != ticket.ActTicketData.PastelID {
			continue
		}
		//iterate through filtered activation tickets
		group.Go(func() error {
			//request art registration tickets
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
func (task *NftSearchingTask) filterRegTicket(regTicket *pastel.RegTicket) (srch *RegTicketSearch, matched bool) {
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
func (task *NftSearchingTask) addMatchedResult(res *RegTicketSearch) {
	task.searchResMutex.Lock()
	defer task.searchResMutex.Unlock()

	task.searchResult = append(task.searchResult, res)
}

// SubscribeSearchResult returns a new search result of the state.
func (task *NftSearchingTask) SubscribeSearchResult() <-chan *RegTicketSearch {
	return task.resultChan
}

// Error returns task err
func (task *NftSearchingTask) Error() error {
	return task.err
}

func (task *NftSearchingTask) removeArtifacts() {
}

// NewNftSearchTask returns a new NftSearchingTask instance.
func NewNftSearchTask(service *NftSearchingService, request *NftSearchingRequest) *NftSearchingTask {
	task := &NftSearchingTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
		resultChan:     make(chan *RegTicketSearch),
	}

	meshHandler := common.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &NftSearchingNodeMaker{},
		service.pastelHandler,
		request.UserPastelID, request.UserPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	task.thumbnail = NewThumbnailHandler(meshHandler)

	return task
}

// NftGetSearchTask helper
type NftGetSearchTask struct {
	*common.WalletNodeTask
	thumbnail *ThumbnailHandler
}

// NewNftGetSearchTask returns a new NftSearchingTask instance.
func NewNftGetSearchTask(service *NftSearchingService, pastelID string, passphrase string) *NftGetSearchTask {
	task := &NftGetSearchTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
	}

	meshHandler := common.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &NftSearchingNodeMaker{},
		service.pastelHandler,
		pastelID, passphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	task.thumbnail = NewThumbnailHandler(meshHandler)

	return task
}
