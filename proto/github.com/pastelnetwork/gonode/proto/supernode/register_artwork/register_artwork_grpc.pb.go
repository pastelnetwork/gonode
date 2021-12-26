// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package register_artwork

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

// RegisterArtworkClient is the client API for RegisterArtwork service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterArtworkClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (RegisterArtwork_SessionClient, error)
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error)
	// SendArtTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendArtTicketSignature(ctx context.Context, in *SendArtTicketSignatureRequest, opts ...grpc.CallOption) (*SendArtTicketSignatureReply, error)
}

type registerArtworkClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterArtworkClient(cc grpc.ClientConnInterface) RegisterArtworkClient {
	return &registerArtworkClient{cc}
}

func (c *registerArtworkClient) Session(ctx context.Context, opts ...grpc.CallOption) (RegisterArtwork_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &RegisterArtwork_ServiceDesc.Streams[0], "/register_artwork.RegisterArtwork/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerArtworkSessionClient{stream}
	return x, nil
}

type RegisterArtwork_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type registerArtworkSessionClient struct {
	grpc.ClientStream
}

func (x *registerArtworkSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *registerArtworkSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerArtworkClient) SendSignedDDAndFingerprints(ctx context.Context, in *SendSignedDDAndFingerprintsRequest, opts ...grpc.CallOption) (*SendSignedDDAndFingerprintsReply, error) {
	out := new(SendSignedDDAndFingerprintsReply)
	err := c.cc.Invoke(ctx, "/register_artwork.RegisterArtwork/SendSignedDDAndFingerprints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerArtworkClient) SendArtTicketSignature(ctx context.Context, in *SendArtTicketSignatureRequest, opts ...grpc.CallOption) (*SendArtTicketSignatureReply, error) {
	out := new(SendArtTicketSignatureReply)
	err := c.cc.Invoke(ctx, "/register_artwork.RegisterArtwork/SendArtTicketSignature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterArtworkServer is the server API for RegisterArtwork service.
// All implementations must embed UnimplementedRegisterArtworkServer
// for forward compatibility
type RegisterArtworkServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants to connect to.
	// The stream is used by the parties to inform each other about the cancellation of the task.
	Session(RegisterArtwork_SessionServer) error
	// SendSignedDDAndFingerprints send its SignedDDAndFingerprints to other node
	SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error)
	// SendArtTicketSignature send signature from supernodes mn2/mn3 for given reg NFT session id to primary supernode
	SendArtTicketSignature(context.Context, *SendArtTicketSignatureRequest) (*SendArtTicketSignatureReply, error)
	mustEmbedUnimplementedRegisterArtworkServer()
}

// UnimplementedRegisterArtworkServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterArtworkServer struct {
}

func (UnimplementedRegisterArtworkServer) Session(RegisterArtwork_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedRegisterArtworkServer) SendSignedDDAndFingerprints(context.Context, *SendSignedDDAndFingerprintsRequest) (*SendSignedDDAndFingerprintsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSignedDDAndFingerprints not implemented")
}
func (UnimplementedRegisterArtworkServer) SendArtTicketSignature(context.Context, *SendArtTicketSignatureRequest) (*SendArtTicketSignatureReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendArtTicketSignature not implemented")
}
func (UnimplementedRegisterArtworkServer) mustEmbedUnimplementedRegisterArtworkServer() {}

// UnsafeRegisterArtworkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterArtworkServer will
// result in compilation errors.
type UnsafeRegisterArtworkServer interface {
	mustEmbedUnimplementedRegisterArtworkServer()
}

func RegisterRegisterArtworkServer(s grpc.ServiceRegistrar, srv RegisterArtworkServer) {
	s.RegisterService(&RegisterArtwork_ServiceDesc, srv)
}

func _RegisterArtwork_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RegisterArtworkServer).Session(&registerArtworkSessionServer{stream})
}

type RegisterArtwork_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type registerArtworkSessionServer struct {
	grpc.ServerStream
}

func (x *registerArtworkSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *registerArtworkSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RegisterArtwork_SendSignedDDAndFingerprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSignedDDAndFingerprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterArtworkServer).SendSignedDDAndFingerprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_artwork.RegisterArtwork/SendSignedDDAndFingerprints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterArtworkServer).SendSignedDDAndFingerprints(ctx, req.(*SendSignedDDAndFingerprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegisterArtwork_SendArtTicketSignature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendArtTicketSignatureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterArtworkServer).SendArtTicketSignature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/register_artwork.RegisterArtwork/SendArtTicketSignature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterArtworkServer).SendArtTicketSignature(ctx, req.(*SendArtTicketSignatureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterArtwork_ServiceDesc is the grpc.ServiceDesc for RegisterArtwork service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterArtwork_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "register_artwork.RegisterArtwork",
	HandlerType: (*RegisterArtworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSignedDDAndFingerprints",
			Handler:    _RegisterArtwork_SendSignedDDAndFingerprints_Handler,
		},
		{
			MethodName: "SendArtTicketSignature",
			Handler:    _RegisterArtwork_SendArtTicketSignature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _RegisterArtwork_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "supernode/register_artwork.proto",
}
