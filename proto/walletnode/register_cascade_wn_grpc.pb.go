// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: register_cascade_wn.proto

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

// RegisterCascadeClient is the client API for RegisterCascade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterCascadeClient interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_SessionClient, error)
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error)
	// MeshNodes informs to SNs other SNs on same meshNodes created for this registration request
	MeshNodes(ctx context.Context, in *MeshNodesRequest, opts ...grpc.CallOption) (*MeshNodesReply, error)
	// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
	SendRegMetadata(ctx context.Context, in *SendRegMetadataRequest, opts ...grpc.CallOption) (*SendRegMetadataReply, error)
	// Upload the asset for storing
	UploadAsset(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_UploadAssetClient, error)
	// SendArtTicket sends a signed art-ticket to the supernode.
	SendSignedActionTicket(ctx context.Context, in *SendSignedCascadeTicketRequest, opts ...grpc.CallOption) (*SendSignedActionTicketReply, error)
	// SendActionAc informs to SN that walletnode activated action_reg
	SendActionAct(ctx context.Context, in *SendActionActRequest, opts ...grpc.CallOption) (*SendActionActReply, error)
}

type registerCascadeClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterCascadeClient(cc grpc.ClientConnInterface) RegisterCascadeClient {
	return &registerCascadeClient{cc}
}

func (c *registerCascadeClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterCascade_ServiceDesc.Streams[0], "/walletnode.RegisterCascade/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerCascadeSessionClient{stream}
	return x, nil
}

type RegisterCascade_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type registerCascadeSessionClient struct {
	grpc.ClientStream
}

func (x *registerCascadeSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerCascadeSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerCascadeClient) AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error) {
	out := new(AcceptedNodesReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/AcceptedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerCascadeClient) ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error) {
	out := new(ConnectToReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/ConnectTo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerCascadeClient) MeshNodes(ctx context.Context, in *MeshNodesRequest, opts ...grpc.CallOption) (*MeshNodesReply, error) {
	out := new(MeshNodesReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/MeshNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerCascadeClient) SendRegMetadata(ctx context.Context, in *SendRegMetadataRequest, opts ...grpc.CallOption) (*SendRegMetadataReply, error) {
	out := new(SendRegMetadataReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/SendRegMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerCascadeClient) UploadAsset(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_UploadAssetClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterCascade_ServiceDesc.Streams[1], "/walletnode.RegisterCascade/UploadAsset", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerCascadeUploadAssetClient{stream}
	return x, nil
}

type RegisterCascade_UploadAssetClient interface {
	Send(*UploadAssetRequest) error
	CloseAndRecv() (*UploadAssetReply, error)
	grpc.ClientStream
}

type registerCascadeUploadAssetClient struct {
	grpc.ClientStream
}

func (x *registerCascadeUploadAssetClient) Send(m *UploadAssetRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerCascadeUploadAssetClient) CloseAndRecv() (*UploadAssetReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadAssetReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerCascadeClient) SendSignedActionTicket(ctx context.Context, in *SendSignedCascadeTicketRequest, opts ...grpc.CallOption) (*SendSignedActionTicketReply, error) {
	out := new(SendSignedActionTicketReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/SendSignedActionTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerCascadeClient) SendActionAct(ctx context.Context, in *SendActionActRequest, opts ...grpc.CallOption) (*SendActionActReply, error) {
	out := new(SendActionActReply)
	err := c.cc.Invoke(ctx, "/walletnode.RegisterCascade/SendActionAct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterCascadeServer is the server API for RegisterCascade service.
// All implementations must embed UnimplementedRegisterCascadeServer
// for forward compatibility
type RegisterCascadeServer interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the task on the supernode. By sending `sessID` in the Metadata.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterCascade_SessionServer) error
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error)
	// MeshNodes informs to SNs other SNs on same meshNodes created for this registration request
	MeshNodes(context.Context, *MeshNodesRequest) (*MeshNodesReply, error)
	// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
	SendRegMetadata(context.Context, *SendRegMetadataRequest) (*SendRegMetadataReply, error)
	// Upload the asset for storing
	UploadAsset(RegisterCascade_UploadAssetServer) error
	// SendArtTicket sends a signed art-ticket to the supernode.
	SendSignedActionTicket(context.Context, *SendSignedCascadeTicketRequest) (*SendSignedActionTicketReply, error)
	// SendActionAc informs to SN that walletnode activated action_reg
	SendActionAct(context.Context, *SendActionActRequest) (*SendActionActReply, error)
	mustEmbedUnimplementedRegisterCascadeServer()
}

// UnimplementedRegisterCascadeServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterCascadeServer struct {
}

func (UnimplementedRegisterCascadeServer) Session(RegisterCascade_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterCascadeServer) AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptedNodes not implemented")
}
func (UnimplementedRegisterCascadeServer) ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectTo not implemented")
}
func (UnimplementedRegisterCascadeServer) MeshNodes(context.Context, *MeshNodesRequest) (*MeshNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MeshNodes not implemented")
}
func (UnimplementedRegisterCascadeServer) SendRegMetadata(context.Context, *SendRegMetadataRequest) (*SendRegMetadataReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRegMetadata not implemented")
}
func (UnimplementedRegisterCascadeServer) UploadAsset(RegisterCascade_UploadAssetServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadAsset not implemented")
}
func (UnimplementedRegisterCascadeServer) SendSignedActionTicket(context.Context, *SendSignedCascadeTicketRequest) (*SendSignedActionTicketReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedActionTicket not implemented")
}
func (UnimplementedRegisterCascadeServer) SendActionAct(context.Context, *SendActionActRequest) (*SendActionActReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendActionAct not implemented")
}
func (UnimplementedRegisterCascadeServer) mustEmbedUnimplementedRegisterCascadeServer() {}

// UnsafeRegisterCascadeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterCascadeServer will
// result in compilation errors.
type UnsafeRegisterCascadeServer interface {
	mustEmbedUnimplementedRegisterCascadeServer()
}

func RegisterRegisterCascadeServer(s grpc.ServiceRegistrar, srv RegisterCascadeServer) {
	s.RegisterService(&RegisterCascade_ServiceDesc, srv)
}

func _RegisterCascade_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterCascadeServer).Session(&registerCascadeSessionServer{stream})
}

type RegisterCascade_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type registerCascadeSessionServer struct {
	grpc.ServerStream
}

func (x *registerCascadeSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerCascadeSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterCascade_AcceptedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptedNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).AcceptedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/AcceptedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).AcceptedNodes(ctx, req.(*AcceptedNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterCascade_ConnectTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectToRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).ConnectTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/ConnectTo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).ConnectTo(ctx, req.(*ConnectToRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterCascade_MeshNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MeshNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).MeshNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/MeshNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).MeshNodes(ctx, req.(*MeshNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterCascade_SendRegMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRegMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).SendRegMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/SendRegMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).SendRegMetadata(ctx, req.(*SendRegMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterCascade_UploadAsset_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterCascadeServer).UploadAsset(&registerCascadeUploadAssetServer{stream})
}

type RegisterCascade_UploadAssetServer interface {
	SendAndClose(*UploadAssetReply) error
	Recv() (*UploadAssetRequest, error)
	grpc.ServerStream
}

type registerCascadeUploadAssetServer struct {
	grpc.ServerStream
}

func (x *registerCascadeUploadAssetServer) SendAndClose(m *UploadAssetReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerCascadeUploadAssetServer) Recv() (*UploadAssetRequest, error) {
	m := new(UploadAssetRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterCascade_SendSignedActionTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedCascadeTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).SendSignedActionTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/SendSignedActionTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).SendSignedActionTicket(ctx, req.(*SendSignedCascadeTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterCascade_SendActionAct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendActionActRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).SendActionAct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.RegisterCascade/SendActionAct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).SendActionAct(ctx, req.(*SendActionActRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterCascade_ServiceDesc is the grpc.ServiceDesc for RegisterCascade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterCascade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "walletnode.RegisterCascade",
	HandlerType: (*RegisterCascadeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptedNodes",
			Handler:    _RegisterCascade_AcceptedNodes_Handler,
		},
		{
			MethodName: "ConnectTo",
			Handler:    _RegisterCascade_ConnectTo_Handler,
		},
		{
			MethodName: "MeshNodes",
			Handler:    _RegisterCascade_MeshNodes_Handler,
		},
		{
			MethodName: "SendRegMetadata",
			Handler:    _RegisterCascade_SendRegMetadata_Handler,
		},
		{
			MethodName: "SendSignedActionTicket",
			Handler:    _RegisterCascade_SendSignedActionTicket_Handler,
		},
		{
			MethodName: "SendActionAct",
			Handler:    _RegisterCascade_SendActionAct_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterCascade_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadAsset",
			Handler:       _RegisterCascade_UploadAsset_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "register_cascade_wn.proto",
}
