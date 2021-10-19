package externalstorage

// Ticket represents artwork registration ticket.
type Ticket struct {
	// The Request_Identifier of the corresponding External_Dupe_Detection_Request
	RequestIdentifier string
	// The SHA3-256 hash of the image file the user wants to run dupe-detection on
	ImageHash string
	// The Maximum_PSL_Fee_for_External_Dupe_Detection
	MaximumFee float64
	// The External_Dupe_Detection_Sending_Address
	SendingAddress string
	// The effective total fee paid by the user in PSL (both paid to Supernodes and burned)
	TotalFee float64
	// The SHA3-256 hash for the JSON file stored in Kademlia which itself contains the SHA3-256 hashes of the set
	// of redundant copies of the JSON files in Kademlia (which contain the results of the dupe detection computation)
	ResultHash string
	// The Pastel signatures of the Supernodes that performed the dupe detection computations for the ticket on the above data
	SuperNodeSignarure [][]byte
	// The Pastel signature of the user on the above data
	Signature []byte
	// The Pastel TXID of this transaction would then serve as the identifier for this ticket
	TxnID string
}
