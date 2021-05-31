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
	writeString(&encodedMsg, &msg.chosenEncryption)
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
	encryptionParams []byte
}

func (msg *ServerHandshakeMessage) marshall() ([]byte, error) {
	encodedMsg := bytes.Buffer{}
	encodedMsg.WriteByte(typeServerHandshakeMsg)
	writeByteArray(&encodedMsg, &msg.pastelID)
	writeByteArray(&encodedMsg, &msg.signedPastelID)
	writeByteArray(&encodedMsg, &msg.pubKey)
	writeByteArray(&encodedMsg, &msg.ctx)
	writeByteArray(&encodedMsg, &msg.encryptionParams)
	return encodedMsg.Bytes(), nil
}

// DecodeServerHandshakeMsg - unmarshall []byte to ServerHandshakeMessage
func DecodeServerHandshakeMsg(msg []byte) (*ServerHandshakeMessage, error) {
	if msg[0] != typeServerHandshakeMsg {
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
		encryptionParams: *encryptionParams,
	}, nil
}
