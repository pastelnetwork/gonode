package conn

// ServerHelloMessage - first server message during handshake
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

// DecodeServerMsg - unmarshall []byte to ServerHelloMessage
func DecodeServerMsg(msg []byte) (*ServerHelloMessage, error) {
	if msg[0] != typeServerHello {
		return nil, ErrWrongFormat
	}

	publicKey := make([]byte, msg[3])

	copy(publicKey, msg[4:])
	return &ServerHelloMessage{
		chosenEncryption: EncryptionScheme(msg[1]),
		chosenSignature:  SignScheme(msg[2]),
	}, nil
}

// ServerHandshakeMessage - second server message during handshake process
type ServerHandshakeMessage struct {
	pastelID       []byte
	signedPastelID []byte
	pubKey         []byte
	ctx            []byte
	aes256key      []byte
}

func (msg *ServerHandshakeMessage) marshall() []byte {
	var pastelIDLen = len(msg.pastelID)
	var signedPastelIDLen = len(msg.signedPastelID)
	var pubKeyLen = len(msg.pubKey)
	var ctxLen = len(msg.pubKey)
	var aes265keyLen = len(msg.aes256key) // should be 32 byte

	encodedMsg := make([]byte, 4+pastelIDLen+signedPastelIDLen+pubKeyLen)

	encodedMsg[0] = typeClientHandshakeMsg
	encodedMsg[1] = byte(pastelIDLen)
	copy(encodedMsg[2:], msg.pastelID)

	shift := 2 + pastelIDLen // msg type + len and array itself
	encodedMsg[shift] = byte(signedPastelIDLen)
	copy(encodedMsg[2:], msg.signedPastelID)

	shift += 1 + signedPastelIDLen // msg type + len and array itself
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
