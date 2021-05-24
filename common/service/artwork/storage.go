package artwork

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
)

// Storage represents a file srorage.
type Storage struct {
	storage.FileStorage

	idCounter int64
	prefix    string
	files     map[string]*File
}

// Run removes all files when the context is canceled.
func (storage *Storage) Run(ctx context.Context) error {
	<-ctx.Done()

	var errs error
	for _, file := range storage.files {
		if err := file.Remove(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

// NewFile returns a new File instance with a unique name.
func (storage *Storage) NewFile() *File {
	id := atomic.AddInt64(&storage.idCounter, 1)
	name := fmt.Sprintf("%s-%d", storage.prefix, id)

	file := NewFile(storage, name)
	storage.files[name] = file

	return file
}

// File returns File by the given name.
func (storage *Storage) File(name string) (*File, error) {
	file, ok := storage.files[name]
	if !ok {
		return nil, errors.New("image not found")
	}
	return file, nil
}

// NewStorage returns a new Storage instance.
func NewStorage(storage storage.FileStorage) *Storage {
	prefix, _ := random.String(8, random.Base62Chars)

	return &Storage{
		FileStorage: storage,

		prefix: prefix,
		files:  make(map[string]*File),
	}
}
