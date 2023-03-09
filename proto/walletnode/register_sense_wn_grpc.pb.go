// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: register_sense_wn.proto

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

// RegisterSenseClient is the client API for RegisterSense service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterSenseClient interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_SessionClient, error)
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error)
	// MeshNodes informs to SNs other SNs on same meshNodes created for this registration request
	MeshNodes(ctx context.Context, in *MeshNodesRequest, opts ...grpc.CallOption) (*MeshNodesReply, error)
	// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
	SendRegMetadata(ctx context.Context, in *SendRegMetadataRequest, opts ...grpc.CallOption) (*SendRegMetadataReply, error)
	// ProbeImage uploads the resampled image compute/burn txid and return a fingerpirnt and MN signature.
	ProbeImage(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_ProbeImageClient, error)
	// SendArtTicket sends a signed art-ticket to the supernode.
	SendSignedActionTicket(ctx context.Context, in *SendSignedSenseTicketRequest, opts ...grpc.CallOption) (*SendSignedActionTicketReply, error)
	// SendActionAc informs to SN that walletnode activated action_reg
	SendActionAct(ctx context.Context, in *SendActionActRequest, opts ...grpc.CallOption) (*SendActionActReply, error)
	// GetDDDatabaseHash returns hash of dupe detection database hash
	GetDDDatabaseHash(ctx context.Context, in *GetDBHashRequest, opts ...grpc.CallOption) (*DBHashReply, error)
}

type registerSenseClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterSenseClient(cc grpc.ClientConnInterface) RegisterSenseClient {
	return &registerSenseClient{cc}
}

func (c *registerSenseClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterSense_ServiceDesc.Streams[0], "/walletnode.RegisterSense/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerSenseSessionClient{stream}
	return x, nil
}

type RegisterSense_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type registerSenseSessionClient struct {
	grpc.ClientStream
}

func (x *registerSenseSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerSenseSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerSenseClient) AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error) {
	out := new(AcceptedNodesReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/AcceptedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error) {
	out := new(ConnectToReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/ConnectTo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) MeshNodes(ctx context.Context, in *MeshNodesRequest, opts ...grpc.CallOption) (*MeshNodesReply, error) {
	out := new(MeshNodesReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/MeshNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) SendRegMetadata(ctx context.Context, in *SendRegMetadataRequest, opts ...grpc.CallOption) (*SendRegMetadataReply, error) {
	out := new(SendRegMetadataReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/SendRegMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) ProbeImage(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_ProbeImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterSense_ServiceDesc.Streams[1], "/walletnode.RegisterSense/ProbeImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerSenseProbeImageClient{stream}
	return x, nil
}

type RegisterSense_ProbeImageClient interface {
	Send(*ProbeImageRequest) error
	CloseAndRecv() (*ProbeImageReply, error)
	grpc.ClientStream
}

type registerSenseProbeImageClient struct {
	grpc.ClientStream
}

func (x *registerSenseProbeImageClient) Send(m *ProbeImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerSenseProbeImageClient) CloseAndRecv() (*ProbeImageReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ProbeImageReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerSenseClient) SendSignedActionTicket(ctx context.Context, in *SendSignedSenseTicketRequest, opts ...grpc.CallOption) (*SendSignedActionTicketReply, error) {
	out := new(SendSignedActionTicketReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/SendSignedActionTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) SendActionAct(ctx context.Context, in *SendActionActRequest, opts ...grpc.CallOption) (*SendActionActReply, error) {
	out := new(SendActionActReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/SendActionAct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) GetDDDatabaseHash(ctx context.Context, in *GetDBHashRequest, opts ...grpc.CallOption) (*DBHashReply, error) {
	out := new(DBHashReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterSense/GetDDDatabaseHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterSenseServer is the server API for RegisterSense service.
// All implementations must embed UnimplementedRegisterSenseServer
// for forward compatibility
type RegisterSenseServer interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterSense_SessionServer) error
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error)
	// MeshNodes informs to SNs other SNs on same meshNodes created for this registration request
	MeshNodes(context.Context, *MeshNodesRequest) (*MeshNodesReply, error)
	// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
	SendRegMetadata(context.Context, *SendRegMetadataRequest) (*SendRegMetadataReply, error)
	// ProbeImage uploads the resampled image compute/burn txid and return a fingerpirnt and MN signature.
	ProbeImage(RegisterSense_ProbeImageServer) error
	// SendArtTicket sends a signed art-ticket to the supernode.
	SendSignedActionTicket(context.Context, *SendSignedSenseTicketRequest) (*SendSignedActionTicketReply, error)
	// SendActionAc informs to SN that walletnode activated action_reg
	SendActionAct(context.Context, *SendActionActRequest) (*SendActionActReply, error)
	// GetDDDatabaseHash returns hash of dupe detection database hash
	GetDDDatabaseHash(context.Context, *GetDBHashRequest) (*DBHashReply, error)
	mustEmbedUnimplementedRegisterSenseServer()
}

// UnimplementedRegisterSenseServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterSenseServer struct {
}

func (UnimplementedRegisterSenseServer) Session(RegisterSense_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterSenseServer) AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptedNodes not implemented")
}
func (UnimplementedRegisterSenseServer) ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectTo not implemented")
}
func (UnimplementedRegisterSenseServer) MeshNodes(context.Context, *MeshNodesRequest) (*MeshNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MeshNodes not implemented")
}
func (UnimplementedRegisterSenseServer) SendRegMetadata(context.Context, *SendRegMetadataRequest) (*SendRegMetadataReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRegMetadata not implemented")
}
func (UnimplementedRegisterSenseServer) ProbeImage(RegisterSense_ProbeImageServer) error {
	return status.Errorf(codes.Unimplemented, "method ProbeImage not implemented")
}
func (UnimplementedRegisterSenseServer) SendSignedActionTicket(context.Context, *SendSignedSenseTicketRequest) (*SendSignedActionTicketReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedActionTicket not implemented")
}
func (UnimplementedRegisterSenseServer) SendActionAct(context.Context, *SendActionActRequest) (*SendActionActReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendActionAct not implemented")
}
func (UnimplementedRegisterSenseServer) GetDDDatabaseHash(context.Context, *GetDBHashRequest) (*DBHashReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDDDatabaseHash not implemented")
}
func (UnimplementedRegisterSenseServer) mustEmbedUnimplementedRegisterSenseServer() {}

// UnsafeRegisterSenseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterSenseServer will
// result in compilation errors.
type UnsafeRegisterSenseServer interface {
	mustEmbedUnimplementedRegisterSenseServer()
}

func RegisterRegisterSenseServer(s grpc.ServiceRegistrar, srv RegisterSenseServer) {
	s.RegisterService(&RegisterSense_ServiceDesc, srv)
}

func _RegisterSense_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterSenseServer).Session(&registerSenseSessionServer{stream})
}

type RegisterSense_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type registerSenseSessionServer struct {
	grpc.ServerStream
}

func (x *registerSenseSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerSenseSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterSense_AcceptedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptedNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).AcceptedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/AcceptedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).AcceptedNodes(ctx, req.(*AcceptedNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_ConnectTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectToRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).ConnectTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/ConnectTo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).ConnectTo(ctx, req.(*ConnectToRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_MeshNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MeshNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).MeshNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/MeshNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).MeshNodes(ctx, req.(*MeshNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_SendRegMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRegMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).SendRegMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/SendRegMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).SendRegMetadata(ctx, req.(*SendRegMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_ProbeImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterSenseServer).ProbeImage(&registerSenseProbeImageServer{stream})
}

type RegisterSense_ProbeImageServer interface {
	SendAndClose(*ProbeImageReply) error
	Recv() (*ProbeImageRequest, error)
	grpc.ServerStream
}

type registerSenseProbeImageServer struct {
	grpc.ServerStream
}

func (x *registerSenseProbeImageServer) SendAndClose(m *ProbeImageReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerSenseProbeImageServer) Recv() (*ProbeImageRequest, error) {
	m := new(ProbeImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterSense_SendSignedActionTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedSenseTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).SendSignedActionTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/SendSignedActionTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).SendSignedActionTicket(ctx, req.(*SendSignedSenseTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_SendActionAct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendActionActRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).SendActionAct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/SendActionAct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).SendActionAct(ctx, req.(*SendActionActRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_GetDDDatabaseHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDBHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).GetDDDatabaseHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterSense/GetDDDatabaseHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).GetDDDatabaseHash(ctx, req.(*GetDBHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterSense_ServiceDesc is the grpc.ServiceDesc for RegisterSense service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterSense_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "walletnode.RegisterSense",
	HandlerType: (*RegisterSenseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptedNodes",
			Handler:    _RegisterSense_AcceptedNodes_Handler,
		},
		{
			MethodName: "ConnectTo",
			Handler:    _RegisterSense_ConnectTo_Handler,
		},
		{
			MethodName: "MeshNodes",
			Handler:    _RegisterSense_MeshNodes_Handler,
		},
		{
			MethodName: "SendRegMetadata",
			Handler:    _RegisterSense_SendRegMetadata_Handler,
		},
		{
			MethodName: "SendSignedActionTicket",
			Handler:    _RegisterSense_SendSignedActionTicket_Handler,
		},
		{
			MethodName: "SendActionAct",
			Handler:    _RegisterSense_SendActionAct_Handler,
		},
		{
			MethodName: "GetDDDatabaseHash",
			Handler:    _RegisterSense_GetDDDatabaseHash_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterSense_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ProbeImage",
			Handler:       _RegisterSense_ProbeImage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "register_sense_wn.proto",
}
