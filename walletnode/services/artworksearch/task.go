package artworksearch

import (
	"context"
	"fmt"
	"sync"

	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch/thumbnail"
)

// Task is the task of searching for artwork.
type Task struct {
	task.Task
	*Service

	searchResult   []*RegTicketSearch
	searchResMutex sync.Mutex

	resultChan      chan *RegTicketSearch
	err             error
	request         *ArtSearchRequest
	thumbnailHelper thumbnail.Helper
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")
	defer close(task.resultChan)
	defer task.thumbnailHelper.Close()

	if err := task.run(ctx); err != nil {
		task.err = err
		task.UpdateStatus(StatusTaskFailure)
		log.WithContext(ctx).WithError(err).Warnf("Task failed")

		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)

	return nil
}

func (task *Task) run(ctx context.Context) error {
	actTickets, err := task.pastelClient.ActTickets(ctx, pastel.ActTicketAll, task.request.MinBlock)
	if err != nil {
		return fmt.Errorf("act ticket: %s", err)
	}

	group, gctx := errgroup.WithContext(ctx)

	for _, ticket := range actTickets {
		ticket := ticket

		if !inIntRange(ticket.Height, nil, task.request.MaxBlock) {
			continue
		}

		if task.request.Artist != nil && *task.request.Artist != ticket.ActTicketData.PastelID {
			continue
		}

		group.Go(func() error {
			regTicket, err := task.RegTicket(gctx, ticket.ActTicketData.RegTXID)
			if err != nil {
				log.WithContext(ctx).WithField("txid", ticket.TXID).WithError(err).Error("Reg Ticket")

				return fmt.Errorf("reg ticket - txid: %s - err: %s", ticket.TXID, err)
			}

			if srch, isMatched := task.filterRegTicket(regTicket); isMatched {
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

	if len(task.searchResult) > task.request.Limit {
		task.searchResult = task.searchResult[:task.request.Limit]
	}

	pastelConnections := 10
	if len(task.searchResult) < pastelConnections {
		pastelConnections = len(task.searchResult)
	}

	if err := task.thumbnailHelper.Connect(ctx, uint(pastelConnections)); err != nil {
		return fmt.Errorf("connect Thumbnail helper : %s", err)
	}

	group, gctx = errgroup.WithContext(ctx)
	for i, res := range task.searchResult {
		res := res
		res.MatchIndex = i

		group.Go(func() error {
			data, err := task.thumbnailHelper.Fetch(gctx, string(res.RegTicket.RegTicketData.ArtTicketData.AppTicketData.PreviewHash))
			if err != nil {
				log.WithContext(ctx).WithField("txid", res.TXID).WithError(err).Error("Fetch Thumbnail")

				return fmt.Errorf("fetch thumbnail: txid: %s - err: %s", res.TXID, err)
			}
			res.Thumbnail = data

			// Post on result channel
			task.resultChan <- res

			log.WithContext(ctx).WithField("search_result", res).Debug("Posted search result")

			return nil
		})
	}

	return group.Wait()
}

// filterRegTicket filters ticket against request params & checks if its a match
func (task *Task) filterRegTicket(regTicket *pastel.RegTicket) (srch *RegTicketSearch, matched bool) {
	if !inIntRange(regTicket.RegTicketData.ArtTicketData.AppTicketData.RarenessScore,
		task.request.MinRarenessScore, task.request.MaxRarenessScore) {
		return srch, false
	}

	if !inIntRange(regTicket.RegTicketData.ArtTicketData.AppTicketData.NSFWScore,
		task.request.MinNsfwScore, task.request.MaxNsfwScore) {
		return srch, false
	}

	if !inIntRange(regTicket.RegTicketData.ArtTicketData.AppTicketData.TotalCopies,
		task.request.MinCopies, task.request.MaxCopies) {
		return srch, false
	}

	regSearch := &RegTicketSearch{
		RegTicket: regTicket,
	}

	return regSearch.Search(task.request)
}

// addMatchedResult adds to search result
func (task *Task) addMatchedResult(res *RegTicketSearch) {
	task.searchResMutex.Lock()
	defer task.searchResMutex.Unlock()

	task.searchResult = append(task.searchResult, res)
}

// Error returns task err
func (task *Task) Error() error {
	return task.err
}

// SubscribeSearchResult returns a new search resultof the state.
func (task *Task) SubscribeSearchResult() <-chan *RegTicketSearch {

	return task.resultChan
}

// NewTask returns a new Task instance.
func NewTask(service *Service, request *ArtSearchRequest) *Task {
	return &Task{
		Task:            task.New(StatusTaskStarted),
		Service:         service,
		request:         request,
		resultChan:      make(chan *RegTicketSearch),
		thumbnailHelper: thumbnail.New(service.pastelClient, service.nodeClient, service.config.ConnectTimeout),
	}
}
