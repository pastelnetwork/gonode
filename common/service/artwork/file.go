package artwork

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
)

// File represents a file.
type File struct {
	fmt.Stringer
	sync.Mutex

	storage *Storage

	isCreated bool
	name      string
	format    Format
}

// Name returns filename.
func (file *File) Name() string {
	return file.name
}

func (file *File) String() string {
	return file.name
}

// SetFormatFromExtension parses and sets image format from filename extension:
// "jpg" (or "jpeg"), "png", "gif" are supported.
func (file *File) SetFormatFromExtension(ext string) error {
	if f, ok := formatExts[strings.ToLower(strings.TrimPrefix(ext, "."))]; ok {
		file.format = f
		return nil
	}
	return ErrUnsupportedFormat
}

// Format returns file extension.
func (file *File) Format() Format {
	return file.format
}

// Open opens a file and returns file descriptor.
// If file is not found, storage.ErrFileNotFound is returned.
func (file *File) Open() (storage.File, error) {
	file.Lock()
	defer file.Unlock()

	return file.storage.Open(file.name)
}

// Create creates a file and returns file descriptor.
func (file *File) Create() (storage.File, error) {
	file.Lock()
	defer file.Unlock()

	fl, err := file.storage.Create(file.name)
	if err != nil {
		return nil, err
	}

	file.isCreated = true
	return fl, nil
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

// Copy creates a copy of the current file.
func (file *File) Copy() (*File, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	newFile := file.storage.NewFile()
	dst, err := newFile.Create()
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, errors.Errorf("failed to copy file: %w", err)
	}
	return newFile, nil
}

// Bytes returns the contents of the file by bytes.
func (file *File) Bytes() ([]byte, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, errors.Errorf("failed to read from file: %w", err)
	}
	return buf.Bytes(), nil
}

// Write writes data to the file.
func (file *File) Write(data []byte) error {
	f, err := file.Create()
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return errors.Errorf("failed to write to file: %w", err)
	}
	return nil
}

// ResizeImage resizes image.
func (file *File) ResizeImage(width, height int) error {
	src, err := file.OpenImage()
	if err != nil {
		return err
	}

	dst := imaging.Resize(src, width, height, imaging.Lanczos)

	return file.SaveImage(dst)
}

// RemoveAfter removes the file after the specified duration.
func (file *File) RemoveAfter(d time.Duration) {
	go func() {
		time.AfterFunc(d, func() { file.Remove() })
	}()
}

// OpenImage opens images from the file.
func (file *File) OpenImage() (image.Image, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Errorf("failed to decode image: %w", err).WithField("filename", f.Name())
	}
	return img, nil
}

// SaveImage saves image to the file.
func (file *File) SaveImage(img image.Image) error {
	f, err := file.Create()
	if err != nil {
		return err
	}
	defer f.Close()

	switch file.format {
	case JPEG:
		if nrgba, ok := img.(*image.NRGBA); ok && nrgba.Opaque() {
			rgba := &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
			if err := jpeg.Encode(f, rgba, nil); err != nil {
				return errors.Errorf("failed to encode jpeg: %w", err).WithField("filename", f.Name())
			}
			return nil
		}
		if err := jpeg.Encode(f, img, nil); err != nil {
			return errors.Errorf("failed to encode jpeg: %w", err).WithField("filename", f.Name())
		}
		return nil

	case PNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := encoder.Encode(f, img); err != nil {
			return errors.Errorf("failed to encode png: %w", err).WithField("filename", f.Name())
		}
		return nil

	case GIF:
		if err := gif.Encode(f, img, nil); err != nil {
			return errors.Errorf("failed to encode gif: %w", err).WithField("filename", f.Name())
		}
		return nil
	}
	return ErrUnsupportedFormat
}

// NewFile returns a newFile File instance.
func NewFile(storage *Storage, name string) *File {
	return &File{
		storage: storage,
		name:    name,
	}
}
