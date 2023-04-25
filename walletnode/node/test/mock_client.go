package test

import (
	"context"
	"fmt"
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
	RegisterNftMethod = "RegisterNftInterface"

	// DownloadNftMethod represent DownloadNftInterface name method
	DownloadNftMethod = "DownloadNftInterface"

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

	SendTicketForSignatureMethod = "SendTicketForSignature"

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

	// GetDupeDetectionDBHashMethod  is method name for GetDupeDetectionDBHash
	GetDupeDetectionDBHashMethod = "GetDupeDetectionDBHash"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.RegisterNftInterface
	*mocks.DownloadNftInterface
	*mocks.ProcessUserdataInterface
	*mocks.RegisterSenseInterface
	*mocks.RegisterCascadeInterface
	*mocks.RegisterCollectionInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                           t,
		ClientInterface:             &mocks.ClientInterface{},
		ConnectionInterface:         &mocks.ConnectionInterface{},
		RegisterNftInterface:        &mocks.RegisterNftInterface{},
		DownloadNftInterface:        &mocks.DownloadNftInterface{},
		ProcessUserdataInterface:    &mocks.ProcessUserdataInterface{},
		RegisterSenseInterface:      &mocks.RegisterSenseInterface{},
		RegisterCascadeInterface:    &mocks.RegisterCascadeInterface{},
		RegisterCollectionInterface: &mocks.RegisterCollectionInterface{},
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

// ListenOnRegisterGetDupeDetectionDBHash listening GetDupeDetectionDBHash call
func (client *Client) ListenOnRegisterGetDupeDetectionDBHash(hash string, err error) *Client {
	client.RegisterNftInterface.On(GetDupeDetectionDBHashMethod, mock.Anything, mock.Anything).Return(hash, err)
	return client
}

// ListenOnSenseGetDupeDetectionDBHash listening GetDupeDetectionDBHash call
func (client *Client) ListenOnSenseGetDupeDetectionDBHash(hash string, err error) *Client {
	client.RegisterSenseInterface.On(GetDupeDetectionDBHashMethod, mock.Anything, mock.Anything).Return(hash, err)
	return client
}

// ListenOnSendSignedTicket listening SendPreBurntFeeTxIdMethod call
func (client *Client) ListenOnSendSignedTicket(id int64, err error) *Client {
	client.RegisterNftInterface.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(id, err)

	client.RegisterSenseInterface.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(fmt.Sprint(id), err)

	client.RegisterCascadeInterface.On(SendSignedTicketMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(fmt.Sprint(id), err)

	return client
}

// ListenOnSendTicketForSignature listening SendTicketForSignature call
func (client *Client) ListenOnSendTicketForSignature(id int64, err error) *Client {
	client.RegisterCollectionInterface.On(SendTicketForSignatureMethod, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(fmt.Sprint(id), err)

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

// ListenOnDownloadNft listening DownloadNftInterface call
func (client *Client) ListenOnDownloadNft() *Client {
	client.ConnectionInterface.On(DownloadNftMethod).Return(client.DownloadNftInterface)
	return client
}

// AssertDownloadNftCall assertion DownloadNftInterface call
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

// ListenOnMeshNodes listening MeshNodes call and returns args value
func (client *Client) ListenOnMeshNodes(arguments ...interface{}) *Client {
	client.RegisterNftInterface.On(MeshNodesMethod, mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnProbeImage(arguments ...interface{}) *Client {
	client.RegisterNftInterface.On(ProbeImageMethod, mock.Anything, mock.Anything).Return(arguments...)
	client.RegisterSenseInterface.On(ProbeImageMethod, mock.Anything, mock.Anything).Return(arguments...)

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

// ListenOnUploadAsset listening UploadAsset call and returns args value
func (client *Client) ListenOnUploadAsset(arguments ...interface{}) *Client {
	client.RegisterCascadeInterface.On(UploadAssetMethod, mock.Anything, mock.Anything).Return(arguments...)
	client.RegisterCascadeInterface.On(UploadAssetMethod, mock.Anything, mock.Anything).Return(arguments...)

	return client
}

// AssertUploadAssetCall assertion UploadAsset call
func (client *Client) AssertUploadAssetCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterCascadeInterface.AssertCalled(client.t, UploadAssetMethod, arguments...)
	}
	client.RegisterCascadeInterface.AssertNumberOfCalls(client.t, UploadAssetMethod, expectedCalls)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterNftInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	client.RegisterSenseInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	client.RegisterCascadeInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	client.RegisterCollectionInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)

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
	client.RegisterSenseInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	client.RegisterCascadeInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	client.RegisterCollectionInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)

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
	client.RegisterSenseInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	client.RegisterCascadeInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	client.RegisterCollectionInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)

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
	client.RegisterSenseInterface.On(SessIDMethod).Return(sessID)
	client.RegisterCascadeInterface.On(SessIDMethod).Return(sessID)
	client.RegisterCollectionInterface.On(SessIDMethod).Return(sessID)

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

// ListenOnDownload listening Download call and returns args value
func (client *Client) ListenOnDownload(arguments ...interface{}) *Client {
	client.DownloadNftInterface.On(DownloadMethod, mock.Anything,
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string("")),
		mock.IsType(string(""))).Return(arguments...)
	return client
}

// ListenOnDownloadThumbnail listening DownloadThumbnail call and returns args value
func (client *Client) ListenOnDownloadThumbnail(arguments ...interface{}) *Client {
	client.DownloadNftInterface.On(DownloadThumbnailMethod, mock.Anything,
		mock.Anything, mock.Anything).Return(arguments...)
	return client
}

// ListenOnDownloadDDAndFP listens for a ListenOnDownloadDDAndFP call and returns args value
func (client *Client) ListenOnDownloadDDAndFP(arguments ...interface{}) *Client {
	client.DownloadNftInterface.On(DownloadDDAndFPMethod, mock.Anything,
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

// ListenOnProcessUserdata listening ProcessUserdataInterface call
func (client *Client) ListenOnProcessUserdata() *Client {
	client.ConnectionInterface.On(ProcessUserdataMethod).Return(client.ProcessUserdataInterface)
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

// ListenOnSessionUserdata listening Session call and returns error from args
func (client *Client) ListenOnSessionUserdata(returnErr error) *Client {
	client.ProcessUserdataInterface.On(SessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCallUserdata assertion Session Call
func (client *Client) AssertSessionCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, SessionMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodesUserdata listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodesUserdata(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.ProcessUserdataInterface.On(AcceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCallUserdata assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, AcceptedNodesMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, AcceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectToUserdata listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectToUserdata(returnErr error) *Client {
	client.ProcessUserdataInterface.On(ConnectToMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}

// AssertConnectToCallUserdata assertion ConnectTo call
func (client *Client) AssertConnectToCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, ConnectToMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, ConnectToMethod, expectedCalls)
	return client
}

// ListenOnSessIDUserdata listening SessID call and returns sessID from args
func (client *Client) ListenOnSessIDUserdata(sessID string) *Client {
	client.ProcessUserdataInterface.On(SessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCallUserdata assertion SessID call
func (client *Client) AssertSessIDCallUserdata(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ProcessUserdataInterface.AssertCalled(client.t, SessIDMethod, arguments...)
	}
	client.ProcessUserdataInterface.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnSendActionAct listening RegisterNftInterface call
func (client *Client) ListenOnSendActionAct(retErr error) *Client {
	client.RegisterSenseInterface.On(SendActionActMethod, mock.Anything, mock.Anything).Return(retErr)
	client.RegisterCascadeInterface.On(SendActionActMethod, mock.Anything, mock.Anything).Return(retErr)

	return client
}
