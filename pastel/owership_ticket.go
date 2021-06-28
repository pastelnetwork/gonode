package pastel

// OwnershipTicket represents ticket's ownership
type OwnershipTicket struct {
	Art   string `json:"art"`   // txid from the request
	Trade string `json:"trade"` // txid from trade ticket
}
