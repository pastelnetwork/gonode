package node

import (
	"context"
)

type SuperNodeInterface interface {
	String() string
	PastelID() string
	Address() string
	SetPrimary(primary bool)
	IsPrimary() bool
	SetActive(active bool)
	IsActive() bool
	SetRemoteStatus() bool
	IsRemoteStatus() bool
	RLock()
	RUnlock()
}

type SuperNodeCollectionInterface interface {
	AddNewNode(client ClientInterface, address string, pastelID string)
	Activate()
	DisconnectInactive()
	DisconnectAll()
	WaitConnClose(ctx context.Context, done <-chan struct{}) error
}
