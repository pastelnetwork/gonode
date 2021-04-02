package models

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Verificationprogress float64 `json:"verificationprogress"`
	Chainwork            string  `json:"chainwork"`
	Pruned               bool    `json:"pruned"`
	Commitments          int     `json:"commitments"`
	ValuePools           []struct {
		ID            string  `json:"id"`
		Monitored     bool    `json:"monitored"`
		ChainValue    float64 `json:"chainValue"`
		ChainValueZat int     `json:"chainValueZat"`
	} `json:"valuePools"`
	Softforks []struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
		Enforce struct {
			Status   bool `json:"status"`
			Found    int  `json:"found"`
			Required int  `json:"required"`
			Window   int  `json:"window"`
		} `json:"enforce"`
		Reject struct {
			Status   bool `json:"status"`
			Found    int  `json:"found"`
			Required int  `json:"required"`
			Window   int  `json:"window"`
		} `json:"reject"`
	} `json:"softforks"`
	Upgrades struct {
		FiveBa81B19 struct {
			Name             string `json:"name"`
			Activationheight int    `json:"activationheight"`
			Status           string `json:"status"`
			Info             string `json:"info"`
		} `json:"5ba81b19"`
		Seven6B809Bb struct {
			Name             string `json:"name"`
			Activationheight int    `json:"activationheight"`
			Status           string `json:"status"`
			Info             string `json:"info"`
		} `json:"76b809bb"`
	} `json:"upgrades"`
	Consensus struct {
		Chaintip  string `json:"chaintip"`
		Nextblock string `json:"nextblock"`
	} `json:"consensus"`
}
