package pastel

import "context"

// Client represents pastel RPC client.
type Client interface {
	TopMasterNodes(ctx context.Context) (MasterNodes, error)
	Getblockchaininfo(ctx context.Context) (*BlockchainInfo, error)
	ListIDTickets(ctx context.Context, idType string) (IDTickets, error)
	FindIDTicket(ctx context.Context, search string) (*IDTicket, error)
	FindIDTickets(ctx context.Context, search string) (IDTickets, error)
	ListPastelIDs(ctx context.Context) (PastelIDs, error)
	GetMNRegFee(ctx context.Context) (int, error)
}
