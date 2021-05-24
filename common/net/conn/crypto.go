package conn

// Crypto interface that should use in handshake and ecrytion process
type Crypto interface {
	Handshake() error
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	About() EncryptionScheme
}

// Handshake interface that should during handshake establishment
type Handshake interface {
	Handshake() error
	IsHandshakeEstablished() bool
	GetCryptoProtocol() EncryptionScheme
}

// ToDo: Implement whole crypto

// CryptoImpl structure that implement all methods of Crypto interface
type CryptoImpl struct {
	handshakeImpl Handshake
}

// NewClientCrypto - creates CryptoImpl
func NewClientCrypto(conn Conn) *CryptoImpl {
	client := NewClient(conn)
	return &CryptoImpl{
		handshakeImpl: client,
	}
}

// Encrypt - encrypt the data
func (crypto *CryptoImpl) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

// Decrypt - decrypt the data
func (crypto *CryptoImpl) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

// About - return information about encryption scheme
func (crypto *CryptoImpl) About() EncryptionScheme {
	return crypto.handshakeImpl.GetCryptoProtocol()
}

// Handshake process itself
func (crypto *CryptoImpl) Handshake() error {
	if crypto.handshakeImpl.IsHandshakeEstablished() {
		return nil
	}

	return crypto.handshakeImpl.Handshake()
}
