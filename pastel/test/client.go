package test

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

// PastelClient implementing pastel.Client mock for testing purpose
type PastelClient struct {
	ClientMock *mocks.Client
}

// NewMockPastelClient create new pastel client mock
func NewMockPastelClient() *PastelClient {
	return &PastelClient{
		ClientMock: &mocks.Client{},
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returning Mn's and error from args
func (p *PastelClient) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *PastelClient {
	p.ClientMock.On("MasterNodesTop", mock.Anything).Return(nodes, err)
	return p
}

// ListenOnStorageFee listening StorageFee call and returning pastel.StorageFee, error form args
func (p *PastelClient) ListenOnStorageFee(fee *pastel.StorageFee, returnErr error) *PastelClient {
	p.ClientMock.On("StorageFee", mock.Anything).Return(fee, returnErr)
	return p
}
