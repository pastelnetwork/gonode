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

// DownloadArtworkClient is the client API for DownloadArtwork service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DownloadArtworkClient interface {
	// Download downloads artwork by given txid, timestamp signature.
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (DownloadArtwork_DownloadClient, error)
}

type downloadArtworkClient struct {
	cc grpc.ClientConnInterface
}

func NewDownloadArtworkClient(cc grpc.ClientConnInterface) DownloadArtworkClient {
	return &downloadArtworkClient{cc}
}

func (c *downloadArtworkClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (DownloadArtwork_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &DownloadArtwork_ServiceDesc.Streams[0], "/walletnode.DownloadArtwork/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &downloadArtworkDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DownloadArtwork_DownloadClient interface {
	Recv() (*DownloadReply, error)
	grpc.ClientStream
}

type downloadArtworkDownloadClient struct {
	grpc.ClientStream
}

func (x *downloadArtworkDownloadClient) Recv() (*DownloadReply, error) {
	m := new(DownloadReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DownloadArtworkServer is the server API for DownloadArtwork service.
// All implementations must embed UnimplementedDownloadArtworkServer
// for forward compatibility
type DownloadArtworkServer interface {
	// Download downloads artwork by given txid, timestamp signature.
	Download(*DownloadRequest, DownloadArtwork_DownloadServer) error
	mustEmbedUnimplementedDownloadArtworkServer()
}

// UnimplementedDownloadArtworkServer must be embedded to have forward compatible implementations.
type UnimplementedDownloadArtworkServer struct {
}

func (UnimplementedDownloadArtworkServer) Download(*DownloadRequest, DownloadArtwork_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedDownloadArtworkServer) mustEmbedUnimplementedDownloadArtworkServer() {}

// UnsafeDownloadArtworkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DownloadArtworkServer will
// result in compilation errors.
type UnsafeDownloadArtworkServer interface {
	mustEmbedUnimplementedDownloadArtworkServer()
}

func RegisterDownloadArtworkServer(s grpc.ServiceRegistrar, srv DownloadArtworkServer) {
	s.RegisterService(&DownloadArtwork_ServiceDesc, srv)
}

func _DownloadArtwork_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DownloadArtworkServer).Download(m, &downloadArtworkDownloadServer{stream})
}

type DownloadArtwork_DownloadServer interface {
	Send(*DownloadReply) error
	grpc.ServerStream
}

type downloadArtworkDownloadServer struct {
	grpc.ServerStream
}

func (x *downloadArtworkDownloadServer) Send(m *DownloadReply) error {
	return x.ServerStream.SendMsg(m)
}

// DownloadArtwork_ServiceDesc is the grpc.ServiceDesc for DownloadArtwork service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DownloadArtwork_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "walletnode.DownloadArtwork",
	HandlerType: (*DownloadArtworkServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _DownloadArtwork_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "download_artwork.proto",
}
