// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: register_cascade_sn.proto

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
	RegisterCascade_Session_FullMethodName                    = "/supernode.RegisterCascade/Session"
	RegisterCascade_SendCascadeTicketSignature_FullMethodName = "/supernode.RegisterCascade/SendCascadeTicketSignature"
)

// RegisterCascadeClient is the client API for RegisterCascade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterCascadeClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_SessionClient, error)
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendCascadeTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error)
}

type registerCascadeClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterCascadeClient(cc grpc.ClientConnInterface) RegisterCascadeClient {
	return &registerCascadeClient{cc}
}

func (c *registerCascadeClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCascade_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterCascade_ServiceDesc.Streams[0], RegisterCascade_Session_FullMethodName, opts...)
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

func (c *registerCascadeClient) SendCascadeTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error) {
	out := new(SendTicketSignatureReply)
	err := c.cc.Invoke(ctx, RegisterCascade_SendCascadeTicketSignature_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterCascadeServer is the server API for RegisterCascade service.
// All implementations must embed UnimplementedRegisterCascadeServer
// for forward compatibility
type RegisterCascadeServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterCascade_SessionServer) error
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendCascadeTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error)
	mustEmbedUnimplementedRegisterCascadeServer()
}

// UnimplementedRegisterCascadeServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterCascadeServer struct {
}

func (UnimplementedRegisterCascadeServer) Session(RegisterCascade_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterCascadeServer) SendCascadeTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCascadeTicketSignature not implemented")
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

func _RegisterCascade_SendCascadeTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCascadeServer).SendCascadeTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegisterCascade_SendCascadeTicketSignature_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCascadeServer).SendCascadeTicketSignature(ctx, req.(*SendTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterCascade_ServiceDesc is the grpc.ServiceDesc for RegisterCascade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterCascade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.RegisterCascade",
	HandlerType: (*RegisterCascadeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCascadeTicketSignature",
			Handler:    _RegisterCascade_SendCascadeTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterCascade_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "register_cascade_sn.proto",
}
