//go:generate mockery --name=FileStorage
//go:generate mockery --name=File

package storage

import (
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	// ErrFileNotFound is returned when file isn't found.
	ErrFileNotFound = errors.New("file not found")
	// ErrFileExists is returned when file already exists.
	ErrFileExists = errors.New("file exists")
)

// FileStorageInterface represents a file storage.
type FileStorageInterface interface {
	// Open opens a file and returns file descriptor.
	// If name is not found, ErrFileNotFound is returned.
	Open(name string) (file FileInterface, err error)

	// Create creates a new file with the given name and returns file descriptor.
	Create(name string) (file FileInterface, err error)

	// Remove removes a file by the given name.
	Remove(name string) error

	// Rename renames oldname to newname.
	Rename(oldname, newname string) error
}

// File represents a file.
type FileInterface interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	Name() string
}
