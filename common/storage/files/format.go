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
	WEBP
)

var formatExts = map[string]Format{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
	"webp": WEBP,
}

var formatNames = map[Format]string{
	JPEG: "jpeg",
	PNG:  "png",
	WEBP: "webp",
}

// Format is an image file format.
type Format int

func (f Format) String() string {
	return formatNames[f]
}
