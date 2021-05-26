package conn

import (
	"io"
	"math/rand"
	"net"
)

type Ed448Transport struct {
	cryptos                map[string]Crypto
	isHandshakeEstablished bool
	chosenEncryptionScheme string
}

func NewEd448Transport(cryptos ...Crypto) *Ed448Transport {
	cryptosMap := make(map[string]Crypto)
	for _, crypto := range cryptos {
		cryptosMap[crypto.About()] = crypto
	}
	return &Ed448Transport{
		cryptos: cryptosMap,
	}
}

// IsHandshakeEstablished - returns flag that says is handshake process finished
func (t *Ed448Transport) IsHandshakeEstablished() bool {
	return t.isHandshakeEstablished
}

func (t *Ed448Transport) readRecord(conn net.Conn) (interface{}, error) {
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small buffer

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		buf = append(buf, tmp[:n]...)

	}
	if _, err := conn.Read(buf); err == nil {
		return nil, err
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

func (t *Ed448Transport) writeRecord(msg message, conn net.Conn) error {
	data, err := msg.marshall()
	if err != nil {
		return err
	}

	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}

func (t *Ed448Transport) initEncryptedConnection(conn net.Conn) {

}

// ToDo: update with appropriate implementation
func (t *Ed448Transport) getPastelID() []byte {
	pastelID := make([]byte, 4)
	rand.Read(pastelID)
	return pastelID
}

// ToDo: update with appropriate implementation
func (t *Ed448Transport) getSignedPastelId() *signedPastelID {
	var pastelID = t.getPastelID()

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
func (t *Ed448Transport) verifyPastelID(_ []byte) bool {
	return true
}

// ToDo: update with external check through cNode
func (t *Ed448Transport) verifySignature(_, _, _, _ []byte) bool {
	return true
}
