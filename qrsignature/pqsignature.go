// pqsignature signs an image using post-quantum cryptography.

package pqsignature

import "image"

type Signature struct {
	image.Image

	qrImages []image.Image
}

func (pq *Signature) Add(data []byte, title, alias string) error {

	return nil
}

func (pq *Signature) SignImage(src image.Image) (image.Image, error) {

	return nil, nil
}

func New(img image.Image) *Signature {
	return &Signature{
		Image: img,
	}
}
