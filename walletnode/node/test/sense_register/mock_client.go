package senseregister

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
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

	// MeshNodesMethod represent MeshNodes name method
	MeshNodesMethod = "MeshNodes"

	// ProbeImageMethod represent ProbeImage name method
	ProbeImageMethod = "ProbeImage"

	// RegisterSenseMethod represent RegisterSense name method
	RegisterSenseMethod = "RegisterSense"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	// SendActionActMethod represent SendActionAction method
	SendActionActMethod = "SendActionAct"

	// SendRegMetadataMethod represent SendRegMetadata method
	SendRegMetadataMethod = "SendRegMetadata"

	// SendSignedTicketMethod represent SendSignedTicket method
	SendSignedTicketMethod = "SendSignedTicket"

	// GetDupeDetectionDBHashMethod  is method name for GetDupeDetectionDBHash
	GetDupeDetectionDBHashMethod = "GetDupeDetectionDBHash"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.RegisterSenseInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                      t,
		ClientInterface:        &mocks.ClientInterface{},
		ConnectionInterface:    &mocks.ConnectionInterface{},
		RegisterSenseInterface: &mocks.RegisterSenseInterface{},
	}
}

// ListenOnRegisterSense listening RegisterSenseInterface call
func (client *Client) ListenOnRegisterSense() *Client {
	client.ConnectionInterface.On(RegisterSenseMethod).Return(client.RegisterSenseInterface)
	return client
}

// ListenOnSendSignedTicket listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendSignedTicket(id string, err error) *Client {
	client.RegisterSenseInterface.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnGetDupeDetectionDBHash listening GetDupeDetectionDBHash call
func (client *Client) ListenOnGetDupeDetectionDBHash(hash string, err error) *Client {
	client.RegisterSenseInterface.On(GetDupeDetectionDBHashMethod, mock.Anything, mock.Anything).Return(hash, err)
	return client
}

// AssertRegisterSenseCall assertion RegisterSenseInterface call
func (client *Client) AssertRegisterSenseCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, RegisterSenseMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, RegisterSenseMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.ClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string("")), mock.Anything).Return(client.ConnectionInterface, returnErr)
	} else {
		client.ClientInterface.On(ConnectMethod, mock.Anything, addr, mock.Anything).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ClientInterface.AssertCalled(client.t, ConnectMethod, arguments...)
	}
	client.ClientInterface.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
	return client
}

// ListenOnClose listening Close call and returns error from args
func (client *Client) ListenOnClose(returnErr error) *Client {
	client.ConnectionInterface.On(CloseMethod).Return(returnErr)
	return client
}

// AssertCloseCall assertion Close call
func (client *Client) AssertCloseCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, CloseMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, CloseMethod, expectedCalls)
	return client
}

// ListenOnDone listening Done call and returns channel from args
func (client *Client) ListenOnDone() *Client {
	client.ConnectionInterface.On(DoneMethod).Return(nil)
	return client
}

// AssertDoneCall assertion Done call
func (client *Client) AssertDoneCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, DoneMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, DoneMethod, expectedCalls)
	return client
}

// ListenOnMeshNodes listening MeshNodes call and returns args value
func (client *Client) ListenOnMeshNodes(arguments ...interface{}) *Client {
	client.RegisterSenseInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnProbeImage(arguments ...interface{}) *Client {
	client.RegisterSenseInterface.On(ProbeImageMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// AssertProbeImageCall assertion ProbeImage call
func (client *Client) AssertProbeImageCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, ProbeImageMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, ProbeImageMethod, expectedCalls)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterSenseInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session Call
func (client *Client) AssertSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterSenseInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterSenseInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterSenseInterface.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnSendRegMetadata listening SendRegMetadata call and returns error from args
func (client *Client) ListenOnSendRegMetadata(returnErr error) *Client {
	client.RegisterSenseInterface.On(SendRegMetadataMethod).Return(returnErr)
	return client
}

// AssertSendRegMetadata assertion SendRegMetadata call
func (client *Client) AssertSendRegMetadata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, SendRegMetadataMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, SendRegMetadataMethod, expectedCalls)
	return client
}

// ListenOnSendActionAct listening SendActionAct call and returns error from args
func (client *Client) ListenOnSendActionAct(returnErr error) *Client {
	client.RegisterSenseInterface.On(SendActionActMethod).Return(returnErr)
	return client
}

// AssertSendActionAct assertion SendActionAct call
func (client *Client) AssertSendActionAct(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSenseInterface.AssertCalled(client.t, SendActionActMethod, arguments...)
	}
	client.RegisterSenseInterface.AssertNumberOfCalls(client.t, SendActionActMethod, expectedCalls)
	return client
}
