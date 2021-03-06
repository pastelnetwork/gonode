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

// ProcessUserdataClient is the client API for ProcessUserdata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProcessUserdataClient interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants
	// to connect to. The stream is used by the parties to inform each other about
	// the cancellation of the task.
	Session(ctx context.Context, opts ...grpc.CallOption) (ProcessUserdata_SessionClient, error)
	// SendUserdataToPrimary send signed userdata from other supernodes to primary
	// supernode
	SendUserdataToPrimary(ctx context.Context, in *SuperNodeRequest, opts ...grpc.CallOption) (*SuperNodeReply, error)
	// SendUserdataToLeader send final userdata to supernode contain rqlite leader
	SendUserdataToLeader(ctx context.Context, in *UserdataRequest, opts ...grpc.CallOption) (*SuperNodeReply, error)
}

type processUserdataClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessUserdataClient(cc grpc.ClientConnInterface) ProcessUserdataClient {
	return &processUserdataClient{cc}
}

func (c *processUserdataClient) Session(ctx context.Context, opts ...grpc.CallOption) (ProcessUserdata_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProcessUserdata_ServiceDesc.Streams[0], "/supernode.ProcessUserdata/Session", opts...)
	if err != nil {
		return nil, err
	}
	x := &processUserdataSessionClient{stream}
	return x, nil
}

type ProcessUserdata_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type processUserdataSessionClient struct {
	grpc.ClientStream
}

func (x *processUserdataSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *processUserdataSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *processUserdataClient) SendUserdataToPrimary(ctx context.Context, in *SuperNodeRequest, opts ...grpc.CallOption) (*SuperNodeReply, error) {
	out := new(SuperNodeReply)
	err := c.cc.Invoke(ctx, "/supernode.ProcessUserdata/SendUserdataToPrimary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processUserdataClient) SendUserdataToLeader(ctx context.Context, in *UserdataRequest, opts ...grpc.CallOption) (*SuperNodeReply, error) {
	out := new(SuperNodeReply)
	err := c.cc.Invoke(ctx, "/supernode.ProcessUserdata/SendUserdataToLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessUserdataServer is the server API for ProcessUserdata service.
// All implementations must embed UnimplementedProcessUserdataServer
// for forward compatibility
type ProcessUserdataServer interface {
	// Session informs primary supernode about its `nodeID` and `sessID` it wants
	// to connect to. The stream is used by the parties to inform each other about
	// the cancellation of the task.
	Session(ProcessUserdata_SessionServer) error
	// SendUserdataToPrimary send signed userdata from other supernodes to primary
	// supernode
	SendUserdataToPrimary(context.Context, *SuperNodeRequest) (*SuperNodeReply, error)
	// SendUserdataToLeader send final userdata to supernode contain rqlite leader
	SendUserdataToLeader(context.Context, *UserdataRequest) (*SuperNodeReply, error)
	mustEmbedUnimplementedProcessUserdataServer()
}

// UnimplementedProcessUserdataServer must be embedded to have forward compatible implementations.
type UnimplementedProcessUserdataServer struct {
}

func (UnimplementedProcessUserdataServer) Session(ProcessUserdata_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedProcessUserdataServer) SendUserdataToPrimary(context.Context, *SuperNodeRequest) (*SuperNodeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUserdataToPrimary not implemented")
}
func (UnimplementedProcessUserdataServer) SendUserdataToLeader(context.Context, *UserdataRequest) (*SuperNodeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUserdataToLeader not implemented")
}
func (UnimplementedProcessUserdataServer) mustEmbedUnimplementedProcessUserdataServer() {}

// UnsafeProcessUserdataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProcessUserdataServer will
// result in compilation errors.
type UnsafeProcessUserdataServer interface {
	mustEmbedUnimplementedProcessUserdataServer()
}

func RegisterProcessUserdataServer(s grpc.ServiceRegistrar, srv ProcessUserdataServer) {
	s.RegisterService(&ProcessUserdata_ServiceDesc, srv)
}

func _ProcessUserdata_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProcessUserdataServer).Session(&processUserdataSessionServer{stream})
}

type ProcessUserdata_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type processUserdataSessionServer struct {
	grpc.ServerStream
}

func (x *processUserdataSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *processUserdataSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProcessUserdata_SendUserdataToPrimary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SuperNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).SendUserdataToPrimary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supernode.ProcessUserdata/SendUserdataToPrimary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).SendUserdataToPrimary(ctx, req.(*SuperNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProcessUserdata_SendUserdataToLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserdataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).SendUserdataToLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supernode.ProcessUserdata/SendUserdataToLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).SendUserdataToLeader(ctx, req.(*UserdataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProcessUserdata_ServiceDesc is the grpc.ServiceDesc for ProcessUserdata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProcessUserdata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.ProcessUserdata",
	HandlerType: (*ProcessUserdataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendUserdataToPrimary",
			Handler:    _ProcessUserdata_SendUserdataToPrimary_Handler,
		},
		{
			MethodName: "SendUserdataToLeader",
			Handler:    _ProcessUserdata_SendUserdataToLeader_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _ProcessUserdata_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "process_userdata_sn.proto",
}
