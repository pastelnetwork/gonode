package synchronizer

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func TestWaitSynchronizationSuccessful(t *testing.T) {
	s := Synchronizer{}

	status := &pastel.MasterNodeStatus{
		Status: "Masternode successfully started",
	}

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(status, nil)
	s.PastelClient = pMock

	err := s.WaitSynchronization(context.Background())
	assert.True(t, err == nil)
}

func TestWaitSynchronizationTimeout(t *testing.T) {
	s := Synchronizer{}

	status := &pastel.MasterNodeStatus{
		Status: "hello",
	}

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(status, nil)
	s.PastelClient = pMock

	err := s.WaitSynchronization(context.Background())
	assert.Equal(t, err.Error(), "timeout expired")
}

func TestWaitSynchronizationError(t *testing.T) {
	s := Synchronizer{}
	errMsg := "timeout expired"

	pMock := pastelMock.NewMockClient(t)
	pMock.ListenOnMasterNodeStatus(nil, errors.New(errMsg))
	s.PastelClient = pMock

	err := s.WaitSynchronization(context.Background())
	assert.Equal(t, err.Error(), errMsg)
}
