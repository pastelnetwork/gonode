package conn

import (
	"bytes"
	"encoding/binary"
)

// ServerHelloMessage - first server message during handshake
type ServerHelloMessage struct {
	chosenEncryption string
	chosenSignature  SignScheme
}

func (msg *ServerHelloMessage) marshall() ([]byte, error) {
	encodedMsg := bytes.Buffer{}
	encodedMsg.WriteByte(typeServerHello)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, len(msg.chosenEncryption)); err != nil {
		return nil, err
	}
	encodedMsg.WriteString(msg.chosenEncryption)
	encodedMsg.WriteByte(byte(msg.chosenSignature))
	return encodedMsg.Bytes(), nil
}

// DecodeServerMsg - unmarshall []byte to ServerHelloMessage
func DecodeServerMsg(msg []byte) (*ServerHelloMessage, error) {
	if msg[0] != typeServerHello {
		return nil, ErrWrongFormat
	}

	reader := bytes.NewReader(msg[1:])
	var chosenEncryptionLen int
	if err := binary.Read(reader, binary.LittleEndian, &chosenEncryptionLen); err != nil {
		return nil, err
	}

	byteString := make([]byte, chosenEncryptionLen)
	if _, err := reader.Read(byteString); err != nil {
		return nil, err
	}

	binarySignature, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ServerHelloMessage{
		chosenEncryption: string(byteString),
		chosenSignature:  SignScheme(binarySignature),
	}, nil
}

// ServerHandshakeMessage - second server message during handshake process
type ServerHandshakeMessage struct {
	pastelID         []byte
	signedPastelID   []byte
	pubKey           []byte
	ctx              []byte
	encryptionParams [][]byte
}

func (msg *ServerHandshakeMessage) marshall() ([]byte, error) {
	var pastelIDLen = len(msg.pastelID)
	var signedPastelIDLen = len(msg.signedPastelID)
	var pubKeyLen = len(msg.pubKey)
	var ctxLen = len(msg.pubKey)
	var encryptionParamsLen = len(msg.encryptionParams)

	encodedMsg := bytes.Buffer{}
	encodedMsg.WriteByte(typeClientHandshakeMsg)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, pastelIDLen); err != nil {
		return nil, err
	}
	encodedMsg.Write(msg.pastelID)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, signedPastelIDLen); err != nil {
		return nil, err
	}
	encodedMsg.Write(msg.signedPastelID)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, pubKeyLen); err != nil {
		return nil, err
	}
	encodedMsg.Write(msg.pubKey)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, ctxLen); err != nil {
		return nil, err
	}
	encodedMsg.Write(msg.ctx)
	if err := binary.Write(&encodedMsg, binary.LittleEndian, encryptionParamsLen); err != nil {
		return nil, err
	}
	// write encryption data
	for _, paramData := range msg.encryptionParams {
		paramDataLen := len(paramData)
		if err := binary.Write(&encodedMsg, binary.LittleEndian, paramDataLen); err != nil {
			return nil, err
		}
		encodedMsg.Write(paramData)
	}

	return encodedMsg.Bytes(), nil
}

// DecodeServerHandshakeMsg - unmarshall []byte to ServerHandshakeMessage
func DecodeServerHandshakeMsg(msg []byte) (*ServerHandshakeMessage, error) {
	if msg[0] != typeClientHandshakeMsg {
		return nil, ErrWrongFormat
	}

	pastelID := make([]byte, msg[1])
	copy(pastelID, msg[2:])

	shift := 2 + msg[1] // msg type + len and array itself
	signedPastelID := make([]byte, msg[shift])
	copy(signedPastelID, msg[shift+1:])

	shift += 1 + msg[1] // len and array itself
	pubKey := make([]byte, msg[shift])
	copy(pubKey, msg[shift+1:])

	shift += 1 + msg[shift] // len and array itself
	ctx := make([]byte, msg[shift])
	copy(ctx, msg[shift+1:])

	shift += 1 + msg[shift] // len and array itself
	aes256key := make([]byte, msg[shift])
	copy(aes256key, msg[shift+1:])

	return &ServerHandshakeMessage{
		pastelID:       pastelID,
		signedPastelID: signedPastelID,
		pubKey:         pubKey,
		ctx:            ctx,
		aes256key:      aes256key,
	}, nil
}
