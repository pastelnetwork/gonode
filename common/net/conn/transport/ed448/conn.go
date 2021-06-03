package ed448

import (
	"encoding/binary"
	// "fmt"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"math"
	"net"
)

// conn represents a secured connection. It implements the net.Conn interface.
type conn struct {
	net.Conn
	crypto             transport.Crypto
	overhead           int
	payloadLengthLimit int
	nextFrame          []byte
	writeBuf           []byte
}

const recordDefaultLength = 4 * 1024     // 4KiB
const writeBufferInitialSize = 32 * 1024 // 32KiB
const writeBufferMaxSize = 512 * 1024    // 512KiB

// NewConn creates a new secure channel instance
func NewConn(c net.Conn, crypto transport.Crypto) net.Conn {
	overhead := msgHeaderSize + crypto.EncryptionOverhead()
	return &conn{
		Conn:               c,
		crypto:             crypto,
		overhead:           overhead,
		payloadLengthLimit: recordDefaultLength - overhead,
		nextFrame:          make([]byte, 0, 2*recordDefaultLength),
		writeBuf:           make([]byte, writeBufferInitialSize),
	}
}

// Read reads and decrypts a frame from the underlying connection, and copies the
// decrypted payload into b
func (c *conn) Read(b []byte) (n int, err error) {
	var framedMsg []byte
	framedMsg, c.nextFrame, err = parseFramedMsg(c.nextFrame, msgMaxSize)
	if err != nil {
		return n, err
	}

	// Check whether a complete frame has been received yet.
	for len(framedMsg) == 0 {
		// frame is full and incomplete increase size of it
		if len(c.nextFrame) == cap(c.nextFrame) {
			tmp := make([]byte, len(c.nextFrame), cap(c.nextFrame)+recordDefaultLength)
			copy(tmp, c.nextFrame)
			c.nextFrame = tmp
		}

		// read last data
		n, err = c.Conn.Read(c.nextFrame[len(c.nextFrame):cap(c.nextFrame)])
		if err != nil {
			return 0, errors.Errorf("error during reading data %s", err)
		}
		c.nextFrame = c.nextFrame[:len(c.nextFrame)+n]
		framedMsg, c.nextFrame, err = parseFramedMsg(c.nextFrame, msgMaxSize)
		if err != nil {
			return 0, err
		}
	}

	// Now we have a complete frame, decrypted it.
	if framedMsg[0] != typeEncryptedMsg {
		return 0, errors.New(ErrWrongFormat)
	}

	decryptedData, err := c.crypto.Decrypt(framedMsg[msgHeaderSize:])
	if err != nil {
		return 0, errors.Errorf("error during decrypting data %w", err)
	}
	n = copy(b, decryptedData)
	return n, nil
}

// Write encrypts, frames, and writes bytes from b to the underlying connection.
func (c *conn) Write(b []byte) (int, error) {
	// fmt.Println("write data")
	msgSize := len(b)

	// calculate the output buffer size with framing and encryption overhead.
	numOfFrames := int(math.Ceil(float64(len(b)) / float64(c.payloadLengthLimit)))
	size := msgSize + numOfFrames*c.overhead

	partialBSize := len(b)
	if size > writeBufferMaxSize {
		size = writeBufferMaxSize
		const numOfFramesInMaxWriteBuf = writeBufferMaxSize / recordDefaultLength
		partialBSize = numOfFramesInMaxWriteBuf * c.payloadLengthLimit
	}
	if len(c.writeBuf) < size {
		c.writeBuf = make([]byte, size)
	}

	for partialBStart := 0; partialBStart < len(b); partialBStart += partialBSize {
		partialBEnd := partialBStart + partialBSize
		if partialBEnd > len(b) {
			partialBEnd = len(b)
		}
		partialB := b[partialBStart:partialBEnd]
		writeBufIndex := 0
		for len(partialB) > 0 {
			payloadLen := len(partialB)
			if payloadLen > c.payloadLengthLimit {
				payloadLen = c.payloadLengthLimit
			}
			buf := partialB[:payloadLen]
			partialB = partialB[payloadLen:]

			// write buffer contains: type, length, payload
			c.writeBuf[writeBufIndex] = typeEncryptedMsg
			// encrypt the payload
			cypherText, err := c.crypto.Encrypt(buf)
			if err != nil {
				return msgSize, err
			}
			// fill in the size field
			binary.LittleEndian.PutUint32(c.writeBuf[writeBufIndex+msgTypeLen:], uint32(len(cypherText)))
			copy(c.writeBuf[writeBufIndex+msgHeaderSize:], cypherText)
			// increase writeBufIndex
			writeBufIndex += len(buf) + c.overhead
		}
		nn, err := c.Conn.Write(c.writeBuf[:writeBufIndex])
		if err != nil {
			// calculate the actual data size that was written.
			// this means we need to remove header, encryption overheads, and any partially-written frame data.
			numOfWrittenFrames := int(math.Floor(float64(nn) / float64(recordDefaultLength)))
			return partialBStart + numOfWrittenFrames*c.payloadLengthLimit, err
		}
	}
	return msgSize, nil
}
