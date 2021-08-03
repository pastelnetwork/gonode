package test

import (
	"github.com/pastelnetwork/gonode/common/storage/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	readMethod  = "Read"
	nameMethod  = "Name"
	closeMethod = "Close"
)

// File implementing storage.File mock for testing purpose
type File struct {
	*mocks.File
}

// NewMockFile new File instance
func NewMockFile() *File {
	return &File{
		File: &mocks.File{},
	}
}

// ListenOnClose listening Open call and returns file and error from args
func (f *File) ListenOnClose(err error) *File {
	f.On(closeMethod).Return(err)

	return f
}

// ListenOnRead listening Read call and returns file and error from args
func (f *File) ListenOnRead(n int, err error) *File {
	f.On(readMethod, mock.Anything).Return(n, err)

	return f
}

// ListenOnName listening Name call and returns file and error from args
func (f *File) ListenOnName(name string) *File {
	f.On(nameMethod).Return(name)

	return f
}
