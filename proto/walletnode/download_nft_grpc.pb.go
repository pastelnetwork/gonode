// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package walletnode

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

// DownloadNftClient is the client API for DownloadNft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DownloadNftClient interface {
	// Download downloads NFT by given txid, timestamp and signature.
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (DownloadNft_DownloadClient, error)
	DownloadThumbnail(ctx context.Context, in *DownloadThumbnailRequest, opts ...grpc.CallOption) (*DownloadThumbnailReply, error)
	DownloadDDAndFingerprints(ctx context.Context, in *DownloadDDAndFingerprintsRequest, opts ...grpc.CallOption) (*DownloadDDAndFingerprintsReply, error)
}

type downloadNftClient struct {
	cc grpc.ClientConnInterface
}

func NewDownloadNftClient(cc grpc.ClientConnInterface) DownloadNftClient {
	return &downloadNftClient{cc}
}

func (c *downloadNftClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (DownloadNft_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &DownloadNft_ServiceDesc.Streams[0], "/walletnode.DownloadNft/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &downloadNftDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DownloadNft_DownloadClient interface {
	Recv() (*DownloadReply, error)
	grpc.ClientStream
}

type downloadNftDownloadClient struct {
	grpc.ClientStream
}

func (x *downloadNftDownloadClient) Recv() (*DownloadReply, error) {
	m := new(DownloadReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *downloadNftClient) DownloadThumbnail(ctx context.Context, in *DownloadThumbnailRequest, opts ...grpc.CallOption) (*DownloadThumbnailReply, error) {
	out := new(DownloadThumbnailReply)
	err := c.cc.Invoke(ctx, "/walletnode.DownloadNft/DownloadThumbnail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *downloadNftClient) DownloadDDAndFingerprints(ctx context.Context, in *DownloadDDAndFingerprintsRequest, opts ...grpc.CallOption) (*DownloadDDAndFingerprintsReply, error) {
	out := new(DownloadDDAndFingerprintsReply)
	err := c.cc.Invoke(ctx, "/walletnode.DownloadNft/DownloadDDAndFingerprints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DownloadNftServer is the server API for DownloadNft service.
// All implementations must embed UnimplementedDownloadNftServer
// for forward compatibility
type DownloadNftServer interface {
	// Download downloads NFT by given txid, timestamp and signature.
	Download(*DownloadRequest, DownloadNft_DownloadServer) error
	DownloadThumbnail(context.Context, *DownloadThumbnailRequest) (*DownloadThumbnailReply, error)
	DownloadDDAndFingerprints(context.Context, *DownloadDDAndFingerprintsRequest) (*DownloadDDAndFingerprintsReply, error)
	mustEmbedUnimplementedDownloadNftServer()
}

// UnimplementedDownloadNftServer must be embedded to have forward compatible implementations.
type UnimplementedDownloadNftServer struct {
}

func (UnimplementedDownloadNftServer) Download(*DownloadRequest, DownloadNft_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedDownloadNftServer) DownloadThumbnail(context.Context, *DownloadThumbnailRequest) (*DownloadThumbnailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadThumbnail not implemented")
}
func (UnimplementedDownloadNftServer) DownloadDDAndFingerprints(context.Context, *DownloadDDAndFingerprintsRequest) (*DownloadDDAndFingerprintsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadDDAndFingerprints not implemented")
}
func (UnimplementedDownloadNftServer) mustEmbedUnimplementedDownloadNftServer() {}

// UnsafeDownloadNftServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DownloadNftServer will
// result in compilation errors.
type UnsafeDownloadNftServer interface {
	mustEmbedUnimplementedDownloadNftServer()
}

func RegisterDownloadNftServer(s grpc.ServiceRegistrar, srv DownloadNftServer) {
	s.RegisterService(&DownloadNft_ServiceDesc, srv)
}

func _DownloadNft_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DownloadNftServer).Download(m, &downloadNftDownloadServer{stream})
}

type DownloadNft_DownloadServer interface {
	Send(*DownloadReply) error
	grpc.ServerStream
}

type downloadNftDownloadServer struct {
	grpc.ServerStream
}

func (x *downloadNftDownloadServer) Send(m *DownloadReply) error {
	return x.ServerStream.SendMsg(m)
}

func _DownloadNft_DownloadThumbnail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadThumbnailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadNftServer).DownloadThumbnail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.DownloadNft/DownloadThumbnail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadNftServer).DownloadThumbnail(ctx, req.(*DownloadThumbnailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DownloadNft_DownloadDDAndFingerprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadDDAndFingerprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadNftServer).DownloadDDAndFingerprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/walletnode.DownloadNft/DownloadDDAndFingerprints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadNftServer).DownloadDDAndFingerprints(ctx, req.(*DownloadDDAndFingerprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DownloadNft_ServiceDesc is the grpc.ServiceDesc for DownloadNft service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DownloadNft_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "walletnode.DownloadNft",
	HandlerType: (*DownloadNftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DownloadThumbnail",
			Handler:    _DownloadNft_DownloadThumbnail_Handler,
		},
		{
			MethodName: "DownloadDDAndFingerprints",
			Handler:    _DownloadNft_DownloadDDAndFingerprints_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _DownloadNft_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "download_nft.proto",
}
