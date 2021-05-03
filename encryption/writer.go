package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// StreamEncrypter is an encrypter for a stream of data with authentication
type StreamEncrypter struct {
	Dist io.Writer
	Block  cipher.Block
	Stream cipher.Stream
	StreamWriter *cipher.StreamWriter
}

func (s *StreamEncrypter) Write(b []byte) (n int, err error){
	return s.StreamWriter.Write(b)
}

func NewStreamEncrypter(key []byte, iv []byte, writer io.Writer) (*StreamEncrypter, error){
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	return &StreamEncrypter{
		Dist: writer,
		Stream: stream,
		StreamWriter: &cipher.StreamWriter{S: stream, W: writer},
	}, nil
}