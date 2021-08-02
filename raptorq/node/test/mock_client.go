package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/raptorq/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// ConnectMethod represents Connect name method
	ConnectMethod = "Connect"

	// CloseMethod represents Close name method
	CloseMethod = "Close"

	// DoneMethod represents Done call
	DoneMethod = "Done"

	// RaptorQMethod represents RaptorQ call
	RaptorQMethod = "RaptorQ"

	// EncodeMethod represents Encode call
	EncodeMethod = "Encode"

	// EncodeInfoMethod represents EncodeInfo call
	EncodeInfoMethod = "EncodeInfo"

	// DecodeMethod represents Decode call
	DecodeMethod = "Decode"
)

// Client implements node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.RaptorQ
}

// NewMockClient creates new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:          t,
		Client:     &mocks.Client{},
		Connection: &mocks.Connection{},
		RaptorQ:    &mocks.RaptorQ{},
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

// ListenOnRaptorQ listening RaptorQ call and returns channel from args
func (client *Client) ListenOnRaptorQ() *Client {
	client.Connection.On(RaptorQMethod, mock.Anything).Return(client.RaptorQ)
	return client
}

// AssertRaptorQCall assertion RaptorQ call
func (client *Client) AssertRaptorQCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, RaptorQMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, RaptorQMethod, expectedCalls)
	return client
}

// ListenOnEncode listening Encode call and returns channel from args
func (client *Client) ListenOnEncode(returnEnc *node.Encode, returnErr error) *Client {
	client.Connection.On(EncodeMethod, mock.Anything, mock.IsType([]byte{})).Return(returnEnc, returnErr)
	return client
}

// AssertEncodeCall assertion Encode call
func (client *Client) AssertEncodeCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, EncodeMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, EncodeMethod, expectedCalls)
	return client
}

// ListenOnEncodeInfo listening EncodeInfo call and returns channel from args
func (client *Client) ListenOnEncodeInfo(returnEncInfo *node.EncodeInfo, returnErr error) *Client {
	client.RaptorQ.On(EncodeInfoMethod, mock.Anything, mock.IsType([]byte{}), mock.IsType(uint32(0)), mock.IsType(string("")), mock.IsType(string(""))).Return(returnEncInfo, returnErr)
	return client
}

// AssertEncodeInfoCall assertion EncodeInfo call
func (client *Client) AssertEncodeInfoCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, EncodeInfoMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, EncodeInfoMethod, expectedCalls)
	return client
}

// ListenOnDecode listening Decode call and returns channel from args
func (client *Client) ListenOnDecode(returnDec *node.Decode, returnErr error) *Client {
	client.Connection.On(DecodeMethod, mock.Anything, mock.Anything).Return(returnDec, returnErr)
	return client
}

// AssertDecodeCall assertion Decode call
func (client *Client) AssertDecodeCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, DecodeMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, DecodeMethod, expectedCalls)
	return client
}
