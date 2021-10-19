// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package walletnode

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExternalDupeDetectionClient is the client API for ExternalDupeDetection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExternalDupeDetectionClient interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_SessionClient, error)
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error)
	// ProbeImage uploads the resampled image compute and return a fingerpirnt.
	ProbeImage(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_ProbeImageClient, error)
	// SendSignedEDDTicket sends a signed art-ticket to the supernode.
	SendSignedEDDTicket(ctx context.Context, in *SendSignedEDDTicketRequest, opts ...grpc.CallOption) (*SendSignedEDDTicketReply, error)
	// SendPreBurnedFeeEDDTxID sends tx_id of 10% burnt transaction fee to the supernode.
	SendPreBurnedFeeEDDTxID(ctx context.Context, in *SendPreBurnedFeeEDDTxIDRequest, opts ...grpc.CallOption) (*SendPreBurnedFeeEDDTxIDReply, error)
	// SendTicket sends a ticket to the supernode.
	SendTicket(ctx context.Context, in *SendTicketRequest, opts ...grpc.CallOption) (*SendTicketReply, error)
	// Upload the image after pq signature is appended along with its thumbnail coordinates
	UploadImage(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_UploadImageClient, error)
}

type externalDupeDetectionClient struct {
	cc grpc.ClientConnInterface
}

func NewExternalDupeDetectionClient(cc grpc.ClientConnInterface) ExternalDupeDetectionClient {
	return &externalDupeDetectionClient{cc}
}

func (c *externalDupeDetectionClient) Session(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExternalDupeDetection_ServiceDesc.Streams[0], "/walletnode.ExternalDupeDetection/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &externalDupeDetectionSessionClient{stream}
	return x, nil
}

type ExternalDupeDetection_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type externalDupeDetectionSessionClient struct {
	grpc.ClientStream
}

func (x *externalDupeDetectionSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *externalDupeDetectionSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *externalDupeDetectionClient) AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error) {
	out := new(AcceptedNodesReply)
	err := c.cc.Invoke(ctx, "/walletnode.ExternalDupeDetection/AcceptedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalDupeDetectionClient) ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error) {
	out := new(ConnectToReply)
	err := c.cc.Invoke(ctx, "/walletnode.ExternalDupeDetection/ConnectTo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalDupeDetectionClient) ProbeImage(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_ProbeImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExternalDupeDetection_ServiceDesc.Streams[1], "/walletnode.ExternalDupeDetection/ProbeImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &externalDupeDetectionProbeImageClient{stream}
	return x, nil
}

type ExternalDupeDetection_ProbeImageClient interface {
	Send(*ProbeImageRequest) error
	CloseAndRecv() (*ProbeImageReply, error)
	grpc.ClientStream
}

type externalDupeDetectionProbeImageClient struct {
	grpc.ClientStream
}

func (x *externalDupeDetectionProbeImageClient) Send(m *ProbeImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *externalDupeDetectionProbeImageClient) CloseAndRecv() (*ProbeImageReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ProbeImageReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *externalDupeDetectionClient) SendSignedEDDTicket(ctx context.Context, in *SendSignedEDDTicketRequest, opts ...grpc.CallOption) (*SendSignedEDDTicketReply, error) {
	out := new(SendSignedEDDTicketReply)
	err := c.cc.Invoke(ctx, "/walletnode.ExternalDupeDetection/SendSignedEDDTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalDupeDetectionClient) SendPreBurnedFeeEDDTxID(ctx context.Context, in *SendPreBurnedFeeEDDTxIDRequest, opts ...grpc.CallOption) (*SendPreBurnedFeeEDDTxIDReply, error) {
	out := new(SendPreBurnedFeeEDDTxIDReply)
	err := c.cc.Invoke(ctx, "/walletnode.ExternalDupeDetection/SendPreBurnedFeeEDDTxID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalDupeDetectionClient) SendTicket(ctx context.Context, in *SendTicketRequest, opts ...grpc.CallOption) (*SendTicketReply, error) {
	out := new(SendTicketReply)
	err := c.cc.Invoke(ctx, "/walletnode.ExternalDupeDetection/SendTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalDupeDetectionClient) UploadImage(ctx context.Context, opts ...grpc.CallOption) (ExternalDupeDetection_UploadImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExternalDupeDetection_ServiceDesc.Streams[2], "/walletnode.ExternalDupeDetection/UploadImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &externalDupeDetectionUploadImageClient{stream}
	return x, nil
}

type ExternalDupeDetection_UploadImageClient interface {
	Send(*UploadImageRequest) error
	CloseAndRecv() (*UploadImageReply, error)
	grpc.ClientStream
}

type externalDupeDetectionUploadImageClient struct {
	grpc.ClientStream
}

func (x *externalDupeDetectionUploadImageClient) Send(m *UploadImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *externalDupeDetectionUploadImageClient) CloseAndRecv() (*UploadImageReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadImageReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExternalDupeDetectionServer is the server API for ExternalDupeDetection service.
// All implementations must embed UnimplementedExternalDupeDetectionServer
// for forward compatibility
type ExternalDupeDetectionServer interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ExternalDupeDetection_SessionServer) error
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error)
	// ProbeImage uploads the resampled image compute and return a fingerpirnt.
	ProbeImage(ExternalDupeDetection_ProbeImageServer) error
	// SendSignedEDDTicket sends a signed art-ticket to the supernode.
	SendSignedEDDTicket(context.Context, *SendSignedEDDTicketRequest) (*SendSignedEDDTicketReply, error)
	// SendPreBurnedFeeEDDTxID sends tx_id of 10% burnt transaction fee to the supernode.
	SendPreBurnedFeeEDDTxID(context.Context, *SendPreBurnedFeeEDDTxIDRequest) (*SendPreBurnedFeeEDDTxIDReply, error)
	// SendTicket sends a ticket to the supernode.
	SendTicket(context.Context, *SendTicketRequest) (*SendTicketReply, error)
	// Upload the image after pq signature is appended along with its thumbnail coordinates
	UploadImage(ExternalDupeDetection_UploadImageServer) error
	mustEmbedUnimplementedExternalDupeDetectionServer()
}

// UnimplementedExternalDupeDetectionServer must be embedded to have forward compatible implementations.
type UnimplementedExternalDupeDetectionServer struct {
}

func (UnimplementedExternalDupeDetectionServer) Session(ExternalDupeDetection_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedExternalDupeDetectionServer) AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptedNodes not implemented")
}
func (UnimplementedExternalDupeDetectionServer) ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectTo not implemented")
}
func (UnimplementedExternalDupeDetectionServer) ProbeImage(ExternalDupeDetection_ProbeImageServer) error {
	return status.Errorf(codes.Unimplemented, "method ProbeImage not implemented")
}
func (UnimplementedExternalDupeDetectionServer) SendSignedEDDTicket(context.Context, *SendSignedEDDTicketRequest) (*SendSignedEDDTicketReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedEDDTicket not implemented")
}
func (UnimplementedExternalDupeDetectionServer) SendPreBurnedFeeEDDTxID(context.Context, *SendPreBurnedFeeEDDTxIDRequest) (*SendPreBurnedFeeEDDTxIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPreBurnedFeeEDDTxID not implemented")
}
func (UnimplementedExternalDupeDetectionServer) SendTicket(context.Context, *SendTicketRequest) (*SendTicketReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTicket not implemented")
}
func (UnimplementedExternalDupeDetectionServer) UploadImage(ExternalDupeDetection_UploadImageServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (UnimplementedExternalDupeDetectionServer) mustEmbedUnimplementedExternalDupeDetectionServer() {}

// UnsafeExternalDupeDetectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExternalDupeDetectionServer will
// result in compilation errors.
type UnsafeExternalDupeDetectionServer interface {
	mustEmbedUnimplementedExternalDupeDetectionServer()
}

func RegisterExternalDupeDetectionServer(s grpc.ServiceRegistrar, srv ExternalDupeDetectionServer) {
	s.RegisterService(&ExternalDupeDetection_ServiceDesc, srv)
}

func _ExternalDupeDetection_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExternalDupeDetectionServer).Session(&externalDupeDetectionSessionServer{stream})
}

type ExternalDupeDetection_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type externalDupeDetectionSessionServer struct {
	grpc.ServerStream
}

func (x *externalDupeDetectionSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *externalDupeDetectionSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExternalDupeDetection_AcceptedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptedNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalDupeDetectionServer).AcceptedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.ExternalDupeDetection/AcceptedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalDupeDetectionServer).AcceptedNodes(ctx, req.(*AcceptedNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalDupeDetection_ConnectTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectToRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalDupeDetectionServer).ConnectTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.ExternalDupeDetection/ConnectTo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalDupeDetectionServer).ConnectTo(ctx, req.(*ConnectToRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalDupeDetection_ProbeImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExternalDupeDetectionServer).ProbeImage(&externalDupeDetectionProbeImageServer{stream})
}

type ExternalDupeDetection_ProbeImageServer interface {
	SendAndClose(*ProbeImageReply) error
	Recv() (*ProbeImageRequest, error)
	grpc.ServerStream
}

type externalDupeDetectionProbeImageServer struct {
	grpc.ServerStream
}

func (x *externalDupeDetectionProbeImageServer) SendAndClose(m *ProbeImageReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *externalDupeDetectionProbeImageServer) Recv() (*ProbeImageRequest, error) {
	m := new(ProbeImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExternalDupeDetection_SendSignedEDDTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedEDDTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalDupeDetectionServer).SendSignedEDDTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.ExternalDupeDetection/SendSignedEDDTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalDupeDetectionServer).SendSignedEDDTicket(ctx, req.(*SendSignedEDDTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalDupeDetection_SendPreBurnedFeeEDDTxID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPreBurnedFeeEDDTxIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalDupeDetectionServer).SendPreBurnedFeeEDDTxID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.ExternalDupeDetection/SendPreBurnedFeeEDDTxID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalDupeDetectionServer).SendPreBurnedFeeEDDTxID(ctx, req.(*SendPreBurnedFeeEDDTxIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalDupeDetection_SendTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalDupeDetectionServer).SendTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.ExternalDupeDetection/SendTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalDupeDetectionServer).SendTicket(ctx, req.(*SendTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalDupeDetection_UploadImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExternalDupeDetectionServer).UploadImage(&externalDupeDetectionUploadImageServer{stream})
}

type ExternalDupeDetection_UploadImageServer interface {
	SendAndClose(*UploadImageReply) error
	Recv() (*UploadImageRequest, error)
	grpc.ServerStream
}

type externalDupeDetectionUploadImageServer struct {
	grpc.ServerStream
}

func (x *externalDupeDetectionUploadImageServer) SendAndClose(m *UploadImageReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *externalDupeDetectionUploadImageServer) Recv() (*UploadImageRequest, error) {
	m := new(UploadImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExternalDupeDetection_ServiceDesc is the grpc.ServiceDesc for ExternalDupeDetection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExternalDupeDetection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "walletnode.ExternalDupeDetection",
	HandlerType: (*ExternalDupeDetectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptedNodes",
			Handler:    _ExternalDupeDetection_AcceptedNodes_Handler,
		},
		{
			MethodName: "ConnectTo",
			Handler:    _ExternalDupeDetection_ConnectTo_Handler,
		},
		{
			MethodName: "SendSignedEDDTicket",
			Handler:    _ExternalDupeDetection_SendSignedEDDTicket_Handler,
		},
		{
			MethodName: "SendPreBurnedFeeEDDTxID",
			Handler:    _ExternalDupeDetection_SendPreBurnedFeeEDDTxID_Handler,
		},
		{
			MethodName: "SendTicket",
			Handler:    _ExternalDupeDetection_SendTicket_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _ExternalDupeDetection_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ProbeImage",
			Handler:       _ExternalDupeDetection_ProbeImage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadImage",
			Handler:       _ExternalDupeDetection_UploadImage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "external_dupe_detection.proto",
}