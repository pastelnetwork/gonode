package conn

type ClientHelloMessage struct {
	supportedEncryptions         []EncryptionScheme
	supportedSignatureAlgorithms []SignScheme
}

func (msg *ClientHelloMessage) marshall() []byte {
	var supportedEncryptionsCount = len(msg.supportedEncryptions)
	var supportedSignatureAlgorithmsCount = len(msg.supportedSignatureAlgorithms)

	// prepare structure to send a message, it should keep supportedEncryptions and supportedSignatureAlgorithms within theirs length and type of msg
	encodedMsg := make([]byte, supportedEncryptionsCount+supportedSignatureAlgorithmsCount+3)

	encodedMsg[0] = typeClientHello
	encodedMsg[1] = byte(supportedEncryptionsCount) // store len of supportedEncryptions
	shift := 2
	copy(encodedMsg[shift:], string(msg.supportedEncryptions))

	shift += supportedSignatureAlgorithmsCount
	encodedMsg[shift] = byte(supportedSignatureAlgorithmsCount)
	copy(encodedMsg[shift:], string(msg.supportedSignatureAlgorithms))

	return encodedMsg
}

func DecodeClientMsg(msg []byte) (*ClientHelloMessage, error) {
	if msg[0] != typeClientHello {
		return nil, WrongFormatErr
	}
	var supportedEncryptionsCount = msg[1]
	supportedEncryptions := make([]EncryptionScheme, supportedEncryptionsCount)

	shift := 2
	for i := 0; i < int(supportedEncryptionsCount); i++ {
		supportedEncryptions[i] = EncryptionScheme(msg[shift+i])
	}

	shift += int(supportedEncryptionsCount)
	var supportedSignatureAlgorithmsCount = msg[shift]
	supportedSignatureAlgorithms := make([]SignScheme, supportedEncryptionsCount)
	for i := 0; i < int(supportedSignatureAlgorithmsCount); i++ {
		supportedSignatureAlgorithms[i] = SignScheme(msg[shift+i])
	}

	return &ClientHelloMessage{
		supportedEncryptions:         supportedEncryptions,
		supportedSignatureAlgorithms: supportedSignatureAlgorithms,
	}, nil
}

type ClientHandshakeMessage struct {
	pastelId       []byte
	signedPastelId []byte
	pubKey         []byte
	ctx            []byte
}

func (msg *ClientHandshakeMessage) marshall() []byte {
	var pastelIdLen = len(msg.signedPastelId)
	var signedPastelIdLen = len(msg.signedPastelId)
	var clientPubKeyLen = len(msg.pubKey)
	var ctxLen = len(msg.pubKey)

	encodedMsg := make([]byte, 4+pastelIdLen+signedPastelIdLen+clientPubKeyLen)

	encodedMsg[0] = typeClientHandshakeMsg
	encodedMsg[1] = byte(pastelIdLen)
	copy(encodedMsg[2:], msg.pastelId)

	shift := 2 + pastelIdLen // msg type + len and array itself
	encodedMsg[shift] = byte(signedPastelIdLen)
	copy(encodedMsg[shift+1:], msg.signedPastelId)

	shift += 1 + signedPastelIdLen // len and array itself
	encodedMsg[shift] = byte(clientPubKeyLen)
	copy(encodedMsg[shift+1:], msg.pubKey)

	shift += 1 + clientPubKeyLen // len and array itself
	encodedMsg[shift] = byte(ctxLen)
	copy(encodedMsg[shift+1:], msg.ctx)

	return encodedMsg
}

func DecodeClientHandshakeMessage(msg []byte) (*ClientHandshakeMessage, error) {
	if msg[0] != typeClientHandshakeMsg {
		return nil, WrongFormatErr
	}

	signedPastelId := make([]byte, msg[1])

	copy(signedPastelId, msg[2:])

	shift := 2 + msg[1] // msg type + len and array itself
	pubKey := make([]byte, msg[shift])
	copy(pubKey, msg[shift+1:])

	shift += 1 + msg[shift] // len and array itself
	ctx := make([]byte, msg[shift])
	copy(ctx, msg[shift+1:])

	return &ClientHandshakeMessage{
		signedPastelId: signedPastelId,
		pubKey:         pubKey,
		ctx:            ctx,
	}, nil
}
