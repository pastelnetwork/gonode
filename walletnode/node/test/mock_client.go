package test

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
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

	// ProbeImageMethod represent ProbeImage name method
	ProbeImageMethod = "ProbeImage"

	// RegisterArtworkMethod represent RegisterArtwork name method
	RegisterArtworkMethod = "RegisterArtwork"

	// DownloadArtworkMethod represent DownloadArtwork name method
	DownloadArtworkMethod = "DownloadArtwork"

	// ProcessUserdataMethod represent ProcessUserdata name method
	ProcessUserdataMethod = "ProcessUserdata"

	// ExternalDupeDetectionMethod represents ExternalDupeDetection name method
	ExternalDupeDetectionMethod = "ExternalDupeDetection"

	// SessionMethod represent Session name method
	SessionMethod = "Session"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	// SendPreBurnedFeeTxidMethod represent SendPreBurnedFeeTxid method
	SendPreBurnedFeeTxidMethod = "SendPreBurnedFeeTxid"

	// SendPreBurnedFeeEDDTxIDMethod represent SendPreBurnedFeeTxid method
	SendPreBurnedFeeEDDTxIDMethod = "SendPreBurnedFeeEDDTxID"

	// SendSignedTicketMethod represent SendSignedTicket method
	SendSignedTicketMethod = "SendSignedTicket"

	// SendSignedEDDTicketMethod represent SendSignedEDDTicket method
	SendSignedEDDTicketMethod = "SendSignedEDDTicket"

	// UploadImageWithThumbnailMethod represent UploadImageWithThumbnail method
	UploadImageWithThumbnailMethod = "UploadImageWithThumbnail"
	// UploadImageMethod represent UploadImage method
	UploadImageMethod = "UploadImage"
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
	*mocks.RegisterArtwork
	*mocks.DownloadArtwork
	*mocks.ProcessUserdata
	*mocks.ExternalDupeDetection
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                     t,
		Client:                &mocks.Client{},
		Connection:            &mocks.Connection{},
		RegisterArtwork:       &mocks.RegisterArtwork{},
		DownloadArtwork:       &mocks.DownloadArtwork{},
		ProcessUserdata:       &mocks.ProcessUserdata{},
		ExternalDupeDetection: &mocks.ExternalDupeDetection{},
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On(RegisterArtworkMethod).Return(client.RegisterArtwork)
	return client
}

// ListenOnExternalDupeDetection listening ExternalDupeDetection call
func (client *Client) ListenOnExternalDupeDetection() *Client {
	client.Connection.On(ExternalDupeDetectionMethod).Return(client.ExternalDupeDetection)
	return client
}

// ListenOnRegisterArtworkSendPreBurnedFeeTxID listening SendPreBurnedFeeTxIDMethod call
func (client *Client) ListenOnRegisterArtworkSendPreBurnedFeeTxID(txid string, err error) *Client {
	client.RegisterArtwork.On(SendPreBurnedFeeTxidMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnRegisterArtworkSendSignedTicket listening SendPreBurntFeeTxIDMethod call
func (client *Client) ListenOnRegisterArtworkSendSignedTicket(id int64, err error) *Client {
	client.RegisterArtwork.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnRegisterArtworkUploadImageWithThumbnail listening UploadImageWithThumbnail call
func (client *Client) ListenOnRegisterArtworkUploadImageWithThumbnail(retPreviewHash []byte,
	retMediumThumbnailHash []byte, retsmallThumbnailHash []byte, retErr error) *Client {

	client.RegisterArtwork.On(UploadImageWithThumbnailMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(retPreviewHash,
		retMediumThumbnailHash, retsmallThumbnailHash, retErr)

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

// ListenOnRegisterArtworkProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnRegisterArtworkProbeImage(arguments ...interface{}) *Client {
	client.RegisterArtwork.On(ProbeImageMethod, mock.Anything, mock.IsType(&artwork.File{})).Return(arguments...)
	return client
}

// AssertRegisterArtworkProbeImageCall assertion ProbeImage call
func (client *Client) AssertRegisterArtworkProbeImageCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, ProbeImageMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, ProbeImageMethod, expectedCalls)
	return client
}

// ListenOnRegisterArtworkSession listening Session call and returns error from args
func (client *Client) ListenOnRegisterArtworkSession(returnErr error) *Client {
	client.RegisterArtwork.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertRegisterArtworkSessionCall assertion Session Call
func (client *Client) AssertRegisterArtworkSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnRegisterArtworkAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnRegisterArtworkAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterArtwork.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertRegisterArtworkAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertRegisterArtworkAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnRegisterArtworkConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnRegisterArtworkConnectTo(returnErr error) *Client {
	client.RegisterArtwork.On(ConnectToMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// AssertRegisterArtworkConnectToCall assertion ConnectTo call
func (client *Client) AssertRegisterArtworkConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnRegisterArtworkSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnRegisterArtworkSessID(sessID string) *Client {
	client.RegisterArtwork.On(SessIDMethod).Return(sessID)
	return client
}

// AssertRegisterArtworkSessIDCall assertion SessID call
func (client *Client) AssertRegisterArtworkSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
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
	client.ProcessUserdata.On(ConnectToMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
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

// ListenOnExternalDupeDetectionSendPreBurnedFeeEDDTxID listening SendPreBurnedFeeEDDTxIDMethod call
func (client *Client) ListenOnExternalDupeDetectionSendPreBurnedFeeEDDTxID(txid string, err error) *Client {
	client.ExternalDupeDetection.On(SendPreBurnedFeeEDDTxIDMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnExternalDupeDetectionSendSignedEDDTicket listening SendSignedEDDTicketMethod call
func (client *Client) ListenOnExternalDupeDetectionSendSignedEDDTicket(id int64, err error) *Client {
	client.ExternalDupeDetection.On(SendSignedEDDTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnExternalDupeDetectionUploadImage listening UploadImage call
func (client *Client) ListenOnExternalDupeDetectionUploadImage(retErr error) *Client {

	client.ExternalDupeDetection.On(UploadImageMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(retErr)

	return client
}

// AssertExternalDupeDetectionCall assertion ExternalDupeDetection call
func (client *Client) AssertExternalDupeDetectionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, ExternalDupeDetectionMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(client.t, ExternalDupeDetectionMethod, expectedCalls)
	return client
}

// ListenOnExternalDupeDetectionProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnExternalDupeDetectionProbeImage(arguments ...interface{}) *Client {
	client.ExternalDupeDetection.On(ProbeImageMethod, mock.Anything, mock.IsType(&artwork.File{})).Return(arguments...)
	return client
}

// AssertExternalDupeDetectionProbeImageCall assertion ProbeImage call
func (client *Client) AssertExternalDupeDetectionProbeImageCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ExternalDupeDetection.AssertCalled(client.t, ProbeImageMethod, arguments...)
	}
	client.ExternalDupeDetection.AssertNumberOfCalls(client.t, ProbeImageMethod, expectedCalls)
	return client
}

// ListenOnExternalDupeDetectionSession listening Session call and returns error from args
func (client *Client) ListenOnExternalDupeDetectionSession(returnErr error) *Client {
	client.ExternalDupeDetection.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertExternalDupeDetectionSessionCall assertion Session Call
func (client *Client) AssertExternalDupeDetectionSessionCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ExternalDupeDetection.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.ExternalDupeDetection.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnExternalDupeDetectionAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnExternalDupeDetectionAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.ExternalDupeDetection.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertExternalDupeDetectionAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertExternalDupeDetectionAcceptedNodesCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ExternalDupeDetection.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.ExternalDupeDetection.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnExternalDupeDetectionConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnExternalDupeDetectionConnectTo(returnErr error) *Client {
	client.ExternalDupeDetection.On(ConnectToMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// AssertExternalDupeDetectionConnectToCall assertion ConnectTo call
func (client *Client) AssertExternalDupeDetectionConnectToCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ExternalDupeDetection.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.ExternalDupeDetection.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnExternalDupeDetectionSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnExternalDupeDetectionSessID(sessID string) *Client {
	client.ExternalDupeDetection.On(SessIDMethod).Return(sessID)
	return client
}

// AssertExternalDupeDetectionSessIDCall assertion SessID call
func (client *Client) AssertExternalDupeDetectionSessIDCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ExternalDupeDetection.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.ExternalDupeDetection.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}
