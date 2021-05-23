package image

import (
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/storage"
)

// File represents a file.
type File struct {
	fmt.Stringer
	sync.Mutex

	storage *Storage

	isCreated bool
	name      string
	ext       string
}

// Name returns filename.
func (file *File) Name() string {
	return file.name
}

func (file *File) String() string {
	return file.name
}

// Ext returns file extension.
func (file *File) Ext() string {
	return file.ext
}

// SetExt set file extension.
func (file *File) SetExt(ext string) {
	file.ext = ext
}

// Create creates a file and returns file descriptor.
func (file *File) Create() (storage.File, error) {
	fl, err := file.storage.Create(file.name)
	if err != nil {
		return nil, err
	}

	file.isCreated = true
	return fl, nil
}

// Open opens a file and returns file descriptor.
// If file is not found, storage.ErrFileNotFound is returned.
func (file *File) Open() (storage.File, error) {
	return file.storage.Open(file.name)
}

// Remove removes the file.
func (file *File) Remove() error {
	file.Lock()
	defer file.Unlock()

	delete(file.storage.files, file.name)

	if !file.isCreated {
		return nil
	}
	file.isCreated = false

	return file.storage.Remove(file.name)
}

// RemoveAfter removes the file after the specified duration.
func (file *File) RemoveAfter(d time.Duration) {
	go func() {
		time.AfterFunc(d, func() { file.Remove() })
	}()
}

// NewFile returns a new File instance.
func NewFile(storage *Storage, name string) *File {
	return &File{
		storage: storage,
		name:    name,
	}
}
