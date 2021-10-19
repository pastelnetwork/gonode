// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// ExternalStorageClient is the client API for ExternalStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExternalStorageClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants
	// to connect to. The stream is used by the parties to inform each other about
	// the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (ExternalStorage_SessionClient, error)
	// SendExternalStorageTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendExternalStorageTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error)
}

type externalStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewExternalStorageClient(cc grpc.ClientConnInterface) ExternalStorageClient {
	return &externalStorageClient{cc}
}

func (c *externalStorageClient) Session(ctx context.Context, opts ...grpc.CallOption) (ExternalStorage_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExternalStorage_ServiceDesc.Streams[0], "/supernode.ExternalStorage/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &externalStorageSessionClient{stream}
	return x, nil
}

type ExternalStorage_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type externalStorageSessionClient struct {
	grpc.ClientStream
}

func (x *externalStorageSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *externalStorageSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *externalStorageClient) SendExternalStorageTicketSignature(ctx context.Context, in *SendTicketSignatureRequest, opts ...grpc.CallOption) (*SendTicketSignatureReply, error) {
	out := new(SendTicketSignatureReply)
	err := c.cc.Invoke(ctx, "/supernode.ExternalStorage/SendExternalStorageTicketSignature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExternalStorageServer is the server API for ExternalStorage service.
// All implementations must embed UnimplementedExternalStorageServer
// for forward compatibility
type ExternalStorageServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants
	// to connect to. The stream is used by the parties to inform each other about
	// the cancellation of the task.
	Session(ExternalStorage_SessionServer) error
	// SendExternalStorageTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendExternalStorageTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error)
	mustEmbedUnimplementedExternalStorageServer()
}

// UnimplementedExternalStorageServer must be embedded to have forward compatible implementations.
type UnimplementedExternalStorageServer struct {
}

func (UnimplementedExternalStorageServer) Session(ExternalStorage_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedExternalStorageServer) SendExternalStorageTicketSignature(context.Context, *SendTicketSignatureRequest) (*SendTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendExternalStorageTicketSignature not implemented")
}
func (UnimplementedExternalStorageServer) mustEmbedUnimplementedExternalStorageServer() {}

// UnsafeExternalStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExternalStorageServer will
// result in compilation errors.
type UnsafeExternalStorageServer interface {
	mustEmbedUnimplementedExternalStorageServer()
}

func RegisterExternalStorageServer(s grpc.ServiceRegistrar, srv ExternalStorageServer) {
	s.RegisterService(&ExternalStorage_ServiceDesc, srv)
}

func _ExternalStorage_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExternalStorageServer).Session(&externalStorageSessionServer{stream})
}

type ExternalStorage_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type externalStorageSessionServer struct {
	grpc.ServerStream
}

func (x *externalStorageSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *externalStorageSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ExternalStorage_SendExternalStorageTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalStorageServer).SendExternalStorageTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supernode.ExternalStorage/SendExternalStorageTicketSignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalStorageServer).SendExternalStorageTicketSignature(ctx, req.(*SendTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExternalStorage_ServiceDesc is the grpc.ServiceDesc for ExternalStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExternalStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.ExternalStorage",
	HandlerType: (*ExternalStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendExternalStorageTicketSignature",
			Handler:    _ExternalStorage_SendExternalStorageTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _ExternalStorage_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "external_storage_request.proto",
}
