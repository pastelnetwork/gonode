package ed448

import "errors"

// ErrWrongFormat - error when structure can't be parsed
var ErrWrongFormat = errors.New("unknown message format")

// ErrWrongSignature - error when signature is incorrect
var ErrWrongSignature = errors.New("wrong signature")

// ErrIncorrectPastelID - error when pastel ID is wrong
var ErrIncorrectPastelID = errors.New("incorrect Pastel Id")

// ErrUnsupportedEncryption - error when encryption doesn't supported by server or client
var ErrUnsupportedEncryption = errors.New("unsupported encryption")

// ErrMaxHandshakes - error when number of concurrent handshakes reached the limit
var ErrMaxHandshakes = errors.New("maximum number of concurrent handshakes is reached")
