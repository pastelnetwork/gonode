package processdata

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

	// RegisterArtworkMethod represent RegisterNft name method
	RegisterArtworkMethod = "RegisterNft"

	// DownloadArtworkMethod represent DownloadNft name method
	DownloadArtworkMethod = "DownloadNft"

	// ProcessUserdataMethod represent ProcessUserdataInterface name method
	ProcessUserdataMethod = "ProcessUserdataInterface"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.ProcessUserdataInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                        t,
		ClientInterface:          &mocks.ClientInterface{},
		ConnectionInterface:      &mocks.ConnectionInterface{},
		ProcessUserdataInterface: &mocks.ProcessUserdataInterface{},
	}
}

// ListenOnRegisterArtwork listening RegisterNft call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.ConnectionInterface.On(RegisterArtworkMethod).Return(nil)
	return client
}

// AssertRegisterArtworkCall assertion RegisterNft call
func (client *Client) AssertRegisterArtworkCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, RegisterArtworkMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, RegisterArtworkMethod, expectedCalls)
	return client
}

// ListenOnDownloadArtwork listening DownloadNft call
func (client *Client) ListenOnDownloadArtwork() *Client {
	client.ConnectionInterface.On(DownloadArtworkMethod).Return(client.DownloadNft)
	return client
}

// AssertDownloadArtworkCall assertion DownloadNft call
func (client *Client) AssertDownloadArtworkCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, DownloadArtworkMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, DownloadArtworkMethod, expectedCalls)
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
	client.ProcessUserdataInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProcessUserdata listening ProcessUserdataInterface call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.ConnectionInterface.On(ProcessUserdataMethod).Return(client.ProcessUserdataInterface)
	return client
}

// AssertProcessUserdataCall assertion ProcessUserdataInterface call
func (client *Client) AssertProcessUserdataCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, ProcessUserdataMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, ProcessUserdataMethod, expectedCalls)
	return client
}

// ListenOnSessionUserdata listening Session call and returns error from args
func (client *Client) ListenOnSessionUserdata(returnErr error) *Client {
	client.ProcessUserdataInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCallUserdata assertion Session Call
func (client *Client) AssertSessionCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodesUserdata listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodesUserdata(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.ProcessUserdataInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCallUserdata assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectToUserdata listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectToUserdata(returnErr error) *Client {
	client.ProcessUserdataInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCallUserdata assertion ConnectTo call
func (client *Client) AssertConnectToCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessIDUserdata listening SessID call and returns sessID from args
func (client *Client) ListenOnSessIDUserdata(sessID string) *Client {
	client.ProcessUserdataInterface.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCallUserdata assertion SessID call
func (client *Client) AssertSessIDCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}
