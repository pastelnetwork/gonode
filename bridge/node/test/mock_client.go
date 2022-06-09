package test

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/bridge/node/mocks"
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

	// RegisterNftMethod represent RegisterNftInterface name method
	RegisterNftMethod = "RegisterNftInterface"

	// DownloadDataMethod represent DownloadNftInterface name method
	DownloadDataMethod = "DownloadData"

	// ProcessUserdataMethod represent ProcessUserdataInterface name method
	ProcessUserdataMethod = "ProcessUserdataInterface"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	// SendPreBurntFeeTxidMethod represent SendPreBurntFeeTxId method
	SendPreBurntFeeTxidMethod = "SendPreBurntFeeTxid"

	// SendSignedTicketMethod represent SendSignedTicket method
	SendSignedTicketMethod = "SendSignedTicket"

	// UploadImageWithThumbnailMethod represent UploadImageWithThumbnail method
	UploadImageWithThumbnailMethod = "UploadImageWithThumbnail"

	// DownloadMethod represent Download name method
	DownloadMethod = "Download"

	// DownloadThumbnailMethod represent DownloadThumbnail name method
	DownloadThumbnailMethod = "DownloadThumbnail"

	// DownloadDDAndFPMethod represent DownloadDDAndFP name method
	DownloadDDAndFPMethod = "DownloadDDAndFingerprints"

	// SendActionActMethod represent SendActionAct method
	SendActionActMethod = "SendActionAct"

	// UploadAssetMethod represent UploadAsset name method
	UploadAssetMethod = "UploadAsset"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	//*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.SNClientInterface
	*mocks.DownloadDataInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                     t,
		SNClientInterface:     &mocks.SNClientInterface{},
		ConnectionInterface:   &mocks.ConnectionInterface{},
		DownloadDataInterface: &mocks.DownloadDataInterface{},
	}
}

// AssertRegisterNftCall assertion RegisterNftInterface call
func (client *Client) AssertRegisterNftCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, RegisterNftMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, RegisterNftMethod, expectedCalls)
	return client
}

// ListenOnDownloadNft listening DownloadNftInterface call
func (client *Client) ListenOnDownloadData() *Client {
	client.ConnectionInterface.On(DownloadDataMethod).Return(client.DownloadDataInterface)
	client.SNClientInterface.On(DownloadDataMethod).Return(client.DownloadDataInterface)
	return client
}

// AssertDownloadNftCall assertion DownloadNftInterface call
func (client *Client) AssertDownloadNftCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, DownloadDataMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, DownloadDataMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		//client.ClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string("")), mock.Anything).Return(client.ConnectionInterface, returnErr)
		client.SNClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string("")), mock.Anything).Return(client.ConnectionInterface, returnErr)
	} else {
		//client.ClientInterface.On(ConnectMethod, mock.Anything, addr, mock.Anything).Return(client.ConnectionInterface, returnErr)
		client.SNClientInterface.On(ConnectMethod, mock.Anything, addr, mock.Anything).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.SNClientInterface.AssertCalled(client.t, ConnectMethod, arguments...)
	}
	client.SNClientInterface.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
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
	client.SNClientInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	client.DownloadDataInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.SNClientInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	client.ConnectionInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	client.DownloadDataInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)

	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.DownloadDataInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)

	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.DownloadDataInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)

	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadDataInterface.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.DownloadDataInterface.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.DownloadDataInterface.On(SessIDMethod).Return(sessID)
	client.DownloadDataInterface.On(SessIDMethod).Return(sessID)
	client.DownloadDataInterface.On(SessIDMethod).Return(sessID)

	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadDataInterface.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.DownloadDataInterface.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnDownloadThumbnail listening DownloadThumbnail call and returns args value
func (client *Client) ListenOnDownloadThumbnail(arguments ...interface{}) *Client {
	client.DownloadDataInterface.On(DownloadThumbnailMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnDownloadDDAndFP listens for a ListenOnDownloadDDAndFP call and returns args value
func (client *Client) ListenOnDownloadDDAndFP(arguments ...interface{}) *Client {
	client.DownloadDataInterface.On(DownloadDDAndFPMethod, mock.Anything,
		mock.Anything).Return(arguments...)
	return client
}

// AssertDownloadCall assertion Download call
func (client *Client) AssertDownloadCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadDataInterface.AssertCalled(client.t, DownloadMethod, arguments...)
	}
	client.DownloadDataInterface.AssertNumberOfCalls(client.t, DownloadMethod, expectedCalls)
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

// ListenOnAcceptedNodesUserdata listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodesUserdata(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.DownloadDataInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCallUserdata assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadDataInterface.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.DownloadDataInterface.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}
