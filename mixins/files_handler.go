package mixins

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
)

// FilesHandler handles files
type FilesHandler struct {
	FileStorage *files.Storage
	FileDb      storage.KeyValue
	fileTTL     time.Duration
}

func NewFilesHandler(
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
	fileTTL time.Duration) *FilesHandler {
	return &FilesHandler{
		FileStorage: files.NewStorage(fileStorage),
		FileDb:      db,
		fileTTL:     fileTTL,
	}
}

// StoreImage stores the image in the FileDb and storage and return image_id
func (st *FilesHandler) StoreFileNameIntoStorage(_ context.Context, fileName *string) (string, string, error) {

	id, err := random.String(8, random.Base62Chars)
	if err != nil {
		return "", "0", err
	}
	if err := st.FileDb.Set(id, []byte(*fileName)); err != nil {
		return "", "0", err
	}

	file, err := st.FileStorage.File(*fileName)
	if err != nil {
		return "", "0", err
	}
	file.RemoveAfter(st.fileTTL)

	return id, time.Now().Add(st.fileTTL).Format(time.RFC3339), nil
}

// GetImgData returns the image data from the storage
func (st *FilesHandler) GetImgData(imageID string) ([]byte, error) {
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
