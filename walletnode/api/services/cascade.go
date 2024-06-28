package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/cascade/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

const (
	maxFileSize                 = 350 * 1024 * 1024 // 350MB in bytes
	maxFileRegistrationAttempts = 3
	downloadDeadline            = 30 * time.Minute
)

// CascadeAPIHandler - CascadeAPIHandler service
type CascadeAPIHandler struct {
	*Common
	register     *cascaderegister.CascadeRegistrationService
	download     *download.NftDownloadingService
	fileMappings *sync.Map // maps unique ID to file path
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *CascadeAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := cascade.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		CascadeUploadAssetDecoderFunc(ctx, service),
		CascadeUploadAssetV2DecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadAssetV2 - Uploads an asset file and return unique file id
func (service *CascadeAPIHandler) UploadAssetV2(ctx context.Context, p *cascade.UploadAssetV2Payload) (res *cascade.AssetV2, err error) {
	if p.Filename == nil {
		log.Error("file not specified")
		return nil, cascade.MakeBadRequest(errors.New("file not specified"))
	}

	fee, preburn, err := service.register.StoreFileMetadata(ctx, filepath.Join(service.config.CascadeFilesDir, *p.Filename), *p.Hash, *p.Size)
	if err != nil {
		log.WithError(err).Error(err)
		return nil, cascade.MakeInternalServerError(err)
	}

	res = &cascade.AssetV2{
		FileID:                            *p.Filename,
		TotalEstimatedFee:                 fee,
		RequiredPreburnTransactionAmounts: preburn,
	}

	return res, nil
}

// UploadAsset - Uploads an asset file and return unique file id
func (service *CascadeAPIHandler) UploadAsset(ctx context.Context, p *cascade.UploadAssetPayload) (res *cascade.Asset, err error) {
	if p.Filename == nil {
		log.Error("file not specified")
		return nil, cascade.MakeBadRequest(errors.New("file not specified"))
	}

	fee, preburn, err := service.register.StoreFileMetadata(ctx, filepath.Join(service.config.CascadeFilesDir, *p.Filename), *p.Hash, *p.Size)
	if err != nil {
		log.WithError(err).Error(err)
		return nil, cascade.MakeInternalServerError(err)
	}

	res = &cascade.Asset{
		FileID:                *p.Filename,
		TotalEstimatedFee:     fee,
		RequiredPreburnAmount: preburn[0],
	}

	return res, nil
}

// StartProcessing - Starts a processing image task
func (service *CascadeAPIHandler) StartProcessing(ctx context.Context, p *cascade.StartProcessingPayload) (res *cascade.StartProcessingResult, err error) {
	if !service.register.ValidateUser(ctx, p.AppPastelID, p.Key) {
		return nil, cascade.MakeUnAuthorized(errors.New("user not authorized: invalid PastelID or Key"))
	}

	// Get related files
	relatedFiles, err := service.register.GetFilesByBaseFileID(p.FileID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}
	log.WithContext(ctx).WithField("total_volumes", len(relatedFiles)).Info("Related volumes retrieved from the base file-id")

	switch {
	case len(relatedFiles) == 1:
		baseFile := relatedFiles.GetBase()
		if baseFile == nil {
			return nil, cascade.MakeInternalServerError(err)
		}

		taskID, err := service.register.ProcessFile(ctx, *baseFile, p)
		if err != nil {
			return nil, cascade.MakeBadRequest(err)
		}

		return &cascade.StartProcessingResult{
			TaskID: taskID,
		}, nil

	case len(relatedFiles) > 1:
		log.WithContext(ctx).Info("multi-volume registration...")

		if len(p.BurnTxids) != len(relatedFiles) {
			log.WithContext(ctx).WithField("related_volumes", len(relatedFiles)).
				WithField("burn_txids", len(p.BurnTxids)).
				Info("no of provided burn txids and volumes are not equal")
			return nil, cascade.MakeBadRequest(errors.New("provided burn txids and no of volumes are not equal"))
		}

		sortedBurnTxids, err := service.register.SortBurnTxIDs(ctx, p.BurnTxids)
		if err != nil {
			return nil, cascade.MakeInternalServerError(err)
		}
		sortedRelatedFiles := service.register.SortFilesWithHigherAmounts(relatedFiles)

		var taskIDs []string
		for index, file := range sortedRelatedFiles {
			burnTxID := sortedBurnTxids[index]
			p.BurnTxid = &burnTxID
			taskID, err := service.register.ProcessFile(ctx, *file, p)
			if err != nil {
				log.WithContext(ctx).WithField("file_id", file.FileID).WithError(err).Error("error processing volume")
				continue
			}

			if taskID != "" {
				taskIDs = append(taskIDs, taskID)
			}
		}

		return &cascade.StartProcessingResult{
			TaskID: strings.Join(taskIDs, ","),
		}, nil
	}

	return res, nil
}

// RegisterTaskState - Registers a task state
func (service *CascadeAPIHandler) RegisterTaskState(ctx context.Context, p *cascade.RegisterTaskStatePayload, stream cascade.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		log.Error("unable to get task")
		return cascade.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &cascade.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return cascade.MakeInternalServerError(err)
			}

			if status.IsFailure() {
				if task.Error() != nil {
					errStr := task.Error()
					log.WithContext(ctx).WithError(errStr).Errorf("error registering cascade")
				}

			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *CascadeAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// Download registered cascade file - also supports multi-volume files
func (service *CascadeAPIHandler) Download(ctx context.Context, p *cascade.DownloadPayload) (*cascade.FileDownloadResult, error) {
	log.WithContext(ctx).WithField("txid", p.Txid).Info("Start downloading")

	if !service.register.ValidateUser(ctx, p.Pid, p.Key) {
		return nil, cascade.MakeUnAuthorized(errors.New("user not authorized: invalid PastelID or Key"))
	}
	defer log.WithContext(ctx).WithField("txid", p.Txid).Info("Finished downloading")

	var txIDs []string
	isMultiVolume := false
	var ticket pastel.CascadeMultiVolumeTicket
	c, err := service.download.CheckForMultiVolumeCascadeTicket(ctx, p.Txid)
	if err == nil {
		ticket, err = c.GetCascadeMultiVolumeMetadataTicket()
		if err == nil {
			for _, volumeTxID := range ticket.Volumes {
				isMultiVolume = true
				txIDs = append(txIDs, volumeTxID)
			}
		} else {
			txIDs = append(txIDs, p.Txid)
		}
	} else {
		txIDs = append(txIDs, p.Txid)
	}

	type DownloadResult struct {
		File     []byte
		Filename string
		Txid     string
	}

	// Channel to control the concurrency of downloads
	sem := make(chan struct{}, 3) // Max 3 concurrent downloads
	taskResults := make(chan *DownloadResult)
	errorsChan := make(chan error)

	ctx, cancel := context.WithTimeout(ctx, downloadDeadline)
	defer cancel()

	// Starting multiple download tasks
	for _, txID := range txIDs {
		go func(txID string) {
			sem <- struct{}{}        // Acquiring the semaphore
			defer func() { <-sem }() // Releasing the semaphore

			taskID := service.download.AddTask(&nft.DownloadPayload{Key: p.Key, Pid: p.Pid, Txid: txID}, pastel.ActionTypeCascade, false)
			task := service.download.GetTask(taskID)
			defer task.Cancel()

			sub := task.SubscribeStatus()

			for {
				select {
				case <-ctx.Done():
					errorsChan <- cascade.MakeBadRequest(errors.Errorf("context done: %w", ctx.Err()))
					return
				case status := <-sub():
					if status.IsFailure() {
						if strings.Contains(utils.SafeErrStr(task.Error()), "validate ticket") {
							errorsChan <- cascade.MakeBadRequest(errors.New("ticket not found. Please make sure you are using correct registration ticket TXID"))
							return
						}

						if strings.Contains(utils.SafeErrStr(task.Error()), "ticket ownership") {
							errorsChan <- cascade.MakeBadRequest(errors.New("failed to verify ownership"))
							return
						}

						errStr := fmt.Errorf("internal processing error: %s", status.String())
						if task.Error() != nil {
							errStr = task.Error()
						}
						errorsChan <- cascade.MakeInternalServerError(errStr)
						return
					}

					if status.IsFinal() {
						taskResults <- &DownloadResult{File: task.File, Filename: task.Filename, Txid: txID}
						return
					}
				}
			}
		}(txID)
	}

	// Create directory with p.Txid
	folderPath := filepath.Join(service.config.StaticFilesDir, p.Txid)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			cancel()
			return nil, err
		}
	}

	var filePath string
	for i := 0; i < len(txIDs); i++ {
		select {
		case res := <-taskResults:
			filePath = filepath.Join(folderPath, res.Filename)
			if err := os.WriteFile(filePath, res.File, 0644); err != nil {
				cancel()
				return nil, cascade.MakeInternalServerError(errors.New("unable to write file"))
			}
		case err := <-errorsChan:
			cancel()
			return nil, err
		}
	}

	if isMultiVolume {
		fsp := common.FileSplitter{PartSizeMB: 300}
		if err := fsp.JoinFiles(folderPath); err != nil {
			return nil, cascade.MakeInternalServerError(errors.New("unable to join files"))
		}

		filePath = filepath.Join(folderPath, ticket.NameOfOriginalFile)

		// Check if the hash of the file matches the hash in the ticket
		hash, err := utils.ComputeSHA256HashOfFile(filePath)
		if err != nil {
			return nil, cascade.MakeInternalServerError(errors.New("unable to compute hash"))
		}

		if hex.EncodeToString(hash) != ticket.SHA3256HashOfOriginalFile {
			return nil, cascade.MakeInternalServerError(errors.New("hash mismatch"))
		}
	}

	// generating a single unique ID
	uniqueID := randIDFunc()
	service.fileMappings.Store(uniqueID, filePath)

	return &cascade.FileDownloadResult{
		FileID: uniqueID,
	}, nil
}

// GetTaskHistory - Gets a task's history
func (service *CascadeAPIHandler) GetTaskHistory(ctx context.Context, p *cascade.GetTaskHistoryPayload) (history []*cascade.TaskHistory, err error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, cascade.MakeInternalServerError(errors.New("error retrieving status"))
	}
	defer store.CloseHistoryDB(ctx)

	statuses, err := store.QueryTaskHistory(p.TaskID)
	if err != nil {
		return nil, cascade.MakeNotFound(errors.New("task not found"))
	}

	for _, entry := range statuses {
		timestamp := entry.CreatedAt.String()
		historyItem := &cascade.TaskHistory{
			Timestamp: &timestamp,
			Status:    entry.Status,
		}

		if entry.Details != nil {
			historyItem.Details = &cascade.Details{
				Message: &entry.Details.Message,
				Fields:  entry.Details.Fields,
			}
		}

		history = append(history, historyItem)
	}

	return history, nil
}

func (service *CascadeAPIHandler) RegistrationDetails(ctx context.Context, rdp *cascade.RegistrationDetailsPayload) (registrationDetail *cascade.Registration, err error) {
	log.WithContext(ctx).WithField("file_id", rdp.FileID).Info("Registration detail api invoked")
	defer log.WithContext(ctx).WithField("file_id", rdp.FileID).Info("Finished registration details")

	baseFile, err := service.register.GetFile(rdp.FileID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	relatedFiles, err := service.register.GetFilesByBaseFileID(baseFile.FileID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	var fileDetails []*cascade.File
	for _, relatedFile := range relatedFiles {
		relatedFileActivationAttempts, err := service.register.GetActivationAttemptsByFileID(relatedFile.FileID)
		if err != nil {
			return nil, cascade.MakeInternalServerError(err)
		}

		relatedFileRegistrationAttempts, err := service.register.GetRegistrationAttemptsByFileID(relatedFile.FileID)
		if err != nil {
			return nil, cascade.MakeInternalServerError(err)
		}

		var activationAttempts []*cascade.ActivationAttempt
		for _, aAttempt := range relatedFileActivationAttempts {
			activationAttempts = append(activationAttempts, &cascade.ActivationAttempt{
				ID:                  aAttempt.ID,
				FileID:              aAttempt.FileID,
				ActivationAttemptAt: aAttempt.ActivationAttemptAt.String(),
				IsSuccessful:        &aAttempt.IsSuccessful,
				ErrorMessage:        &aAttempt.ErrorMessage,
			})
		}

		var registrationAttempts []*cascade.RegistrationAttempt
		for _, rAttempt := range relatedFileRegistrationAttempts {
			registrationAttempts = append(registrationAttempts, &cascade.RegistrationAttempt{
				ID:           rAttempt.ID,
				FileID:       rAttempt.FileID,
				RegStartedAt: rAttempt.RegStartedAt.String(),
				FinishedAt:   rAttempt.FinishedAt.String(),
				IsSuccessful: &rAttempt.IsSuccessful,
				ErrorMessage: &rAttempt.ErrorMessage,
				ProcessorSns: &rAttempt.ProcessorSNS,
			})
		}

		fileDetails = append(fileDetails, &cascade.File{
			FileID:                       relatedFile.FileID,
			UploadTimestamp:              relatedFile.UploadTimestamp.String(),
			Path:                         &relatedFile.Path,
			FileIndex:                    &relatedFile.FileIndex,
			BaseFileID:                   relatedFile.BaseFileID,
			TaskID:                       relatedFile.TaskID,
			RegTxid:                      &relatedFile.RegTxid,
			ActivationTxid:               &relatedFile.ActivationTxid,
			ReqBurnTxnAmount:             relatedFile.ReqBurnTxnAmount,
			BurnTxnID:                    &relatedFile.BurnTxnID,
			ReqAmount:                    relatedFile.ReqAmount,
			IsConcluded:                  &relatedFile.IsConcluded,
			CascadeMetadataTicketID:      relatedFile.CascadeMetadataTicketID,
			UUIDKey:                      &relatedFile.UUIDKey,
			HashOfOriginalBigFile:        relatedFile.HashOfOriginalBigFile,
			NameOfOriginalBigFileWithExt: relatedFile.NameOfOriginalBigFileWithExt,
			SizeOfOriginalBigFile:        relatedFile.SizeOfOriginalBigFile,
			DataTypeOfOriginalBigFile:    relatedFile.DataTypeOfOriginalBigFile,
			StartBlock:                   &relatedFile.StartBlock,
			DoneBlock:                    &relatedFile.DoneBlock,

			ActivationAttempts:   activationAttempts,
			RegistrationAttempts: registrationAttempts,
		})
	}

	return &cascade.Registration{
		Files: fileDetails,
	}, nil
}

// NewCascadeAPIHandler returns the swagger OpenAPI implementation.
func NewCascadeAPIHandler(config *Config, filesMap *sync.Map, register *cascaderegister.CascadeRegistrationService, download *download.NftDownloadingService) *CascadeAPIHandler {
	return &CascadeAPIHandler{
		Common:       NewCommon(config),
		register:     register,
		download:     download,
		fileMappings: filesMap,
	}
}
