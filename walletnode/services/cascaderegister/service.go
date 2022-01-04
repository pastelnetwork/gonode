package cascaderegister

import (
	"bytes"
	"context"
	"image"
	"time"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Same with png.
	_ "image/jpeg"
	_ "image/png"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix       = "sense"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
	db           storage.KeyValue

	imageTTL        time.Duration
	registrationFee int64
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Storage.Run(ctx)
	})

	// Run worker service
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*Task {
	var tasks []*Task
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*Task))
	}
	return tasks
}

// Task returns the task of the registration artwork by the given id.
func (service *Service) Task(id string) *Task {
	if t := service.Worker.Task(id); t != nil {
		return t.(*Task)
	}
	return nil
}

// StoreImage stores the image in the db and storage and return image_id
func (service *Service) StoreImage(ctx context.Context, p *cascade.UploadImagePayload) (res *cascade.Image, err error) {
	// Detect image format
	imgReader := bytes.NewReader(p.Bytes)
	_, format, err := image.Decode(imgReader)
	if err != nil {
		return nil, errors.Errorf("decode image: %w", err)
	}
	log.WithContext(ctx).Debugf("image format: %s", format)

	// Store image content to storage
	image := service.Storage.NewFile()
	if err := image.SetFormatFromExtension(format); err != nil {
		return nil, errors.Errorf("set image format: %w", err)
	}

	log.WithContext(ctx).Debugf("writing image: %s to storage", image.Name())

	fl, err := image.Create()
	if err != nil {
		return nil, errors.Errorf("create image file: %w", err)
	}
	defer fl.Close()

	if _, err := fl.Write(p.Bytes); err != nil {
		return nil, errors.Errorf("write image file: %w", err)
	}

	// Generate unique image_id and write to db
	id, _ := random.String(8, random.Base62Chars)
	if err := service.db.Set(id, []byte(image.Name())); err != nil {
		return nil, errors.Errorf("store image in db: %w", err)
	}

	// Set storage expire time
	image.RemoveAfter(service.imageTTL)

	res = &cascade.Image{
		ImageID:   id,
		ExpiresIn: time.Now().Add(service.imageTTL).Format(time.RFC3339),
	}

	return res, nil
}

// GetImgData returns the image data from the storage
func (service *Service) GetImgData(imageID string) ([]byte, error) {
	// get image filename from storage based on image_id
	filename, err := service.db.Get(imageID)
	if err != nil {
		return nil, errors.Errorf("get image filename from db: %w", err)
	}

	// get image data from storage
	file, err := service.Storage.File(string(filename))
	if err != nil {
		return nil, errors.Errorf("get image fd: %v", err)
	}

	imgData, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("read image data: %w", err)
	}

	return imgData, nil
}

// AddTask create ticket request and start a new task with the given payload
func (service *Service) AddTask(p *cascade.StartProcessingPayload) (string, error) {
	ticket := &Request{
		BurnTxID:              p.BurnTxid,
		AppPastelID:           p.AppPastelID,
		AppPastelIDPassphrase: p.AppPastelidPassphrase,
	}

	// get image filename from storage based on image_id
	filename, err := service.db.Get(p.ImageID)
	if err != nil {
		return "", errors.Errorf("get image filename from storage: %w", err)
	}

	// get image data from storage
	file, err := service.Storage.File(string(filename))
	if err != nil {
		return "", errors.Errorf("get image data: %v", err)
	}
	ticket.Image = file

	task := NewTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// VerifyImageSignature verifies the signature of the image
func (service *Service) VerifyImageSignature(ctx context.Context, imgData []byte, signature string, pastelID string) error {
	ok, err := service.pastelClient.Verify(ctx, imgData, signature, pastelID, pastel.SignAlgorithmED448)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("signature verification failed")
	}

	return nil
}

// GetEstimatedFee returns the estimated sense fee for the given image
func (service *Service) GetEstimatedFee(ctx context.Context, ImgSizeInMb int64) (float64, error) {
	actionFees, err := service.pastelClient.GetActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return 0, err
	}

	service.registrationFee = (int64)(actionFees.SenseFee)

	return actionFees.SenseFee, nil
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	fileStorage storage.FileStorage,
	pastelClient pastel.Client,
	nodeClient node.Client,
	db storage.KeyValue,
) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
		imageTTL:     defaultImageTTL,
	}
}
