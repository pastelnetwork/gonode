package cascaderegister

import (
	"context"
	"database/sql"
	"encoding/base64"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/storage/ticketstore"
	"github.com/pastelnetwork/gonode/common/types"

	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Same with png.
	_ "image/jpeg"
	_ "image/png"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

const (
	logPrefix                   = "cascade"
	defaultImageTTL             = time.Second * 3600 // 1 hour
	maxFileRegistrationAttempts = 3
)

// CascadeRegistrationService represents a service for Cascade Open API
type CascadeRegistrationService struct {
	*task.Worker
	config *Config

	task          *CascadeRegistrationTask
	ImageHandler  *mixins.FilesHandler
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
	restore       *RestoreService

	rqClient rqnode.ClientInterface

	downloadHandler download.NftDownloadingService
	historyDB       queries.LocalStoreInterface
	ticketDB        ticketstore.TicketStorageInterface
}

// Run starts worker.
func (service *CascadeRegistrationService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.ImageHandler.FileStorage.Run(ctx)
	})

	// Run worker service
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})

	// Run restore volume worker
	group.Go(func() error {
		return service.restore.Run(ctx, *service)
	})

	return group.Wait()
}

// Tasks returns all tasks.
func (service *CascadeRegistrationService) Tasks() []*CascadeRegistrationTask {
	var tasks []*CascadeRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*CascadeRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the Cascade OpenAPI by the given id.
func (service *CascadeRegistrationService) GetTask(id string) *CascadeRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*CascadeRegistrationTask)
	}
	return nil
}

// ValidateUser validates user
func (service *CascadeRegistrationService) ValidateUser(ctx context.Context, id string, pass string) bool {
	return common.ValidateUser(ctx, service.pastelHandler.PastelClient, id, pass) &&
		common.IsPastelIDTicketRegistered(ctx, service.pastelHandler.PastelClient, id)
}

// AddTask create ticket request and start a new task with the given payload
func (service *CascadeRegistrationService) AddTask(p *common.AddTaskPayload, regAttemptID int64, filename string, baseFileID string) (string, error) {
	request := FromTaskPayload(p)
	request.RegAttemptID = regAttemptID
	request.FileID = filename
	request.BaseFileID = baseFileID
	request.FileName = filename

	// get image filename from storage based on image_id
	filePath := filepath.Join(service.config.CascadeFilesDir, baseFileID, filename)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	file := service.ImageHandler.FileStorage.NewFile()
	if file == nil {
		return "", errors.Errorf("unable to create new file instance for %s", filename)
	}

	// Write the []byte data to the file
	if _, err := file.Write(fileData); err != nil {
		return "", errors.Errorf("write image data to file: %v", err)
	}

	// Assign the newly created File instance to the request
	request.Image = file

	task := NewCascadeRegisterTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// StoreFile stores file into walletnode file storage
func (service *CascadeRegistrationService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.ImageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// CalculateFee stores file into walletnode file storage
func (service *CascadeRegistrationService) CalculateFee(ctx context.Context, fileID string) (float64, error) {
	fileData, err := service.ImageHandler.GetImgData(fileID)
	if err != nil {
		return 0.0, err
	}

	return service.pastelHandler.GetEstimatedCascadeFee(ctx, utils.GetFileSizeInMB(fileData))
}

type FileMetadata struct {
	TaskID            string
	TotalEstimatedFee float64
	ReqPreBurnAmount  float64
	UploadAssetReq    *cascade.UploadAssetPayload
}

// StoreFileMetadata stores file metadata into the ticket db
func (service *CascadeRegistrationService) StoreFileMetadata(ctx context.Context, dir string, hash string, size int64) (fee float64, preburn []float64, err error) {
	blockCount, err := service.pastelHandler.PastelClient.GetBlockCount(ctx)
	if err != nil {
		return fee, preburn, errors.Errorf("cannot get block count: %w", err)
	}

	files, err := os.ReadDir(dir)

	type fileMetadata struct {
		Name string
		Fee  float64
		Hash []byte
	}
	var metadataList []fileMetadata

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return fee, preburn, err
		}

		fee, err := service.pastelHandler.GetEstimatedCascadeFee(ctx, utils.GetFileSizeInMB(fileData))
		if err != nil {
			return fee, preburn, err
		}

		hash, err := utils.Sha3256hash(fileData)
		if err != nil {
			return fee, preburn, err
		}

		metadataList = append(metadataList, fileMetadata{
			Name: file.Name(),
			Fee:  fee,
			Hash: hash,
		})
	}

	now := time.Now().UTC()
	// Print the metadata list or handle it as needed
	for i, metadata := range metadataList {
		localFee := metadata.Fee + 10.0
		localPreburn := localFee * 0.2

		fee += localFee
		preburn = append(preburn, localPreburn)

		err = service.ticketDB.UpsertFile(types.File{
			FileID:                       metadata.Name,
			UploadTimestamp:              now,
			FileIndex:                    strconv.Itoa(i),
			BaseFileID:                   filepath.Base(dir),
			TaskID:                       "",
			ReqBurnTxnAmount:             localPreburn,
			ReqAmount:                    localFee,
			UUIDKey:                      uuid.NewString(),
			HashOfOriginalBigFile:        hash,
			NameOfOriginalBigFileWithExt: strings.Split(metadata.Name, ".7z")[0],
			SizeOfOriginalBigFile:        float64(size),
			StartBlock:                   blockCount,
		})
		if err != nil {
			return fee, preburn, errors.Errorf("error upsert for file data: %w", err)
		}
	}

	return fee, preburn, nil
}

func (service *CascadeRegistrationService) GetFileByTaskID(taskID string) (*types.File, error) {
	file, err := service.ticketDB.GetFileByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (service *CascadeRegistrationService) GetFilesByBaseFileID(fileID string) (types.Files, error) {
	files, err := service.ticketDB.GetFilesByBaseFileID(fileID)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (service *CascadeRegistrationService) GetActivationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.ActivationAttempt, error) {
	activationAttempts, err := service.ticketDB.GetActivationAttemptsByFileIDAndBaseFileID(fileID, baseFileID)
	if err != nil {
		return nil, err
	}
	return activationAttempts, nil
}

func (service *CascadeRegistrationService) GetRegistrationAttemptsByFileIDAndBaseFileID(fileID, baseFileID string) ([]*types.RegistrationAttempt, error) {
	registrationAttempts, err := service.ticketDB.GetRegistrationAttemptsByFileIDAndBaseFileID(fileID, baseFileID)
	if err != nil {
		return nil, err
	}
	return registrationAttempts, nil
}

func (service *CascadeRegistrationService) GetRegistrationAttemptsByID(attemptID int) (*types.RegistrationAttempt, error) {
	registrationAttempt, err := service.ticketDB.GetRegistrationAttemptByID(attemptID)
	if err != nil {
		return nil, err
	}
	return registrationAttempt, nil
}

func (service *CascadeRegistrationService) GetActivationAttemptByID(attemptID int) (*types.ActivationAttempt, error) {
	actAttempt, err := service.ticketDB.GetActivationAttemptByID(attemptID)
	if err != nil {
		return nil, err
	}
	return actAttempt, nil
}

func (service *CascadeRegistrationService) InsertRegistrationAttempts(regAttempt types.RegistrationAttempt) (int64, error) {
	id, err := service.ticketDB.InsertRegistrationAttempt(regAttempt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *CascadeRegistrationService) UpdateRegistrationAttempts(regAttempt types.RegistrationAttempt) (int64, error) {
	id, err := service.ticketDB.UpdateRegistrationAttempt(regAttempt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *CascadeRegistrationService) InsertActivationAttempt(actAttempt types.ActivationAttempt) (int64, error) {
	id, err := service.ticketDB.InsertActivationAttempt(actAttempt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *CascadeRegistrationService) UpdateActivationAttempts(actAttempt types.ActivationAttempt) (int64, error) {
	id, err := service.ticketDB.UpdateActivationAttempt(actAttempt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *CascadeRegistrationService) UpsertFile(file types.File) error {
	err := service.ticketDB.UpsertFile(file)
	if err != nil {
		return err
	}
	return nil
}

func (service *CascadeRegistrationService) InsertMultiVolCascadeTicketTxIDMap(m types.MultiVolCascadeTicketTxIDMap) error {
	err := service.ticketDB.InsertMultiVolCascadeTicketTxIDMap(m)
	if err != nil {
		return err
	}
	return nil
}

func (service *CascadeRegistrationService) GetMultiVolTicketTxIDMap(baseFileID string) (*types.MultiVolCascadeTicketTxIDMap, error) {
	txIDMapping, err := service.ticketDB.GetMultiVolCascadeTicketTxIDMap(baseFileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return txIDMapping, nil
}

func (service *CascadeRegistrationService) HandleTaskRegistrationErrorAttempts(ctx context.Context, taskID, regTxid, actTxid string, regAttemptID int64, actAttemptID int64, taskError error) error {
	doneBlock, err := service.pastelHandler.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("task_id", taskID).WithError(err).Error("error retrieving block count")
		return err
	}

	switch {
	case regTxid == "":
		ra, err := service.GetRegistrationAttemptsByID(int(regAttemptID))
		if err != nil {
			log.WithContext(ctx).WithField("task_id", taskID).WithError(err).Error("error retrieving file reg attempt")
			return err
		}

		ra.FinishedAt = time.Now().UTC()
		ra.IsSuccessful = false
		ra.ErrorMessage = taskError.Error()
		_, err = service.UpdateRegistrationAttempts(*ra)
		if err != nil {
			log.WithContext(ctx).WithField("task_id", taskID).WithError(err).Error("error updating file reg attempt")
			return err
		}

		return nil
	case actTxid == "":
		err := service.HandleTaskActError(ctx, taskID, regTxid, int(actAttemptID), doneBlock, taskError)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (service *CascadeRegistrationService) HandleTaskActError(ctx context.Context, taskID string, regTxid string, actAttemptID int, doneBlock int32, err error) error {
	file, err := service.GetFileByTaskID(taskID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving file")
		return nil
	}

	file.DoneBlock = int(doneBlock)
	file.RegTxid = regTxid
	file.IsConcluded = false
	err = service.ticketDB.UpsertFile(*file)
	if err != nil {
		log.Errorf("Error in file upsert: %v", err.Error())
		return nil
	}

	actAttempt, err := service.GetActivationAttemptByID(int(actAttemptID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("Error retrieving file act attempt: %v", err.Error())
		return err
	}

	if actAttempt == nil {
		return nil
	}

	actAttempt.IsSuccessful = false
	actAttempt.ActivationAttemptAt = time.Now().UTC()
	actAttempt.ErrorMessage = err.Error()
	_, err = service.UpdateActivationAttempts(*actAttempt)
	if err != nil {
		log.Errorf("Error in activation attempts upsert: %v", err.Error())
		return err
	}

	return err
}

func (service *CascadeRegistrationService) HandleTaskRegSuccess(ctx context.Context, taskID string, regTxid string, actTxid string, regAttemptID int, actAttemptID int) error {
	doneBlock, err := service.pastelHandler.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving block-count")
		return err
	}

	file, err := service.GetFileByTaskID(taskID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving file")
		return nil
	}

	file.DoneBlock = int(doneBlock)
	file.RegTxid = regTxid
	file.ActivationTxid = actTxid
	file.IsConcluded = true
	err = service.ticketDB.UpsertFile(*file)
	if err != nil {
		log.Errorf("Error in file upsert: %v", err.Error())
		return nil
	}

	// Upsert registration attempts

	ra, err := service.GetRegistrationAttemptsByID(regAttemptID)
	if err != nil {
		log.Errorf("Error retrieving file reg attempt: %v", err.Error())
		return nil
	}

	ra.IsConfirmed = true
	_, err = service.UpdateRegistrationAttempts(*ra)
	if err != nil {
		log.Errorf("Error in registration attempts upsert: %v", err.Error())
		return nil
	}

	actAttempt, err := service.GetActivationAttemptByID(actAttemptID)
	if err != nil {
		log.Errorf("Error retrieving file act attempt: %v", err.Error())
		return nil
	}

	actAttempt.IsConfirmed = true
	_, err = service.UpdateActivationAttempts(*actAttempt)
	if err != nil {
		log.Errorf("Error in activation attempts upsert: %v", err.Error())
		return nil
	}

	return nil
}

func (service *CascadeRegistrationService) ProcessFile(ctx context.Context, file types.File, p *common.AddTaskPayload) (taskID string, err error) {
	log.WithContext(ctx).WithField("file_id", file.FileID).WithField("base_file_id", file.BaseFileID).Infof("Processing started for file")

	switch {
	case file.IsConcluded:
		return "", cascade.MakeInternalServerError(errors.New("ticket has already been registered & activated"))

	case file.RegTxid == "":
		baseFileRegistrationAttempts, err := service.GetRegistrationAttemptsByFileIDAndBaseFileID(file.FileID, file.BaseFileID)
		if err != nil {
			return "", cascade.MakeInternalServerError(err)
		}

		switch {
		case len(baseFileRegistrationAttempts) > maxFileRegistrationAttempts:
			return "", cascade.MakeInternalServerError(errors.New("ticket registration attempts have been exceeded"))
		default:
			regAttemptID, err := service.InsertRegistrationAttempts(types.RegistrationAttempt{
				BaseFileID:   file.BaseFileID,
				FileID:       file.FileID,
				RegStartedAt: time.Now().UTC(),
			})
			if err != nil {
				log.WithContext(ctx).WithField("file_id", file.FileID).WithError(err).Error("error inserting registration attempt")
				return "", err
			}

			taskID, err = service.AddTask(p, regAttemptID, file.FileID, file.BaseFileID)
			if err != nil {
				log.WithError(err).Error("unable to add task")
				return "", cascade.MakeInternalServerError(err)
			}

			file.TaskID = taskID
			file.BurnTxnID = *p.BurnTxid
			file.PastelID = p.AppPastelID
			file.Passphrase = base64.StdEncoding.EncodeToString([]byte(p.Key))
			err = service.UpsertFile(file)
			if err != nil {
				log.WithField("task_id", taskID).WithField("file_id", p.FileID).
					WithField("burn_txid", p.BurnTxid).
					WithField("file_name", file.FileID).
					Errorf("Error in file upsert: %v", err.Error())
				return "", cascade.MakeInternalServerError(errors.New("Error in file upsert"))
			}
		}

		return taskID, nil
	default:
		// Activation code
	}

	return taskID, nil
}

func (service *CascadeRegistrationService) SortBurnTxIDs(ctx context.Context, burnTxIDs []string) ([]string, error) {
	type Txn struct {
		ID     string
		Amount float64
	}
	var txns []Txn

	for _, burnTxID := range burnTxIDs {
		txnDetails, err := service.pastelHandler.PastelClient.GetTransaction(ctx, burnTxID)
		if err != nil {
			return []string{}, err
		}
		txns = append(txns, Txn{ID: burnTxID, Amount: txnDetails.Amount})
	}

	sort.SliceStable(txns, func(i, j int) bool {
		return math.Abs(txns[i].Amount) > math.Abs(txns[j].Amount)
	})

	sortedBurnTxIDs := make([]string, len(txns))
	for i, txn := range txns {
		sortedBurnTxIDs[i] = txn.ID
	}

	return sortedBurnTxIDs, nil
}

func (service *CascadeRegistrationService) SortFilesWithHigherAmounts(files types.Files) types.Files {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].ReqBurnTxnAmount > files[j].ReqBurnTxnAmount
	})

	return files
}

func (service *CascadeRegistrationService) GetConcludedVolumesByBaseFileID(BaseFileID string) (types.Files, error) {
	concludedVolumes, err := service.ticketDB.GetFilesByBaseFileIDAndConcludedCheck(BaseFileID, true)
	if err != nil {
		return nil, err
	}

	return concludedVolumes, nil
}

func (service *CascadeRegistrationService) GetUnConcludedVolumesByBaseFileID(BaseFileID string) (types.Files, error) {

	UnConcludedVolumes, err := service.ticketDB.GetFilesByBaseFileIDAndConcludedCheck(BaseFileID, false)
	if err != nil {
		return nil, err
	}

	return UnConcludedVolumes, nil
}

func (service *CascadeRegistrationService) GetBurnTxIdByAmount(ctx context.Context, amount int64) (string, error) {
	txID, err := service.pastelHandler.PastelClient.SendToAddress(ctx, service.pastelHandler.GetBurnAddress(), amount)
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (service *CascadeRegistrationService) ActivateActionTicketAndRegisterVolumeTicket(ctx context.Context, aar pastel.ActivateActionRequest, baseFile types.File, actAttemptId int64) error {
	regTicket, err := service.pastelHandler.PastelClient.ActionRegTicket(ctx, aar.RegTxID)
	if err != nil {
		return err
	}
	aar.BlockNum = regTicket.Height

	actTxId, err := service.pastelHandler.PastelClient.ActivateActionTicket(ctx, aar)
	if err != nil {
		return err
	}

	baseFile.ActivationTxid = actTxId
	err = service.UpsertFile(baseFile)
	if err != nil {
		return err
	}

	actAttempt, err := service.GetActivationAttemptByID(int(actAttemptId))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("Error retrieving file act attempt: %v", err.Error())
		return err
	}

	actAttempt.IsSuccessful = true
	actAttempt.ActivationAttemptAt = time.Now().UTC()
	_, err = service.UpdateActivationAttempts(*actAttempt)
	if err != nil {
		log.Errorf("Error in activation attempts upsert: %v", err.Error())
		return err
	}

	return nil
}

func (service *CascadeRegistrationService) RegisterVolumeTicket(ctx context.Context, baseFileID string) (string, error) {
	log.WithContext(ctx).WithField("base_file_id", baseFileID).Info("Register volumes ticket started")

	multiVolCascadeTicketMap, err := service.GetMultiVolTicketTxIDMap(baseFileID)
	if err != nil {
		return "", errors.New("errors retrieving multi-vol-ticket txid map")
	}

	if multiVolCascadeTicketMap != nil {
		log.WithContext(ctx).WithField("multi_vol_cascade_ticket_txid", multiVolCascadeTicketMap.MultiVolCascadeTicketTxid).
			Info("multi-vol-ticket has already been created")
		return "", nil
	}

	relatedFiles, err := service.GetFilesByBaseFileID(baseFileID)
	if err != nil {
		return "", err
	}

	if len(relatedFiles) <= 1 {
		return "", errors.New("related files must be greater than 1")
	}

	concludedCount := 0
	requiredConcludedCount := len(relatedFiles)
	for _, rf := range relatedFiles {
		if rf.IsConcluded {
			concludedCount += 1
		}
	}

	if concludedCount < requiredConcludedCount || concludedCount > requiredConcludedCount {
		log.WithContext(ctx).WithField("base_file_id", baseFileID).
			WithField("related_files", relatedFiles.Names()).WithField("no_of_volumes_registered", concludedCount).
			WithField("total_volumes_count", requiredConcludedCount).Info("all volumes are not registered or activated yet")
		return "", nil
	}

	volumes := make(map[int]string)
	var sizeOfOriginalFileMB int
	var nameOfOriginalFile string
	var sha3256HashOfOriginalFile string

	nameOfOriginalFile = relatedFiles[0].NameOfOriginalBigFileWithExt
	sizeOfOriginalFileMB = int(relatedFiles[0].SizeOfOriginalBigFile)
	sha3256HashOfOriginalFile = relatedFiles[0].HashOfOriginalBigFile

	for _, relatedFile := range relatedFiles {
		fileIndexInt, err := strconv.Atoi(relatedFile.FileIndex)
		if err != nil {
			return "", errors.Errorf("error converting file index to int from string")
		}

		volumes[fileIndexInt] = relatedFile.RegTxid
	}

	txID, err := service.pastelHandler.PastelClient.RegisterCascadeMultiVolumeTicket(ctx, pastel.CascadeMultiVolumeTicket{
		NameOfOriginalFile:        nameOfOriginalFile,
		SizeOfOriginalFileMB:      sizeOfOriginalFileMB,
		SHA3256HashOfOriginalFile: sha3256HashOfOriginalFile,
		Volumes:                   volumes,
	})
	if err != nil {
		log.WithContext(ctx).WithField("base_file_id", baseFileID).
			WithField("related_files", relatedFiles.Names()).WithField("no_of_volumes_registered", concludedCount).
			WithField("total_volumes_count", requiredConcludedCount).Error("Error in pastel register volume cascade ticket")
		return txID, err
	}

	err = service.InsertMultiVolCascadeTicketTxIDMap(types.MultiVolCascadeTicketTxIDMap{
		MultiVolCascadeTicketTxid: txID,
		BaseFileID:                baseFileID,
	})
	if err != nil {
		log.WithContext(ctx).WithField("base_file_id", baseFileID).
			WithField("related_files", relatedFiles.Names()).WithField("no_of_volumes_registered", concludedCount).
			WithField("total_volumes_count", requiredConcludedCount).Error("Error in storing txId in db after pastel volume registration")
		return txID, err
	}
	log.WithContext(ctx).WithField("base_file_id", baseFileID).WithField("txID", txID).Info("cascade multi-volume registration has been completed")

	return txID, nil
}

func (service *CascadeRegistrationService) CheckBurnTxIDTicketDuplication(ctx context.Context, burnTxID string) error {
	err := service.pastelHandler.CheckBurnTxID(ctx, burnTxID)
	if err != nil {
		return err
	}

	return nil
}

func (service *CascadeRegistrationService) ValidBurnTxnAmount(ctx context.Context, burnTxID string, estimatedFee float64) error {
	txnDetail, err := service.pastelHandler.PastelClient.GetTransaction(ctx, burnTxID)
	if err != nil {
		return err
	}

	burnTxnAmount := math.Round(math.Abs(txnDetail.Amount)*100) / 100
	estimatedFeeAmount := math.Round(estimatedFee*100) / 100

	if burnTxnAmount >= estimatedFeeAmount {
		return nil
	}

	return errors.New("the amounts are not equal")
}

func (service *CascadeRegistrationService) ValidateBurnTxn(ctx context.Context, burnTxID string, estimatedFee float64) error {
	return service.pastelHandler.ValidateBurnTxID(ctx, burnTxID, estimatedFee)
}

func (service *CascadeRegistrationService) GetUnCompletedFiles() ([]*types.File, error) {
	uCFiles, err := service.ticketDB.GetUnCompletedFiles()
	if err != nil {
		return []*types.File{}, err
	}

	return uCFiles, nil
}

func (service *CascadeRegistrationService) RestoreFile(ctx context.Context, p *cascade.RestorePayload) (res *cascade.RestoreFile, err error) {
	volumes, err := service.GetFilesByBaseFileID(p.BaseFileID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	if len(volumes) == 0 {
		return nil, cascade.MakeInternalServerError(errors.New("no volumes associated with the base-file-id to recover"))
	}

	log.WithContext(ctx).WithField("total_volumes", len(volumes)).
		Info("total volumes retrieved from the base file-id")

	unConcludedVolumes, err := volumes.GetUnconcludedFiles()
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	// check if volumes already concluded
	if len(unConcludedVolumes) == 0 {
		log.WithContext(ctx).
			Info("All volumes are already concluded")
		return &cascade.RestoreFile{
			TotalVolumes:                   len(volumes),
			RegisteredVolumes:              len(volumes),
			VolumesWithPendingRegistration: 0,
			VolumesRegistrationInProgress:  0,
			ActivatedVolumes:               len(volumes),
			VolumesActivatedInRecoveryFlow: 0,
		}, nil

	}

	var (
		registeredVolumes              int
		volumesWithInProgressRegCount  int
		volumesWithPendingRegistration int
		volumesActivatedInRecoveryFlow int
		activatedVolumes               int
	)
	logger := log.WithContext(ctx).WithField("base_file_id", p.BaseFileID)

	runningTasks := service.Worker.Tasks()
	runningTaskIDs := make(map[string]struct{}, len(runningTasks))
	for _, rt := range runningTasks {
		runningTaskIDs[rt.ID()] = struct{}{}
	}

	for _, v := range volumes {
		if _, exists := runningTaskIDs[v.TaskID]; exists {
			log.WithContext(ctx).WithField("task_id", v.TaskID).WithField("base_file_id", v.BaseFileID).
				Info("current task is already in-progress can't execute the recovery-flow")
			continue
		}

		if v.RegTxid == "" {
			volumesWithPendingRegistration++
			logger.WithField("volume_name", v.FileID).Info("find a volume with no registration, trying again...")

			var burnTxId string
			if v.BurnTxnID != "" && service.IsBurnTxIDValidForRecovery(ctx, v.BurnTxnID, v.ReqAmount-10) {
				log.WithContext(ctx).WithField("burn_txid", v.BurnTxnID).Info("existing burn-txid is valid")
				burnTxId = v.BurnTxnID
			} else {
				log.WithContext(ctx).WithField("burn_txid", v.BurnTxnID).Info("existing burn-txid is not valid, burning the new txid")

				burnTxId, err = service.GetBurnTxIdByAmount(ctx, int64(v.ReqBurnTxnAmount))
				if err != nil {
					log.WithContext(ctx).WithField("amount", int64(v.ReqBurnTxnAmount)).WithError(err).Error("error getting burn TxId for amount")
					return nil, cascade.MakeInternalServerError(err)
				}

				logger.WithField("volume_name", v.FileID).Info("estimated fee has been burned, sending for registration")
			}

			addTaskPayload := &common.AddTaskPayload{
				FileID:                 v.FileID,
				BurnTxid:               &burnTxId,
				AppPastelID:            p.AppPastelID,
				MakePubliclyAccessible: p.MakePubliclyAccessible,
				Key:                    p.Key,
			}

			if p.SpendableAddress != nil {
				addTaskPayload.SpendableAddress = p.SpendableAddress
			}

			_, err = service.ProcessFile(ctx, *v, addTaskPayload)
			if err != nil {
				log.WithContext(ctx).WithField("file_id", v.FileID).WithError(err).Error("error processing un-registered volume")
				continue
			}
			logger.WithField("volume_name", v.FileID).Info("task added for registration recovery")

			volumesWithInProgressRegCount += 1
		} else if v.ActivationTxid == "" {
			logger.WithField("volume_name", v.FileID).Info("find a volume with no activation, trying again...")

			// activation code
			actAttemptId, err := service.InsertActivationAttempt(types.ActivationAttempt{
				FileID:              v.FileID,
				ActivationAttemptAt: time.Now().UTC(),
			})
			if err != nil {
				log.WithContext(ctx).WithField("file_id", v.FileID).WithError(err).Error("error inserting activation attempt")
				continue
			}

			activateActReq := pastel.ActivateActionRequest{
				RegTxID:    v.RegTxid,
				Fee:        int64(v.ReqAmount) - 10,
				PastelID:   p.AppPastelID,
				Passphrase: p.Key,
			}

			err = service.ActivateActionTicketAndRegisterVolumeTicket(ctx, activateActReq, *v, actAttemptId)
			if err != nil {
				log.WithContext(ctx).WithField("file_id", v.FileID).WithError(err).Error("error activating or registering un-concluded volume")
				continue
			}
			logger.WithField("volume_name", v.FileID).Info("request has been sent for activation")

			volumesActivatedInRecoveryFlow += 1
		} else {
			if v.RegTxid != "" {
				registeredVolumes += 1
			}
			if v.ActivationTxid != "" {
				activatedVolumes += 1
			}
		}
	}

	// only set base file txId return by pastel register all else remains nil
	_, err = service.RegisterVolumeTicket(ctx, p.BaseFileID)
	if err != nil {
		return nil, err
	}

	return &cascade.RestoreFile{
		TotalVolumes:                   len(volumes),
		RegisteredVolumes:              registeredVolumes,
		VolumesWithPendingRegistration: volumesWithPendingRegistration,
		VolumesRegistrationInProgress:  volumesWithInProgressRegCount,
		ActivatedVolumes:               activatedVolumes,
		VolumesActivatedInRecoveryFlow: volumesActivatedInRecoveryFlow,
	}, nil
}

func (service *CascadeRegistrationService) IsBurnTxIDValidForRecovery(ctx context.Context, burnTxID string, estimatedFee float64) bool {
	if err := service.CheckBurnTxIDTicketDuplication(ctx, burnTxID); err != nil {
		return false
	}

	err := service.ValidateBurnTxn(ctx, burnTxID, estimatedFee)
	if err != nil {
		return false
	}

	return true
}

// NewService returns a new Service instance
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface, db storage.KeyValue,
	downloadService download.NftDownloadingService,
	historyDB queries.LocalStoreInterface,
	ticketDB ticketstore.TicketStorageInterface,
) *CascadeRegistrationService {
	return &CascadeRegistrationService{
		Worker:          task.NewWorker(),
		config:          config,
		nodeClient:      nodeClient,
		ImageHandler:    mixins.NewFilesHandler(fileStorage, db, defaultImageTTL),
		pastelHandler:   mixins.NewPastelHandler(pastelClient),
		rqClient:        rqgrpc.NewClient(),
		restore:         NewRestoreService(),
		downloadHandler: downloadService,
		historyDB:       historyDB,
		ticketDB:        ticketDB,
	}
}
