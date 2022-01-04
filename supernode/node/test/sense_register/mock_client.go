package test

import (
	"context"
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

	// ProbeImageMethod represent ProbeImage name method
	ProbeImageMethod = "ProbeImage"

	// RegisterSenseMethod represent RegisterSense name method
	RegisterSenseMethod = "RegisterSense"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	// UploadImageWithThumbnailMethod represent UploadImageWithThumbnail method
	UploadImageWithThumbnailMethod = "UploadImageWithThumbnail"
	// DownloadMethod represent Download name method
	DownloadMethod = "Download"

	// SendArtTicketSignatureMethod represent SendArtTicketSignature method
	SendArtTicketSignatureMethod = "SendArtTicketSignature"

	// SendSignedDDAndFingerprintsMethod represent SendSignedDDAndFingerprints method
	SendSignedDDAndFingerprintsMethod = "SendSignedDDAndFingerprints"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.RegisterSense
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:             t,
		Client:        &mocks.Client{},
		Connection:    &mocks.Connection{},
		RegisterSense: &mocks.RegisterSense{},
	}
}

// ListenOnRegisterSense listening RegisterSense call
func (client *Client) ListenOnRegisterSense() *Client {
	client.Connection.On(RegisterSenseMethod).Return(client.RegisterSense)
	return client
}

// AssertRegisterSenseCall assertion RegisterSense call
func (client *Client) AssertRegisterSenseCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, RegisterSenseMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, RegisterSenseMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.Client.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.Connection, returnErr)
	} else {
		client.Client.On(ConnectMethod, mock.Anything, addr).Return(client.Connection, returnErr)
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

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterSense.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// ListenOnSendArtTicketSignature listens on send art ticket signature
func (client *Client) ListenOnSendArtTicketSignature(returnErr error) *Client {
	client.RegisterSense.On(SendArtTicketSignatureMethod, mock.Anything, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// ListenOnSendSignedDDAndFingerprints listens on send art ticket signature
func (client *Client) ListenOnSendSignedDDAndFingerprints(returnErr error) *Client {
	client.RegisterSense.On(SendSignedDDAndFingerprintsMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session Call
func (client *Client) AssertSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSense.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.RegisterSense.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterSense.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSense.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.RegisterSense.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterSense.On(ConnectToMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSense.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.RegisterSense.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterSense.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterSense.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.RegisterSense.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}