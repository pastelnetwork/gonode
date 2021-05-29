package pastel

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test/mocks"
	"github.com/stretchr/testify/mock"
)

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
