package nftsearch

import (
	"context"
	"fmt"
	"sync"

	json "github.com/json-iterator/go"

	"sort"

	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// NftSearchingTask is the task of searching for nft.
type NftSearchingTask struct {
	*common.WalletNodeTask

	thumbnail *ThumbnailHandler
	ddAndFP   *DDFPHandler

	service *NftSearchingService
	// request is search request from API call
	request *NftSearchingRequest

	searchResult   []*RegTicketSearch
	resultChan     chan *RegTicketSearch
	searchResMutex sync.Mutex
}

// Run starts the task
func (task *NftSearchingTask) Run(ctx context.Context) error {
	defer close(task.resultChan)
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *NftSearchingTask) run(ctx context.Context) error {
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go task.connectThumbnailAndDDHelper(newCtx, cancel)

	if err := task.search(ctx); err != nil {
		return errors.Errorf("search tickets: %w", err)
	}

	if err := task.thumbnail.FetchMultiple(newCtx, task.searchResult, &task.resultChan); err != nil {
		return errors.Errorf("fetch multiple thumbnails: %w", err)
	}

	return task.thumbnail.CloseAll(newCtx)
}

func (task *NftSearchingTask) search(ctx context.Context) error {
	actTickets, err := task.service.pastelHandler.PastelClient.ActTickets(ctx, pastel.ActTicketAll, task.request.MinBlock)
	if err != nil {
		return fmt.Errorf("act ticket: %s", err)
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, ticket := range actTickets {
		ticket := ticket
		//filter list of activation tickets by blocknum if provided
		if !common.InIntRange(ticket.Height, nil, task.request.MaxBlock) {
			continue
		}
		//filter list of activation tickets by artist pastelid if artist is provided
		if task.request.Artist != nil && *task.request.Artist != ticket.ActTicketData.PastelID {
			continue
		}

		//iterate through filtered activation tickets

		//request art registration tickets
		regTicket, err := task.service.pastelHandler.PastelClient.RegTicket(ctx, ticket.ActTicketData.RegTXID)
		if err != nil {
			log.WithContext(ctx).WithField("txid", ticket.ActTicketData.RegTXID).WithError(err).Error("Reg request")
			continue
		}

		group.Go(func() error {
			nftData, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
			if err != nil {
				log.WithContext(gctx).WithError(err).Error("Failed to decode reg ticket")
				return nil
			}
			regTicket.RegTicketData.NFTTicketData = *nftData

			if srch, isMatched := task.filterRegTicket(ctx, &regTicket); isMatched {
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

	if len(task.searchResult) == 0 {
		log.WithContext(ctx).WithField("request", task.request).Debug("No matching results")
	}

	return nil
}

// filterRegTicket filters ticket against request params & checks if its a match
func (task *NftSearchingTask) filterRegTicket(ctx context.Context, regTicket *pastel.RegTicket) (srch *RegTicketSearch, matched bool) {
	var ddAndFpData []byte
	var err error

	ddAndFpData, err = task.ddAndFP.Fetch(ctx, regTicket.TXID)
	if err != nil {
		log.WithContext(ctx).WithField("request", task.request).WithField("txid", regTicket.TXID).Warn("Could not get dd and fp for this txid in search.")
		return srch, false
	}

	ddAndFpStruct := &pastel.DDAndFingerprints{}
	if err := json.Unmarshal(ddAndFpData, ddAndFpStruct); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to  unmarshalled dd and fp struct")
	}

	//rareness score
	if !(task.request.MinRarenessScore == float64(0) && task.request.MinRarenessScore == task.request.MaxRarenessScore) && !common.InFloatRange(float64(ddAndFpStruct.OverallRarenessScore),
		&task.request.MinRarenessScore, &task.request.MaxRarenessScore) {
		log.WithContext(ctx).WithField("task.request.minrarenessscore", task.request.MinRarenessScore).
			WithField("task.request.maxrarenessscore", task.request.MaxRarenessScore).
			WithField("overallaveragerarnessscore", ddAndFpStruct.OverallRarenessScore).Debug("rareness score outside of range")
		return srch, false
	}

	//opennsfw score
	if !(task.request.MinNsfwScore == float64(0) && task.request.MinNsfwScore == task.request.MaxNsfwScore) && !common.InFloatRange(float64(ddAndFpStruct.OpenNSFWScore),
		&task.request.MinNsfwScore, &task.request.MaxNsfwScore) {
		log.WithContext(ctx).WithField("task.request.minnsfwscore", task.request.MinNsfwScore).WithField("task.request.maxnsfwscore", task.request.MaxNsfwScore).WithField("opennsfscore", ddAndFpStruct.OpenNSFWScore).Println("nsfw score outside of range")
		return srch, false
	}

	//Is likely dupe
	if task.request.IsLikelyDupe != ddAndFpStruct.IsLikelyDupe {
		log.WithContext(ctx).WithField("task.request.islikelydupe", task.request.IsLikelyDupe).WithField("ddandfpstruct.islikelydupe", ddAndFpStruct.IsLikelyDupe).Println("IsLikelyDupe not match")
		return srch, false
	}

	if !common.InIntRange(regTicket.RegTicketData.NFTTicketData.AppTicketData.TotalCopies,
		task.request.MinCopies, task.request.MaxCopies) {
		return srch, false
	}

	regSearch := &RegTicketSearch{
		RegTicket:     regTicket,
		RarenessScore: ddAndFpStruct.OverallRarenessScore,
		OpenNSFWScore: ddAndFpStruct.OpenNSFWScore,
		IsLikelyDupe:  ddAndFpStruct.IsLikelyDupe,
	}
	//performs fuzzy matching on string portions of search
	return regSearch.Search(task.request)
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
	return task.WalletNodeTask.Error()
}

func (task *NftSearchingTask) removeArtifacts() {
}

// NewNftSearchTask returns a new NftSearchingTask instance.
func NewNftSearchTask(service *NftSearchingService, request *NftSearchingRequest) *NftSearchingTask {
	task := common.NewWalletNodeTask(logPrefix, service.historyDB)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &NftSearchingNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     service.config.AcceptNodesTimeout,
			MinSNs:                 service.config.NumberSuperNodes,
			PastelID:               request.UserPastelID,
			Passphrase:             request.UserPassphrase,
		},
	}

	return &NftSearchingTask{
		WalletNodeTask: task,
		service:        service,
		request:        request,
		resultChan:     make(chan *RegTicketSearch),
		thumbnail:      NewThumbnailHandler(common.NewMeshHandler(meshHandlerOpts)),
		ddAndFP:        NewDDFPHandler(common.NewMeshHandler(meshHandlerOpts)),
	}
}

// NftGetSearchTask helper
type NftGetSearchTask struct {
	*common.WalletNodeTask
	thumbnail *ThumbnailHandler
	ddAndFP   *DDFPHandler
}

// NewNftGetSearchTask returns a new NftSearchingTask instance.
func NewNftGetSearchTask(service *NftSearchingService, pastelID string, passphrase string) *NftGetSearchTask {
	task := common.NewWalletNodeTask(logPrefix, service.historyDB)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &NftSearchingNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay:        service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:          service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:            service.config.AcceptNodesTimeout,
			MinSNs:                        service.config.NumberSuperNodes,
			PastelID:                      pastelID,
			Passphrase:                    passphrase,
			RequireSNAgreementOnMNTopList: false,
		},
	}

	return &NftGetSearchTask{
		WalletNodeTask: task,
		thumbnail:      NewThumbnailHandler(common.NewMeshHandler(meshHandlerOpts)),
		ddAndFP:        NewDDFPHandler(common.NewMeshHandler(meshHandlerOpts)),
	}
}

func (task *NftSearchingTask) connectThumbnailAndDDHelper(ctx context.Context, cancel context.CancelFunc) {
	pastelConnections := task.service.config.NumberSuperNodes

	if err := task.ddAndFP.Connect(ctx, pastelConnections, cancel); err != nil {
		log.WithContext(ctx).WithError(err).Error("connect and setup fetchers")
	}

	if err := task.thumbnail.Connect(ctx, pastelConnections, cancel); err != nil {
		log.WithContext(ctx).WithError(err).Error("connect and setup fetchers")
	}
}
