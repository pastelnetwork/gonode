package files

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
)

// Storage represents a file srorage.
type Storage struct {
	storage.FileStorageInterface

	idCounter int64
	prefix    string
	filesMap  map[string]*File
}

// Run removes all files when the context is canceled.
func (storage *Storage) Run(ctx context.Context) error {
	<-ctx.Done()

	var errs error
	for _, file := range storage.filesMap {
		if err := file.Remove(); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// NewFile returns a new File instance with a unique name.
func (storage *Storage) NewFile() *File {
	id := atomic.AddInt64(&storage.idCounter, 1)
	name := fmt.Sprintf("%s-%d", storage.prefix, id)

	file := NewFile(storage, name)
	storage.filesMap[name] = file

	return file
}

// File returns File by the given name.
func (storage *Storage) File(name string) (*File, error) {
	file, ok := storage.filesMap[name]
	if !ok {
		return nil, errors.New("image not found")
	}
	return file, nil
}

// Update changes the key to identify a *File to a new key
func (storage *Storage) Update(oldname, newname string, file *File) error {
	f, ok := storage.filesMap[oldname]
	if !ok {
		return errors.New("file not found")
	}

	if f != file {
		return errors.New("not the same file")
	}

	delete(storage.filesMap, oldname)
	storage.filesMap[newname] = file
	return nil
}

// NewStorage returns a new Storage instance.
func NewStorage(storage storage.FileStorageInterface) *Storage {
	prefix, _ := random.String(8, random.Base62Chars)

	return &Storage{
		FileStorageInterface: storage,

		prefix:   prefix,
		filesMap: make(map[string]*File),
	}
}
