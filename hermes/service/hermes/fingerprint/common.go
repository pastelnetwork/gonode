package fingerprint

func typeMapper(val string) string {
	if val == "action-reg" {
		return "SENSE"
	}
	if val == "nft-reg" {
		return "NFT"
	}

	return val
}
