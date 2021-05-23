package fs

import (
	"os"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
)

const (
	logPrefix = "storage"
)

// FS represents file sysmte storage.
type FS struct {
	dir string
}

// Open implements storage.FileStorage.Open
func (fs *FS) Open(filename string) (storage.File, error) {
	filename = filepath.Join(fs.dir, filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, storage.ErrFileNotFound
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Errorf("failed to open file %q: %w", filename, err)
	}
	return file, nil
}

// Create implements storage.FileStorage.Create
func (fs *FS) Create(filename string) (storage.File, error) {
	filename = filepath.Join(fs.dir, filename)

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return nil, storage.ErrFileExists
	}

	log.WithPrefix(logPrefix).Debugf("Create file %q", filename)

	file, err := os.Create(filename)
	if err != nil {
		return nil, errors.Errorf("failed to create file %q: %w", filename, err)
	}
	return file, nil
}

// Remove implements storage.FileStorage.Remove
func (fs *FS) Remove(filename string) error {
	filename = filepath.Join(fs.dir, filename)

	log.WithPrefix(logPrefix).Debugf("Remove file %q", filename)

	if err := os.Remove(filename); err != nil {
		return errors.Errorf("failed to remove file %q: %w", filename, err)
	}
	return nil
}

// NewFileStorage returns new FS instance. Where `dir` is the path for storing files.
func NewFileStorage(dir string) storage.FileStorage {
	return &FS{
		dir: dir,
	}
}
