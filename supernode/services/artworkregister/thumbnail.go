package artworkregister

type thumbnailType int

const (
	// previewThumbnail is best effort webp image from the original artwork
	previewThumbnail thumbnailType = iota

	// mediumThumbnail is wepb imnage from cropped image
	mediumThumbnail

	// smallhumbnail is snaller webp imnage from cropped image
	smallThumbnail
)

// thumbnails defines the target size of each thumbnail type and their quality
var thumbnails = map[thumbnailType]struct {
	targetSize int
	quality    float32
}{
	previewThumbnail: {targetSize: 400000, quality: 85},
	mediumThumbnail:  {targetSize: 25000, quality: 25},
	smallThumbnail:   {targetSize: 10000, quality: 10},
}
