package conn

type Crypto interface {
	Handshake() error
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	About() EncryptionScheme
}

type Handshake interface {
	Handshake() error
	IsHandshakeEstablished() bool
	GetCryptoProtocol() EncryptionScheme
}

// ToDo: Implement whole crypto

type CryptoImpl struct {
	handshakeImpl Handshake
}

func NewClientCrypto(conn Conn) *CryptoImpl {
	client := NewClient(conn)
	return &CryptoImpl{
		handshakeImpl: client,
	}
}

func (crypto *CryptoImpl) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (crypto *CryptoImpl) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (crypto *CryptoImpl) About() EncryptionScheme {
	return crypto.handshakeImpl.GetCryptoProtocol()
}

func (crypto *CryptoImpl) Handshake() error {
	if crypto.handshakeImpl.IsHandshakeEstablished() {
		return nil
	}

	return crypto.handshakeImpl.Handshake()
}
