package pastel

import "context"

// Client represents pastel RPC client.
type Client interface {
	MasterNodesTop(ctx context.Context) (MasterNodes, error)
	MasterNodeStatus(ctx context.Context) (*MasterNodeStatus, error)

	StorageFee(ctx context.Context) (*StorageFee, error)

	IDTickets(ctx context.Context, idType IDTicketType) (IDTickets, error)
}

type IDTicketType string

const (
	IDTicketAll      IDTicketType = "all"
	IDTicketMine     IDTicketType = "mine"
	IDTicketMN       IDTicketType = "mn"
	IDTicketPersonal IDTicketType = "personal"
)
