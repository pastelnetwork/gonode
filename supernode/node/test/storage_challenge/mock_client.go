package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/common/types"
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

	// ConnectSNMethod represent ConnectSN name method
	ConnectSNMethod = "ConnectSN"

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

	//VerifyEvaluationResultMethod intercepts the VerifyEvaluationResult Method
	VerifyEvaluationResultMethod = "VerifyEvaluationResult"

	//StorageChallengeInterfaceMethod intercepts the StorageChallenge interface
	StorageChallengeInterfaceMethod = "StorageChallenge"

	//BroadcastStorageChallengeResultMethod intercepts the BroadcastStorageChallengeResult Method
	BroadcastStorageChallengeResultMethod = "BroadcastStorageChallengeResult"
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

// ListenOnConnectSN listening Connect call and returns error from args
func (client *Client) ListenOnConnectSN(addr string, returnErr error) *Client {
	if addr == "" {
		client.ClientInterface.On(ConnectSNMethod, mock.Anything, mock.IsType(string(""))).Return(client.ConnectionInterface, returnErr)
	} else {
		client.ClientInterface.On(ConnectSNMethod, mock.Anything, addr).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// ListenOnProcessStorageChallengeFunc returns returnErr
func (client *Client) ListenOnProcessStorageChallengeFunc(returnErr error) *Client {
	client.StorageChallengeInterface.On(ProcessStorageChallengeMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}

// ListenOnVerifyStorageChallengeFunc returns returnErr
func (client *Client) ListenOnVerifyStorageChallengeFunc(returnErr error) *Client {
	client.StorageChallengeInterface.On(VerifyStorageChallengeMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}

// ListenOnVerifyEvaluationResultFunc returns affirmations & returnErr
func (client *Client) ListenOnVerifyEvaluationResultFunc(nodeID string, data types.Message, returnErr error) *Client {
	data.Sender = nodeID
	client.StorageChallengeInterface.On(VerifyEvaluationResultMethod, mock.Anything, mock.Anything).Times(1).Return(data, returnErr)

	return client
}

// ListenOnBroadcastStorageChallengeResultFunc returns error
func (client *Client) ListenOnBroadcastStorageChallengeResultFunc(returnErr error) *Client {
	client.StorageChallengeInterface.On(BroadcastStorageChallengeResultMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}
