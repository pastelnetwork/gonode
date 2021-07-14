package pastel

// TxIDType represents a type of transaction id
type TxIDType string

// GetInfoResult models the data returned by the server getinfo
// command.
// Refer https://github.com/pastelnetwork/pastel/blob/d8c0520c590693f4282f72e8553d79beeaf9da27/src/rpc/misc.cpp#L45
type GetInfoResult struct {
	Version         int32   `json:"version"`
	ProtocolVersion int32   `json:"protocolversion"`
	WalletVersion   int32   `json:"walletversion"`
	Balance         float64 `json:"balance"`
	Blocks          int32   `json:"blocks"`
	TimeOffset      int64   `json:"timeoffset"`
	Connections     int32   `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float64 `json:"difficulty"`
	TestNet         bool    `json:"testnet"`
	KeypoolOldest   int64   `json:"keypoololdest"`
	KeypoolSize     int32   `json:"keypoolsize"`
	UnlockedUntil   int64   `json:"unlocked_until"`
	PaytxFee        float64 `json:"paytxfee"`
	RelayFee        float64 `json:"relayfee"`
	Errors          string  `json:"errors"`
}

type GetTransactionDetailsResult struct {
	Account           string   `json:"account"`
	Address           string   `json:"address,omitempty"`
	Amount            float64  `json:"amount"`
	Category          string   `json:"category"`
	InvolvesWatchOnly bool     `json:"involveswatchonly,omitempty"`
	Fee               *float64 `json:"fee,omitempty"`
	Vout              uint32   `json:"vout"`
}

// GetTransactionResult models the data from the gettransaction command.
type GetTransactionResult struct {
	Amount          float64                       `json:"amount"`
	Fee             float64                       `json:"fee,omitempty"`
	Confirmations   int64                         `json:"confirmations"`
	BlockHash       string                        `json:"blockhash"`
	BlockIndex      int64                         `json:"blockindex"`
	BlockTime       int64                         `json:"blocktime"`
	TxID            string                        `json:"txid"`
	WalletConflicts []string                      `json:"walletconflicts"`
	Time            int64                         `json:"time"`
	TimeReceived    int64                         `json:"timereceived"`
	Details         []GetTransactionDetailsResult `json:"details"`
	Hex             string                        `json:"hex"`
}

// GetBlockVerboseResult models the data from the getblock command when the
// verbose flag is set to 1 - getblock returns an object
// whose tx field is an array of transaction hashes
type GetBlockVerbose1Result struct {
	Hash             string   `json:"hash"`
	Confirmations    int64    `json:"confirmations"`
	Size             int32    `json:"size"`
	Height           int64    `json:"height"`
	Version          int32    `json:"version"`
	MerkleRoot       string   `json:"merkleroot"`
	FinalSaplingRoot string   `json:"finalsaplingroot"`
	Tx               []string `json:"tx,omitempty"`
	Time             int64    `json:"time"`
	Nonce            uint32   `json:"nonce"`
	Bits             string   `json:"bits"`
	Difficulty       float64  `json:"difficulty"`
	PreviousHash     string   `json:"previousblockhash"`
	NextHash         string   `json:"nextblockhash,omitempty"`
}

func (t TxIDType) Valid() bool {
	return t != ""
}
