// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common_bridge.proto

package bridge

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SessionRequest struct {
	IsPrimary            bool     `protobuf:"varint,1,opt,name=is_primary,json=isPrimary,proto3" json:"is_primary,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionRequest) Reset()         { *m = SessionRequest{} }
func (m *SessionRequest) String() string { return proto.CompactTextString(m) }
func (*SessionRequest) ProtoMessage()    {}
func (*SessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{0}
}

func (m *SessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionRequest.Unmarshal(m, b)
}
func (m *SessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionRequest.Marshal(b, m, deterministic)
}
func (m *SessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionRequest.Merge(m, src)
}
func (m *SessionRequest) XXX_Size() int {
	return xxx_messageInfo_SessionRequest.Size(m)
}
func (m *SessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SessionRequest proto.InternalMessageInfo

func (m *SessionRequest) GetIsPrimary() bool {
	if m != nil {
		return m.IsPrimary
	}
	return false
}

type SessionReply struct {
	SessID               string   `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionReply) Reset()         { *m = SessionReply{} }
func (m *SessionReply) String() string { return proto.CompactTextString(m) }
func (*SessionReply) ProtoMessage()    {}
func (*SessionReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{1}
}

func (m *SessionReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionReply.Unmarshal(m, b)
}
func (m *SessionReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionReply.Marshal(b, m, deterministic)
}
func (m *SessionReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionReply.Merge(m, src)
}
func (m *SessionReply) XXX_Size() int {
	return xxx_messageInfo_SessionReply.Size(m)
}
func (m *SessionReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionReply.DiscardUnknown(m)
}

var xxx_messageInfo_SessionReply proto.InternalMessageInfo

func (m *SessionReply) GetSessID() string {
	if m != nil {
		return m.SessID
	}
	return ""
}

type AcceptedNodesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptedNodesRequest) Reset()         { *m = AcceptedNodesRequest{} }
func (m *AcceptedNodesRequest) String() string { return proto.CompactTextString(m) }
func (*AcceptedNodesRequest) ProtoMessage()    {}
func (*AcceptedNodesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{2}
}

func (m *AcceptedNodesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptedNodesRequest.Unmarshal(m, b)
}
func (m *AcceptedNodesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptedNodesRequest.Marshal(b, m, deterministic)
}
func (m *AcceptedNodesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptedNodesRequest.Merge(m, src)
}
func (m *AcceptedNodesRequest) XXX_Size() int {
	return xxx_messageInfo_AcceptedNodesRequest.Size(m)
}
func (m *AcceptedNodesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptedNodesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptedNodesRequest proto.InternalMessageInfo

type AcceptedNodesReply struct {
	Peers                []*AcceptedNodesReply_Peer `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *AcceptedNodesReply) Reset()         { *m = AcceptedNodesReply{} }
func (m *AcceptedNodesReply) String() string { return proto.CompactTextString(m) }
func (*AcceptedNodesReply) ProtoMessage()    {}
func (*AcceptedNodesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{3}
}

func (m *AcceptedNodesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptedNodesReply.Unmarshal(m, b)
}
func (m *AcceptedNodesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptedNodesReply.Marshal(b, m, deterministic)
}
func (m *AcceptedNodesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptedNodesReply.Merge(m, src)
}
func (m *AcceptedNodesReply) XXX_Size() int {
	return xxx_messageInfo_AcceptedNodesReply.Size(m)
}
func (m *AcceptedNodesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptedNodesReply.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptedNodesReply proto.InternalMessageInfo

func (m *AcceptedNodesReply) GetPeers() []*AcceptedNodesReply_Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type AcceptedNodesReply_Peer struct {
	NodeID               string   `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptedNodesReply_Peer) Reset()         { *m = AcceptedNodesReply_Peer{} }
func (m *AcceptedNodesReply_Peer) String() string { return proto.CompactTextString(m) }
func (*AcceptedNodesReply_Peer) ProtoMessage()    {}
func (*AcceptedNodesReply_Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{3, 0}
}

func (m *AcceptedNodesReply_Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptedNodesReply_Peer.Unmarshal(m, b)
}
func (m *AcceptedNodesReply_Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptedNodesReply_Peer.Marshal(b, m, deterministic)
}
func (m *AcceptedNodesReply_Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptedNodesReply_Peer.Merge(m, src)
}
func (m *AcceptedNodesReply_Peer) XXX_Size() int {
	return xxx_messageInfo_AcceptedNodesReply_Peer.Size(m)
}
func (m *AcceptedNodesReply_Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptedNodesReply_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptedNodesReply_Peer proto.InternalMessageInfo

func (m *AcceptedNodesReply_Peer) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

type ConnectToRequest struct {
	SessID               string   `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID               string   `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectToRequest) Reset()         { *m = ConnectToRequest{} }
func (m *ConnectToRequest) String() string { return proto.CompactTextString(m) }
func (*ConnectToRequest) ProtoMessage()    {}
func (*ConnectToRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{4}
}

func (m *ConnectToRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectToRequest.Unmarshal(m, b)
}
func (m *ConnectToRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectToRequest.Marshal(b, m, deterministic)
}
func (m *ConnectToRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectToRequest.Merge(m, src)
}
func (m *ConnectToRequest) XXX_Size() int {
	return xxx_messageInfo_ConnectToRequest.Size(m)
}
func (m *ConnectToRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectToRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectToRequest proto.InternalMessageInfo

func (m *ConnectToRequest) GetSessID() string {
	if m != nil {
		return m.SessID
	}
	return ""
}

func (m *ConnectToRequest) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

type ConnectToReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectToReply) Reset()         { *m = ConnectToReply{} }
func (m *ConnectToReply) String() string { return proto.CompactTextString(m) }
func (*ConnectToReply) ProtoMessage()    {}
func (*ConnectToReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{5}
}

func (m *ConnectToReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectToReply.Unmarshal(m, b)
}
func (m *ConnectToReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectToReply.Marshal(b, m, deterministic)
}
func (m *ConnectToReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectToReply.Merge(m, src)
}
func (m *ConnectToReply) XXX_Size() int {
	return xxx_messageInfo_ConnectToReply.Size(m)
}
func (m *ConnectToReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectToReply.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectToReply proto.InternalMessageInfo

type MeshNodesRequest struct {
	Nodes                []*MeshNodesRequest_Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *MeshNodesRequest) Reset()         { *m = MeshNodesRequest{} }
func (m *MeshNodesRequest) String() string { return proto.CompactTextString(m) }
func (*MeshNodesRequest) ProtoMessage()    {}
func (*MeshNodesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{6}
}

func (m *MeshNodesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MeshNodesRequest.Unmarshal(m, b)
}
func (m *MeshNodesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MeshNodesRequest.Marshal(b, m, deterministic)
}
func (m *MeshNodesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeshNodesRequest.Merge(m, src)
}
func (m *MeshNodesRequest) XXX_Size() int {
	return xxx_messageInfo_MeshNodesRequest.Size(m)
}
func (m *MeshNodesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MeshNodesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MeshNodesRequest proto.InternalMessageInfo

func (m *MeshNodesRequest) GetNodes() []*MeshNodesRequest_Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type MeshNodesRequest_Node struct {
	SessID               string   `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID               string   `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MeshNodesRequest_Node) Reset()         { *m = MeshNodesRequest_Node{} }
func (m *MeshNodesRequest_Node) String() string { return proto.CompactTextString(m) }
func (*MeshNodesRequest_Node) ProtoMessage()    {}
func (*MeshNodesRequest_Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{6, 0}
}

func (m *MeshNodesRequest_Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MeshNodesRequest_Node.Unmarshal(m, b)
}
func (m *MeshNodesRequest_Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MeshNodesRequest_Node.Marshal(b, m, deterministic)
}
func (m *MeshNodesRequest_Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeshNodesRequest_Node.Merge(m, src)
}
func (m *MeshNodesRequest_Node) XXX_Size() int {
	return xxx_messageInfo_MeshNodesRequest_Node.Size(m)
}
func (m *MeshNodesRequest_Node) XXX_DiscardUnknown() {
	xxx_messageInfo_MeshNodesRequest_Node.DiscardUnknown(m)
}

var xxx_messageInfo_MeshNodesRequest_Node proto.InternalMessageInfo

func (m *MeshNodesRequest_Node) GetSessID() string {
	if m != nil {
		return m.SessID
	}
	return ""
}

func (m *MeshNodesRequest_Node) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

type MeshNodesReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MeshNodesReply) Reset()         { *m = MeshNodesReply{} }
func (m *MeshNodesReply) String() string { return proto.CompactTextString(m) }
func (*MeshNodesReply) ProtoMessage()    {}
func (*MeshNodesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9bfa44c8dd3b80f, []int{7}
}

func (m *MeshNodesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MeshNodesReply.Unmarshal(m, b)
}
func (m *MeshNodesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MeshNodesReply.Marshal(b, m, deterministic)
}
func (m *MeshNodesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeshNodesReply.Merge(m, src)
}
func (m *MeshNodesReply) XXX_Size() int {
	return xxx_messageInfo_MeshNodesReply.Size(m)
}
func (m *MeshNodesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MeshNodesReply.DiscardUnknown(m)
}

var xxx_messageInfo_MeshNodesReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SessionRequest)(nil), "bridge.SessionRequest")
	proto.RegisterType((*SessionReply)(nil), "bridge.SessionReply")
	proto.RegisterType((*AcceptedNodesRequest)(nil), "bridge.AcceptedNodesRequest")
	proto.RegisterType((*AcceptedNodesReply)(nil), "bridge.AcceptedNodesReply")
	proto.RegisterType((*AcceptedNodesReply_Peer)(nil), "bridge.AcceptedNodesReply.Peer")
	proto.RegisterType((*ConnectToRequest)(nil), "bridge.ConnectToRequest")
	proto.RegisterType((*ConnectToReply)(nil), "bridge.ConnectToReply")
	proto.RegisterType((*MeshNodesRequest)(nil), "bridge.MeshNodesRequest")
	proto.RegisterType((*MeshNodesRequest_Node)(nil), "bridge.MeshNodesRequest.Node")
	proto.RegisterType((*MeshNodesReply)(nil), "bridge.MeshNodesReply")
}

func init() { proto.RegisterFile("common_bridge.proto", fileDescriptor_e9bfa44c8dd3b80f) }

var fileDescriptor_e9bfa44c8dd3b80f = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xc1, 0x4b, 0x3a, 0x41,
	0x14, 0xc7, 0xf1, 0xf7, 0x53, 0xc9, 0x57, 0x88, 0x6c, 0x21, 0x22, 0x58, 0xb2, 0x87, 0xf0, 0x10,
	0xb3, 0x90, 0xd4, 0x3d, 0xeb, 0xd2, 0xa1, 0x90, 0xad, 0x53, 0x17, 0xd1, 0xdd, 0x87, 0x0e, 0xee,
	0xce, 0x9b, 0xe6, 0x8d, 0xc4, 0x9e, 0xfa, 0xd7, 0x63, 0x1c, 0xb7, 0x36, 0xa3, 0x4b, 0xc7, 0xcf,
	0xf0, 0x7d, 0x5f, 0x3e, 0x5f, 0x06, 0x8e, 0x13, 0xca, 0x73, 0x52, 0xb3, 0x85, 0x91, 0xe9, 0x12,
	0x85, 0x36, 0x64, 0x29, 0x68, 0x7a, 0x0a, 0x23, 0x68, 0x3f, 0x21, 0xb3, 0x24, 0x15, 0xe3, 0xeb,
	0x06, 0xd9, 0x06, 0x03, 0x00, 0xc9, 0x33, 0x6d, 0x64, 0x3e, 0x37, 0x45, 0xaf, 0x36, 0xac, 0x8d,
	0x0e, 0xe2, 0x96, 0xe4, 0xa9, 0x7f, 0x08, 0xcf, 0xe1, 0xe8, 0xf3, 0x40, 0x67, 0x45, 0xd0, 0x85,
	0x26, 0x23, 0xf3, 0xfd, 0xdd, 0x36, 0xda, 0x8a, 0x77, 0x14, 0x76, 0xe1, 0xe4, 0x26, 0x49, 0x50,
	0x5b, 0x4c, 0x1f, 0x29, 0x45, 0xde, 0xd5, 0x87, 0x6b, 0x08, 0xf6, 0xde, 0x5d, 0xcb, 0x15, 0x34,
	0x34, 0xa2, 0xe1, 0x5e, 0x6d, 0xf8, 0x7f, 0x74, 0x78, 0x79, 0x26, 0x76, 0xb2, 0x3f, 0xa3, 0x62,
	0x8a, 0x68, 0x62, 0x9f, 0xee, 0x9f, 0x42, 0xdd, 0xa1, 0x93, 0x50, 0x94, 0xe2, 0x97, 0x84, 0xa7,
	0x70, 0x02, 0x9d, 0x5b, 0x52, 0x0a, 0x13, 0xfb, 0x4c, 0xe5, 0xbe, 0x5f, 0x84, 0x2b, 0x1d, 0xff,
	0xbe, 0x75, 0x74, 0xa0, 0x5d, 0xe9, 0xd0, 0x59, 0x11, 0xbe, 0x43, 0xe7, 0x01, 0x79, 0x55, 0x9d,
	0x15, 0x8c, 0xa1, 0xe1, 0xf2, 0xe5, 0x80, 0x41, 0x39, 0x60, 0x3f, 0x28, 0x1c, 0xc4, 0x3e, 0xdb,
	0xbf, 0x86, 0xba, 0xc3, 0xbf, 0x28, 0x55, 0x7a, 0x75, 0x56, 0x4c, 0xc4, 0xcb, 0xc5, 0x52, 0xda,
	0xd5, 0x66, 0x21, 0x12, 0xca, 0x23, 0x3d, 0x67, 0x8b, 0x99, 0x42, 0xfb, 0x46, 0x66, 0x1d, 0x2d,
	0xc9, 0x9d, 0x45, 0xdb, 0x7f, 0x8f, 0xbc, 0xd6, 0xa2, 0xb9, 0xa5, 0xf1, 0x47, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xa7, 0xd2, 0x93, 0x02, 0x1c, 0x02, 0x00, 0x00,
}
