package files

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
)

var formatExts = map[string]Format{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
}

var formatNames = map[Format]string{
	JPEG: "jpeg",
	PNG:  "png",
}

// Format is an image file format.
type Format int

func (f Format) String() string {
	return formatNames[f]
}
