// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package process_userdata

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
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the
	// task on the supernode. By sending `sessID` in the Metadata. The stream is
	// used by the parties to inform each other about the cancellation of the
	// task.
	Session(ctx context.Context, opts ...grpc.CallOption) (ProcessUserdata_SessionClient, error)
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error)
	// SendUserdata send the user info and return the operation is success or
	// detail on error.
	SendUserdata(ctx context.Context, in *UserdataRequest, opts ...grpc.CallOption) (*UserdataReply, error)
	// ReceiveUserdata receive the user info and return the operation is success
	// or detail on error.
	ReceiveUserdata(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*UserdataRequest, error)
}

type processUserdataClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessUserdataClient(cc grpc.ClientConnInterface) ProcessUserdataClient {
	return &processUserdataClient{cc}
}

func (c *processUserdataClient) Session(ctx context.Context, opts ...grpc.CallOption) (ProcessUserdata_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProcessUserdata_ServiceDesc.Streams[0], "/process_userdata.ProcessUserdata/Session", opts...)
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

func (c *processUserdataClient) AcceptedNodes(ctx context.Context, in *AcceptedNodesRequest, opts ...grpc.CallOption) (*AcceptedNodesReply, error) {
	out := new(AcceptedNodesReply)
	err := c.cc.Invoke(ctx, "/process_userdata.ProcessUserdata/AcceptedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processUserdataClient) ConnectTo(ctx context.Context, in *ConnectToRequest, opts ...grpc.CallOption) (*ConnectToReply, error) {
	out := new(ConnectToReply)
	err := c.cc.Invoke(ctx, "/process_userdata.ProcessUserdata/ConnectTo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processUserdataClient) SendUserdata(ctx context.Context, in *UserdataRequest, opts ...grpc.CallOption) (*UserdataReply, error) {
	out := new(UserdataReply)
	err := c.cc.Invoke(ctx, "/process_userdata.ProcessUserdata/SendUserdata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processUserdataClient) ReceiveUserdata(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*UserdataRequest, error) {
	out := new(UserdataRequest)
	err := c.cc.Invoke(ctx, "/process_userdata.ProcessUserdata/ReceiveUserdata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessUserdataServer is the server API for ProcessUserdata service.
// All implementations must embed UnimplementedProcessUserdataServer
// for forward compatibility
type ProcessUserdataServer interface {
	// Session informs the supernode its position (primary/secondary).
	// Returns `SessID` that are used by all other rpc methods to identify the
	// task on the supernode. By sending `sessID` in the Metadata. The stream is
	// used by the parties to inform each other about the cancellation of the
	// task.
	Session(ProcessUserdata_SessionServer) error
	// AcceptedNodes returns peers of the secondary supernodes connected to it.
	AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error)
	// ConnectTo requests to connect to the primary supernode.
	ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error)
	// SendUserdata send the user info and return the operation is success or
	// detail on error.
	SendUserdata(context.Context, *UserdataRequest) (*UserdataReply, error)
	// ReceiveUserdata receive the user info and return the operation is success
	// or detail on error.
	ReceiveUserdata(context.Context, *RetrieveRequest) (*UserdataRequest, error)
	mustEmbedUnimplementedProcessUserdataServer()
}

// UnimplementedProcessUserdataServer must be embedded to have forward compatible implementations.
type UnimplementedProcessUserdataServer struct {
}

func (UnimplementedProcessUserdataServer) Session(ProcessUserdata_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedProcessUserdataServer) AcceptedNodes(context.Context, *AcceptedNodesRequest) (*AcceptedNodesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptedNodes not implemented")
}
func (UnimplementedProcessUserdataServer) ConnectTo(context.Context, *ConnectToRequest) (*ConnectToReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectTo not implemented")
}
func (UnimplementedProcessUserdataServer) SendUserdata(context.Context, *UserdataRequest) (*UserdataReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUserdata not implemented")
}
func (UnimplementedProcessUserdataServer) ReceiveUserdata(context.Context, *RetrieveRequest) (*UserdataRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveUserdata not implemented")
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

func _ProcessUserdata_AcceptedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptedNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).AcceptedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/process_userdata.ProcessUserdata/AcceptedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).AcceptedNodes(ctx, req.(*AcceptedNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProcessUserdata_ConnectTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectToRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).ConnectTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/process_userdata.ProcessUserdata/ConnectTo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).ConnectTo(ctx, req.(*ConnectToRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProcessUserdata_SendUserdata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserdataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).SendUserdata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/process_userdata.ProcessUserdata/SendUserdata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).SendUserdata(ctx, req.(*UserdataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProcessUserdata_ReceiveUserdata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessUserdataServer).ReceiveUserdata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/process_userdata.ProcessUserdata/ReceiveUserdata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessUserdataServer).ReceiveUserdata(ctx, req.(*RetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProcessUserdata_ServiceDesc is the grpc.ServiceDesc for ProcessUserdata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProcessUserdata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "process_userdata.ProcessUserdata",
	HandlerType: (*ProcessUserdataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptedNodes",
			Handler:    _ProcessUserdata_AcceptedNodes_Handler,
		},
		{
			MethodName: "ConnectTo",
			Handler:    _ProcessUserdata_ConnectTo_Handler,
		},
		{
			MethodName: "SendUserdata",
			Handler:    _ProcessUserdata_SendUserdata_Handler,
		},
		{
			MethodName: "ReceiveUserdata",
			Handler:    _ProcessUserdata_ReceiveUserdata_Handler,
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
	Metadata: "process_userdata_walletnode.proto",
}
