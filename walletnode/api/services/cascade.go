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

	"github.com/pastelnetwork/gonode/common/types"

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
	maxFileRegistrationAttempts = 3
	downloadDeadline            = 120 * time.Minute
	downloadConcurrency         = 2
	statusCompleted             = "Completed"
	statusInProgress            = "In Progress"
	statusFailed                = "Failed"
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

	if isProcessingStarted(relatedFiles) {
		log.WithContext(ctx).Error("task registration is already in progress")
		return nil, cascade.MakeInternalServerError(errors.New("task registration is already in progress"))
	}

	switch {
	case len(relatedFiles) == 1:
		baseFile := relatedFiles.GetBase()
		if baseFile == nil {
			return nil, cascade.MakeInternalServerError(err)
		}

		burnTxID, err := validateBurnTxID(p.BurnTxid, p.BurnTxids)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error validating burn-txid from the payload")
			return nil, cascade.MakeInternalServerError(err)
		}

		addTaskPayload := &common.AddTaskPayload{
			FileID:                 p.FileID,
			BurnTxid:               &burnTxID,
			AppPastelID:            p.AppPastelID,
			MakePubliclyAccessible: p.MakePubliclyAccessible,
			Key:                    p.Key,
		}

		if p.SpendableAddress != nil {
			addTaskPayload.SpendableAddress = p.SpendableAddress
		}

		taskID, err := service.register.ProcessFile(ctx, *baseFile, addTaskPayload)
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

		err = service.checkBurnTxIDsValidForRegistration(ctx, sortedBurnTxids, sortedRelatedFiles)
		if err != nil {
			log.WithContext(ctx).WithField("related_volumes", len(relatedFiles)).
				WithError(err).
				WithField("burn_txids", len(p.BurnTxids)).
				Error("given burn-tx-ids are not valid")
			return nil, cascade.MakeBadRequest(errors.New("given burn txids are not valid"))
		}

		var taskIDs []string
		for index, file := range sortedRelatedFiles {
			burnTxID := sortedBurnTxids[index]
			p.BurnTxid = &burnTxID

			addTaskPayload := &common.AddTaskPayload{
				FileID:                 file.FileID,
				BurnTxid:               p.BurnTxid,
				AppPastelID:            p.AppPastelID,
				MakePubliclyAccessible: p.MakePubliclyAccessible,
				Key:                    p.Key,
			}

			if p.SpendableAddress != nil {
				addTaskPayload.SpendableAddress = p.SpendableAddress
			}

			taskID, err := service.register.ProcessFile(ctx, *file, addTaskPayload)
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

func isDuplicateExists(burnTxIDs []string) bool {
	duplicateBurnTxIDMap := make(map[string]int)

	for _, burnTxID := range burnTxIDs {
		duplicateBurnTxIDMap[burnTxID]++
	}

	for _, count := range duplicateBurnTxIDMap {
		if count > 1 {
			return true
		}
	}

	return false
}

func isProcessingStarted(files types.Files) bool {
	for _, file := range files {
		if file.TaskID != "" {
			return true
		}
	}

	return false
}

func validateBurnTxID(burnTxID *string, burnTxIDs []string) (string, error) {
	if burnTxID != nil {
		return *burnTxID, nil
	}

	if len(burnTxIDs) > 0 {
		return burnTxIDs[0], nil
	}

	return "", errors.New("burn txid should either passed in burn_txid field or as an item in burn_txids array")
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

func (service *CascadeAPIHandler) getTxIDs(ctx context.Context, txID string) (txIDs []string, ticket pastel.CascadeMultiVolumeTicket) {
	c, err := service.download.CheckForMultiVolumeCascadeTicket(ctx, txID)
	if err == nil {
		ticket, err = c.GetCascadeMultiVolumeMetadataTicket()
		if err == nil {
			for _, volumeTxID := range ticket.Volumes {
				txIDs = append(txIDs, volumeTxID)
			}

			return
		} else {
			filename, size, err := service.download.GetFilenameAndSize(ctx, txID)
			if err != nil {
				return nil, ticket
			}

			ticket.NameOfOriginalFile = filename
			ticket.SizeOfOriginalFileMB = utils.BytesIntToMB(size)

		}
	}
	txIDs = append(txIDs, txID)

	return
}

type DownloadResult struct {
	File     []byte
	Filename string
	Txid     string
}

// Download registered cascade file - also supports multi-volume files
func (service *CascadeAPIHandler) Download(ctx context.Context, p *cascade.DownloadPayload) (*cascade.FileDownloadResult, error) {
	log.WithContext(ctx).WithField("txid", p.Txid).Info("Start downloading")

	if !service.register.ValidateUser(ctx, p.Pid, p.Key) {
		return nil, cascade.MakeUnAuthorized(errors.New("user not authorized: invalid PastelID or Key"))
	}
	defer log.WithContext(ctx).WithField("txid", p.Txid).Info("Finished downloading")

	txIDs, ticket := service.getTxIDs(ctx, p.Txid)
	isMultiVolume := len(txIDs) > 1
	if isMultiVolume {
		return nil, cascade.MakeBadRequest(errors.New("multi-volume download not supported. Please use /v2/download endpoint to download multi-volume files"))
	}

	// Create directory with p.Txid
	folderPath := filepath.Join(service.config.StaticFilesDir, p.Txid)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// Channel to control the concurrency of downloads
	sem := make(chan struct{}, downloadConcurrency) // Max 3 concurrent downloads
	taskResults := make(chan *DownloadResult)
	errorsChan := make(chan error)

	ctx, cancel := context.WithTimeout(ctx, downloadDeadline)
	defer cancel()

	// Starting multiple download tasks
	tasks := 0
	for _, txID := range txIDs {

		filename, size, err := service.download.GetFilenameAndSize(ctx, txID)
		if err != nil {
			return nil, err
		}

		// check if file already exists
		filePath := filepath.Join(folderPath, filename)
		if fileinfo, err := os.Stat(filePath); err == nil {
			if fileinfo.Size() == int64(size) {
				log.WithContext(ctx).WithField("filename", filename).WithField("txid", txID).Info("skipping file download as it already exists")
				continue
			} else {
				log.WithContext(ctx).WithField("filename", filename).WithField("txid", txID).Warn("file exists but downloading again due to size mismatch")
			}
		}

		tasks++
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

	var filePath string
	for i := 0; i < tasks; i++ {
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
		fsp := common.FileSplitter{PartSizeMB: partSizeMB}
		if err := fsp.JoinFiles(folderPath); err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to join files")
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

func getTxIDsKey(id string) string {
	return fmt.Sprintf("%s_txids", id)
}

func getPrimaryTxIDKey(id string) string {
	return fmt.Sprintf("%s_primary_txid", id)
}

// Download registered cascade file - also supports multi-volume files
func (service *CascadeAPIHandler) DownloadV2(ctx context.Context, p *cascade.DownloadPayload) (*cascade.FileDownloadV2Result, error) {
	log.WithContext(ctx).WithField("txid", p.Txid).Info("Download Request Received (V2)")
	if !service.register.ValidateUser(ctx, p.Pid, p.Key) {
		return nil, cascade.MakeUnAuthorized(errors.New("user not authorized: invalid PastelID or Key"))
	}

	txIDs, ticket := service.getTxIDs(ctx, p.Txid)
	isMultiVolume := len(txIDs) > 1

	// Create directory with p.Txid
	folderPath := filepath.Join(service.config.StaticFilesDir, p.Txid)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// generating a single unique ID
	uniqueID := randIDFunc()
	service.fileMappings.Store(uniqueID, folderPath)
	service.fileMappings.Store(getPrimaryTxIDKey(uniqueID), p.Txid)
	// Start asynchronous download process
	go service.manageDownloads(ctx, p, txIDs, folderPath, isMultiVolume, ticket.NameOfOriginalFile, ticket.SHA3256HashOfOriginalFile, uniqueID)

	return &cascade.FileDownloadV2Result{
		FileID: uniqueID,
	}, nil
}

func (service *CascadeAPIHandler) manageDownloads(_ context.Context, p *cascade.DownloadPayload, txIDs []string, folderPath string, isMulti bool, name string, sha string, id string) {
	sem := make(chan struct{}, downloadConcurrency) // Control Max concurrent downloads
	taskResults := make(chan *DownloadResult)
	errorsChan := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), downloadDeadline)
	defer cancel()

	service.fileMappings.Store(getTxIDsKey(id), txIDs)

	// Starting multiple download tasks
	tasks := 0
	for _, txID := range txIDs {
		filename, size, err := service.download.GetFilenameAndSize(ctx, txID)
		if err != nil {
			return
		}

		// check if file already exists
		filePath := filepath.Join(folderPath, filename)
		if fileinfo, err := os.Stat(filePath); err == nil {
			if fileinfo.Size() == int64(size) {
				log.WithContext(ctx).WithField("filename", filename).WithField("txid", txID).Info("skipping file download as it already exists")
				if len(txIDs) == 1 {
					service.fileMappings.Store(id, filePath)
				}

				service.fileMappings.Store(txID, txID)
				continue
			} else {
				log.WithContext(ctx).WithField("filename", filename).WithField("txid", txID).Warn("file exists but downloading again due to size mismatch")
			}
		}
		tasks++

		go func(txID string) {
			sem <- struct{}{}        // Acquiring the semaphore
			defer func() { <-sem }() // Releasing the semaphore

			taskID := service.download.AddTask(&nft.DownloadPayload{Key: p.Key, Pid: p.Pid, Txid: txID}, pastel.ActionTypeCascade, false)
			service.fileMappings.Store(txID, taskID)

			task := service.download.GetTask(taskID)
			defer task.Cancel()

			sub := task.SubscribeStatus()

			for {
				select {
				case <-ctx.Done():
					errorsChan <- cascade.MakeBadRequest(errors.Errorf("context done: %w", ctx.Err()))
					log.WithContext(ctx).Info("context done down v2")
					return
				case status := <-sub():
					if status.IsFailure() {
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
	log.WithContext(ctx).WithField("tasks-count", tasks).Info("Download tasks started")

	// Processing results
	service.processDownloadResults(ctx, taskResults, errorsChan, folderPath, isMulti, tasks, sha, name, id)
}

func (service *CascadeAPIHandler) GetDownloadTaskState(ctx context.Context, p *cascade.GetDownloadTaskStatePayload) (th *cascade.DownloadTaskStatus, err error) {
	filePath, ok := service.fileMappings.Load(p.FileID)
	if !ok {
		return nil, cascade.MakeNotFound(errors.New("file not found"))
	}

	filePathStr, ok := filePath.(string)
	if !ok {
		return nil, cascade.MakeInternalServerError(errors.New("unable to cast file path"))
	}

	// Check if the file exists
	primaryTxID, ok := service.fileMappings.Load(getPrimaryTxIDKey(p.FileID))
	if !ok {
		return nil, cascade.MakeNotFound(errors.New("file not found - please try requesting download again"))
	}

	_, ticket := service.getTxIDs(ctx, primaryTxID.(string))
	if _, err := os.Stat(filepath.Join(filePathStr, ticket.NameOfOriginalFile)); err == nil {
		th = &cascade.DownloadTaskStatus{
			TaskStatus:              statusCompleted,
			SizeOfTheFileMegabytes:  ticket.SizeOfOriginalFileMB,
			DataDownloadedMegabytes: ticket.SizeOfOriginalFileMB,
		}

		return th, nil
	}

	val, ok := service.fileMappings.Load(getTxIDsKey(p.FileID))
	if !ok {
		return nil, cascade.MakeNotFound(errors.New("No such file download in progress. Please double check the file id - Please request download again if the app was restarted"))
	}

	details := cascade.Details{
		Fields: make(map[string]any),
	}

	type txidDetails struct {
		Txid   string
		TaskID string
		Status string
		Error  string
	}

	txids := val.([]string)
	var success int
	var failed int
	var pending int
	var dataDownloaded int
	for _, txid := range txids {
		_, ticket := service.getTxIDs(ctx, txid)
		td := txidDetails{
			Txid: txid,
		}

		taskID, ok := service.fileMappings.Load(txid)
		if !ok || taskID == "" {
			pending++
			td.Status = "Pending"
			details.Fields[txid] = td
			continue

		}

		if taskID == txid {
			td.Status = common.StatusDownloaded.String()
			details.Fields[txid] = td
			dataDownloaded += ticket.SizeOfOriginalFileMB
			success++
			continue
		}

		task := service.download.GetTask(taskID.(string))
		if task == nil {
			if _, err := os.Stat(filepath.Join(filePathStr, ticket.NameOfOriginalFile)); err == nil {
				td.Status = common.StatusDownloaded.String()
				details.Fields[txid] = td
				success++
				dataDownloaded += ticket.SizeOfOriginalFileMB
				continue
			} else {
				return nil, cascade.MakeInternalServerError(fmt.Errorf("no such task found for txid %s - please try requesting download again", txid))
			}
		}
		td.TaskID = taskID.(string)

		statuses := task.StatusHistory()
		for _, status := range statuses {
			if status.IsFailure() {
				failed++
				errStr := fmt.Errorf("internal processing error: %s", status.String())
				if task.Error() != nil {
					errStr = task.Error()
				}
				td.Error = errStr.Error()
				details.Fields[txid] = td

				continue
			}

			if status.Is(common.StatusTaskCompleted) || status.Is(common.StatusDownloaded) {
				success++
				td.Status = common.StatusDownloaded.String()
				dataDownloaded += ticket.SizeOfOriginalFileMB
				details.Fields[txid] = td
				continue
			}
			td.Status = statusInProgress
			details.Fields[txid] = td
		}
	}

	total := len(txids)

	th = &cascade.DownloadTaskStatus{
		TotalVolumes:              total,
		DownloadedVolumes:         success,
		VolumesDownloadFailed:     failed,
		VolumesPendingDownload:    pending,
		VolumesDownloadInProgress: total - success - failed - pending,
		SizeOfTheFileMegabytes:    ticket.SizeOfOriginalFileMB,
		DataDownloadedMegabytes:   dataDownloaded,
	}

	if total == success {
		th.TaskStatus = statusCompleted
	} else if failed == 0 {
		th.TaskStatus = statusInProgress
	} else {
		th.TaskStatus = statusFailed
		msg := fmt.Sprintf("Download failed for %d volumes - please try requesting download again", failed)
		th.Message = &msg
	}

	th.Details = &details

	return th, nil
}

func (service *CascadeAPIHandler) processDownloadResults(ctx context.Context, taskResults chan *DownloadResult, errorsChan chan error, folderPath string,
	isMultiVolume bool, tasks int, sha string, name string, id string) error {

	filePath := ""
	for i := 0; i < tasks; i++ {
		select {
		case res := <-taskResults:
			filePath = filepath.Join(folderPath, res.Filename)
			if err := os.WriteFile(filePath, res.File, 0644); err != nil {
				log.WithContext(ctx).WithError(err).Error("unable to write file")
				return fmt.Errorf("unable to write file: %w", err)
			}
		case err := <-errorsChan:
			log.WithContext(ctx).WithError(err).Error("error during download process")
			return fmt.Errorf("error during download process: %w", err)
		case <-ctx.Done():
			log.WithContext(ctx).Info("download process context expired")
			return fmt.Errorf("download process context expired: %w", ctx.Err())
		}
	}

	if isMultiVolume {
		fsp := common.FileSplitter{PartSizeMB: partSizeMB}
		if err := fsp.JoinFiles(folderPath); err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to join files")
			return cascade.MakeInternalServerError(errors.New("unable to join files"))
		}

		filePath = filepath.Join(folderPath, name)

		// Check if the hash of the file matches the hash in the ticket
		hash, err := utils.ComputeSHA256HashOfFile(filePath)
		if err != nil {
			return fmt.Errorf("unable to compute hash: %w", err)
		}

		if hex.EncodeToString(hash) != sha {
			return fmt.Errorf("hash mismatch: %w", errors.New("hash mismatch"))
		}

		if err := deleteUnmatchedFiles(folderPath, name); err != nil {
			log.WithContext(ctx).WithError(err).Error("unable to delete volume files")
		}
	}
	if filePath != "" {
		service.fileMappings.Store(id, filePath)
	}
	log.WithContext(ctx).WithField("file_id", id).Info("file downloaded successfully")

	return nil
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
	log.WithContext(ctx).WithField("base_file_id", rdp.BaseFileID).Info("Registration detail api invoked")
	defer log.WithContext(ctx).WithField("base_file_id", rdp.BaseFileID).Info("Finished registration details")

	relatedFiles, err := service.register.GetFilesByBaseFileID(rdp.BaseFileID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	var metadataTicketTxID string
	cascadeMultiVolTicketTxIDMap, err := service.register.GetMultiVolTicketTxIDMap(rdp.BaseFileID)
	if cascadeMultiVolTicketTxIDMap != nil {
		metadataTicketTxID = cascadeMultiVolTicketTxIDMap.MultiVolCascadeTicketTxid
	}

	var fileDetails []*cascade.File
	for _, relatedFile := range relatedFiles {
		relatedFileActivationAttempts, err := service.register.GetActivationAttemptsByFileIDAndBaseFileID(relatedFile.FileID, relatedFile.BaseFileID)
		if err != nil {
			return nil, cascade.MakeInternalServerError(err)
		}

		relatedFileRegistrationAttempts, err := service.register.GetRegistrationAttemptsByFileIDAndBaseFileID(relatedFile.FileID, relatedFile.BaseFileID)
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
			FileIndex:                    &relatedFile.FileIndex,
			BaseFileID:                   relatedFile.BaseFileID,
			TaskID:                       relatedFile.TaskID,
			RegTxid:                      &relatedFile.RegTxid,
			ActivationTxid:               &relatedFile.ActivationTxid,
			ReqBurnTxnAmount:             relatedFile.ReqBurnTxnAmount,
			BurnTxnID:                    &relatedFile.BurnTxnID,
			ReqAmount:                    relatedFile.ReqAmount,
			IsConcluded:                  &relatedFile.IsConcluded,
			CascadeMetadataTicketID:      metadataTicketTxID,
			UUIDKey:                      &relatedFile.UUIDKey,
			HashOfOriginalBigFile:        relatedFile.HashOfOriginalBigFile,
			NameOfOriginalBigFileWithExt: relatedFile.NameOfOriginalBigFileWithExt,
			SizeOfOriginalBigFile:        relatedFile.SizeOfOriginalBigFile,
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

// Restore - to recover files registration and activation that was missed in initial file registration.
func (service *CascadeAPIHandler) Restore(ctx context.Context, p *cascade.RestorePayload) (res *cascade.RestoreFile, err error) {
	if !service.register.ValidateUser(ctx, p.AppPastelID, p.Key) {
		return nil, cascade.MakeUnAuthorized(errors.New("user not authorized: invalid PastelID or Key"))
	}

	restoreFile, err := service.register.RestoreFile(ctx, p)
	if err != nil {
		return nil, err
	}

	return restoreFile, nil
}

func (service *CascadeAPIHandler) checkBurnTxIDsValidForRegistration(ctx context.Context, burnTxIDs []string, files types.Files) error {
	if isDuplicateExists(burnTxIDs) {
		return errors.New("duplicate burn-tx-ids provided")
	}

	for _, bid := range burnTxIDs {
		if err := service.register.CheckBurnTxIDTicketDuplication(ctx, bid); err != nil {
			return err
		}
	}

	for i := 0; i < len(files); i++ {
		err := service.register.ValidBurnTxnAmount(ctx, burnTxIDs[i], files[i].ReqBurnTxnAmount)
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteUnmatchedFiles deletes all files in the specified folderPath that do not match the given filename.
func deleteUnmatchedFiles(folderPath, name string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("unable to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() != name {
			// Check if the directory entry is a file before attempting to delete
			if !entry.IsDir() {
				filePath := filepath.Join(folderPath, entry.Name())
				if err := os.Remove(filePath); err != nil {
					log.WithContext(context.Background()).WithError(err).Errorf("unable to delete file: %s", filePath)
					// Continue with the loop even if there's an error deleting a file
				}
			}
		}
	}
	return nil
}

// NewCascadeAPIHandler returns the swagger OpenAPI implementation.
func NewCascadeAPIHandler(config *Config, filesMap *sync.Map, register *cascaderegister.CascadeRegistrationService, download *download.NftDownloadingService) *CascadeAPIHandler {
	partSizeMB = config.MultiVolumeChunkSize
	chunkSize = config.MultiVolumeChunkSize * 1024 * 1024

	return &CascadeAPIHandler{
		Common:       NewCommon(config),
		register:     register,
		download:     download,
		fileMappings: filesMap,
	}
}
