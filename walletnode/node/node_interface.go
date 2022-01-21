package node

import (
	"context"
)

type SuperNodeInterface interface {
	String() string
	PastelID() string
	Address() string
	SetPrimary(bool)
	IsPrimary() bool
	SetActive(bool)
	IsActive() bool
	SetRemoteState(bool)
	IsRemoteState() bool
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
