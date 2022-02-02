package nftdownload

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

	// RegisterNftMethod represent RegisterNft name method
	RegisterNftMethod = "RegisterNft"

	// DownloadNftMethod represent DownloadNft name method
	DownloadNftMethod = "DownloadNft"

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
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.DownloadNftInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                    t,
		ClientInterface:      &mocks.ClientInterface{},
		ConnectionInterface:  &mocks.ConnectionInterface{},
		DownloadNftInterface: &mocks.DownloadNftInterface{},
	}
}

// ListenOnRegisterNft listening RegisterNft call
func (client *Client) ListenOnRegisterNft() *Client {
	client.ConnectionInterface.On(RegisterNftMethod).Return(nil)
	return client
}

// ListenOnDownloadNft listening DownloadNft call
func (client *Client) ListenOnDownloadNft() *Client {
	client.ConnectionInterface.On(DownloadNftMethod).Return(client.DownloadNft)
	return client
}

// AssertDownloadNftCall assertion DownloadNft call
func (client *Client) AssertDownloadNftCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, DownloadNftMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, DownloadNftMethod, expectedCalls)
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

// ListenOnDownload listening Download call and returns args value
func (client *Client) ListenOnDownload(arguments ...interface{}) *Client {
	client.DownloadNftInterface.On(DownloadMethod, mock.Anything,
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string(""))).Return(arguments...)
	return client
}

// ListenOnDownloadThumbnail listening DownloadThumbnail call and returns args value
func (client *Client) ListenOnDownloadThumbnail(arguments ...interface{}) *Client {
	client.DownloadNftInterface.On(DownloadThumbnailMethod, mock.Anything,
		mock.Anything).Return(arguments...)
	return client
}

// AssertDownloadCall assertion Download call
func (client *Client) AssertDownloadCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadNftInterface.AssertCalled(client.t, DownloadMethod, arguments...)
	}
	client.DownloadNftInterface.AssertNumberOfCalls(client.t, DownloadMethod, expectedCalls)
	return client
}

// ListenOnProcessUserdata listening ProcessUserdata call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.ConnectionInterface.On(ProcessUserdataMethod).Return(nil)
	return client
}
