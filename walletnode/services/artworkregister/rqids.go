package artworkregister

// RQIDSList is list of RQIDS
type RQIDSList []*RQIDS

// RQIDS define a symbol id file
type RQIDS struct {
	// Id = bash64(hash(file + signature))
	Id string

	// Content is the byte serialization of the RaptorQ symbols file in the following format
	// RANDOM-GUID1..X
	// BLOCK_HASH
	// raptroq id1
	// ...
	// raptroq idN
	// signature
	Content []byte
}

// Returns a list of identifiers
func (rqidsList RQIDSList) Identifiers() []string {
	var identifier []string
	for _, rqids := range rqidsList {
		identifier = append(identifier, rqids.Id)
	}

	return identifier
}

// Returns a map from symbold id file id to it's content
func (rqidsList RQIDSList) ToMap() map[string][]byte {
	m := make(map[string][]byte, len(rqidsList))
	for _, rqids := range rqidsList {
		m[rqids.Id] = rqids.Content
	}
	return m
}
