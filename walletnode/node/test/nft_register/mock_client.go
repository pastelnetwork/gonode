package nftregister

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

	// RegisterNftMethod represent RegisterNftInterface name method
	RegisterNftMethod = "RegisterNft"

	// DownloadNftMethod represent DownloadNft name method
	DownloadNftMethod = "DownloadNft"

	// ProcessUserdataMethod represent ProcessUserdata name method
	ProcessUserdataMethod = "ProcessUserdata"

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

	// GetDupeDetectionDBHashMethod  is method name for GetDupeDetectionDBHash
	GetDupeDetectionDBHashMethod = "GetDupeDetectionDBHash"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.RegisterNftInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                    t,
		ClientInterface:      &mocks.ClientInterface{},
		ConnectionInterface:  &mocks.ConnectionInterface{},
		RegisterNftInterface: &mocks.RegisterNftInterface{},
	}
}

// ListenOnRegisterNft listening RegisterNftInterface call
func (client *Client) ListenOnRegisterNft() *Client {
	client.ConnectionInterface.On(RegisterNftMethod).Return(client.RegisterNftInterface)
	return client
}

// ListenOnSendPreBurntFeeTxID listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendPreBurntFeeTxID(txid string, err error) *Client {
	client.RegisterNftInterface.On(SendPreBurntFeeTxidMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnSendSignedTicket listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendSignedTicket(id int64, err error) *Client {
	client.RegisterNftInterface.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnUploadImageWithThumbnail listening UploadImageWithThumbnail call
func (client *Client) ListenOnUploadImageWithThumbnail(retPreviewHash []byte,
	retMediumThumbnailHash []byte, retsmallThumbnailHash []byte, retErr error) *Client {

	client.RegisterNftInterface.On(UploadImageWithThumbnailMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(retPreviewHash,
		retMediumThumbnailHash, retsmallThumbnailHash, retErr)

	return client
}

// AssertRegisterNftCall assertion RegisterNftInterface call
func (client *Client) AssertRegisterNftCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, RegisterNftMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, RegisterNftMethod, expectedCalls)
	return client
}

// ListenOnDownloadNft listening DownloadNft call
func (client *Client) ListenOnDownloadNft() *Client {
	client.ConnectionInterface.On(DownloadNftMethod).Return(nil)
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

// ListenOnGetDupeDetectionDBHas listening GetDupeDetectionDBHash call
func (client *Client) ListenOnGetDupeDetectionDBHash(hash string, err error) *Client {
	client.RegisterNftInterface.On(GetDupeDetectionDBHashMethod, mock.Anything, mock.Anything).Return(hash, err)
	return client
}

// ListenOnMeshNodes listening MeshNodes call and returns args value
func (client *Client) ListenOnMeshNodes(arguments ...interface{}) *Client {
	client.RegisterNftInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnProbeImage(arguments ...interface{}) *Client {
	client.RegisterNftInterface.On(ProbeImageMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// AssertProbeImageCall assertion ProbeImage call
func (client *Client) AssertProbeImageCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNftInterface.AssertCalled(client.t, ProbeImageMethod, arguments...)
	}
	client.RegisterNftInterface.AssertNumberOfCalls(client.t, ProbeImageMethod, expectedCalls)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterNftInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session Call
func (client *Client) AssertSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNftInterface.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.RegisterNftInterface.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterNftInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNftInterface.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.RegisterNftInterface.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterNftInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNftInterface.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.RegisterNftInterface.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterNftInterface.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNftInterface.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.RegisterNftInterface.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnProcessUserdata listening ProcessUserdata call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.ConnectionInterface.On(ProcessUserdataMethod).Return(nil)
	return client
}

// AssertProcessUserdataCall assertion ProcessUserdata call
func (client *Client) AssertProcessUserdataCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, ProcessUserdataMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, ProcessUserdataMethod, expectedCalls)
	return client
}
