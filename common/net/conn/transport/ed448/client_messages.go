package ed448

import (
	"bytes"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
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
	writeInt(&encodedMsg, supportedEncryptionsCount)
	for _, encryption := range msg.supportedEncryptions {
		writeString(&encodedMsg, &encryption)
	}

	writeInt(&encodedMsg, supportedSignatureAlgorithmsCount)

	for _, signature := range msg.supportedSignatureAlgorithms {
		encodedMsg.WriteByte(byte(signature))
	}
	return encodedMsg.Bytes(), nil
}

// DecodeClientMsg - unmarshall []byte to ClientHelloMessage
func DecodeClientMsg(msg []byte) (*ClientHelloMessage, error) {
	if msg[0] != typeClientHello {
		fmt.Println("DecodeClientMsg")
		return nil, errors.New(ErrWrongFormat)
	}
	reader := bytes.NewReader(msg[1:])

	supportedEncryptionsCount, err := readInt(reader)
	if err != nil {
		return nil, errors.Errorf("can not read int value %w", err)
	}

	supportedEncryptions := make([]string, supportedEncryptionsCount)
	for i := 0; i < supportedEncryptionsCount; i++ {
		encryption, err := readString(reader)
		if err != nil {
			return nil, errors.Errorf("can not read string value %w", err)
		}
		supportedEncryptions[i] = *encryption
	}

	supportedSignatureAlgorithmsCount, err := readInt(reader)
	if err != nil {
		return nil, errors.Errorf("can not read int value %w", err)
	}
	supportedSignatureAlgorithms := make([]SignScheme, supportedEncryptionsCount)
	for i := 0; i < int(supportedSignatureAlgorithmsCount); i++ {
		signScheme, err := reader.ReadByte()
		if err != nil {
			return nil, errors.Errorf("can not read SingScheme %w", err)
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
	writeByteArray(&encodedMsg, &msg.pastelID)
	writeByteArray(&encodedMsg, &msg.signedPastelID)
	writeByteArray(&encodedMsg, &msg.pubKey)
	writeByteArray(&encodedMsg, &msg.ctx)

	return encodedMsg.Bytes(), nil
}

// DecodeClientHandshakeMessage - unmarshall []byte to ClientHandshakeMessage
func DecodeClientHandshakeMessage(msg []byte) (*ClientHandshakeMessage, error) {
	if msg[0] != typeClientHandshakeMsg {
		return nil, errors.New(ErrWrongFormat)
	}

	reader := bytes.NewReader(msg[1:])
	pastelID, err := readByteArray(reader)
	if err != nil {
		return nil, errors.Errorf("can not read PastelID %w", err)
	}

	signedPastelID, err := readByteArray(reader)
	if err != nil {
		return nil, errors.Errorf("can not read signed PastelID %w", err)
	}

	pubKey, err := readByteArray(reader)
	if err != nil {
		return nil, errors.Errorf("can not read public key %w", err)
	}

	ctx, err := readByteArray(reader)
	if err != nil {
		return nil, errors.Errorf("can not read context data %w", err)
	}

	return &ClientHandshakeMessage{
		pastelID:       *pastelID,
		signedPastelID: *signedPastelID,
		pubKey:         *pubKey,
		ctx:            *ctx,
	}, nil
}
