package qrsignature

// List of available payloads
const (
	Fingerprint DataType = iota
	PostQuantumSignature
	PostQuantumPubKey
	Ed448Signature
	Ed448PubKey
)

var dataTypeNames = map[DataType]string{
	Fingerprint:          "Fingerprint",
	PostQuantumSignature: "PQ Signature",
	PostQuantumPubKey:    "PQ Pub Key",
	Ed448Signature:       "Ed448 Signature",
	Ed448PubKey:          "Ed448 Pub Key",
}

// DataType represents data that can be stored inside QR code.
type DataType byte

func (dataType DataType) String() string {
	if name, ok := dataTypeNames[dataType]; ok {
		return name
	}
	return ""
}

// NewFingerprint returns a new instance of Data with preset Fingerpirnt datatype.
func NewFingerprint(raw []byte) *Data {
	return NewData(raw, Fingerprint)
}

// NewPostQuantumData returns a new instance of Data with preset PostQuantumData datatype.
func NewPostQuantumData(raw []byte) *Data {
	return NewData(raw, PostQuantumSignature)
}

// NewPostQuantumPubKey returns a new instance of Data with preset PostQuantumPubKey datatype.
func NewPostQuantumPubKey(raw []byte) *Data {
	return NewData(raw, PostQuantumPubKey)
}

// NewEd448Data returns a new instance of Data with preset Ed448Data datatype.
func NewEd448Data(raw []byte) *Data {
	return NewData(raw, Ed448Signature)
}

// NewEd448PubKey returns a new instance of Data with preset Ed448PubKey datatype.
func NewEd448PubKey(raw []byte) *Data {
	return NewData(raw, Ed448PubKey)
}
