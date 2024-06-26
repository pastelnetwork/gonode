package cascaderegister

import (
	"context"
	"database/sql"
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
	maxFileActAttempts          = 3
)

// CascadeRegistrationService represents a service for Cascade Open API
type CascadeRegistrationService struct {
	*task.Worker
	config *Config

	task          *CascadeRegistrationTask
	ImageHandler  *mixins.FilesHandler
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface

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
func (service *CascadeRegistrationService) AddTask(p *cascade.StartProcessingPayload, regAttemptID int64, filename string, baseFileID string) (string, error) {
	request := FromStartProcessingPayload(p)
	request.RegAttemptID = regAttemptID
	request.FileID = filename
	request.BaseFileID = baseFileID
	request.FileName = filename

	// get image filename from storage based on image_id
	filePath := filepath.Join(service.config.CascadeFilesDir, p.FileID, filename)
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

func (service *CascadeRegistrationService) GetFile(fileID string) (*types.File, error) {
	file, err := service.ticketDB.GetFileByID(fileID)
	if err != nil {
		return nil, err
	}
	return file, nil
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

func (service *CascadeRegistrationService) GetActivationAttemptsByFileID(fileID string) ([]*types.ActivationAttempt, error) {
	activationAttempts, err := service.ticketDB.GetActivationAttemptsByFileID(fileID)
	if err != nil {
		return nil, err
	}
	return activationAttempts, nil
}

func (service *CascadeRegistrationService) GetRegistrationAttemptsByFileID(fileID string) ([]*types.RegistrationAttempt, error) {
	registrationAttempts, err := service.ticketDB.GetRegistrationAttemptsByFileID(fileID)
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

	ra.FinishedAt = time.Now().UTC()
	ra.IsSuccessful = true
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

	actAttempt.IsSuccessful = true
	_, err = service.UpdateActivationAttempts(*actAttempt)
	if err != nil {
		log.Errorf("Error in activation attempts upsert: %v", err.Error())
		return nil
	}

	return nil
}

func (service *CascadeRegistrationService) ProcessFile(ctx context.Context, file types.File, p *cascade.StartProcessingPayload) (taskID string, err error) {
	log.WithContext(ctx).WithField("file_id", file.FileID).WithField("base_file_id", file.BaseFileID).Infof("Processing started for file")

	switch {
	case file.IsConcluded:
		return "", cascade.MakeInternalServerError(errors.New("ticket has already been registered & activated"))

	case file.RegTxid == "":
		baseFileRegistrationAttempts, err := service.GetRegistrationAttemptsByFileID(file.FileID)
		if err != nil {
			return "", cascade.MakeInternalServerError(err)
		}

		switch {
		case len(baseFileRegistrationAttempts) > maxFileRegistrationAttempts:
			return "", cascade.MakeInternalServerError(errors.New("ticket registration attempts have been exceeded"))
		default:
			regAttemptID, err := service.InsertRegistrationAttempts(types.RegistrationAttempt{
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
		downloadHandler: downloadService,
		historyDB:       historyDB,
		ticketDB:        ticketDB,
	}
}
