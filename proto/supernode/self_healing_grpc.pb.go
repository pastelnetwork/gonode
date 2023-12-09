// Copyright (c) 2021-2021 The Pastel Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: self_healing.proto

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

const (
	SelfHealing_Session_FullMethodName                     = "/supernode.SelfHealing/Session"
	SelfHealing_Ping_FullMethodName                        = "/supernode.SelfHealing/Ping"
	SelfHealing_ProcessSelfHealingChallenge_FullMethodName = "/supernode.SelfHealing/ProcessSelfHealingChallenge"
	SelfHealing_VerifySelfHealingChallenge_FullMethodName  = "/supernode.SelfHealing/VerifySelfHealingChallenge"
)

// SelfHealingClient is the client API for SelfHealing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SelfHealingClient interface {
	Session(ctx context.Context, opts ...grpc.CallOption) (SelfHealing_SessionClient, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	ProcessSelfHealingChallenge(ctx context.Context, in *ProcessSelfHealingChallengeRequest, opts ...grpc.CallOption) (*ProcessSelfHealingChallengeReply, error)
	VerifySelfHealingChallenge(ctx context.Context, in *VerifySelfHealingChallengeRequest, opts ...grpc.CallOption) (*VerifySelfHealingChallengeReply, error)
}

type selfHealingClient struct {
	cc grpc.ClientConnInterface
}

func NewSelfHealingClient(cc grpc.ClientConnInterface) SelfHealingClient {
	return &selfHealingClient{cc}
}

func (c *selfHealingClient) Session(ctx context.Context, opts ...grpc.CallOption) (SelfHealing_SessionClient, error) {
	stream, err := c.cc.NewStream(ctx, &SelfHealing_ServiceDesc.Streams[0], SelfHealing_Session_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &selfHealingSessionClient{stream}
	return x, nil
}

type SelfHealing_SessionClient interface {
	Send(*SessionRequest) error
	Recv() (*SessionReply, error)
	grpc.ClientStream
}

type selfHealingSessionClient struct {
	grpc.ClientStream
}

func (x *selfHealingSessionClient) Send(m *SessionRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *selfHealingSessionClient) Recv() (*SessionReply, error) {
	m := new(SessionReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *selfHealingClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, SelfHealing_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *selfHealingClient) ProcessSelfHealingChallenge(ctx context.Context, in *ProcessSelfHealingChallengeRequest, opts ...grpc.CallOption) (*ProcessSelfHealingChallengeReply, error) {
	out := new(ProcessSelfHealingChallengeReply)
	err := c.cc.Invoke(ctx, SelfHealing_ProcessSelfHealingChallenge_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *selfHealingClient) VerifySelfHealingChallenge(ctx context.Context, in *VerifySelfHealingChallengeRequest, opts ...grpc.CallOption) (*VerifySelfHealingChallengeReply, error) {
	out := new(VerifySelfHealingChallengeReply)
	err := c.cc.Invoke(ctx, SelfHealing_VerifySelfHealingChallenge_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SelfHealingServer is the server API for SelfHealing service.
// All implementations must embed UnimplementedSelfHealingServer
// for forward compatibility
type SelfHealingServer interface {
	Session(SelfHealing_SessionServer) error
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	ProcessSelfHealingChallenge(context.Context, *ProcessSelfHealingChallengeRequest) (*ProcessSelfHealingChallengeReply, error)
	VerifySelfHealingChallenge(context.Context, *VerifySelfHealingChallengeRequest) (*VerifySelfHealingChallengeReply, error)
	mustEmbedUnimplementedSelfHealingServer()
}

// UnimplementedSelfHealingServer must be embedded to have forward compatible implementations.
type UnimplementedSelfHealingServer struct {
}

func (UnimplementedSelfHealingServer) Session(SelfHealing_SessionServer) error {
	return status.Errorf(codes.Unimplemented, "method Session not implemented")
}
func (UnimplementedSelfHealingServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSelfHealingServer) ProcessSelfHealingChallenge(context.Context, *ProcessSelfHealingChallengeRequest) (*ProcessSelfHealingChallengeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessSelfHealingChallenge not implemented")
}
func (UnimplementedSelfHealingServer) VerifySelfHealingChallenge(context.Context, *VerifySelfHealingChallengeRequest) (*VerifySelfHealingChallengeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifySelfHealingChallenge not implemented")
}
func (UnimplementedSelfHealingServer) mustEmbedUnimplementedSelfHealingServer() {}

// UnsafeSelfHealingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SelfHealingServer will
// result in compilation errors.
type UnsafeSelfHealingServer interface {
	mustEmbedUnimplementedSelfHealingServer()
}

func RegisterSelfHealingServer(s grpc.ServiceRegistrar, srv SelfHealingServer) {
	s.RegisterService(&SelfHealing_ServiceDesc, srv)
}

func _SelfHealing_Session_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SelfHealingServer).Session(&selfHealingSessionServer{stream})
}

type SelfHealing_SessionServer interface {
	Send(*SessionReply) error
	Recv() (*SessionRequest, error)
	grpc.ServerStream
}

type selfHealingSessionServer struct {
	grpc.ServerStream
}

func (x *selfHealingSessionServer) Send(m *SessionReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *selfHealingSessionServer) Recv() (*SessionRequest, error) {
	m := new(SessionRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SelfHealing_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SelfHealingServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SelfHealing_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SelfHealingServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SelfHealing_ProcessSelfHealingChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessSelfHealingChallengeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SelfHealingServer).ProcessSelfHealingChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SelfHealing_ProcessSelfHealingChallenge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SelfHealingServer).ProcessSelfHealingChallenge(ctx, req.(*ProcessSelfHealingChallengeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SelfHealing_VerifySelfHealingChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifySelfHealingChallengeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SelfHealingServer).VerifySelfHealingChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SelfHealing_VerifySelfHealingChallenge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SelfHealingServer).VerifySelfHealingChallenge(ctx, req.(*VerifySelfHealingChallengeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SelfHealing_ServiceDesc is the grpc.ServiceDesc for SelfHealing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SelfHealing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "supernode.SelfHealing",
	HandlerType: (*SelfHealingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _SelfHealing_Ping_Handler,
		},
		{
			MethodName: "ProcessSelfHealingChallenge",
			Handler:    _SelfHealing_ProcessSelfHealingChallenge_Handler,
		},
		{
			MethodName: "VerifySelfHealingChallenge",
			Handler:    _SelfHealing_VerifySelfHealingChallenge_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Session",
			Handler:       _SelfHealing_Session_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "self_healing.proto",
}
