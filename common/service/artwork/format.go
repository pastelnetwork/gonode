package artwork

import (
	"github.com/pastelnetwork/gonode/common/errors"
)

// ErrUnsupportedFormat means the given image format is not supported.
var ErrUnsupportedFormat = errors.New("imaging: unsupported image format")

// Image file formats.
const (
	JPEG Format = iota
	PNG
	GIF
	TIFF
	BMP
)

var formatExts = map[string]Format{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
	"gif":  GIF,
	"tif":  TIFF,
	"tiff": TIFF,
	"bmp":  BMP,
}

var formatNames = map[Format]string{
	JPEG: "JPEG",
	PNG:  "PNG",
	GIF:  "GIF",
	TIFF: "TIFF",
	BMP:  "BMP",
}

// Format is an image file format.
type Format int

func (f Format) String() string {
	return formatNames[f]
}
