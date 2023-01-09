package test

import (
	"testing"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
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

	//ProcessSelfHealingChallengeMethod intercepts the ProcessSelfHealingChallenge method
	ProcessSelfHealingChallengeMethod = "ProcessSelfHealingChallenge"

	//VerifySelfHealingChallengeMethod intercepts the VerifySelfHealingChallenge method
	VerifySelfHealingChallengeMethod = "VerifySelfHealingChallenge"

	//SelfHealingChallengeInterfaceMethod intercepts the SelfHealing interface
	SelfHealingChallengeInterfaceMethod = "SelfHealingChallenge"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.SelfHealingChallengeInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                             t,
		ClientInterface:               &mocks.ClientInterface{},
		ConnectionInterface:           &mocks.ConnectionInterface{},
		SelfHealingChallengeInterface: &mocks.SelfHealingChallengeInterface{},
	}
}

// ListenOnSelfHealingChallengeInterface listening SelfHealingChallengeInterface call
func (client *Client) ListenOnSelfHealingChallengeInterface() *Client {
	client.ConnectionInterface.On(SelfHealingChallengeInterfaceMethod).Return(client.SelfHealingChallengeInterface)
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

// ListenOnProcessSelfHealingChallengeFunc returns returnErr
func (client *Client) ListenOnProcessSelfHealingChallengeFunc(returnErr error) *Client {
	client.SelfHealingChallengeInterface.On(ProcessSelfHealingChallengeMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}

// ListenOnVerifySelfHealingChallengeFunc returns returnErr
func (client *Client) ListenOnVerifySelfHealingChallengeFunc(data *pb.SelfHealingData, returnErr error) *Client {
	client.SelfHealingChallengeInterface.On(VerifySelfHealingChallengeMethod, mock.Anything, mock.Anything).Return(data, returnErr).Times(5)

	return client
}
