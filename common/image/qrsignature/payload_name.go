package qrsignature

// List of available payloads
const (
	PayloadFingerprint PayloadName = iota
	PayloadPostQuantumSignature
	PayloadPostQuantumPubKey
	PayloadEd448Signature
	PayloadEd448PubKey
)

var payloadNames = map[PayloadName]string{
	PayloadFingerprint:          "Fingerprint",
	PayloadPostQuantumSignature: "PQ Signature",
	PayloadPostQuantumPubKey:    "PQ Pub Key",
	PayloadEd448Signature:       "Ed448 Signature",
	PayloadEd448PubKey:          "Ed448 Pub Key",
}

// PayloadName represents data that can be stored inside QR code.
type PayloadName byte

func (dataType PayloadName) String() string {
	if name, ok := payloadNames[dataType]; ok {
		return name
	}
	return ""
}

// Fingerprint returns a new instance of Payload with preset Fingerpirnt datatype.
func Fingerprint(raw []byte) *Payload {
	return NewPayload(raw, PayloadFingerprint)
}

// PostQuantumSignature returns a new instance of Payload with preset PostQuantumData datatype.
func PostQuantumSignature(raw []byte) *Payload {
	return NewPayload(raw, PayloadPostQuantumSignature)
}

// PostQuantumPubKey returns a new instance of Payload with preset PayloadPostQuantumPubKey datatype.
func PostQuantumPubKey(raw []byte) *Payload {
	return NewPayload(raw, PayloadPostQuantumPubKey)
}

// Ed448Signature returns a new instance of Payload with preset Ed448Data datatype.
func Ed448Signature(raw []byte) *Payload {
	return NewPayload(raw, PayloadEd448Signature)
}

// Ed448PubKey returns a new instance of Payload with preset PayloadEd448PubKey datatype.
func Ed448PubKey(raw []byte) *Payload {
	return NewPayload(raw, PayloadEd448PubKey)
}
