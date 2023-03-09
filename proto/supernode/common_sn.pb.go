// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common_sn.proto

package supernode

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
	NodeID               string   `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionRequest) Reset()         { *m = SessionRequest{} }
func (m *SessionRequest) String() string { return proto.CompactTextString(m) }
func (*SessionRequest) ProtoMessage()    {}
func (*SessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f4ceffd962018273, []int{0}
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

func (m *SessionRequest) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
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
	return fileDescriptor_f4ceffd962018273, []int{1}
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

type SendSignedDDAndFingerprintsRequest struct {
	SessID                    string   `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID                    string   `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	ZstdCompressedFingerprint []byte   `protobuf:"bytes,3,opt,name=zstd_compressed_fingerprint,json=zstdCompressedFingerprint,proto3" json:"zstd_compressed_fingerprint,omitempty"`
	XXX_NoUnkeyedLiteral      struct{} `json:"-"`
	XXX_unrecognized          []byte   `json:"-"`
	XXX_sizecache             int32    `json:"-"`
}

func (m *SendSignedDDAndFingerprintsRequest) Reset()         { *m = SendSignedDDAndFingerprintsRequest{} }
func (m *SendSignedDDAndFingerprintsRequest) String() string { return proto.CompactTextString(m) }
func (*SendSignedDDAndFingerprintsRequest) ProtoMessage()    {}
func (*SendSignedDDAndFingerprintsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f4ceffd962018273, []int{2}
}

func (m *SendSignedDDAndFingerprintsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendSignedDDAndFingerprintsRequest.Unmarshal(m, b)
}
func (m *SendSignedDDAndFingerprintsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendSignedDDAndFingerprintsRequest.Marshal(b, m, deterministic)
}
func (m *SendSignedDDAndFingerprintsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendSignedDDAndFingerprintsRequest.Merge(m, src)
}
func (m *SendSignedDDAndFingerprintsRequest) XXX_Size() int {
	return xxx_messageInfo_SendSignedDDAndFingerprintsRequest.Size(m)
}
func (m *SendSignedDDAndFingerprintsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendSignedDDAndFingerprintsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendSignedDDAndFingerprintsRequest proto.InternalMessageInfo

func (m *SendSignedDDAndFingerprintsRequest) GetSessID() string {
	if m != nil {
		return m.SessID
	}
	return ""
}

func (m *SendSignedDDAndFingerprintsRequest) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

func (m *SendSignedDDAndFingerprintsRequest) GetZstdCompressedFingerprint() []byte {
	if m != nil {
		return m.ZstdCompressedFingerprint
	}
	return nil
}

type SendSignedDDAndFingerprintsReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendSignedDDAndFingerprintsReply) Reset()         { *m = SendSignedDDAndFingerprintsReply{} }
func (m *SendSignedDDAndFingerprintsReply) String() string { return proto.CompactTextString(m) }
func (*SendSignedDDAndFingerprintsReply) ProtoMessage()    {}
func (*SendSignedDDAndFingerprintsReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_f4ceffd962018273, []int{3}
}

func (m *SendSignedDDAndFingerprintsReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendSignedDDAndFingerprintsReply.Unmarshal(m, b)
}
func (m *SendSignedDDAndFingerprintsReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendSignedDDAndFingerprintsReply.Marshal(b, m, deterministic)
}
func (m *SendSignedDDAndFingerprintsReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendSignedDDAndFingerprintsReply.Merge(m, src)
}
func (m *SendSignedDDAndFingerprintsReply) XXX_Size() int {
	return xxx_messageInfo_SendSignedDDAndFingerprintsReply.Size(m)
}
func (m *SendSignedDDAndFingerprintsReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendSignedDDAndFingerprintsReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendSignedDDAndFingerprintsReply proto.InternalMessageInfo

type SendTicketSignatureRequest struct {
	NodeID               string   `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTicketSignatureRequest) Reset()         { *m = SendTicketSignatureRequest{} }
func (m *SendTicketSignatureRequest) String() string { return proto.CompactTextString(m) }
func (*SendTicketSignatureRequest) ProtoMessage()    {}
func (*SendTicketSignatureRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f4ceffd962018273, []int{4}
}

func (m *SendTicketSignatureRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTicketSignatureRequest.Unmarshal(m, b)
}
func (m *SendTicketSignatureRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTicketSignatureRequest.Marshal(b, m, deterministic)
}
func (m *SendTicketSignatureRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTicketSignatureRequest.Merge(m, src)
}
func (m *SendTicketSignatureRequest) XXX_Size() int {
	return xxx_messageInfo_SendTicketSignatureRequest.Size(m)
}
func (m *SendTicketSignatureRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTicketSignatureRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendTicketSignatureRequest proto.InternalMessageInfo

func (m *SendTicketSignatureRequest) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

func (m *SendTicketSignatureRequest) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type SendTicketSignatureReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTicketSignatureReply) Reset()         { *m = SendTicketSignatureReply{} }
func (m *SendTicketSignatureReply) String() string { return proto.CompactTextString(m) }
func (*SendTicketSignatureReply) ProtoMessage()    {}
func (*SendTicketSignatureReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_f4ceffd962018273, []int{5}
}

func (m *SendTicketSignatureReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTicketSignatureReply.Unmarshal(m, b)
}
func (m *SendTicketSignatureReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTicketSignatureReply.Marshal(b, m, deterministic)
}
func (m *SendTicketSignatureReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTicketSignatureReply.Merge(m, src)
}
func (m *SendTicketSignatureReply) XXX_Size() int {
	return xxx_messageInfo_SendTicketSignatureReply.Size(m)
}
func (m *SendTicketSignatureReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTicketSignatureReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendTicketSignatureReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SessionRequest)(nil), "supernode.SessionRequest")
	proto.RegisterType((*SessionReply)(nil), "supernode.SessionReply")
	proto.RegisterType((*SendSignedDDAndFingerprintsRequest)(nil), "supernode.SendSignedDDAndFingerprintsRequest")
	proto.RegisterType((*SendSignedDDAndFingerprintsReply)(nil), "supernode.SendSignedDDAndFingerprintsReply")
	proto.RegisterType((*SendTicketSignatureRequest)(nil), "supernode.SendTicketSignatureRequest")
	proto.RegisterType((*SendTicketSignatureReply)(nil), "supernode.SendTicketSignatureReply")
}

func init() { proto.RegisterFile("common_sn.proto", fileDescriptor_f4ceffd962018273) }

var fileDescriptor_f4ceffd962018273 = []byte{
	// 271 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x41, 0x4b, 0xf4, 0x30,
	0x10, 0x86, 0xe9, 0xf7, 0xc1, 0x42, 0x43, 0x51, 0xe8, 0x41, 0xea, 0xea, 0xa1, 0xe4, 0x20, 0x3d,
	0x6d, 0x11, 0xef, 0x82, 0x5a, 0x04, 0xaf, 0xa9, 0x27, 0x2f, 0x65, 0xb7, 0x19, 0x6b, 0xd8, 0x36,
	0x13, 0x33, 0x29, 0x52, 0x7f, 0x87, 0x3f, 0x58, 0x52, 0xb5, 0xed, 0x61, 0xd1, 0xe3, 0x4c, 0x9e,
	0xbc, 0xcf, 0x30, 0xc3, 0x8e, 0x6b, 0xec, 0x3a, 0xd4, 0x15, 0xe9, 0x8d, 0xb1, 0xe8, 0x30, 0x0e,
	0xa9, 0x37, 0x60, 0x35, 0x4a, 0xe0, 0x19, 0x3b, 0x2a, 0x81, 0x48, 0xa1, 0x16, 0xf0, 0xda, 0x03,
	0xb9, 0xf8, 0x84, 0xad, 0xfc, 0xcb, 0x43, 0x91, 0x04, 0x69, 0x90, 0x85, 0xe2, 0xbb, 0xe2, 0x17,
	0x2c, 0x9a, 0x48, 0xd3, 0x0e, 0x9e, 0x23, 0x20, 0x9a, 0xb9, 0xaf, 0x8a, 0x7f, 0x04, 0x8c, 0x97,
	0xa0, 0x65, 0xa9, 0x1a, 0x0d, 0xb2, 0x28, 0x6e, 0xb4, 0xbc, 0x57, 0xba, 0x01, 0x6b, 0xac, 0xd2,
	0x8e, 0x16, 0x9a, 0x43, 0xdf, 0x17, 0xfa, 0x7f, 0x4b, 0x7d, 0x7c, 0xcd, 0xce, 0xde, 0xc9, 0xc9,
	0xaa, 0xc6, 0xce, 0x58, 0x20, 0x02, 0x59, 0x3d, 0xcf, 0xb1, 0xc9, 0xff, 0x34, 0xc8, 0x22, 0x71,
	0xea, 0x91, 0xbb, 0x89, 0x58, 0x78, 0x39, 0x67, 0xe9, 0xaf, 0x53, 0x99, 0x76, 0xe0, 0x82, 0xad,
	0x3d, 0xf3, 0xa8, 0xea, 0x3d, 0x38, 0x4f, 0x6e, 0x5d, 0x6f, 0xe1, 0x8f, 0xc5, 0xc4, 0xe7, 0x2c,
	0xa4, 0x1f, 0x76, 0x1c, 0x3a, 0x12, 0x73, 0x83, 0xaf, 0x59, 0x72, 0x30, 0xd3, 0xb4, 0xc3, 0xed,
	0xe5, 0x53, 0xde, 0x28, 0xf7, 0xd2, 0xef, 0x36, 0x35, 0x76, 0xb9, 0xd9, 0x92, 0x83, 0x56, 0x83,
	0x7b, 0x43, 0xbb, 0xcf, 0x1b, 0xf4, 0xf1, 0xf9, 0x78, 0xb0, 0x7c, 0xba, 0xd7, 0x6e, 0x35, 0x36,
	0xae, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x33, 0x95, 0x80, 0xd4, 0x01, 0x00, 0x00,
}
