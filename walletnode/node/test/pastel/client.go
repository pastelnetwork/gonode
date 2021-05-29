package pastel

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing PastelClient for testing purpose
type Client struct {
	PastelMock *mocks.PastelClient
}

// NewMockPastelClient create new pastel client mock
func NewMockPastelClient() *Client {
	return &Client{
		PastelMock: &mocks.PastelClient{},
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returning Mn's and error from args
func (p *Client) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *Client {
	p.PastelMock.On("MasterNodesTop", mock.Anything).Return(nodes, err)
	return p
}

// ListenOnStorageFee listening StorageFee call and returning pastel.StorageFee, error form args
func (p *Client) ListenOnStorageFee(fee *pastel.StorageFee, returnErr error) *Client {
	p.PastelMock.On("StorageFee", mock.Anything).Return(fee, returnErr)
	return p
}
