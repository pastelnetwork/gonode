package process_data

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

	// RegisterArtworkMethod represent RegisterArtwork name method
	RegisterArtworkMethod = "RegisterArtwork"

	// DownloadArtworkMethod represent DownloadArtwork name method
	DownloadArtworkMethod = "DownloadArtwork"

	// ProcessUserdataMethod represent ProcessUserdata name method
	ProcessUserdataMethod = "ProcessUserdata"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.ProcessUserdata
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:               t,
		Client:          &mocks.Client{},
		Connection:      &mocks.Connection{},
		ProcessUserdata: &mocks.ProcessUserdata{},
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On(RegisterArtworkMethod).Return(nil)
	return client
}

// AssertRegisterArtworkCall assertion RegisterArtwork call
func (client *Client) AssertRegisterArtworkCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, RegisterArtworkMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, RegisterArtworkMethod, expectedCalls)
	return client
}

// ListenOnDownloadArtwork listening DownloadArtwork call
func (client *Client) ListenOnDownloadArtwork() *Client {
	client.Connection.On(DownloadArtworkMethod).Return(client.DownloadArtwork)
	return client
}

// AssertDownloadArtworkCall assertion DownloadArtwork call
func (client *Client) AssertDownloadArtworkCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, DownloadArtworkMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, DownloadArtworkMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.Client.On(ConnectMethod, mock.Anything, mock.IsType(string("")), mock.Anything).Return(client.Connection, returnErr)
	} else {
		client.Client.On(ConnectMethod, mock.Anything, addr, mock.Anything).Return(client.Connection, returnErr)
	}

	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Client.AssertCalled(client.t, ConnectMethod, arguments...)
	}
	client.Client.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
	return client
}

// ListenOnClose listening Close call and returns error from args
func (client *Client) ListenOnClose(returnErr error) *Client {
	client.Connection.On(CloseMethod).Return(returnErr)
	return client
}

// AssertCloseCall assertion Close call
func (client *Client) AssertCloseCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, CloseMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, CloseMethod, expectedCalls)
	return client
}

// ListenOnDone listening Done call and returns channel from args
func (client *Client) ListenOnDone() *Client {
	client.Connection.On(DoneMethod).Return(nil)
	return client
}

// AssertDoneCall assertion Done call
func (client *Client) AssertDoneCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, DoneMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, DoneMethod, expectedCalls)
	return client
}

// ListenOnMeshNodes listening MeshNodes call and returns args value
func (client *Client) ListenOnMeshNodes(arguments ...interface{}) *Client {
	client.ProcessUserdata.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProcessUserdata listening ProcessUserdata call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.Connection.On(ProcessUserdataMethod).Return(client.ProcessUserdata)
	return client
}

// AssertProcessUserdataCall assertion ProcessUserdata call
func (client *Client) AssertProcessUserdataCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, ProcessUserdataMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, ProcessUserdataMethod, expectedCalls)
	return client
}

// ListenOnSessionUserdata listening Session call and returns error from args
func (client *Client) ListenOnSessionUserdata(returnErr error) *Client {
	client.ProcessUserdata.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCallUserdata assertion Session Call
func (client *Client) AssertSessionCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdata.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.ProcessUserdata.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodesUserdata listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodesUserdata(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.ProcessUserdata.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCallUserdata assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdata.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.ProcessUserdata.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectToUserdata listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectToUserdata(returnErr error) *Client {
	client.ProcessUserdata.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCallUserdata assertion ConnectTo call
func (client *Client) AssertConnectToCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdata.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.ProcessUserdata.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessIDUserdata listening SessID call and returns sessID from args
func (client *Client) ListenOnSessIDUserdata(sessID string) *Client {
	client.ProcessUserdata.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCallUserdata assertion SessID call
func (client *Client) AssertSessIDCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdata.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.ProcessUserdata.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}
