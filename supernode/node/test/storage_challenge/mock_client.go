package test

import (
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

	//ProcessStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error
	ProcessStorageChallengeMethod = "ProcessStorageChallenge"

	//VerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error
	VerifyStorageChallengeMethod = "VerifyStorageChallenge"

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

// ListenOnRegisterNft listening RegisterNftInterface call
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

// ListenOnProcessStorageChallengeFunc returns returnErr
func (client *Client) ListenOnVerifyStorageChallengeFunc(returnErr error) *Client {
	client.StorageChallengeInterface.On(VerifyStorageChallengeMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}
