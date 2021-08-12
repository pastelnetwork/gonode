package conn

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

// ALTSRecordCrypto is the interface for gRPC ALTS record protocol.
type ALTSRecordCrypto interface {
	// Encrypt encrypts the plaintext and computes the tag (if any) of dst
	// and plaintext. dst and plaintext may fully overlap or not at all.
	Encrypt(dst, plaintext []byte) ([]byte, error)
	// EncryptionOverhead returns the tag size (if any) in bytes.
	EncryptionOverhead() int
	// Decrypt decrypts ciphertext and verify the tag (if any). dst and
	// ciphertext may alias exactly or not at all. To reuse ciphertext's
	// storage for the decrypted output, use ciphertext[:0] as dst.
	Decrypt(dst, ciphertext []byte) ([]byte, error)
}

// ALTSRecordFunc is a function type for factory functions that create
// ALTSRecordCrypto instances.
type ALTSRecordFunc func(s alts.Side, keyData []byte) (ALTSRecordCrypto, error)

const (
	// frameLenFieldSize is the byte size of the frame length field of a
	// framed message.
	frameLenFieldSize = 4
	// The byte size of the message type field of a framed message.
	frameTypeFieldSize = 4
	// The bytes size limit for a ALTS record message.
	altsRecordSizeLimit = 1024 * 1024 // 1 MiB
	// The default bytes size of a ALTS record message.
	altsRecordDefaultSize = 4 * 1024 // 4KiB
	// Message type value included in ALTS record framing.
	altsRecordMsgType = uint32(0x06)
	// The initial write buffer size.
	altsWriteBufferInitialSize = 32 * 1024 // 32KiB
	// The maximum write buffer size. This *must* be multiple of
	// altsRecordDefaultSize.
	altsWriteBufferMaxSize = 512 * 1024 // 512KiB
)

var (
	protocols = make(map[string]ALTSRecordFunc)
)

// RegisterProtocol register a ALTS record encryption protocol.
func RegisterProtocol(protocol string, f ALTSRecordFunc) error {
	if _, ok := protocols[protocol]; ok {
		return fmt.Errorf("protocol %v is already registered", protocol)
	}
	protocols[protocol] = f
	return nil
}

// conn represents a secured connection. It implements the net.Conn interface.
type conn struct {
	net.Conn

	// record crypto
	crypto ALTSRecordCrypto

	// buf holds data that has been read from the connection and decrypted,
	// but has not yet been returned by Read.
	buf          []byte
	payloadLimit int
	// protected holds data read from the network but have not yet been
	// decrypted. This data might not compose a complete frame.
	protected []byte
	// writeBuf is a buffer used to contain encrypted frames before being
	// written to the network.
	writeBuf []byte
	// nextFrame stores the next frame (in protected buffer) info.
	nextFrame []byte
	// overhead is the calculated overhead of each frame.
	overhead int

	// the secret key for connection's encryption
	secret []byte
}

// NewConn creates a new secure channel instance given the other party role and
// handshaking result.
func NewConn(side alts.Side, c net.Conn, protocol string, secret []byte) (net.Conn, error) {
	newCrypto := protocols[protocol]
	if newCrypto == nil {
		return nil, fmt.Errorf("unknown next_protocol %q", protocol)
	}
	crypto, err := newCrypto(side, secret)
	if err != nil {
		return nil, fmt.Errorf("new crypto for %q: %v", protocol, err)
	}

	overhead := frameLenFieldSize + frameTypeFieldSize + crypto.EncryptionOverhead()
	payloadLimit := altsRecordDefaultSize - overhead

	// We pre-allocate protected to be of size
	// 2*altsRecordDefaultLength-1 during initialization. We only
	// read from the network into protected when protected does not
	// contain a complete frame, which is at most
	// altsRecordDefaultLength-1 (bytes). And we read at most
	// altsRecordDefaultLength (bytes) data into protected at one
	// time. Therefore, 2*altsRecordDefaultLength-1 is large enough
	// to buffer data read from the network.
	protected := make([]byte, 0, 2*altsRecordDefaultSize-1)

	altsConn := &conn{
		Conn:         c,
		crypto:       crypto,
		payloadLimit: payloadLimit,
		protected:    protected,
		writeBuf:     make([]byte, altsWriteBufferInitialSize),
		nextFrame:    protected,
		overhead:     overhead,
		secret:       secret,
	}
	return altsConn, nil
}

// Read reads and decrypts a frame from the underlying connection, and copies the
// decrypted payload into b. If the size of the payload is greater than len(b),
// Read retains the remaining bytes in an internal buffer, and subsequent calls
// to Read will read from this buffer until it is exhausted.
func (s *conn) Read(b []byte) (n int, err error) {
	if len(s.buf) == 0 {
		var frame []byte
		frame, s.nextFrame, err = ParseFrame(s.nextFrame, altsRecordSizeLimit)
		if err != nil {
			return n, err
		}
		// Check whether the next frame to be decrypted has been
		// completely received yet.
		if len(frame) == 0 {
			copy(s.protected, s.nextFrame)
			s.protected = s.protected[:len(s.nextFrame)]
			// Always copy next incomplete frame to the beginning of
			// the protected buffer and reset nextFrame to it.
			s.nextFrame = s.protected
		}
		// Check whether a complete frame has been received yet.
		for len(frame) == 0 {
			if len(s.protected) == cap(s.protected) {
				tmp := make([]byte, len(s.protected), cap(s.protected)+altsRecordDefaultSize)
				copy(tmp, s.protected)
				s.protected = tmp
			}
			n, err = s.Conn.Read(s.protected[len(s.protected):min(cap(s.protected), len(s.protected)+altsRecordDefaultSize)])
			if err != nil {
				return 0, err
			}
			s.protected = s.protected[:len(s.protected)+n]
			frame, s.nextFrame, err = ParseFrame(s.protected, altsRecordSizeLimit)
			if err != nil {
				return 0, err
			}
		}

		// Now we have a complete frame, decrypted it.
		msg := frame[frameLenFieldSize:]
		msgType := binary.LittleEndian.Uint32(msg[:frameTypeFieldSize])
		if msgType&0xff != altsRecordMsgType {
			return 0, fmt.Errorf("received frame with incorrect message type %v, expected lower byte %v", msgType, altsRecordMsgType)
		}
		ciphertext := msg[frameTypeFieldSize:]

		// Decrypt requires that if the dst and ciphertext alias, they
		// must alias exactly. Code here used to use msg[:0], but msg
		// starts MsgLenFieldSize+msgTypeFieldSize bytes earlier than
		// ciphertext, so they alias inexactly. Using ciphertext[:0]
		// arranges the appropriate aliasing without needing to copy
		// ciphertext or use a separate destination buffer. For more info
		// check: https://golang.org/pkg/crypto/cipher/#AEAD.
		// decrypt the record payload
		s.buf, err = s.crypto.Decrypt(ciphertext[:0], ciphertext)
		if err != nil {
			return 0, err
		}
	}
	n = copy(b, s.buf)
	s.buf = s.buf[n:]

	return n, nil
}

// Write encrypts, frames, and writes bytes from b to the underlying connection.
func (s *conn) Write(b []byte) (n int, err error) {
	n = len(b)

	// Calculate the output buffer size with framing and encryption overhead.
	framesCount := int(math.Ceil(float64(len(b)) / float64(s.payloadLimit)))

	size := len(b) + framesCount*s.overhead
	// If writeBuf is too small, increase its size up to the maximum size.
	partialSize := len(b)
	if size > altsWriteBufferMaxSize {
		size = altsWriteBufferMaxSize
		maxCount := altsWriteBufferMaxSize / altsRecordDefaultSize
		partialSize = maxCount * s.payloadLimit
	}
	if len(s.writeBuf) < size {
		s.writeBuf = make([]byte, size)
	}

	for start := 0; start < len(b); start += partialSize {
		end := start + partialSize
		if end > len(b) {
			end = len(b)
		}
		partial := b[start:end]

		writeBufIndex := 0
		for len(partial) > 0 {
			payloadLen := len(partial)
			if payloadLen > s.payloadLimit {
				payloadLen = s.payloadLimit
			}
			buf := partial[:payloadLen]
			partial = partial[payloadLen:]

			// Write buffer contains: length, type, payload, and tag
			// if any.

			// 1. Fill in type field.
			msg := s.writeBuf[writeBufIndex+frameLenFieldSize:]
			binary.LittleEndian.PutUint32(msg, altsRecordMsgType)

			// 2. Encrypt the payload and create a tag if any.
			msg, err = s.crypto.Encrypt(msg[:frameTypeFieldSize], buf)
			if err != nil {
				return n, err
			}

			// 3. Fill in the size field.
			binary.LittleEndian.PutUint32(s.writeBuf[writeBufIndex:], uint32(len(msg)))

			// 4. Increase writeBufIndex.
			writeBufIndex += len(buf) + s.overhead
		}
		nn, err := s.Conn.Write(s.writeBuf[:writeBufIndex])
		if err != nil {
			// We need to calculate the actual data size that was
			// written. This means we need to remove header,
			// encryption overheads, and any partially-written
			// frame data.
			numOfWrittenFrames := int(math.Floor(float64(nn) / float64(altsRecordDefaultSize)))
			return start + numOfWrittenFrames*s.payloadLimit, err
		}
	}
	return n, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
