package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/nft/server"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
	"github.com/pastelnetwork/gonode/walletnode/services/nftregister"
	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

const (
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// NftAPIHandler represents services for nfts endpoints.
type NftAPIHandler struct {
	*Common
	register *nftregister.NftRegistrationService
	search   *nftsearch.NftSearchingService
	download *download.NftDownloadingService
	db       storage.KeyValue
	imageTTL time.Duration
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *NftAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// Mount configures the mux to serve the nfts endpoints.
func (service *NftAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := nft.NewEndpoints(service)
	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		NftRegUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// --- Register NFT ---

// UploadImage uploads an image and return unique image id.
func (service *NftAPIHandler) UploadImage(ctx context.Context, p *nft.UploadImagePayload) (res *nft.ImageRes, err error) {
	if p.Filename == nil {
		log.Error("filename not specified")
		return nil, nft.MakeBadRequest(errors.New("file not specified"))
	}

	id, expiry, err := service.register.StoreFile(ctx, p.Filename)
	if err != nil {
		log.WithError(err).Error("error storing file")
		return nil, nft.MakeInternalServerError(err)
	}

	log.Infof("file has been uploaded: %s", id)

	fee, err := service.register.CalculateFee(ctx, id)
	if err != nil {
		log.WithError(err).Error("error calculating fee")
		return nil, cascade.MakeInternalServerError(err)
	}

	log.Infof("estimated fee has been calculated: %f", fee)

	res = &nft.ImageRes{
		ImageID:      id,
		ExpiresIn:    expiry,
		EstimatedFee: fee + (fee * 0.1),
	}

	return res, nil
}

// Register runs registers process for the new NFT.
func (service *NftAPIHandler) Register(_ context.Context, p *nft.RegisterPayload) (res *nft.RegisterResult, err error) {
	taskID, err := service.register.AddTask(p)
	if err != nil {
		log.WithError(err).Error("unable to add task")
		return nil, sense.MakeInternalServerError(err)
	}

	res = &nft.RegisterResult{
		TaskID: taskID,
	}

	log.Infof("task has been added: %s", taskID)

	return res, nil
}

// GetTaskHistory - Gets a task's history
func (service *NftAPIHandler) GetTaskHistory(ctx context.Context, p *nft.GetTaskHistoryPayload) (res []*nft.TaskHistory, err error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, nft.MakeInternalServerError(errors.New("error retrieving status"))
	}
	defer store.CloseHistoryDB(ctx)

	statuses, err := store.QueryTaskHistory(p.TaskID)
	if err != nil {
		log.WithError(err).Error("error retrieving task history")
		return nil, nft.MakeNotFound(errors.New("task not found"))
	}

	history := []*nft.TaskHistory{}
	for _, entry := range statuses {
		timestamp := entry.CreatedAt.String()
		historyItem := &nft.TaskHistory{
			Timestamp: &timestamp,
			Status:    entry.Status,
		}

		if entry.Details != nil {
			historyItem.Details = &nft.Details{
				Message: &entry.Details.Message,
				Fields:  entry.Details.Fields,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}

// RegisterTaskState streams the state of the registration process.
func (service *NftAPIHandler) RegisterTaskState(ctx context.Context, p *nft.RegisterTaskStatePayload, stream nft.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		log.Error("unable to get task")
		return nft.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &nft.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return nft.MakeInternalServerError(err)
			}

			if status.IsFailure() {
				if task.Error() != nil {
					errStr := task.Error()
					log.WithContext(ctx).WithError(errStr).Error("error registering NFT")
				}
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// RegisterTask returns a single task.
func (service *NftAPIHandler) RegisterTask(_ context.Context, p *nft.RegisterTaskPayload) (res *nft.Task, err error) {
	task := service.register.GetTask(p.TaskID)
	if task == nil {
		log.Error("error retrieving task")
		return nil, nft.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	res = &nft.Task{
		ID:     p.TaskID,
		Status: task.Status().String(),
		Ticket: nftregister.ToNftRegisterTicket(task.Request),
		States: toNftStates(task.StatusHistory()),
	}
	return res, nil
}

// RegisterTasks returns list of all tasks.
func (service *NftAPIHandler) RegisterTasks(_ context.Context) (res nft.TaskCollection, err error) {
	tasks := service.register.Tasks()
	for _, task := range tasks {
		res = append(res, &nft.Task{
			ID:     task.ID(),
			Status: task.Status().String(),
			Ticket: nftregister.ToNftRegisterTicket(task.Request),
		})
	}
	return res, nil
}

// Download registered NFT
func (service *NftAPIHandler) Download(ctx context.Context, p *nft.DownloadPayload) (res *nft.FileDownloadResult, err error) {
	log.WithContext(ctx).WithField("txid", p.Txid).Info("Start downloading")
	defer log.WithContext(ctx).WithField("txid", p.Txid).Info("Finished downloading")
	taskID := service.download.AddTask(p, "")
	task := service.download.GetTask(taskID)
	defer task.Cancel()

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil, nft.MakeBadRequest(errors.Errorf("context done: %w", ctx.Err()))
		case status := <-sub():
			if status.IsFailure() {
				if strings.Contains(utils.SafeErrStr(task.Error()), "validate ticket") {
					return nil, nft.MakeBadRequest(errors.New("ticket not found. Please make sure you are using correct registration ticket TXID"))
				}

				if strings.Contains(utils.SafeErrStr(task.Error()), "ticket ownership") {
					return nil, nft.MakeBadRequest(errors.New("failed to verify ownership"))
				}

				errStr := fmt.Errorf("internal processing error: %s", status.String())
				if task.Error() != nil {
					errStr = task.Error()
				}
				return nil, nft.MakeInternalServerError(errStr)
			}

			if status.IsFinal() {
				if len(task.File) == 0 {
					return nil, nft.MakeInternalServerError(errors.New("unable to download file"))
				}

				log.WithContext(ctx).WithField("size in KB", len(task.File)/1000).WithField("txid", p.Txid).Info("File downloaded")

				err := os.WriteFile(filepath.Join(service.config.StaticFilesDir, task.Filename), task.File, 0644)
				if err != nil {
					return nil, nft.MakeInternalServerError(errors.New("unable to write file"))
				}

				return &nft.FileDownloadResult{
					FileID: task.Filename,
				}, nil
			}
		}
	}
}

// --- Search registered NFTs and return details ---

// NftSearch searches for NFT & streams the result based on filters
func (service *NftAPIHandler) NftSearch(ctx context.Context, p *nft.NftSearchPayload, stream nft.NftSearchServerStream) error {
	defer stream.Close()
	taskID := service.search.AddTask(p)
	task := service.search.GetTask(taskID)

	resultChan := task.SubscribeSearchResult()
	for {
		select {
		case <-ctx.Done():
			return nil
		case search, ok := <-resultChan:
			if !ok {
				if task.Status().IsFailure() {
					log.WithContext(ctx).Error("error finding NFT")
					return nft.MakeInternalServerError(task.Error())
				}

				return nil
			}

			res := toNftSearchResult(search)
			if err := stream.Send(res); err != nil {
				return nft.MakeInternalServerError(err)
			}
		}
	}
}

// NftGet returns NFT detail
func (service *NftAPIHandler) NftGet(ctx context.Context, p *nft.NftGetPayload) (res *nft.NftDetail, err error) {
	ticket, err := service.search.RegTicket(ctx, p.Txid)
	if err != nil {
		log.WithError(err).Error("error retrieving ticket")
		return nil, nft.MakeBadRequest(err)
	}

	res = toNftDetail(ticket)
	data, err := service.search.GetThumbnail(ctx, ticket, p.Pid, p.Key)
	if err != nil {
		log.WithError(err).Error("error retrieving thumbnail")

		return nil, nft.MakeInternalServerError(err)
	}
	res.PreviewThumbnail = data

	ddAndFpData, err := service.search.GetDDAndFP(ctx, ticket, p.Pid, p.Key)
	if err != nil {
		log.WithError(err).Error("error retrieving DD&FP")
		return nil, nft.MakeInternalServerError(err)
	}
	ddAndFpStruct := &pastel.DDAndFingerprints{}
	json.Unmarshal(ddAndFpData, ddAndFpStruct)

	res.NsfwScore = ddAndFpStruct.OpenNSFWScore
	res.HentaiNsfwScore = &ddAndFpStruct.AlternativeNSFWScores.Hentai
	res.PornNsfwScore = &ddAndFpStruct.AlternativeNSFWScores.Porn
	res.SexyNsfwScore = &ddAndFpStruct.AlternativeNSFWScores.Sexy
	res.DrawingNsfwScore = &ddAndFpStruct.AlternativeNSFWScores.Drawings
	res.NeutralNsfwScore = &ddAndFpStruct.AlternativeNSFWScores.Neutral

	res.RarenessScore = ddAndFpStruct.OverallRarenessScore
	res.IsLikelyDupe = ddAndFpStruct.IsLikelyDupe
	res.IsRareOnInternet = ddAndFpStruct.IsRareOnInternet

	res.RareOnInternetSummaryTableJSONB64 = &ddAndFpStruct.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64
	res.RareOnInternetGraphJSONB64 = &ddAndFpStruct.InternetRareness.RareOnInternetGraphJSONCompressedB64
	res.AltRareOnInternetDictJSONB64 = &ddAndFpStruct.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64
	res.MinNumExactMatchesOnPage = &ddAndFpStruct.InternetRareness.MinNumberOfExactMatchesInPage
	res.EarliestDateOfResults = &ddAndFpStruct.InternetRareness.EarliestAvailableDateOfInternetResults

	return res, nil
}

// DdServiceOutputFileDetail returns Dupe detection output file details
func (service *NftAPIHandler) DdServiceOutputFileDetail(ctx context.Context, p *nft.DownloadPayload) (res *nft.DDServiceOutputFileResult, err error) {
	ticket, err := service.search.RegTicket(ctx, p.Txid)
	if err != nil {
		log.WithError(err).Error("error retrieving ticket")
		return nil, nft.MakeBadRequest(err)
	}

	ddAndFpData, err := service.search.GetDDAndFP(ctx, ticket, p.Pid, p.Key)
	if err != nil {
		log.WithError(err).Error("error retrieving DD&FP")
		return nil, nft.MakeInternalServerError(err)
	}
	ddAndFpStruct := &pastel.DDAndFingerprints{}
	json.Unmarshal(ddAndFpData, ddAndFpStruct)

	res = &nft.DDServiceOutputFileResult{}
	res = translateNftSummary(res, ticket)

	res = translateDDServiceOutputFile(res, ddAndFpStruct)

	return res, nil
}

// DdServiceOutputFile returns Dupe detection output file
func (service *NftAPIHandler) DdServiceOutputFile(ctx context.Context, p *nft.DownloadPayload) (res *nft.DDFPResultFile, err error) {
	ticket, err := service.search.RegTicket(ctx, p.Txid)
	if err != nil {
		log.WithError(err).Error("error retrieving ticket")
		return nil, nft.MakeBadRequest(err)
	}

	ddAndFpData, err := service.search.GetDDAndFP(ctx, ticket, p.Pid, p.Key)
	if err != nil {
		log.WithError(err).Error("error retrieving DD&FP")
		return nil, nft.MakeInternalServerError(err)
	}
	ddAndFpStruct := &pastel.DDAndFingerprints{}
	json.Unmarshal(ddAndFpData, ddAndFpStruct)

	DDFPFileDetails := toNFTDDServiceFile(ticket.RegTicketData.NFTTicketData.AppTicketData, ddAndFpStruct)

	fileBytes, err := json.Marshal(DDFPFileDetails)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting file details to bytes")
	}

	// Convert json data to base64
	base64Data := base64.StdEncoding.EncodeToString(fileBytes)

	res = &nft.DDFPResultFile{
		File: base64Data,
	}
	return res, nil
}

// NewNftAPIHandler returns the nft NftAPIHandler implementation.
func NewNftAPIHandler(config *Config, register *nftregister.NftRegistrationService, search *nftsearch.NftSearchingService, download *download.NftDownloadingService) *NftAPIHandler {
	return &NftAPIHandler{
		Common:   NewCommon(config),
		register: register,
		search:   search,
		download: download,
		db:       memory.NewKeyValue(),
		imageTTL: defaultImageTTL,
	}
}
