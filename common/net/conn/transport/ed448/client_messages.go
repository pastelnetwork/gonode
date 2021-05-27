package ed448

import (
	"bytes"
	"encoding/binary"
)

// ClientHelloMessage - first Client message during handshake
type ClientHelloMessage struct {
	supportedEncryptions         []string
	supportedSignatureAlgorithms []SignScheme
}

func (msg *ClientHelloMessage) marshall() ([]byte, error) {
	var encodedMsg = bytes.Buffer{}
	var supportedEncryptionsCount = len(msg.supportedEncryptions)
	var supportedSignatureAlgorithmsCount = len(msg.supportedSignatureAlgorithms)

	encodedMsg.WriteByte(typeClientHello)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, supportedEncryptionsCount); err != nil {
		return nil, err
	}
	for _, encryption := range msg.supportedEncryptions {
		if err := writeString(&encodedMsg, &encryption); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(&encodedMsg, binary.LittleEndian, supportedSignatureAlgorithmsCount); err != nil {
		return nil, err
	}

	for _, signature := range msg.supportedSignatureAlgorithms {
		encodedMsg.WriteByte(byte(signature))
	}
	return encodedMsg.Bytes(), nil
}

// DecodeClientMsg - unmarshall []byte to ClientHelloMessage
func DecodeClientMsg(msg []byte) (*ClientHelloMessage, error) {
	if msg[0] != typeClientHello {
		return nil, ErrWrongFormat
	}
	reader := bytes.NewReader(msg[1:])

	var supportedEncryptionsCount int
	if err := binary.Read(reader, binary.LittleEndian, &supportedEncryptionsCount); err != nil {
		return nil, err
	}

	supportedEncryptions := make([]string, supportedEncryptionsCount)
	for i := 0; i < supportedEncryptionsCount; i++ {
		encryption, err := readString(reader)
		if err != nil {
			return nil, err
		}
		supportedEncryptions[i] = *encryption
	}

	var supportedSignatureAlgorithmsCount int
	if err := binary.Read(reader, binary.LittleEndian, &supportedSignatureAlgorithmsCount); err != nil {
		return nil, err
	}
	supportedSignatureAlgorithms := make([]SignScheme, supportedEncryptionsCount)
	for i := 0; i < int(supportedSignatureAlgorithmsCount); i++ {
		signScheme, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		supportedSignatureAlgorithms[i] = SignScheme(signScheme)
	}

	return &ClientHelloMessage{
		supportedEncryptions:         supportedEncryptions,
		supportedSignatureAlgorithms: supportedSignatureAlgorithms,
	}, nil
}

// ClientHandshakeMessage - second Client message during handshake
type ClientHandshakeMessage struct {
	pastelID       []byte
	signedPastelID []byte
	pubKey         []byte
	ctx            []byte
}

func (msg *ClientHandshakeMessage) marshall() ([]byte, error) {
	var encodedMsg = bytes.Buffer{}
	encodedMsg.WriteByte(typeClientHandshakeMsg)
	if err := writeByteArray(&encodedMsg, &msg.pastelID); err != nil {
		return nil, err
	}
	if err := writeByteArray(&encodedMsg, &msg.signedPastelID); err != nil {
		return nil, err
	}
	if err := writeByteArray(&encodedMsg, &msg.pubKey); err != nil {
		return nil, err
	}
	if err := writeByteArray(&encodedMsg, &msg.ctx); err != nil {
		return nil, err
	}

	return encodedMsg.Bytes(), nil
}

// DecodeClientHandshakeMessage - unmarshall []byte to ClientHandshakeMessage
func DecodeClientHandshakeMessage(msg []byte) (*ClientHandshakeMessage, error) {
	if msg[0] != typeClientHandshakeMsg {
		return nil, ErrWrongFormat
	}

	reader := bytes.NewReader(msg[1:])
	pastelID, err := readByteArray(reader)
	if err != nil {
		return nil, err
	}

	signedPastelID, err := readByteArray(reader)
	if err != nil {
		return nil, err
	}

	pubKey, err := readByteArray(reader)
	if err != nil {
		return nil, err
	}

	ctx, err := readByteArray(reader)
	if err != nil {
		return nil, err
	}

	return &ClientHandshakeMessage{
		pastelID:       *pastelID,
		signedPastelID: *signedPastelID,
		pubKey:         *pubKey,
		ctx:            *ctx,
	}, nil
}
