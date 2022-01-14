package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"time"
)

type NodeInterface interface {
	String() (string)
	PastelID() (string)
	Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) (error)
	Address() (string)
	SetPrimary(primary bool)
	IsPrimary() (bool)
}
