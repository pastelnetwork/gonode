package artwork

type Image struct {
	Data []byte
	Size int
}

func NewImage(data []byte) *Image {
	return &Image{
		Data: data,
		Size: len(data),
		//ExpiresIn: time.Now().Add(imageTTL),
	}
}
