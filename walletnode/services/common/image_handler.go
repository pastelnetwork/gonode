package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"time"
)

type ImageHandler struct {
	FileStorage *files.Storage
	FileDb      storage.KeyValue
	imageTTL    time.Duration
}

func NewImageHandler(
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
	imageTTL time.Duration) *ImageHandler {
	return &ImageHandler{
		FileStorage: files.NewStorage(fileStorage),
		FileDb:      db,
		imageTTL:    imageTTL,
	}
}

// StoreImage stores the image in the FileDb and storage and return image_id
func (st *ImageHandler) StoreFileNameIntoStorage(ctx context.Context, fileName *string) (string, string, error) {

	id, err := random.String(8, random.Base62Chars)
	if err := st.FileDb.Set(id, []byte(*fileName)); err != nil {
		return "", "0", err
	}

	file, err := st.FileStorage.File(*fileName)
	if err != nil {
		return "", "0", err
	}
	file.RemoveAfter(st.imageTTL)

	return id, time.Now().Add(st.imageTTL).Format(time.RFC3339), nil
}

// GetImgData returns the image data from the storage
func (st *ImageHandler) GetImgData(imageID string) ([]byte, error) {
	// get image filename from storage based on image_id
	filename, err := st.FileDb.Get(imageID)
	if err != nil {
		return nil, errors.Errorf("get image filename from FileDb: %w", err)
	}

	// get image data from storage
	file, err := st.FileStorage.File(string(filename))
	if err != nil {
		return nil, errors.Errorf("get image fd: %v", err)
	}

	imgData, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("read image data: %w", err)
	}

	return imgData, nil
}
