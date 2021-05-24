package conn

type ServerHelloMessage struct {
	chosenEncryption EncryptionScheme
	chosenSignature  SignScheme
	publicKey        []byte
}

func (msg *ServerHelloMessage) marshall() []byte {
	// prepare structure to send a message, it should keep supportedEncryptions and supportedSignatureAlgorithms within theirs length
	encodedMsg := make([]byte, 4+len(msg.publicKey))

	encodedMsg[0] = typeServerHello
	encodedMsg[1] = byte(msg.chosenEncryption) // store len of supportedEncryptions
	encodedMsg[2] = byte(msg.chosenSignature)
	encodedMsg[3] = byte(len(msg.publicKey))

	copy(encodedMsg[4:], msg.publicKey)
	return encodedMsg
}

func DecodeServerMsg(msg []byte) (*ServerHelloMessage, error) {
	if msg[0] != typeServerHello {
		return nil, WrongFormatErr
	}

	publicKey := make([]byte, msg[3])

	copy(publicKey, msg[4:])
	return &ServerHelloMessage{
		chosenEncryption: EncryptionScheme(msg[1]),
		chosenSignature:  SignScheme(msg[2]),
	}, nil
}

type ServerHandshakeMessage struct {
	pastelId       []byte
	signedPastelId []byte
	pubKey         []byte
	ctx            []byte
	aes256key      []byte
}

func (msg *ServerHandshakeMessage) marshall() []byte {
	var pastelIdLen = len(msg.pastelId)
	var signedPastelIdLen = len(msg.signedPastelId)
	var pubKeyLen = len(msg.pubKey)
	var ctxLen = len(msg.pubKey)
	var aes265keyLen = len(msg.aes256key) // should be 32 byte

	encodedMsg := make([]byte, 4+pastelIdLen+signedPastelIdLen+pubKeyLen)

	encodedMsg[0] = typeClientHandshakeMsg
	encodedMsg[1] = byte(pastelIdLen)
	copy(encodedMsg[2:], msg.pastelId)

	shift := 2 + pastelIdLen // msg type + len and array itself
	encodedMsg[shift] = byte(signedPastelIdLen)
	copy(encodedMsg[2:], msg.signedPastelId)

	shift += 1 + signedPastelIdLen // msg type + len and array itself
	encodedMsg[shift] = byte(pubKeyLen)
	copy(encodedMsg[shift+1:], msg.pubKey)

	shift += 1 + pubKeyLen // len and array itself
	encodedMsg[shift] = byte(ctxLen)
	copy(encodedMsg[shift+1:], msg.ctx)

	shift += 1 + ctxLen // len and array itself
	encodedMsg[shift] = byte(aes265keyLen)
	copy(encodedMsg[shift+1:], msg.aes256key)

	return encodedMsg
}

func DecodeServerHandshakeMsg(msg []byte) (*ServerHandshakeMessage, error) {
	if msg[0] != typeClientHandshakeMsg {
		return nil, WrongFormatErr
	}

	pastelId := make([]byte, msg[1])
	copy(pastelId, msg[2:])

	shift := 2 + msg[1] // msg type + len and array itself
	signedPastelId := make([]byte, msg[shift])
	copy(signedPastelId, msg[shift+1:])

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
		pastelId:       pastelId,
		signedPastelId: signedPastelId,
		pubKey:         pubKey,
		ctx:            ctx,
		aes256key:      aes256key,
	}, nil
}
