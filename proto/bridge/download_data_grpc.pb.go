// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package bridge

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

// DownloadDataClient is the client API for DownloadData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DownloadDataClient interface {
	DownloadThumbnail(ctx context.Context, in *DownloadThumbnailRequest, opts ...grpc.CallOption) (*DownloadThumbnailReply, error)
	DownloadDDAndFingerprints(ctx context.Context, in *DownloadDDAndFingerprintsRequest, opts ...grpc.CallOption) (*DownloadDDAndFingerprintsReply, error)
}

type downloadDataClient struct {
	cc grpc.ClientConnInterface
}

func NewDownloadDataClient(cc grpc.ClientConnInterface) DownloadDataClient {
	return &downloadDataClient{cc}
}

func (c *downloadDataClient) DownloadThumbnail(ctx context.Context, in *DownloadThumbnailRequest, opts ...grpc.CallOption) (*DownloadThumbnailReply, error) {
	out := new(DownloadThumbnailReply)
	err := c.cc.Invoke(ctx, "/bridge.DownloadData/DownloadThumbnail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *downloadDataClient) DownloadDDAndFingerprints(ctx context.Context, in *DownloadDDAndFingerprintsRequest, opts ...grpc.CallOption) (*DownloadDDAndFingerprintsReply, error) {
	out := new(DownloadDDAndFingerprintsReply)
	err := c.cc.Invoke(ctx, "/bridge.DownloadData/DownloadDDAndFingerprints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DownloadDataServer is the server API for DownloadData service.
// All implementations must embed UnimplementedDownloadDataServer
// for forward compatibility
type DownloadDataServer interface {
	DownloadThumbnail(context.Context, *DownloadThumbnailRequest) (*DownloadThumbnailReply, error)
	DownloadDDAndFingerprints(context.Context, *DownloadDDAndFingerprintsRequest) (*DownloadDDAndFingerprintsReply, error)
	mustEmbedUnimplementedDownloadDataServer()
}

// UnimplementedDownloadDataServer must be embedded to have forward compatible implementations.
type UnimplementedDownloadDataServer struct {
}

func (UnimplementedDownloadDataServer) DownloadThumbnail(context.Context, *DownloadThumbnailRequest) (*DownloadThumbnailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadThumbnail not implemented")
}
func (UnimplementedDownloadDataServer) DownloadDDAndFingerprints(context.Context, *DownloadDDAndFingerprintsRequest) (*DownloadDDAndFingerprintsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadDDAndFingerprints not implemented")
}
func (UnimplementedDownloadDataServer) mustEmbedUnimplementedDownloadDataServer() {}

// UnsafeDownloadDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DownloadDataServer will
// result in compilation errors.
type UnsafeDownloadDataServer interface {
	mustEmbedUnimplementedDownloadDataServer()
}

func RegisterDownloadDataServer(s grpc.ServiceRegistrar, srv DownloadDataServer) {
	s.RegisterService(&DownloadData_ServiceDesc, srv)
}

func _DownloadData_DownloadThumbnail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadThumbnailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadDataServer).DownloadThumbnail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bridge.DownloadData/DownloadThumbnail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadDataServer).DownloadThumbnail(ctx, req.(*DownloadThumbnailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DownloadData_DownloadDDAndFingerprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadDDAndFingerprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadDataServer).DownloadDDAndFingerprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bridge.DownloadData/DownloadDDAndFingerprints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadDataServer).DownloadDDAndFingerprints(ctx, req.(*DownloadDDAndFingerprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DownloadData_ServiceDesc is the grpc.ServiceDesc for DownloadData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DownloadData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bridge.DownloadData",
	HandlerType: (*DownloadDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DownloadThumbnail",
			Handler:    _DownloadData_DownloadThumbnail_Handler,
		},
		{
			MethodName: "DownloadDDAndFingerprints",
			Handler:    _DownloadData_DownloadDDAndFingerprints_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "download_data.proto",
}
