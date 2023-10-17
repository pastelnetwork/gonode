// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: register_sense_sn.proto

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

const (
	RegisterSense_Session_FullMethodName                     = "/supernode.RegisterSense/Session"
	RegisterSense_SendSignedDDAndFingerprints_FullMethodName = "/supernode.RegisterSense/SendSignedDDAndFingerprints"
	RegisterSense_SendSenseTicketSignature_FullMethodName    = "/supernode.RegisterSense/SendSenseTicketSignature"
)

// RegisterSenseClient is the client API for RegisterSense service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterSenseClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_SessionClient, error)
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error)
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendSenseTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error)
}

type registerSenseClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterSenseClient(cc grpc.ClientConnInterface) RegisterSenseClient {
	return &registerSenseClient{cc}
}

func (c *registerSenseClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterSense_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterSense_ServiceDesc.Streams[0], RegisterSense_Session_FullMethodName, opts...)
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

func (c *registerSenseClient) SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error) {
	out := new(SendSignedDDAndFingerprintsReply)
	err := c.cc.Invoke(ctx, RegisterSense_SendSignedDDAndFingerprints_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerSenseClient) SendSenseTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error) {
	out := new(SendTicketSignatureReply)
	err := c.cc.Invoke(ctx, RegisterSense_SendSenseTicketSignature_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterSenseServer is the server API for RegisterSense service.
// All implementations must embed UnimplementedRegisterSenseServer
// for forward compatibility
type RegisterSenseServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterSense_SessionServer) error
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error)
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendSenseTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error)
	mustEmbedUnimplementedRegisterSenseServer()
}

// UnimplementedRegisterSenseServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterSenseServer struct {
}

func (UnimplementedRegisterSenseServer) Session(RegisterSense_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterSenseServer) SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedDDAndFingerprints not implemented")
}
func (UnimplementedRegisterSenseServer) SendSenseTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSenseTicketSignature not implemented")
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

func _RegisterSense_SendSignedDDAndFingerprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedDDAndFingerprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).SendSignedDDAndFingerprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegisterSense_SendSignedDDAndFingerprints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).SendSignedDDAndFingerprints(ctx, req.(*SendSignedDDAndFingerprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterSense_SendSenseTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterSenseServer).SendSenseTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegisterSense_SendSenseTicketSignature_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterSenseServer).SendSenseTicketSignature(ctx, req.(*SendTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterSense_ServiceDesc is the grpc.ServiceDesc for RegisterSense service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterSense_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.RegisterSense",
	HandlerType: (*RegisterSenseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSignedDDAndFingerprints",
			Handler:    _RegisterSense_SendSignedDDAndFingerprints_Handler,
		},
		{
			MethodName: "SendSenseTicketSignature",
			Handler:    _RegisterSense_SendSenseTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterSense_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "register_sense_sn.proto",
}
