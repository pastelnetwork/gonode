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
	Nonce            string   `json:"nonce"`
	Bits             string   `json:"bits"`
	Difficulty       float64  `json:"difficulty"`
	PreviousHash     string   `json:"previousblockhash"`
	NextHash         string   `json:"nextblockhash,omitempty"`
}

// {
//     "id": "opid-39c1cdc7-6dbc-4ad9-be0a-4693918d5d77",
//     "status": "success",
//     "creation_time": 1626973263,
//     "result": {
//       "txid": "8ec1c6f5c4ef15b749191f5cb0d95b8c4de9a497842be8690c5285e0adab2ec9"
//     },
//     "execution_secs": 0.003088334,
//     "method": "z_sendmany",
//     "params": {
//       "fromaddress": "tPWLZyFoXe5EdQPhV2G7Y45gZGRQzn1wRHp",
//       "amounts": [
//         {
//           "address": "tPce2T47TFPcHPj3sKiKdFLhQpmzJ1tkzjr",
//           "amount": 6.12
//         }
//       ],
//       "minconf": 1,
//       "fee": 0.1
//     }
//   }
//

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

type GetOperationStatusResult struct {
	Id            string                  `json:"id,omitempty"`
	Status        string                  `json:"status,omitemptyi"`
	Error         GetOperationStatusError `json:"error,omitempty"`
	Result        GetOperationStatusTxid  `json:"result"`
	ExecutionSecs float64                 `json:"execution_secs"`
	Method        string                  `json:"method"`
	Params        GetOperationStatusParam `json:"params"`
	MinConf       int                     `json:"min_conf"`
	Fee           float64                 `json:"fee"`
}
