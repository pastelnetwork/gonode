//go:generate mockery --name=PastelClient

package test

import (
	"context"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test/mocks"
	"github.com/stretchr/testify/mock"
)

// TODO: Ashadi
// PastelClient mock pastel.Client. redefine pastel.Client interface
// go:generate does not work with accros module on current CircleCI build
// not sure how to fix it
type PastelClient interface {
	// MasterNodesTop returns a result of the `masternode top`.
	MasterNodesTop(ctx context.Context) (pastel.MasterNodes, error)

	// MasterNodeStatus returns a result of the `masternode status`.
	MasterNodeStatus(ctx context.Context) (*pastel.MasterNodeStatus, error)

	// MasterNodeConfig returns a result of the `masternode list-conf`.
	MasterNodeConfig(ctx context.Context) (*pastel.MasterNodeConfig, error)

	// StorageFee returns a result of the `storagefee getnetworkfee`.
	StorageFee(ctx context.Context) (*pastel.StorageFee, error)

	// IDTickets returns a result of the `tickets list id`.
	IDTickets(ctx context.Context, idType pastel.IDTicketType) (pastel.IDTickets, error)
}

// PastelClientMock implementing PastelClient for testing purpose
type PastelClientMock struct {
	ClientMock *mocks.PastelClient
}

// NewMockPastelClient create new pastel client mock
func NewMockPastelClient() *PastelClientMock {
	return &PastelClientMock{
		ClientMock: &mocks.PastelClient{},
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returning Mn's and error from args
func (p *PastelClientMock) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *PastelClientMock {
	p.ClientMock.On("MasterNodesTop", mock.Anything).Return(nodes, err)
	return p
}

// ListenOnStorageFee listening StorageFee call and returning pastel.StorageFee, error form args
func (p *PastelClientMock) ListenOnStorageFee(fee *pastel.StorageFee, returnErr error) *PastelClientMock {
	p.ClientMock.On("StorageFee", mock.Anything).Return(fee, returnErr)
	return p
}
