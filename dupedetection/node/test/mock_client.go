package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/dupedetection/node"
	"github.com/pastelnetwork/gonode/dupedetection/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// ConnectMethod represents Connect name method
	ConnectMethod = "Connect"

	// CloseMethod represents Close name method
	CloseMethod = "Close"

	// DoneMethod represents Done call
	DoneMethod = "Done"

	// DupedetectionMethod represents RaptorQ call
	DupedetectionMethod = "Dupedetection"

	// ImageRarenessScoreMethod represents Encode call
	ImageRarenessScoreMethod = "ImageRarenessScore"
)

// Client implements node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.Dupedetection
}

// NewMockClient creates new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:             t,
		Client:        &mocks.Client{},
		Connection:    &mocks.Connection{},
		Dupedetection: &mocks.Dupedetection{},
	}
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(returnErr error) *Client {
	client.Client.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.Connection, returnErr)
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

// ListenOnDupedetection listening Dupedetection call and returns channel from args
func (client *Client) ListenOnDupedetection() *Client {
	client.Connection.On(DupedetectionMethod, mock.Anything).Return(client.Dupedetection)
	return client
}

// AssertRaptorQCall assertion RaptorQ call
func (client *Client) AssertRaptorQCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, DupedetectionMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, DupedetectionMethod, expectedCalls)
	return client
}

// ListenOnImageRarenessScore listening ImageRarenessScore call and returns channel from args
func (client *Client) ListenOnImageRarenessScore(returnEnc *node.DupeDetection, returnErr error) *Client {
	client.Dupedetection.On(ImageRarenessScoreMethod, mock.Anything, mock.IsType([]byte{}), mock.Anything).Return(returnEnc, returnErr)
	return client
}

// AssertEncodeCall assertion ImageRarenessScore call
func (client *Client) AssertImageRarenessScoreCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, ImageRarenessScoreMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, ImageRarenessScoreMethod, expectedCalls)
	return client
}
