package qrsignature

// List of available payloads
const (
	PayloadFingerprint PayloadName = iota + 1
	PayloadPostQuantumSignature
	PayloadPostQuantumPubKey
	PayloadEd448Signature
	PayloadEd448PubKey
	PayloadBlockHash
)

var payloadNames = map[PayloadName]string{
	PayloadFingerprint:          "Fingerprint",
	PayloadPostQuantumSignature: "PQ Signature",
	PayloadPostQuantumPubKey:    "PQ Pub Key",
	PayloadEd448Signature:       "Ed448 Signature",
	PayloadEd448PubKey:          "Ed448 Pub Key",
	PayloadBlockHash:            "Block Hash",
}

func payloadName(subText string) PayloadName {
	for name, text := range payloadNames {
		if text == subText {
			return name
		}
	}
	return PayloadName(0)
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

// BlockHash returns a new instance of Payload with preset BlockHash datatype.
func BlockHash(raw []byte) *Payload {
	return NewPayload(raw, PayloadBlockHash)
}

// Fingerprint returns Fingerprint data.
func (payloads *Payloads) Fingerprint() []byte {
	return payloads.Raw(PayloadFingerprint)
}

// BlockHash returns BlockHash data.
func (payloads *Payloads) BlockHash() []byte {
	return payloads.Raw(PayloadBlockHash)
}

// PostQuantumSignature returns PostQuantumSignature data.
func (payloads *Payloads) PostQuantumSignature() []byte {
	return payloads.Raw(PayloadPostQuantumSignature)
}

// PostQuantumPubKey returns PostQuantumPubKey data.
func (payloads *Payloads) PostQuantumPubKey() []byte {
	return payloads.Raw(PayloadPostQuantumPubKey)
}

// Ed448Signature returns Ed448Signature data.
func (payloads *Payloads) Ed448Signature() []byte {
	return payloads.Raw(PayloadEd448Signature)
}

// Ed448PubKey returns Ed448PubKey data.
func (payloads *Payloads) Ed448PubKey() []byte {
	return payloads.Raw(PayloadEd448PubKey)
}
