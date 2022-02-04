// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: register_nft_sn.proto

package supernode

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

// RegisterNftClient is the client API for RegisterNft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterNftClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterNft_SessionClient, error)
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error)
	// SendNftTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendNftTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error)
}

type registerNftClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterNftClient(cc grpc.ClientConnInterface) RegisterNftClient {
	return &registerNftClient{cc}
}

func (c *registerNftClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterNft_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterNft_ServiceDesc.Streams[0], "/supernode.RegisterNft/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerNftSessionClient{stream}
	return x, nil
}

type RegisterNft_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type registerNftSessionClient struct {
	grpc.ClientStream
}

func (x *registerNftSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerNftSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerNftClient) SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error) {
	out := new(SendSignedDDAndFingerprintsReply)
	err := c.cc.Invoke(ctx, "/supernode.RegisterNft/SendSignedDDAndFingerprints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerNftClient) SendNftTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error) {
	out := new(SendTicketSignatureReply)
	err := c.cc.Invoke(ctx, "/supernode.RegisterNft/SendNftTicketSignature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterNftServer is the server API for RegisterNft service.
// All implementations must embed UnimplementedRegisterNftServer
// for forward compatibility
type RegisterNftServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterNft_SessionServer) error
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error)
	// SendNftTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendNftTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error)
	mustEmbedUnimplementedRegisterNftServer()
}

// UnimplementedRegisterNftServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterNftServer struct {
}

func (UnimplementedRegisterNftServer) Session(RegisterNft_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterNftServer) SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedDDAndFingerprints not implemented")
}
func (UnimplementedRegisterNftServer) SendNftTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNftTicketSignature not implemented")
}
func (UnimplementedRegisterNftServer) mustEmbedUnimplementedRegisterNftServer() {}

// UnsafeRegisterNftServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterNftServer will
// result in compilation errors.
type UnsafeRegisterNftServer interface {
	mustEmbedUnimplementedRegisterNftServer()
}

func RegisterRegisterNftServer(s grpc.ServiceRegistrar, srv RegisterNftServer) {
	s.RegisterService(&RegisterNft_ServiceDesc, srv)
}

func _RegisterNft_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterNftServer).Session(&registerNftSessionServer{stream})
}

type RegisterNft_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type registerNftSessionServer struct {
	grpc.ServerStream
}

func (x *registerNftSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerNftSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterNft_SendSignedDDAndFingerprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedDDAndFingerprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterNftServer).SendSignedDDAndFingerprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supernode.RegisterNft/SendSignedDDAndFingerprints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterNftServer).SendSignedDDAndFingerprints(ctx, req.(*SendSignedDDAndFingerprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterNft_SendNftTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterNftServer).SendNftTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supernode.RegisterNft/SendNftTicketSignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterNftServer).SendNftTicketSignature(ctx, req.(*SendTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterNft_ServiceDesc is the grpc.ServiceDesc for RegisterNft service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterNft_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.RegisterNft",
	HandlerType: (*RegisterNftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSignedDDAndFingerprints",
			Handler:    _RegisterNft_SendSignedDDAndFingerprints_Handler,
		},
		{
			MethodName: "SendNftTicketSignature",
			Handler:    _RegisterNft_SendNftTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterNft_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "register_nft_sn.proto",
}
