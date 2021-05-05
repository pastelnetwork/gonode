package api

// Request represents an API request.
type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  []string    `json:"params"`
}
