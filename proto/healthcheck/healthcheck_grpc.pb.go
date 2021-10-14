// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package healthcheck

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

// HealthCheckClient is the client API for HealthCheck service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthCheckClient interface {
	// rpc get overview status of service
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error)
	// Store a data to p2p
	P2PStore(ctx context.Context, in *P2PStoreRequest, opts ...grpc.CallOption) (*P2PStoreReply, error)
	// Retrieve a key from p2p
	P2PRetrieve(ctx context.Context, in *P2PRetrieveRequest, opts ...grpc.CallOption) (*P2PRetrieveReply, error)
	// Do a query on rqlite
	QueryRqlite(ctx context.Context, in *QueryRqliteRequest, opts ...grpc.CallOption) (*QueryRqliteReply, error)
}

type healthCheckClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthCheckClient(cc grpc.ClientConnInterface) HealthCheckClient {
	return &healthCheckClient{cc}
}

func (c *healthCheckClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthCheckClient) P2PStore(ctx context.Context, in *P2PStoreRequest, opts ...grpc.CallOption) (*P2PStoreReply, error) {
	out := new(P2PStoreReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/P2PStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthCheckClient) P2PRetrieve(ctx context.Context, in *P2PRetrieveRequest, opts ...grpc.CallOption) (*P2PRetrieveReply, error) {
	out := new(P2PRetrieveReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/P2PRetrieve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthCheckClient) QueryRqlite(ctx context.Context, in *QueryRqliteRequest, opts ...grpc.CallOption) (*QueryRqliteReply, error) {
	out := new(QueryRqliteReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/QueryRqlite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthCheckServer is the server API for HealthCheck service.
// All implementations must embed UnimplementedHealthCheckServer
// for forward compatibility
type HealthCheckServer interface {
	// rpc get overview status of service
	Status(context.Context, *StatusRequest) (*StatusReply, error)
	// Store a data to p2p
	P2PStore(context.Context, *P2PStoreRequest) (*P2PStoreReply, error)
	// Retrieve a key from p2p
	P2PRetrieve(context.Context, *P2PRetrieveRequest) (*P2PRetrieveReply, error)
	// Do a query on rqlite
	QueryRqlite(context.Context, *QueryRqliteRequest) (*QueryRqliteReply, error)
	mustEmbedUnimplementedHealthCheckServer()
}

// UnimplementedHealthCheckServer must be embedded to have forward compatible implementations.
type UnimplementedHealthCheckServer struct {
}

func (UnimplementedHealthCheckServer) Status(context.Context, *StatusRequest) (*StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedHealthCheckServer) P2PStore(context.Context, *P2PStoreRequest) (*P2PStoreReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method P2PStore not implemented")
}
func (UnimplementedHealthCheckServer) P2PRetrieve(context.Context, *P2PRetrieveRequest) (*P2PRetrieveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method P2PRetrieve not implemented")
}
func (UnimplementedHealthCheckServer) QueryRqlite(context.Context, *QueryRqliteRequest) (*QueryRqliteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRqlite not implemented")
}
func (UnimplementedHealthCheckServer) mustEmbedUnimplementedHealthCheckServer() {}

// UnsafeHealthCheckServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthCheckServer will
// result in compilation errors.
type UnsafeHealthCheckServer interface {
	mustEmbedUnimplementedHealthCheckServer()
}

func RegisterHealthCheckServer(s grpc.ServiceRegistrar, srv HealthCheckServer) {
	s.RegisterService(&HealthCheck_ServiceDesc, srv)
}

func _HealthCheck_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthCheck_P2PStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(P2PStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).P2PStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/P2PStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).P2PStore(ctx, req.(*P2PStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthCheck_P2PRetrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(P2PRetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).P2PRetrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/P2PRetrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).P2PRetrieve(ctx, req.(*P2PRetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthCheck_QueryRqlite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRqliteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).QueryRqlite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/QueryRqlite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).QueryRqlite(ctx, req.(*QueryRqliteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HealthCheck_ServiceDesc is the grpc.ServiceDesc for HealthCheck service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HealthCheck_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "healthcheck.HealthCheck",
	HandlerType: (*HealthCheckServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _HealthCheck_Status_Handler,
		},
		{
			MethodName: "P2PStore",
			Handler:    _HealthCheck_P2PStore_Handler,
		},
		{
			MethodName: "P2PRetrieve",
			Handler:    _HealthCheck_P2PRetrieve_Handler,
		},
		{
			MethodName: "QueryRqlite",
			Handler:    _HealthCheck_QueryRqlite_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "healthcheck/healthcheck.proto",
}
