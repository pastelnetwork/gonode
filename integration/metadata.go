package integration

const (
	// SN1BaseURI of SN1 Server
	SN1BaseURI = "http://localhost:19090"
	// SN2BaseURI of SN2 Server
	SN2BaseURI = "http://localhost:19091"
	// SN3BaseURI of SN3 Server
	SN3BaseURI = "http://localhost:19092"

	// WNBaseURI of WN
	WNBaseURI = "http://localhost:18089"
)

// RetrieveResponse indicates response structure of retrieve request
type RetrieveResponse struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// StoreRequest indicates request structure of store request
type StoreRequest struct {
	Value []byte `json:"value"`
}

// StoreReply indicates reply structure of store request
type StoreReply struct {
	Key string `json:"key"`
}
