// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: register_collection_sn.proto

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
	RegisterCollection_Session_FullMethodName                       = "/supernode.RegisterCollection/Session"
	RegisterCollection_SendCollectionTicketSignature_FullMethodName = "/supernode.RegisterCollection/SendCollectionTicketSignature"
)

// RegisterCollectionClient is the client API for RegisterCollection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterCollectionClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCollection_SessionClient, error)
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendCollectionTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error)
}

type registerCollectionClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterCollectionClient(cc grpc.ClientConnInterface) RegisterCollectionClient {
	return &registerCollectionClient{cc}
}

func (c *registerCollectionClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterCollection_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterCollection_ServiceDesc.Streams[0], RegisterCollection_Session_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &registerCollectionSessionClient{stream}
	return x, nil
}

type RegisterCollection_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type registerCollectionSessionClient struct {
	grpc.ClientStream
}

func (x *registerCollectionSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerCollectionSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerCollectionClient) SendCollectionTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error) {
	out := new(SendTicketSignatureReply)
	err := c.cc.Invoke(ctx, RegisterCollection_SendCollectionTicketSignature_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterCollectionServer is the server API for RegisterCollection service.
// All implementations must embed UnimplementedRegisterCollectionServer
// for forward compatibility
type RegisterCollectionServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterCollection_SessionServer) error
	// SendSenseTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendCollectionTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error)
	mustEmbedUnimplementedRegisterCollectionServer()
}

// UnimplementedRegisterCollectionServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterCollectionServer struct {
}

func (UnimplementedRegisterCollectionServer) Session(RegisterCollection_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterCollectionServer) SendCollectionTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCollectionTicketSignature not implemented")
}
func (UnimplementedRegisterCollectionServer) mustEmbedUnimplementedRegisterCollectionServer() {}

// UnsafeRegisterCollectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterCollectionServer will
// result in compilation errors.
type UnsafeRegisterCollectionServer interface {
	mustEmbedUnimplementedRegisterCollectionServer()
}

func RegisterRegisterCollectionServer(s grpc.ServiceRegistrar, srv RegisterCollectionServer) {
	s.RegisterService(&RegisterCollection_ServiceDesc, srv)
}

func _RegisterCollection_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterCollectionServer).Session(&registerCollectionSessionServer{stream})
}

type RegisterCollection_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type registerCollectionSessionServer struct {
	grpc.ServerStream
}

func (x *registerCollectionSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerCollectionSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterCollection_SendCollectionTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterCollectionServer).SendCollectionTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegisterCollection_SendCollectionTicketSignature_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterCollectionServer).SendCollectionTicketSignature(ctx, req.(*SendTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterCollection_ServiceDesc is the grpc.ServiceDesc for RegisterCollection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterCollection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.RegisterCollection",
	HandlerType: (*RegisterCollectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCollectionTicketSignature",
			Handler:    _RegisterCollection_SendCollectionTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterCollection_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "register_collection_sn.proto",
}
