package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	MasterNodesExtraMethod = "MasterNodesExtra"
	GetBlockCountMethod    = "GetBlockCount"
	GetBlockVerbose1Method = "GetBlockVerbose1"
)

type PastelClient struct {
	t *testing.T
	*mocks.Client
}

func (p *PastelClient) ListenOnGetBlockCount(blkCount int32, err error) *PastelClient {
	p.On(GetBlockCountMethod, mock.Anything).Return(blkCount, err)
	return p
}

func (p *PastelClient) ListenOnGetBlockVerbose1(rs *pastel.GetBlockVerbose1Result, err error) *PastelClient {
	p.On(GetBlockVerbose1Method, mock.Anything, mock.Anything).Return(rs, err)
	return p
}

func (p *PastelClient) ListenOnMasterNodesExtra(mns pastel.MasterNodes, err error) *PastelClient {
	p.On(MasterNodesExtraMethod, mock.Anything).Return(mns, err)
	return p
}

func (p *PastelClient) AssertGetBlockCountCall(expectedCalls int, arguments ...interface{}) *PastelClient {
	if expectedCalls > 0 {
		p.AssertCalled(p.t, GetBlockCountMethod, arguments...)
	}

	p.AssertNumberOfCalls(p.t, GetBlockCountMethod, expectedCalls)
	return p
}

func (p *PastelClient) AssertGetBlockVerbose1Call(expectedCalls int, arguments ...interface{}) *PastelClient {
	if expectedCalls > 0 {
		p.AssertCalled(p.t, GetBlockVerbose1Method, arguments...)
	}

	p.AssertNumberOfCalls(p.t, GetBlockVerbose1Method, expectedCalls)
	return p
}

func (p *PastelClient) AssertMasterNodesExtraCall(expectedCalls int, arguments ...interface{}) *PastelClient {
	if expectedCalls > 0 {
		p.AssertCalled(p.t, MasterNodesExtraMethod, arguments...)
	}

	p.AssertNumberOfCalls(p.t, MasterNodesExtraMethod, expectedCalls)
	return p
}

func NewMockPastelClient(t *testing.T) *PastelClient {
	return &PastelClient{
		t:      t,
		Client: &mocks.Client{},
	}
}
