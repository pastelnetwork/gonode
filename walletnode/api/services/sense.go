package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// SenseAPIHandler - SenseAPIHandler service
type SenseAPIHandler struct {
	*Common
	register *senseregister.SenseRegistrationService
	download *download.NftDownloadingService
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *SenseAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := sense.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		SenseUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadImage - Uploads an image and return unique image id
func (service *SenseAPIHandler) UploadImage(ctx context.Context, p *sense.UploadImagePayload) (res *sense.Image, err error) {
	if p.Filename == nil {
		log.Error("filename not specified")
		return nil, sense.MakeBadRequest(errors.New("file not specified"))
	}

	id, expiry, err := service.register.StoreFile(ctx, p.Filename)
	if err != nil {
		log.WithError(err).Error("error storing file")
		return nil, sense.MakeInternalServerError(err)
	}
	log.Infof("file has been uploaded: %s", id)

	fee, err := service.register.CalculateFee(ctx, id)
	if err != nil {
		log.WithError(err).Error("error calculating fee")
		return nil, sense.MakeInternalServerError(err)
	}
	log.Infof("estimated fee has been calculated: %f", fee)

	res = &sense.Image{
		ImageID:      id,
		ExpiresIn:    expiry,
		EstimatedFee: fee,
	}

	return res, nil
}

// StartProcessing - Starts a processing image task
func (service *SenseAPIHandler) StartProcessing(_ context.Context, p *sense.StartProcessingPayload) (res *sense.StartProcessingResult, err error) {
	taskID, err := service.register.AddTask(p)
	if err != nil {
		log.WithError(err).Error("unable to add task")
		return nil, sense.MakeInternalServerError(err)
	}

	res = &sense.StartProcessingResult{
		TaskID: taskID,
	}

	log.Infof("task has been added: %s", taskID)
	return res, nil
}

// RegisterTaskState - Registers a task state
func (service *SenseAPIHandler) RegisterTaskState(ctx context.Context, p *sense.RegisterTaskStatePayload, stream sense.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		log.Error("unable to get task")
		return sense.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}
	log.Info("task has been retrieved")

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &sense.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return sense.MakeInternalServerError(err)
			}

			if status.IsFailure() {
				if task.Error() != nil {
					errStr := task.Error()
					log.WithContext(ctx).WithError(errStr).Errorf("error registering sense")
				}

			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// GetTaskHistory - Gets a task's history
func (service *SenseAPIHandler) GetTaskHistory(ctx context.Context, p *sense.GetTaskHistoryPayload) (res []*sense.TaskHistory, err error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		return nil, sense.MakeInternalServerError(errors.New("error retrieving status"))
	}
	defer store.CloseHistoryDB(ctx)

	statuses, err := store.QueryTaskHistory(p.TaskID)
	if err != nil {
		return nil, sense.MakeNotFound(errors.New("task not found"))
	}

	history := []*sense.TaskHistory{}
	for _, entry := range statuses {
		timestamp := entry.CreatedAt.String()
		historyItem := &sense.TaskHistory{
			Timestamp: &timestamp,
			Status:    entry.Status,
		}

		if entry.Details != nil {
			historyItem.Details = &sense.Details{
				Message: &entry.Details.Message,
				Fields:  entry.Details.Fields,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *SenseAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// Download registered NFT
func (service *SenseAPIHandler) Download(ctx context.Context, p *sense.DownloadPayload) (res *sense.DownloadResult, err error) {
	log.Info("Start downloading")
	defer log.Info("Finished downloading")
	taskID := service.download.AddTask(&nft.DownloadPayload{Txid: p.Txid, Pid: p.Pid, Key: p.Key}, pastel.ActionTypeSense)
	task := service.download.GetTask(taskID)
	defer task.Cancel()

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil, sense.MakeBadRequest(errors.Errorf("context done: %w", ctx.Err()))
		case status := <-sub():
			if status.IsFailure() {
				if strings.Contains(utils.SafeErrStr(task.Error()), "ticket ownership") {
					return nil, sense.MakeBadRequest(errors.New("failed to verify ownership"))
				}

				errStr := fmt.Errorf("internal processing error: %s", status.String())
				if task.Error() != nil {
					errStr = task.Error()
				}
				return nil, sense.MakeInternalServerError(errStr)
			}

			if status.IsFinal() {
				log.WithContext(ctx).WithField("size", fmt.Sprintf("%d bytes", len(task.File))).Info("NFT downloaded")
				res = &sense.DownloadResult{
					File: task.File,
				}

				return res, nil
			}
		}
	}
}

// NewSenseAPIHandler returns the swagger OpenAPI implementation.
func NewSenseAPIHandler(register *senseregister.SenseRegistrationService, download *download.NftDownloadingService) *SenseAPIHandler {
	return &SenseAPIHandler{
		Common:   NewCommon(),
		register: register,
		download: download,
	}
}
