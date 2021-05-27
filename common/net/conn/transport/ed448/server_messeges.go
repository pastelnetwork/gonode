package ed448

import (
	"bytes"
)

// ServerHelloMessage - first server message during handshake
type ServerHelloMessage struct {
	chosenEncryption string
	chosenSignature  SignScheme
}

func (msg *ServerHelloMessage) marshall() ([]byte, error) {
	encodedMsg := bytes.Buffer{}
	encodedMsg.WriteByte(typeServerHello)
	if err := writeString(&encodedMsg, &msg.chosenEncryption); err != nil {
		return nil, err
	}
	encodedMsg.WriteByte(byte(msg.chosenSignature))
	return encodedMsg.Bytes(), nil
}

// DecodeServerMsg - unmarshall []byte to ServerHelloMessage
func DecodeServerMsg(msg []byte) (*ServerHelloMessage, error) {
	if msg[0] != typeServerHello {
		return nil, ErrWrongFormat
	}

	reader := bytes.NewReader(msg[1:])
	chosenEncryption, err := readString(reader)
	if err != nil {
		return nil, err
	}

	binarySignature, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ServerHelloMessage{
		chosenEncryption: *chosenEncryption,
		chosenSignature:  SignScheme(binarySignature),
	}, nil
}

// ServerHandshakeMessage - second server message during handshake process
type ServerHandshakeMessage struct {
	pastelID         []byte
	signedPastelID   []byte
	pubKey           []byte
	ctx              []byte
	encryptionParams string
}

func (msg *ServerHandshakeMessage) marshall() ([]byte, error) {
	encodedMsg := bytes.Buffer{}
	encodedMsg.WriteByte(typeClientHandshakeMsg)
	if err := writeByteArray(&encodedMsg, &msg.pastelID); err != nil {
		return nil, err
	}
	if err := writeByteArray(&encodedMsg, &msg.pubKey); err != nil {
		return nil, err
	}
	if err := writeByteArray(&encodedMsg, &msg.ctx); err != nil {
		return nil, err
	}
	if err := writeString(&encodedMsg, &msg.encryptionParams); err != nil {
		return nil, err
	}
	return encodedMsg.Bytes(), nil
}

// DecodeServerHandshakeMsg - unmarshall []byte to ServerHandshakeMessage
func DecodeServerHandshakeMsg(msg []byte) (*ServerHandshakeMessage, error) {
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

	encryptionParams, err := readByteArray(reader)
	if err != nil {
		return nil, err
	}

	return &ServerHandshakeMessage{
		pastelID:         *pastelID,
		signedPastelID:   *signedPastelID,
		pubKey:           *pubKey,
		ctx:              *ctx,
		encryptionParams: string(*encryptionParams),
	}, nil
}
