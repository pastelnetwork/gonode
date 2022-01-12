package files

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
)

// File represents a file.
type File struct {
	fmt.Stringer
	sync.Mutex

	storage.FileInterface
	storage *Storage

	// if a file was created during the process, it should be deleted at the end.
	isCreated bool

	// unique name within the storage.
	name string

	// file format, png, jpg, etc.
	format Format
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
	if format, ok := formatExts[strings.ToLower(strings.TrimPrefix(ext, "."))]; ok {
		return file.SetFormat(format)
	}
	return ErrUnsupportedFormat
}

// SetFormat sets file extension.
func (file *File) SetFormat(format Format) error {
	file.format = format

	newname := fmt.Sprintf("%s.%s", strings.TrimSuffix(file.name, filepath.Ext(file.name)), format)
	oldname := file.name
	file.name = newname

	if err := file.storage.Update(oldname, newname, file); err != nil {
		return err
	}

	if !file.isCreated {
		return nil
	}
	return file.storage.Rename(oldname, newname)
}

// Format returns file extension.
func (file *File) Format() Format {
	return file.format
}

// Open opens a file and returns file descriptor.
// If file is not found, storage.ErrFileNotFound is returned.
func (file *File) Open() (storage.FileInterface, error) {
	file.Lock()
	defer file.Unlock()

	return file.storage.Open(file.Name())
}

// Create creates a file and returns file descriptor.
func (file *File) Create() (storage.FileInterface, error) {
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

	delete(file.storage.filesMap, file.name)

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
	if err := newFile.SetFormat(file.format); err != nil {
		return nil, err
	}

	dst, err := newFile.Create()
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, errors.Errorf("copy file: %w", err)
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
		return nil, errors.Errorf("read file: %w", err)
	}

	return buf.Bytes(), nil
}

// Write writes data to the file.
func (file *File) Write(data []byte) (n int, err error) {
	f, err := file.Create()
	if err != nil {
		return 0, errors.Errorf("create file: %w", err)
	}
	defer f.Close()

	n, err = f.Write(data)
	if err != nil {
		return n, errors.Errorf("write file: %w", err)
	}

	return
}

// ResizeImage resizes image.
func (file *File) ResizeImage(width, height int) error {
	src, err := file.LoadImage()
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

// LoadImage opens images from the file.
func (file *File) LoadImage() (image.Image, error) {
	f, err := file.Open()
	if err != nil {
		log.Debug("Failed to open")
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Errorf("decode image(%s): %w", f.Name(), err)
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
				return errors.Errorf("encode jpeg rgba(%s): %w", f.Name(), err)
			}
			return nil
		}
		if err := jpeg.Encode(f, img, nil); err != nil {
			return errors.Errorf("encode jpeg(%s): %w", f.Name(), err)
		}
		return nil

	case PNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := encoder.Encode(f, img); err != nil {
			return errors.Errorf("encode png(%s): %w", f.Name(), err)
		}
		return nil

	case GIF:
		if err := gif.Encode(f, img, nil); err != nil {
			return errors.Errorf("encode gif(%s): %w", f.Name(), err)
		}
		return nil
	}
	return ErrUnsupportedFormat
}

// Thumbnail creates a thumbnail file from the artwork file and store in to starage layer
func (file *File) Thumbnail(coordinate ThumbnailCoordinate) (*File, error) {
	f := NewFile(file.storage, "thumbnail-of-"+file.name)
	if f == nil {
		return nil, errors.Errorf("create new file for thumbnail-of-%q", file.Name())
	}
	if err := f.SetFormat(file.Format()); err != nil {
		return nil, errors.Errorf("set format for thumbnail-of-%q", file.Name())
	}

	img, err := file.LoadImage()
	if err != nil {
		return nil, errors.Errorf("load image from file(%s): %w", file.Name(), err)
	}

	rect := image.Rect(int(coordinate.TopLeftX), int(coordinate.TopLeftY), int(coordinate.BottomRightX), int(coordinate.BottomRightY))
	thumbnail := imaging.Crop(img, rect)
	if thumbnail == nil {
		return nil, errors.Errorf("generate thumbnail(%s): %w", file.Name(), err)
	}

	if err := f.SaveImage(thumbnail); err != nil {
		return nil, errors.Errorf("save thumbnail(%s): %w", file.Name(), err)
	}

	return f, nil
}

// UpdateFormat updates file format
func (file *File) UpdateFormat() error {
	f, err := file.Open()
	if err != nil {
		log.Debug("Failed to open")
		return err
	}
	defer f.Close()

	_, format, err := image.Decode(f)
	if err != nil {
		return errors.Errorf("decode image(%s): %w", file.Name(), err)
	}
	err = file.SetFormatFromExtension(format)
	if err != nil {
		return errors.Errorf("set file format(%s): %w", file.Name(), err)
	}
	return nil
}

// Encoder represents an image encoder.
type Encoder interface {
	Encode(img image.Image) (image.Image, error)
}

// Encode encodes the image by the given encoder.
func (file *File) Encode(enc Encoder) error {
	img, err := file.LoadImage()
	if err != nil {
		return err
	}

	encImg, err := enc.Encode(img)
	if err != nil {
		return err
	}
	return file.SaveImage(encImg)
}

// Decoder represents an image decoder.
type Decoder interface {
	Decode(img image.Image) error
}

// Decode decodes the image by the given decoder.
func (file *File) Decode(dec Decoder) error {
	img, err := file.LoadImage()
	if err != nil {
		return err
	}
	return dec.Decode(img)
}

// NewFile returns a newFile File instance.
func NewFile(storage *Storage, name string) *File {
	return &File{
		storage: storage,
		name:    name,
	}
}
