package ed448

import (
	"encoding/binary"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"io"
	"math/rand"
	"net"
	"sync"
)

// it contains packet type 1 byte
// 4 bytes - packet size
// 12 bytes - additional data for the future extends
const (
	ed448PacketHeaderLen = 17
	maxPendingHandshakes = 100
)

var (
	mu                   sync.Mutex
	concurrentHandshakes = int64(0)
)

// Ed448 is transport layer that implements ED448 handshake protocol
type Ed448 struct {
	cryptos                map[string]transport.Crypto
	isHandshakeEstablished bool
	chosenEncryptionScheme string
}

func acquire() bool {
	mu.Lock()
	n := int64(1)
	success := maxPendingHandshakes-concurrentHandshakes >= n
	if success {
		concurrentHandshakes += n
	}
	mu.Unlock()
	return success
}

func release() {
	mu.Lock()
	concurrentHandshakes -= int64(1)
	if concurrentHandshakes < 0 {
		mu.Unlock()
		panic("bad release")
	}
	mu.Unlock()
}

// New creates Ed448 transport layer
func New(cryptos ...transport.Crypto) transport.Transport {
	cryptosMap := make(map[string]transport.Crypto)
	for _, crypto := range cryptos {
		cryptosMap[crypto.About()] = crypto
	}
	return &Ed448{
		cryptos: cryptosMap,
	}
}

// GetEncryptionInfo returns chosen encryption scheme
func (trans *Ed448) GetEncryptionInfo() string {
	return trans.chosenEncryptionScheme
}

// Clone creates a clone of current instance
func (trans *Ed448) Clone() transport.Transport {
	cryptos := []transport.Crypto{}
	for _, value := range trans.cryptos {
		crypto, err := value.Clone()
		if err != nil {
			panic(err)
		}
		cryptos = append(cryptos, crypto)
	}
	return New(cryptos...)
}

func (trans *Ed448) readRecord(conn net.Conn) (interface{}, error) {
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small buffer
	packetHeader := make([]byte, ed448PacketHeaderLen)

	read, err := conn.Read(packetHeader)
	if err != nil {
		return nil, errors.Errorf("error during reading pkg header %w", err)
	}

	if read != ed448PacketHeaderLen {
		return nil, errors.New(ErrWrongFormat)
	}

	expected := int(binary.LittleEndian.Uint32(packetHeader[1:5]))

	for expected != 0 {
		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			return nil, errors.Errorf("error during reading pkg body %w", err)
		}
		buf = append(buf, tmp[:n]...)
		expected -= n
	}

	// trying to decrypt message
	switch buf[0] {
	case typeClientHello:
		return DecodeClientMsg(buf)
	case typeServerHello:
		return DecodeServerMsg(buf)
	case typeClientHandshakeMsg:
		return DecodeClientHandshakeMessage(buf)
	case typeServerHandshakeMsg:
		return DecodeServerHandshakeMsg(buf)
	default:
		return nil, ErrWrongFormat
	}
}

func (trans *Ed448) writeRecord(msg message, conn net.Conn) error {
	data, err := msg.marshall()
	if err != nil {
		return err
	}

	size := len(data)
	var packetHeader = make([]byte, ed448PacketHeaderLen) // header + packet size + 4 additional bytes for the future
	packetHeader[0] = typeED448Msg
	binary.LittleEndian.PutUint32(packetHeader[1:], uint32(size))

	if _, err := conn.Write(packetHeader); err != nil {
		return errors.Errorf("errors during writing packet header %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return errors.Errorf("errors during writing data %w", err)
	}

	return nil
}

func (trans *Ed448) initEncryptedConnection(conn net.Conn, cryptoAlias string, params []byte) (net.Conn, error) {
	crypto := trans.cryptos[cryptoAlias]
	if crypto == nil {
		return nil, errors.New(ErrUnsupportedEncryption)
	}

	if err := crypto.Configure(params); err != nil {
		return nil, errors.Errorf("error during configuration %w", err)
	}

	return NewConn(conn, crypto), nil
}

// ToDo: update with appropriate implementation
func (trans *Ed448) getPastelID() []byte {
	pastelID := make([]byte, 4)
	rand.Read(pastelID)
	return pastelID
}

// ToDo: update with appropriate implementation
func (trans *Ed448) getSignedPastelID() *signedPastelID {
	var pastelID = trans.getPastelID()

	signPastelID := make([]byte, 4)
	rand.Read(pastelID)

	pubKey := make([]byte, 4)
	rand.Read(pubKey)

	ctx := make([]byte, 4)
	rand.Read(ctx)

	return &signedPastelID{
		pastelID:       pastelID,
		signedPastelID: signPastelID,
		pubKey:         pubKey,
		ctx:            ctx,
	}
}

// ToDo: update with appropriate implementation
func (trans *Ed448) verifyPastelID(_ []byte) bool {
	return true
}

// ToDo: update with external check through cNode
func (trans *Ed448) verifySignature(_, _, _, _ []byte) bool {
	return true
}
