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
	// Set a data to p2p
	P2PSet(ctx context.Context, in *P2PSetRequest, opts ...grpc.CallOption) (*P2PSetReply, error)
	// Get a key from p2p
	P2PGet(ctx context.Context, in *P2PGetRequest, opts ...grpc.CallOption) (*P2PGetReply, error)
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

func (c *healthCheckClient) P2PSet(ctx context.Context, in *P2PSetRequest, opts ...grpc.CallOption) (*P2PSetReply, error) {
	out := new(P2PSetReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/P2PSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthCheckClient) P2PGet(ctx context.Context, in *P2PGetRequest, opts ...grpc.CallOption) (*P2PGetReply, error) {
	out := new(P2PGetReply)
	err := c.cc.Invoke(ctx, "/healthcheck.HealthCheck/P2PGet", in, out, opts...)
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
	// Set a data to p2p
	P2PSet(context.Context, *P2PSetRequest) (*P2PSetReply, error)
	// Get a key from p2p
	P2PGet(context.Context, *P2PGetRequest) (*P2PGetReply, error)
	mustEmbedUnimplementedHealthCheckServer()
}

// UnimplementedHealthCheckServer must be embedded to have forward compatible implementations.
type UnimplementedHealthCheckServer struct {
}

func (UnimplementedHealthCheckServer) Status(context.Context, *StatusRequest) (*StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedHealthCheckServer) P2PSet(context.Context, *P2PSetRequest) (*P2PSetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method P2PSet not implemented")
}
func (UnimplementedHealthCheckServer) P2PGet(context.Context, *P2PGetRequest) (*P2PGetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method P2PGet not implemented")
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

func _HealthCheck_P2PSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(P2PSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).P2PSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/P2PSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).P2PSet(ctx, req.(*P2PSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthCheck_P2PGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(P2PGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServer).P2PGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.HealthCheck/P2PGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServer).P2PGet(ctx, req.(*P2PGetRequest))
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
			MethodName: "P2PSet",
			Handler:    _HealthCheck_P2PSet_Handler,
		},
		{
			MethodName: "P2PGet",
			Handler:    _HealthCheck_P2PGet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "healthcheck/healthcheck.proto",
}
