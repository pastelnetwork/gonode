// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.18.1
// source: external_dupe_detection.proto

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

type SendSignedEDDTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EddTicket         []byte             `protobuf:"bytes,1,opt,name=edd_ticket,json=eddTicket,proto3" json:"edd_ticket,omitempty"`
	CreatetorPastelId string             `protobuf:"bytes,2,opt,name=createtor_pastel_id,json=createtorPastelId,proto3" json:"createtor_pastel_id,omitempty"`
	CreatorSignature  []byte             `protobuf:"bytes,3,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	Key1              string             `protobuf:"bytes,4,opt,name=key1,proto3" json:"key1,omitempty"`
	Key2              string             `protobuf:"bytes,5,opt,name=key2,proto3" json:"key2,omitempty"`
	EncodeParameters  *EncoderParameters `protobuf:"bytes,6,opt,name=encode_parameters,json=encodeParameters,proto3" json:"encode_parameters,omitempty"`
	EncodeFiles       map[string][]byte  `protobuf:"bytes,7,rep,name=encode_files,json=encodeFiles,proto3" json:"encode_files,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *SendSignedEDDTicketRequest) Reset() {
	*x = SendSignedEDDTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_external_dupe_detection_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedEDDTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedEDDTicketRequest) ProtoMessage() {}

func (x *SendSignedEDDTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_external_dupe_detection_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedEDDTicketRequest.ProtoReflect.Descriptor instead.
func (*SendSignedEDDTicketRequest) Descriptor() ([]byte, []int) {
	return file_external_dupe_detection_proto_rawDescGZIP(), []int{0}
}

func (x *SendSignedEDDTicketRequest) GetEddTicket() []byte {
	if x != nil {
		return x.EddTicket
	}
	return nil
}

func (x *SendSignedEDDTicketRequest) GetCreatetorPastelId() string {
	if x != nil {
		return x.CreatetorPastelId
	}
	return ""
}

func (x *SendSignedEDDTicketRequest) GetCreatorSignature() []byte {
	if x != nil {
		return x.CreatorSignature
	}
	return nil
}

func (x *SendSignedEDDTicketRequest) GetKey1() string {
	if x != nil {
		return x.Key1
	}
	return ""
}

func (x *SendSignedEDDTicketRequest) GetKey2() string {
	if x != nil {
		return x.Key2
	}
	return ""
}

func (x *SendSignedEDDTicketRequest) GetEncodeParameters() *EncoderParameters {
	if x != nil {
		return x.EncodeParameters
	}
	return nil
}

func (x *SendSignedEDDTicketRequest) GetEncodeFiles() map[string][]byte {
	if x != nil {
		return x.EncodeFiles
	}
	return nil
}

type SendSignedEDDTicketReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RegistrationFee int64 `protobuf:"varint,1,opt,name=registration_fee,json=registrationFee,proto3" json:"registration_fee,omitempty"`
}

func (x *SendSignedEDDTicketReply) Reset() {
	*x = SendSignedEDDTicketReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_external_dupe_detection_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedEDDTicketReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedEDDTicketReply) ProtoMessage() {}

func (x *SendSignedEDDTicketReply) ProtoReflect() protoreflect.Message {
	mi := &file_external_dupe_detection_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedEDDTicketReply.ProtoReflect.Descriptor instead.
func (*SendSignedEDDTicketReply) Descriptor() ([]byte, []int) {
	return file_external_dupe_detection_proto_rawDescGZIP(), []int{1}
}

func (x *SendSignedEDDTicketReply) GetRegistrationFee() int64 {
	if x != nil {
		return x.RegistrationFee
	}
	return 0
}

type SendPreBurnedFeeEDDTxIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid string `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
}

func (x *SendPreBurnedFeeEDDTxIDRequest) Reset() {
	*x = SendPreBurnedFeeEDDTxIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_external_dupe_detection_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendPreBurnedFeeEDDTxIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendPreBurnedFeeEDDTxIDRequest) ProtoMessage() {}

func (x *SendPreBurnedFeeEDDTxIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_external_dupe_detection_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendPreBurnedFeeEDDTxIDRequest.ProtoReflect.Descriptor instead.
func (*SendPreBurnedFeeEDDTxIDRequest) Descriptor() ([]byte, []int) {
	return file_external_dupe_detection_proto_rawDescGZIP(), []int{2}
}

func (x *SendPreBurnedFeeEDDTxIDRequest) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

type SendPreBurnedFeeEDDTxIDReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EddRegTxid string `protobuf:"bytes,1,opt,name=edd_reg_txid,json=eddRegTxid,proto3" json:"edd_reg_txid,omitempty"`
}

func (x *SendPreBurnedFeeEDDTxIDReply) Reset() {
	*x = SendPreBurnedFeeEDDTxIDReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_external_dupe_detection_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendPreBurnedFeeEDDTxIDReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendPreBurnedFeeEDDTxIDReply) ProtoMessage() {}

func (x *SendPreBurnedFeeEDDTxIDReply) ProtoReflect() protoreflect.Message {
	mi := &file_external_dupe_detection_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendPreBurnedFeeEDDTxIDReply.ProtoReflect.Descriptor instead.
func (*SendPreBurnedFeeEDDTxIDReply) Descriptor() ([]byte, []int) {
	return file_external_dupe_detection_proto_rawDescGZIP(), []int{3}
}

func (x *SendPreBurnedFeeEDDTxIDReply) GetEddRegTxid() string {
	if x != nil {
		return x.EddRegTxid
	}
	return ""
}

var File_external_dupe_detection_proto protoreflect.FileDescriptor

var file_external_dupe_detection_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x64, 0x75, 0x70, 0x65, 0x5f,
	0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x16, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x72, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa8, 0x03, 0x0a, 0x1a, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x45, 0x44, 0x44, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x64, 0x64, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x65, 0x64, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x74, 0x6f, 0x72, 0x5f, 0x70,
	0x61, 0x73, 0x74, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x74, 0x6f, 0x72, 0x50, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49,
	0x64, 0x12, 0x2b, 0x0a, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x65, 0x79, 0x31, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65,
	0x79, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x32, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6b, 0x65, 0x79, 0x32, 0x12, 0x4a, 0x0a, 0x11, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x45,
	0x6e, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73,
	0x52, 0x10, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x12, 0x5a, 0x0a, 0x0c, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x45, 0x44, 0x44, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0b, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x3e,
	0x0a, 0x10, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x45,
	0x0a, 0x18, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x45, 0x44, 0x44, 0x54,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x10, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x65, 0x65, 0x22, 0x34, 0x0a, 0x1e, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x65,
	0x42, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x46, 0x65, 0x65, 0x45, 0x44, 0x44, 0x54, 0x78, 0x49, 0x44,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x22, 0x40, 0x0a, 0x1c, 0x53,
	0x65, 0x6e, 0x64, 0x50, 0x72, 0x65, 0x42, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x46, 0x65, 0x65, 0x45,
	0x44, 0x44, 0x54, 0x78, 0x49, 0x44, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x20, 0x0a, 0x0c, 0x65,
	0x64, 0x64, 0x5f, 0x72, 0x65, 0x67, 0x5f, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x65, 0x64, 0x64, 0x52, 0x65, 0x67, 0x54, 0x78, 0x69, 0x64, 0x32, 0xb1, 0x05,
	0x0a, 0x15, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x44, 0x75, 0x70, 0x65, 0x44, 0x65,
	0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x43, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x51, 0x0a, 0x0d,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x20, 0x2e,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70,
	0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x45, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x12, 0x1c, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54,
	0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4a, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x28, 0x01, 0x12, 0x63, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x45, 0x44, 0x44, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x26, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x45, 0x44, 0x44, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x24, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x45, 0x44, 0x44, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x6f, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x50,
	0x72, 0x65, 0x42, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x46, 0x65, 0x65, 0x45, 0x44, 0x44, 0x54, 0x78,
	0x49, 0x44, 0x12, 0x2a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x65, 0x42, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x46, 0x65, 0x65,
	0x45, 0x44, 0x44, 0x54, 0x78, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x50, 0x72, 0x65, 0x42, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x46, 0x65, 0x65, 0x45, 0x44, 0x44, 0x54,
	0x78, 0x49, 0x44, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x48, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64,
	0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x4d, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x12, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28,
	0x01, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f,
	0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_external_dupe_detection_proto_rawDescOnce sync.Once
	file_external_dupe_detection_proto_rawDescData = file_external_dupe_detection_proto_rawDesc
)

func file_external_dupe_detection_proto_rawDescGZIP() []byte {
	file_external_dupe_detection_proto_rawDescOnce.Do(func() {
		file_external_dupe_detection_proto_rawDescData = protoimpl.X.CompressGZIP(file_external_dupe_detection_proto_rawDescData)
	})
	return file_external_dupe_detection_proto_rawDescData
}

var file_external_dupe_detection_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_external_dupe_detection_proto_goTypes = []interface{}{
	(*SendSignedEDDTicketRequest)(nil),     // 0: walletnode.SendSignedEDDTicketRequest
	(*SendSignedEDDTicketReply)(nil),       // 1: walletnode.SendSignedEDDTicketReply
	(*SendPreBurnedFeeEDDTxIDRequest)(nil), // 2: walletnode.SendPreBurnedFeeEDDTxIDRequest
	(*SendPreBurnedFeeEDDTxIDReply)(nil),   // 3: walletnode.SendPreBurnedFeeEDDTxIDReply
	nil,                                    // 4: walletnode.SendSignedEDDTicketRequest.EncodeFilesEntry
	(*EncoderParameters)(nil),              // 5: walletnode.EncoderParameters
	(*SessionRequest)(nil),                 // 6: walletnode.SessionRequest
	(*AcceptedNodesRequest)(nil),           // 7: walletnode.AcceptedNodesRequest
	(*ConnectToRequest)(nil),               // 8: walletnode.ConnectToRequest
	(*ProbeImageRequest)(nil),              // 9: walletnode.ProbeImageRequest
	(*SendTicketRequest)(nil),              // 10: walletnode.SendTicketRequest
	(*UploadImageRequest)(nil),             // 11: walletnode.UploadImageRequest
	(*SessionReply)(nil),                   // 12: walletnode.SessionReply
	(*AcceptedNodesReply)(nil),             // 13: walletnode.AcceptedNodesReply
	(*ConnectToReply)(nil),                 // 14: walletnode.ConnectToReply
	(*ProbeImageReply)(nil),                // 15: walletnode.ProbeImageReply
	(*SendTicketReply)(nil),                // 16: walletnode.SendTicketReply
	(*UploadImageReply)(nil),               // 17: walletnode.UploadImageReply
}
var file_external_dupe_detection_proto_depIdxs = []int32{
	5,  // 0: walletnode.SendSignedEDDTicketRequest.encode_parameters:type_name -> walletnode.EncoderParameters
	4,  // 1: walletnode.SendSignedEDDTicketRequest.encode_files:type_name -> walletnode.SendSignedEDDTicketRequest.EncodeFilesEntry
	6,  // 2: walletnode.ExternalDupeDetection.Session:input_type -> walletnode.SessionRequest
	7,  // 3: walletnode.ExternalDupeDetection.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	8,  // 4: walletnode.ExternalDupeDetection.ConnectTo:input_type -> walletnode.ConnectToRequest
	9,  // 5: walletnode.ExternalDupeDetection.ProbeImage:input_type -> walletnode.ProbeImageRequest
	0,  // 6: walletnode.ExternalDupeDetection.SendSignedEDDTicket:input_type -> walletnode.SendSignedEDDTicketRequest
	2,  // 7: walletnode.ExternalDupeDetection.SendPreBurnedFeeEDDTxID:input_type -> walletnode.SendPreBurnedFeeEDDTxIDRequest
	10, // 8: walletnode.ExternalDupeDetection.SendTicket:input_type -> walletnode.SendTicketRequest
	11, // 9: walletnode.ExternalDupeDetection.UploadImage:input_type -> walletnode.UploadImageRequest
	12, // 10: walletnode.ExternalDupeDetection.Session:output_type -> walletnode.SessionReply
	13, // 11: walletnode.ExternalDupeDetection.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	14, // 12: walletnode.ExternalDupeDetection.ConnectTo:output_type -> walletnode.ConnectToReply
	15, // 13: walletnode.ExternalDupeDetection.ProbeImage:output_type -> walletnode.ProbeImageReply
	1,  // 14: walletnode.ExternalDupeDetection.SendSignedEDDTicket:output_type -> walletnode.SendSignedEDDTicketReply
	3,  // 15: walletnode.ExternalDupeDetection.SendPreBurnedFeeEDDTxID:output_type -> walletnode.SendPreBurnedFeeEDDTxIDReply
	16, // 16: walletnode.ExternalDupeDetection.SendTicket:output_type -> walletnode.SendTicketReply
	17, // 17: walletnode.ExternalDupeDetection.UploadImage:output_type -> walletnode.UploadImageReply
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_external_dupe_detection_proto_init() }
func file_external_dupe_detection_proto_init() {
	if File_external_dupe_detection_proto != nil {
		return
	}
	file_register_artwork_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_external_dupe_detection_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedEDDTicketRequest); i {
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
		file_external_dupe_detection_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedEDDTicketReply); i {
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
		file_external_dupe_detection_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendPreBurnedFeeEDDTxIDRequest); i {
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
		file_external_dupe_detection_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendPreBurnedFeeEDDTxIDReply); i {
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
			RawDescriptor: file_external_dupe_detection_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_external_dupe_detection_proto_goTypes,
		DependencyIndexes: file_external_dupe_detection_proto_depIdxs,
		MessageInfos:      file_external_dupe_detection_proto_msgTypes,
	}.Build()
	File_external_dupe_detection_proto = out.File
	file_external_dupe_detection_proto_rawDesc = nil
	file_external_dupe_detection_proto_goTypes = nil
	file_external_dupe_detection_proto_depIdxs = nil
}