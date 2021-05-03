package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type StreamDecrypter struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
}

func (s *StreamDecrypter) Read(p []byte) (int, error) {
	n, readErr := s.Source.Read(p)
	if n > 0 {
		s.Stream.XORKeyStream(p[:n], p[:n])
		return n, readErr
	}
	return 0, io.EOF
}

// NewStreamDecrypter creates a new stream decrypter
func NewStreamDecrypter(key []byte, iv []byte, reader io.Reader) (*StreamDecrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	return &StreamDecrypter{
		Source: reader,
		Block:  block,
		Stream: stream,
	}, nil
}
