package node

import (
	"context"
)

type NodeInterface interface {
	String() string
	PastelID() string
	Address() string
	SetPrimary(primary bool)
	IsPrimary() bool
	SetActive(active bool)
	IsActive() bool
	RLock()
	RUnlock()
}

type NodeCollectionInterface interface {
	Add(node *NodeInterface)
	Activate()
	DisconnectInactive()
	DisconnectAll()
	WaitConnClose(ctx context.Context, done <-chan struct{}) error
	FindByPastelID(id string) *NodeInterface
}
