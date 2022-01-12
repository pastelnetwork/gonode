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

// Open implements storage.FileStorageInterface.Open
func (fs *FS) Open(filename string) (storage.FileInterface, error) {
	filename = filepath.Join(fs.dir, filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, storage.ErrFileNotFound
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Errorf("open file %q: %w", filename, err)
	}
	return file, nil
}

// Create implements storage.FileStorageInterface.Create
func (fs *FS) Create(filename string) (storage.FileInterface, error) {
	filename = filepath.Join(fs.dir, filename)

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		log.WithPrefix(logPrefix).Debugf("Rewrite file %q", filename)
	} else {
		log.WithPrefix(logPrefix).Debugf("Create file %q", filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, errors.Errorf("create file %q: %w", filename, err)
	}
	return file, nil
}

// Remove implements storage.FileStorageInterface.Remove
func (fs *FS) Remove(filename string) error {
	filename = filepath.Join(fs.dir, filename)

	log.WithPrefix(logPrefix).Debugf("Remove file %q", filename)

	if err := os.Remove(filename); err != nil {
		return errors.Errorf("remove file %q: %w", filename, err)
	}
	return nil
}

// Rename renames oldName to newName.
func (fs *FS) Rename(oldname, newname string) error {
	if oldname == newname {
		return nil
	}

	oldname = filepath.Join(fs.dir, oldname)
	newname = filepath.Join(fs.dir, newname)

	log.WithPrefix(logPrefix).Debugf("Rename file %q to %q", oldname, newname)

	if err := os.Rename(oldname, newname); err != nil {
		return errors.Errorf("rename file %q to %q: %w", oldname, newname, err)
	}
	return nil
}

// NewFileStorage returns new FS instance. Where `dir` is the path for storing files.
func NewFileStorage(dir string) storage.FileStorageInterface {
	return &FS{
		dir: dir,
	}
}
