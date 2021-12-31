package artwork_download

import (
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

	// RegisterArtworkMethod represent RegisterArtwork name method
	RegisterArtworkMethod = "RegisterArtwork"

	// DownloadArtworkMethod represent DownloadArtwork name method
	DownloadArtworkMethod = "DownloadArtwork"

	// ProcessUserdataMethod represent ProcessUserdata name method
	ProcessUserdataMethod = "ProcessUserdata"

	// DownloadMethod represent Download name method
	DownloadMethod = "Download"
	// DownloadThumbnailMethod represent DownloadThumbnail name method
	DownloadThumbnailMethod = "DownloadThumbnail"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.DownloadArtwork
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:               t,
		Client:          &mocks.Client{},
		Connection:      &mocks.Connection{},
		DownloadArtwork: &mocks.DownloadArtwork{},
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On(RegisterArtworkMethod).Return(nil)
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

// ListenOnDownload listening Download call and returns args value
func (client *Client) ListenOnDownload(arguments ...interface{}) *Client {
	client.DownloadArtwork.On(DownloadMethod, mock.Anything,
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string(""))).Return(arguments...)
	return client
}

// ListenOnDownloadThumbnail listening DownloadThumbnail call and returns args value
func (client *Client) ListenOnDownloadThumbnail(arguments ...interface{}) *Client {
	client.DownloadArtwork.On(DownloadThumbnailMethod, mock.Anything,
		mock.Anything).Return(arguments...)
	return client
}

// AssertDownloadCall assertion Download call
func (client *Client) AssertDownloadCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadArtwork.AssertCalled(client.t, DownloadMethod, arguments...)
	}
	client.DownloadArtwork.AssertNumberOfCalls(client.t, DownloadMethod, expectedCalls)
	return client
}

// ListenOnProcessUserdata listening ProcessUserdata call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.Connection.On(ProcessUserdataMethod).Return(nil)
	return client
}
