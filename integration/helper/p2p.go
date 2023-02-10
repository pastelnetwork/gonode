package helper

// RetrieveResponse indicates response structure of retrieve request
type RetrieveResponse struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// StoreRequest indicates request structure of store request
type StoreRequest struct {
	Value []byte `json:"value"`
}

// LocalStoreRequest indicates request structure of store request
type LocalStoreRequest struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// StoreReply indicates reply structure of store request
type StoreReply struct {
	Key string `json:"key"`
}

// DeleteRequest indicates request structure of store request
type DeleteRequest struct {
	Key string `json:"key"`
}
