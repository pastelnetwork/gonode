package artworksearch

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"sync"

	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch/thumbnail"
)

// NftSearchTask is the task of searching for artwork.
type NftSearchTask struct {
	*common.WalletNodeTask

	service *NftSearchService
	Request *ArtSearchRequest

	searchResult    []*RegTicketSearch
	searchResMutex  sync.Mutex
	resultChan      chan *RegTicketSearch
	err             error
	thumbnailHelper thumbnail.Helper
}

// Run starts the task
func (task *NftSearchTask) Run(ctx context.Context) error {
	defer close(task.resultChan)
	defer task.thumbnailHelper.Close()
	task.err = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
}

func (task *NftSearchTask) run(ctx context.Context) error {
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
		return nil
	}

	pastelConnections := 10
	if len(task.searchResult) < pastelConnections {
		pastelConnections = len(task.searchResult)
	}

	if err := task.thumbnailHelper.Connect(ctx, uint(pastelConnections), &alts.SecInfo{
		PastelID:   task.Request.UserPastelID,
		PassPhrase: task.Request.UserPassphrase,
		Algorithm:  "ed448",
	}); err != nil {
		return fmt.Errorf("connect Thumbnail helper : %s", err)
	}

	group, gctx = errgroup.WithContext(ctx)
	for i, res := range task.searchResult {
		res := res
		res.MatchIndex = i

		group.Go(func() error {
			tgroup, tgctx := errgroup.WithContext(ctx)
			var t1Data, t2Data []byte
			tgroup.Go(func() (err error) {
				t1Data, err = task.thumbnailHelper.Fetch(tgctx, res.RegTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail1Hash)
				return err
			})

			tgroup.Go(func() (err error) {
				t2Data, err = task.thumbnailHelper.Fetch(tgctx, res.RegTicket.RegTicketData.NFTTicketData.AppTicketData.Thumbnail2Hash)
				return err
			})

			if tgroup.Wait() != nil {
				log.WithContext(ctx).WithField("txid", res.TXID).WithError(err).Error("fetch Thumbnail")

				return fmt.Errorf("fetch thumbnail: txid: %s - err: %s", res.TXID, err)
			}

			res.Thumbnail = t1Data
			res.ThumbnailSecondry = t2Data
			// Post on result channel
			task.resultChan <- res

			log.WithContext(ctx).WithField("search_result", res).Debug("Posted search result")

			return nil
		})
	}

	return group.Wait()
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

// Error returns task err
func (task *NftSearchTask) Error() error {
	return task.err
}

// SubscribeSearchResult returns a new search resultof the state.
func (task *NftSearchTask) SubscribeSearchResult() <-chan *RegTicketSearch {

	return task.resultChan
}

func (task *NftSearchTask) removeArtifacts() {
}

// NewNftSearchTask returns a new NftSearchTask instance.
func NewNftSearchTask(service *NftSearchService, request *ArtSearchRequest) *NftSearchTask {
	return &NftSearchTask{
		WalletNodeTask:  common.NewWalletNodeTask(logPrefix),
		service:         service,
		Request:         request,
		resultChan:      make(chan *RegTicketSearch),
		thumbnailHelper: thumbnail.New(service.pastelHandler.PastelClient, service.nodeClient, service.config.ConnectToNodeTimeout),
	}
}
