// Copyright (c) 2021-2021 The Pastel Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: dd-server.proto

package dupedetection

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
	DupeDetectionServer_ImageRarenessScore_FullMethodName = "/dupedetection.DupeDetectionServer/ImageRarenessScore"
	DupeDetectionServer_GetStatus_FullMethodName          = "/dupedetection.DupeDetectionServer/GetStatus"
)

// DupeDetectionServerClient is the client API for DupeDetectionServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DupeDetectionServerClient interface {
	ImageRarenessScore(ctx context.Context, in *RarenessScoreRequest, opts ...grpc.CallOption) (*ImageRarenessScoreReply, error)
	GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error)
}

type dupeDetectionServerClient struct {
	cc grpc.ClientConnInterface
}

func NewDupeDetectionServerClient(cc grpc.ClientConnInterface) DupeDetectionServerClient {
	return &dupeDetectionServerClient{cc}
}

func (c *dupeDetectionServerClient) ImageRarenessScore(ctx context.Context, in *RarenessScoreRequest, opts ...grpc.CallOption) (*ImageRarenessScoreReply, error) {
	out := new(ImageRarenessScoreReply)
	err := c.cc.Invoke(ctx, DupeDetectionServer_ImageRarenessScore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dupeDetectionServerClient) GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error) {
	out := new(GetStatusResponse)
	err := c.cc.Invoke(ctx, DupeDetectionServer_GetStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DupeDetectionServerServer is the server API for DupeDetectionServer service.
// All implementations must embed UnimplementedDupeDetectionServerServer
// for forward compatibility
type DupeDetectionServerServer interface {
	ImageRarenessScore(context.Context, *RarenessScoreRequest) (*ImageRarenessScoreReply, error)
	GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error)
	mustEmbedUnimplementedDupeDetectionServerServer()
}

// UnimplementedDupeDetectionServerServer must be embedded to have forward compatible implementations.
type UnimplementedDupeDetectionServerServer struct {
}

func (UnimplementedDupeDetectionServerServer) ImageRarenessScore(context.Context, *RarenessScoreRequest) (*ImageRarenessScoreReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImageRarenessScore not implemented")
}
func (UnimplementedDupeDetectionServerServer) GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedDupeDetectionServerServer) mustEmbedUnimplementedDupeDetectionServerServer() {}

// UnsafeDupeDetectionServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DupeDetectionServerServer will
// result in compilation errors.
type UnsafeDupeDetectionServerServer interface {
	mustEmbedUnimplementedDupeDetectionServerServer()
}

func RegisterDupeDetectionServerServer(s grpc.ServiceRegistrar, srv DupeDetectionServerServer) {
	s.RegisterService(&DupeDetectionServer_ServiceDesc, srv)
}

func _DupeDetectionServer_ImageRarenessScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RarenessScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DupeDetectionServerServer).ImageRarenessScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DupeDetectionServer_ImageRarenessScore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DupeDetectionServerServer).ImageRarenessScore(ctx, req.(*RarenessScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DupeDetectionServer_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DupeDetectionServerServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DupeDetectionServer_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DupeDetectionServerServer).GetStatus(ctx, req.(*GetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DupeDetectionServer_ServiceDesc is the grpc.ServiceDesc for DupeDetectionServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DupeDetectionServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dupedetection.DupeDetectionServer",
	HandlerType: (*DupeDetectionServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ImageRarenessScore",
			Handler:    _DupeDetectionServer_ImageRarenessScore_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _DupeDetectionServer_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dd-server.proto",
}
