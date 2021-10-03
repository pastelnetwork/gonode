package test

import (
	"testing"

	"ddclient"
	"ddclient/mocks"

	"github.com/stretchr/testify/mock"
)

const (

	// ImageRarenessScoreMethod represents Encode call
	ImageRarenessScoreMethod = "ImageRarenessScore"
)

// Client implements node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.DDServerClient
}

// NewMockClient creates new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:              t,
		DDServerClient: &mocks.DDServerClient{},
	}
}

// ListenOnImageRarenessScore listening ImageRarenessScore call and returns channel from args
func (client *Client) ListenOnImageRarenessScore(returnEnc *ddclient.DupeDetection, returnErr error) *Client {
	client.DDServerClient.On(ImageRarenessScoreMethod, mock.Anything, mock.IsType([]byte{}), mock.Anything).Return(returnEnc, returnErr)
	return client
}

// AssertEncodeCall assertion ImageRarenessScore call
func (client *Client) AssertImageRarenessScoreCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DDServerClient.AssertCalled(client.t, ImageRarenessScoreMethod, arguments...)
	}
	client.DDServerClient.AssertNumberOfCalls(client.t, ImageRarenessScoreMethod, expectedCalls)
	return client
}
