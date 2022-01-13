// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: register_sense_wn.proto

package walletnode

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SendSignedActionTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionTicket     []byte `protobuf:"bytes,1,opt,name=action_ticket,json=actionTicket,proto3" json:"action_ticket,omitempty"`
	CreatorSignature []byte `protobuf:"bytes,2,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	DdFpFiles        []byte `protobuf:"bytes,3,opt,name=dd_fp_files,json=ddFpFiles,proto3" json:"dd_fp_files,omitempty"`
}

func (x *SendSignedActionTicketRequest) Reset() {
	*x = SendSignedActionTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_sense_wn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedActionTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedActionTicketRequest) ProtoMessage() {}

func (x *SendSignedActionTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_sense_wn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedActionTicketRequest.ProtoReflect.Descriptor instead.
func (*SendSignedActionTicketRequest) Descriptor() ([]byte, []int) {
	return file_register_sense_wn_proto_rawDescGZIP(), []int{0}
}

func (x *SendSignedActionTicketRequest) GetActionTicket() []byte {
	if x != nil {
		return x.ActionTicket
	}
	return nil
}

func (x *SendSignedActionTicketRequest) GetCreatorSignature() []byte {
	if x != nil {
		return x.CreatorSignature
	}
	return nil
}

func (x *SendSignedActionTicketRequest) GetDdFpFiles() []byte {
	if x != nil {
		return x.DdFpFiles
	}
	return nil
}

type SendSignedActionTicketReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionRegTxid string `protobuf:"bytes,1,opt,name=action_reg_txid,json=actionRegTxid,proto3" json:"action_reg_txid,omitempty"`
}

func (x *SendSignedActionTicketReply) Reset() {
	*x = SendSignedActionTicketReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_sense_wn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedActionTicketReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedActionTicketReply) ProtoMessage() {}

func (x *SendSignedActionTicketReply) ProtoReflect() protoreflect.Message {
	mi := &file_register_sense_wn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedActionTicketReply.ProtoReflect.Descriptor instead.
func (*SendSignedActionTicketReply) Descriptor() ([]byte, []int) {
	return file_register_sense_wn_proto_rawDescGZIP(), []int{1}
}

func (x *SendSignedActionTicketReply) GetActionRegTxid() string {
	if x != nil {
		return x.ActionRegTxid
	}
	return ""
}

type SendActionActRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionRegTxid string `protobuf:"bytes,1,opt,name=action_reg_txid,json=actionRegTxid,proto3" json:"action_reg_txid,omitempty"`
}

func (x *SendActionActRequest) Reset() {
	*x = SendActionActRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_sense_wn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendActionActRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendActionActRequest) ProtoMessage() {}

func (x *SendActionActRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_sense_wn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendActionActRequest.ProtoReflect.Descriptor instead.
func (*SendActionActRequest) Descriptor() ([]byte, []int) {
	return file_register_sense_wn_proto_rawDescGZIP(), []int{2}
}

func (x *SendActionActRequest) GetActionRegTxid() string {
	if x != nil {
		return x.ActionRegTxid
	}
	return ""
}

type SendActionActReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendActionActReply) Reset() {
	*x = SendActionActReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_sense_wn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendActionActReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendActionActReply) ProtoMessage() {}

func (x *SendActionActReply) ProtoReflect() protoreflect.Message {
	mi := &file_register_sense_wn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendActionActReply.ProtoReflect.Descriptor instead.
func (*SendActionActReply) Descriptor() ([]byte, []int) {
	return file_register_sense_wn_proto_rawDescGZIP(), []int{3}
}

var File_register_sense_wn_proto protoreflect.FileDescriptor

var file_register_sense_wn_proto_rawDesc = []byte{
	0x0a, 0x17, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x65,
	0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x77, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x1d, 0x53, 0x65, 0x6e, 0x64, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0c, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x2b, 0x0a,
	0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1e, 0x0a, 0x0b, 0x64, 0x64,
	0x5f, 0x66, 0x70, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x64, 0x64, 0x46, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x45, 0x0a, 0x1b, 0x53, 0x65,
	0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x67, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x67, 0x54, 0x78, 0x69,
	0x64, 0x22, 0x3e, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x67, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x67, 0x54, 0x78, 0x69,
	0x64, 0x22, 0x14, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41,
	0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x9b, 0x05, 0x0a, 0x0d, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x07, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x18, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x51,
	0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12,
	0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x45, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x12, 0x1c,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x54, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x68,
	0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x57, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x22, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4a, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x62,
	0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x28, 0x01, 0x12, 0x6c, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x29,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x51, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x41, 0x63, 0x74, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_register_sense_wn_proto_rawDescOnce sync.Once
	file_register_sense_wn_proto_rawDescData = file_register_sense_wn_proto_rawDesc
)

func file_register_sense_wn_proto_rawDescGZIP() []byte {
	file_register_sense_wn_proto_rawDescOnce.Do(func() {
		file_register_sense_wn_proto_rawDescData = protoimpl.X.CompressGZIP(file_register_sense_wn_proto_rawDescData)
	})
	return file_register_sense_wn_proto_rawDescData
}

var file_register_sense_wn_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_register_sense_wn_proto_goTypes = []interface{}{
	(*SendSignedActionTicketRequest)(nil), // 0: walletnode.SendSignedActionTicketRequest
	(*SendSignedActionTicketReply)(nil),   // 1: walletnode.SendSignedActionTicketReply
	(*SendActionActRequest)(nil),          // 2: walletnode.SendActionActRequest
	(*SendActionActReply)(nil),            // 3: walletnode.SendActionActReply
	(*SessionRequest)(nil),                // 4: walletnode.SessionRequest
	(*AcceptedNodesRequest)(nil),          // 5: walletnode.AcceptedNodesRequest
	(*ConnectToRequest)(nil),              // 6: walletnode.ConnectToRequest
	(*MeshNodesRequest)(nil),              // 7: walletnode.MeshNodesRequest
	(*SendRegMetadataRequest)(nil),        // 8: walletnode.SendRegMetadataRequest
	(*ProbeImageRequest)(nil),             // 9: walletnode.ProbeImageRequest
	(*SessionReply)(nil),                  // 10: walletnode.SessionReply
	(*AcceptedNodesReply)(nil),            // 11: walletnode.AcceptedNodesReply
	(*ConnectToReply)(nil),                // 12: walletnode.ConnectToReply
	(*MeshNodesReply)(nil),                // 13: walletnode.MeshNodesReply
	(*SendRegMetadataReply)(nil),          // 14: walletnode.SendRegMetadataReply
	(*ProbeImageReply)(nil),               // 15: walletnode.ProbeImageReply
}
var file_register_sense_wn_proto_depIdxs = []int32{
	4,  // 0: walletnode.RegisterSense.Session:input_type -> walletnode.SessionRequest
	5,  // 1: walletnode.RegisterSense.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	6,  // 2: walletnode.RegisterSense.ConnectTo:input_type -> walletnode.ConnectToRequest
	7,  // 3: walletnode.RegisterSense.MeshNodes:input_type -> walletnode.MeshNodesRequest
	8,  // 4: walletnode.RegisterSense.SendRegMetadata:input_type -> walletnode.SendRegMetadataRequest
	9,  // 5: walletnode.RegisterSense.ProbeImage:input_type -> walletnode.ProbeImageRequest
	0,  // 6: walletnode.RegisterSense.SendSignedActionTicket:input_type -> walletnode.SendSignedActionTicketRequest
	2,  // 7: walletnode.RegisterSense.SendActionAct:input_type -> walletnode.SendActionActRequest
	10, // 8: walletnode.RegisterSense.Session:output_type -> walletnode.SessionReply
	11, // 9: walletnode.RegisterSense.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	12, // 10: walletnode.RegisterSense.ConnectTo:output_type -> walletnode.ConnectToReply
	13, // 11: walletnode.RegisterSense.MeshNodes:output_type -> walletnode.MeshNodesReply
	14, // 12: walletnode.RegisterSense.SendRegMetadata:output_type -> walletnode.SendRegMetadataReply
	15, // 13: walletnode.RegisterSense.ProbeImage:output_type -> walletnode.ProbeImageReply
	1,  // 14: walletnode.RegisterSense.SendSignedActionTicket:output_type -> walletnode.SendSignedActionTicketReply
	3,  // 15: walletnode.RegisterSense.SendActionAct:output_type -> walletnode.SendActionActReply
	8,  // [8:16] is the sub-list for method output_type
	0,  // [0:8] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_register_sense_wn_proto_init() }
func file_register_sense_wn_proto_init() {
	if File_register_sense_wn_proto != nil {
		return
	}
	file_common_wn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_register_sense_wn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedActionTicketRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_register_sense_wn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedActionTicketReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_register_sense_wn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendActionActRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_register_sense_wn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendActionActReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_register_sense_wn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_register_sense_wn_proto_goTypes,
		DependencyIndexes: file_register_sense_wn_proto_depIdxs,
		MessageInfos:      file_register_sense_wn_proto_msgTypes,
	}.Build()
	File_register_sense_wn_proto = out.File
	file_register_sense_wn_proto_rawDesc = nil
	file_register_sense_wn_proto_goTypes = nil
	file_register_sense_wn_proto_depIdxs = nil
}