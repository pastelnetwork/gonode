package pastel

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

// GetTransactionDetailsResult represents a transaction details
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

// GetBlockVerbose1Result models the data from the getblock command when the
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
	Nonce            string   `json:"nonce"`
	Bits             string   `json:"bits"`
	Difficulty       float64  `json:"difficulty"`
	PreviousHash     string   `json:"previousblockhash"`
	NextHash         string   `json:"nextblockhash,omitempty"`
}

// Amount sent to Address
type Amount struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

// GetOperationStatusParam describes  the params of operation
type GetOperationStatusParam struct {
	FromAddress string   `json:"from_address"`
	Amounts     []Amount `json:"amounts"`
	MinConf     int      `json:"min_conf"`
	Fee         float64  `json:"fee"`
}

// GetOperationStatusTxid describes the result transaction id of the operation
type GetOperationStatusTxid struct {
	Txid string `json:"txid"`
}

// GetOperationStatusError describes the reason why the operation failed
type GetOperationStatusError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// GetOperationStatusResult describes the result of an operation identified by Id
type GetOperationStatusResult struct {
	ID            string                  `json:"id,omitempty"`
	Status        string                  `json:"status,omitempty"`
	Error         GetOperationStatusError `json:"error,omitempty"`
	Result        GetOperationStatusTxid  `json:"result"`
	ExecutionSecs float64                 `json:"execution_secs"`
	Method        string                  `json:"method"`
	Params        GetOperationStatusParam `json:"params"`
	MinConf       int                     `json:"min_conf"`
	Fee           float64                 `json:"fee"`
}

// GetRawTransactionVerbose1Result describes the result of "getrawtranasction txid 1"
type GetRawTransactionVerbose1Result struct {
	// Other information are omitted here because we don't care
	// They can be added later if there are business requirements
	Txid          string          `json:"txid"`
	Confirmations int64           `json:"confirmations"`
	Vout          []VoutTxnResult `json:"vout"`
	Time          int64           `json:"time"`
}

// VoutTxnResult is detail of txn
type VoutTxnResult struct {
	Value        float64      `json:"value"`
	ValuePat     int64        `json:"valuePat"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// ScriptPubKey lists addresses of txn
type ScriptPubKey struct {
	Addresses []string `json:"addresses"`
}

// GetActionFeesResult describes the result of "getactionfee getactionfees"
type GetActionFeesResult struct {
	SenseFee   float64 `json:"sensefee"`
	CascadeFee float64 `json:"cascadefee"`
}

// GetTotalBalanceResponse ...
type GetTotalBalanceResponse struct {
	Transparent float64 `json:"transparent"`
	Private     float64 `json:"private"`
	Total       float64 `json:"total"`
}

// NFTStorageFeeEstimate ...
type NFTStorageFeeEstimate struct {
	EstimatedNftStorageFeeMin     float64 `json:"estimatedNftStorageFeeMin"`
	EstimatedNftStorageFeeAverage float64 `json:"estimatedNftStorageFeeAverage"`
	EstimatedNftStorageFeeMax     float64 `json:"estimatedNftStorageFeeMax"`
}
