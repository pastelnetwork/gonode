package test

import (
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/mocks"
	"github.com/stretchr/testify/mock"
)

// FileStorage implementing storage.FileStorage mock for testing purpose
type FileStorage struct {
	FileStorageMock *mocks.FileStorage
}

// NewMockFileStorage new FileStorage instance
func NewMockFileStorage() *FileStorage {
	return &FileStorage{
		FileStorageMock: &mocks.FileStorage{},
	}
}

// ListenOnOpen listening Open call and returns file and error from args
func (f *FileStorage) ListenOnOpen(file storage.File, err error) *FileStorage {
	f.FileStorageMock.On("Open", mock.AnythingOfType("string")).Return(file, err)
	return f
}

// ListenOnCreate listening Create call and returns file and error from args
func (f *FileStorage) ListenOnCreate(file storage.File, err error) *FileStorage {
	f.FileStorageMock.On("Create", mock.AnythingOfType("string")).Return(file, err)
	return f
}

func (f *FileStorage) ListenOnRemove(err error) *FileStorage {
	f.FileStorageMock.On("Remove", mock.AnythingOfType("string")).Return(err)
	return f
}
