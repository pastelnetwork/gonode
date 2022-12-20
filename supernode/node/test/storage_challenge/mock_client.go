package test

import (
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"testing"

	"github.com/pastelnetwork/gonode/supernode/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// AcceptedNodesMethod represent AcceptedNodes name method
	AcceptedNodesMethod = "AcceptedNodes"

	// CloseMethod represent Close name method
	CloseMethod = "Close"

	// ConnectMethod represent Connect name method
	ConnectMethod = "Connect"

	// ConnectToMethod represent Connect name method
	ConnectToMethod = "ConnectTo"

	// DoneMethod represent Done call
	DoneMethod = "Done"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	//ProcessStorageChallengeMethod intercepts the ProcessStorageChallenge method
	ProcessStorageChallengeMethod = "ProcessStorageChallenge"

	//VerifyStorageChallengeMethod intercepts the VerifyStorageChallenge method
	VerifyStorageChallengeMethod = "VerifyStorageChallenge"

	//StorageChallengeInterfaceMethod intercepts the StorageChallenge interface
	StorageChallengeInterfaceMethod = "StorageChallenge"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.StorageChallengeInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                         t,
		ClientInterface:           &mocks.ClientInterface{},
		ConnectionInterface:       &mocks.ConnectionInterface{},
		StorageChallengeInterface: &mocks.StorageChallengeInterface{},
	}
}

// ListenOnStorageChallengeInterface listening RegisterNftInterface call
func (client *Client) ListenOnStorageChallengeInterface() *Client {
	client.ConnectionInterface.On(StorageChallengeInterfaceMethod).Return(client.StorageChallengeInterface)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.ClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.ConnectionInterface, returnErr)
	} else {
		client.ClientInterface.On(ConnectMethod, mock.Anything, addr).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// ListenOnProcessStorageChallengeFunc returns returnErr
func (client *Client) ListenOnProcessStorageChallengeFunc(returnErr error) *Client {
	client.StorageChallengeInterface.On(ProcessStorageChallengeMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}

// ListenOnVerifyStorageChallengeFunc returns returnErr
func (client *Client) ListenOnVerifyStorageChallengeFunc(data *pb.StorageChallengeData, returnErr error) *Client {
	client.StorageChallengeInterface.On(VerifyStorageChallengeMethod, mock.Anything, mock.Anything).Return(data, returnErr)

	return client
}
