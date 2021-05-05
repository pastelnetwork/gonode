package net

type Handshaker interface {
	IsClient() bool
	ClientHello() ([]byte, error)           // returns artistID
	ServerHello() ([]byte, error)           // returns OK
	ClientKeyExchange() ([]byte, error)     // returns signature+publicKey
	ServerKeyVerify([]byte) (Cipher, error) // verify client's signature
	ServerKeyExchange() ([]byte, error)     // returns signature+publicKey
	ClientKeyVerify([]byte) (Cipher, error) // verify serverMock's signature
}

type Cipher interface {
	Encrypt(dst, plaintext []byte) ([]byte, error)
	MaxOverhead() int
	Decrypt(dst, ciphertext []byte) ([]byte, error)
}
