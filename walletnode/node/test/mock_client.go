package test

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

	// RegisterArtworkMethod represent RegisterNft name method
	RegisterArtworkMethod = "RegisterNft"

	// DownloadArtworkMethod represent DownloadNft name method
	DownloadArtworkMethod = "DownloadNft"

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
	*mocks.RegisterNft
	*mocks.DownloadNft
	*mocks.ProcessUserdata
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:               t,
		Client:          &mocks.Client{},
		Connection:      &mocks.Connection{},
		RegisterNft:     &mocks.RegisterNft{},
		DownloadNft:     &mocks.DownloadNft{},
		ProcessUserdata: &mocks.ProcessUserdata{},
	}
}

// ListenOnRegisterArtwork listening RegisterNft call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On(RegisterArtworkMethod).Return(client.RegisterNft)
	return client
}

// ListenOnSendPreBurntFeeTxID listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendPreBurntFeeTxID(txid string, err error) *Client {
	client.RegisterNft.On(SendPreBurntFeeTxidMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnSendSignedTicket listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendSignedTicket(id int64, err error) *Client {
	client.RegisterNft.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnUploadImageWithThumbnail listening UploadImageWithThumbnail call
func (client *Client) ListenOnUploadImageWithThumbnail(retPreviewHash []byte,
	retMediumThumbnailHash []byte, retsmallThumbnailHash []byte, retErr error) *Client {

	client.RegisterNft.On(UploadImageWithThumbnailMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(retPreviewHash,
		retMediumThumbnailHash, retsmallThumbnailHash, retErr)

	return client
}

// AssertRegisterArtworkCall assertion RegisterNft call
func (client *Client) AssertRegisterArtworkCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, RegisterArtworkMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, RegisterArtworkMethod, expectedCalls)
	return client
}

// ListenOnDownloadArtwork listening DownloadNft call
func (client *Client) ListenOnDownloadArtwork() *Client {
	client.Connection.On(DownloadArtworkMethod).Return(client.DownloadNft)
	return client
}

// AssertDownloadArtworkCall assertion DownloadNft call
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
	client.RegisterNft.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnProbeImage(arguments ...interface{}) *Client {
	client.RegisterNft.On(ProbeImageMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// AssertProbeImageCall assertion ProbeImage call
func (client *Client) AssertProbeImageCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNft.AssertCalled(client.t, ProbeImageMethod, arguments...)
	}
	client.RegisterNft.AssertNumberOfCalls(client.t, ProbeImageMethod, expectedCalls)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterNft.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session Call
func (client *Client) AssertSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNft.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.RegisterNft.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterNft.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNft.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.RegisterNft.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterNft.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNft.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.RegisterNft.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterNft.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterNft.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.RegisterNft.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnDownload listening Download call and returns args value
func (client *Client) ListenOnDownload(arguments ...interface{}) *Client {
	client.DownloadNft.On(DownloadMethod, mock.Anything,
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string(""))).Return(arguments...)
	return client
}

// ListenOnDownloadThumbnail listening DownloadThumbnail call and returns args value
func (client *Client) ListenOnDownloadThumbnail(arguments ...interface{}) *Client {
	client.DownloadNft.On(DownloadThumbnailMethod, mock.Anything,
		mock.Anything).Return(arguments...)
	return client
}

// AssertDownloadCall assertion Download call
func (client *Client) AssertDownloadCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.DownloadNft.AssertCalled(client.t, DownloadMethod, arguments...)
	}
	client.DownloadNft.AssertNumberOfCalls(client.t, DownloadMethod, expectedCalls)
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
