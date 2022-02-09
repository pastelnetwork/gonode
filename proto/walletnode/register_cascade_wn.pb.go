// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: register_cascade_wn.proto

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

type SendSignedCascadeTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionTicket     []byte             `protobuf:"bytes,1,opt,name=action_ticket,json=actionTicket,proto3" json:"action_ticket,omitempty"`
	CreatorSignature []byte             `protobuf:"bytes,2,opt,name=creator_signature,json=creatorSignature,proto3" json:"creator_signature,omitempty"`
	EncodeParameters *EncoderParameters `protobuf:"bytes,3,opt,name=encode_parameters,json=encodeParameters,proto3" json:"encode_parameters,omitempty"`
	RqFiles          []byte             `protobuf:"bytes,4,opt,name=rq_files,json=rqFiles,proto3" json:"rq_files,omitempty"`
}

func (x *SendSignedCascadeTicketRequest) Reset() {
	*x = SendSignedCascadeTicketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_cascade_wn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendSignedCascadeTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendSignedCascadeTicketRequest) ProtoMessage() {}

func (x *SendSignedCascadeTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_cascade_wn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendSignedCascadeTicketRequest.ProtoReflect.Descriptor instead.
func (*SendSignedCascadeTicketRequest) Descriptor() ([]byte, []int) {
	return file_register_cascade_wn_proto_rawDescGZIP(), []int{0}
}

func (x *SendSignedCascadeTicketRequest) GetActionTicket() []byte {
	if x != nil {
		return x.ActionTicket
	}
	return nil
}

func (x *SendSignedCascadeTicketRequest) GetCreatorSignature() []byte {
	if x != nil {
		return x.CreatorSignature
	}
	return nil
}

func (x *SendSignedCascadeTicketRequest) GetEncodeParameters() *EncoderParameters {
	if x != nil {
		return x.EncodeParameters
	}
	return nil
}

func (x *SendSignedCascadeTicketRequest) GetRqFiles() []byte {
	if x != nil {
		return x.RqFiles
	}
	return nil
}

type UploadAssetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Payload:
	//	*UploadAssetRequest_AssetPiece
	//	*UploadAssetRequest_MetaData_
	Payload isUploadAssetRequest_Payload `protobuf_oneof:"payload"`
}

func (x *UploadAssetRequest) Reset() {
	*x = UploadAssetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_cascade_wn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadAssetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadAssetRequest) ProtoMessage() {}

func (x *UploadAssetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_register_cascade_wn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadAssetRequest.ProtoReflect.Descriptor instead.
func (*UploadAssetRequest) Descriptor() ([]byte, []int) {
	return file_register_cascade_wn_proto_rawDescGZIP(), []int{1}
}

func (m *UploadAssetRequest) GetPayload() isUploadAssetRequest_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *UploadAssetRequest) GetAssetPiece() []byte {
	if x, ok := x.GetPayload().(*UploadAssetRequest_AssetPiece); ok {
		return x.AssetPiece
	}
	return nil
}

func (x *UploadAssetRequest) GetMetaData() *UploadAssetRequest_MetaData {
	if x, ok := x.GetPayload().(*UploadAssetRequest_MetaData_); ok {
		return x.MetaData
	}
	return nil
}

type isUploadAssetRequest_Payload interface {
	isUploadAssetRequest_Payload()
}

type UploadAssetRequest_AssetPiece struct {
	AssetPiece []byte `protobuf:"bytes,1,opt,name=asset_piece,json=assetPiece,proto3,oneof"`
}

type UploadAssetRequest_MetaData_ struct {
	// Should be included in the last piece of image only
	MetaData *UploadAssetRequest_MetaData `protobuf:"bytes,2,opt,name=meta_data,json=metaData,proto3,oneof"`
}

func (*UploadAssetRequest_AssetPiece) isUploadAssetRequest_Payload() {}

func (*UploadAssetRequest_MetaData_) isUploadAssetRequest_Payload() {}

type UploadAssetReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadAssetReply) Reset() {
	*x = UploadAssetReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_cascade_wn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadAssetReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadAssetReply) ProtoMessage() {}

func (x *UploadAssetReply) ProtoReflect() protoreflect.Message {
	mi := &file_register_cascade_wn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadAssetReply.ProtoReflect.Descriptor instead.
func (*UploadAssetReply) Descriptor() ([]byte, []int) {
	return file_register_cascade_wn_proto_rawDescGZIP(), []int{2}
}

type UploadAssetRequest_MetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// size of the data
	Size int64 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	// format of data: such as jpeg, png ...
	Format string `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
	// sha3-256 hash of the data
	Hash []byte `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *UploadAssetRequest_MetaData) Reset() {
	*x = UploadAssetRequest_MetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_cascade_wn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadAssetRequest_MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadAssetRequest_MetaData) ProtoMessage() {}

func (x *UploadAssetRequest_MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_register_cascade_wn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadAssetRequest_MetaData.ProtoReflect.Descriptor instead.
func (*UploadAssetRequest_MetaData) Descriptor() ([]byte, []int) {
	return file_register_cascade_wn_proto_rawDescGZIP(), []int{1, 0}
}

func (x *UploadAssetRequest_MetaData) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *UploadAssetRequest_MetaData) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *UploadAssetRequest_MetaData) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

var File_register_cascade_wn_proto protoreflect.FileDescriptor

var file_register_cascade_wn_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x61, 0x73, 0x63, 0x61,
	0x64, 0x65, 0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f,
	0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd9, 0x01, 0x0a, 0x1e, 0x53, 0x65, 0x6e,
	0x64, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x2b, 0x0a, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x4a, 0x0a,
	0x11, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x52, 0x10, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x71, 0x5f,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x71, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x22, 0xd6, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41,
	0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0b, 0x61,
	0x73, 0x73, 0x65, 0x74, 0x5f, 0x70, 0x69, 0x65, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x00, 0x52, 0x0a, 0x61, 0x73, 0x73, 0x65, 0x74, 0x50, 0x69, 0x65, 0x63, 0x65, 0x12, 0x46,
	0x0a, 0x09, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x27, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x4a, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x12, 0x0a,
	0x10, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x32, 0xa1, 0x05, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x61,
	0x73, 0x63, 0x61, 0x64, 0x65, 0x12, 0x43, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x51, 0x0a, 0x0d, 0x41, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65,
	0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70,
	0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a,
	0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d,
	0x65, 0x73, 0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x65, 0x73,
	0x68, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x57, 0x0a, 0x0f, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x22,
	0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x4d, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x28, 0x01, 0x12, 0x6d, 0x0a, 0x16, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x2a, 0x2e,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x43, 0x61, 0x73, 0x63, 0x61, 0x64, 0x65, 0x54, 0x69, 0x63, 0x6b,
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
	file_register_cascade_wn_proto_rawDescOnce sync.Once
	file_register_cascade_wn_proto_rawDescData = file_register_cascade_wn_proto_rawDesc
)

func file_register_cascade_wn_proto_rawDescGZIP() []byte {
	file_register_cascade_wn_proto_rawDescOnce.Do(func() {
		file_register_cascade_wn_proto_rawDescData = protoimpl.X.CompressGZIP(file_register_cascade_wn_proto_rawDescData)
	})
	return file_register_cascade_wn_proto_rawDescData
}

var file_register_cascade_wn_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_register_cascade_wn_proto_goTypes = []interface{}{
	(*SendSignedCascadeTicketRequest)(nil), // 0: walletnode.SendSignedCascadeTicketRequest
	(*UploadAssetRequest)(nil),             // 1: walletnode.UploadAssetRequest
	(*UploadAssetReply)(nil),               // 2: walletnode.UploadAssetReply
	(*UploadAssetRequest_MetaData)(nil),    // 3: walletnode.UploadAssetRequest.MetaData
	(*EncoderParameters)(nil),              // 4: walletnode.EncoderParameters
	(*SessionRequest)(nil),                 // 5: walletnode.SessionRequest
	(*AcceptedNodesRequest)(nil),           // 6: walletnode.AcceptedNodesRequest
	(*ConnectToRequest)(nil),               // 7: walletnode.ConnectToRequest
	(*MeshNodesRequest)(nil),               // 8: walletnode.MeshNodesRequest
	(*SendRegMetadataRequest)(nil),         // 9: walletnode.SendRegMetadataRequest
	(*SendActionActRequest)(nil),           // 10: walletnode.SendActionActRequest
	(*SessionReply)(nil),                   // 11: walletnode.SessionReply
	(*AcceptedNodesReply)(nil),             // 12: walletnode.AcceptedNodesReply
	(*ConnectToReply)(nil),                 // 13: walletnode.ConnectToReply
	(*MeshNodesReply)(nil),                 // 14: walletnode.MeshNodesReply
	(*SendRegMetadataReply)(nil),           // 15: walletnode.SendRegMetadataReply
	(*SendSignedActionTicketReply)(nil),    // 16: walletnode.SendSignedActionTicketReply
	(*SendActionActReply)(nil),             // 17: walletnode.SendActionActReply
}
var file_register_cascade_wn_proto_depIdxs = []int32{
	4,  // 0: walletnode.SendSignedCascadeTicketRequest.encode_parameters:type_name -> walletnode.EncoderParameters
	3,  // 1: walletnode.UploadAssetRequest.meta_data:type_name -> walletnode.UploadAssetRequest.MetaData
	5,  // 2: walletnode.RegisterCascade.Session:input_type -> walletnode.SessionRequest
	6,  // 3: walletnode.RegisterCascade.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	7,  // 4: walletnode.RegisterCascade.ConnectTo:input_type -> walletnode.ConnectToRequest
	8,  // 5: walletnode.RegisterCascade.MeshNodes:input_type -> walletnode.MeshNodesRequest
	9,  // 6: walletnode.RegisterCascade.SendRegMetadata:input_type -> walletnode.SendRegMetadataRequest
	1,  // 7: walletnode.RegisterCascade.UploadAsset:input_type -> walletnode.UploadAssetRequest
	0,  // 8: walletnode.RegisterCascade.SendSignedActionTicket:input_type -> walletnode.SendSignedCascadeTicketRequest
	10, // 9: walletnode.RegisterCascade.SendActionAct:input_type -> walletnode.SendActionActRequest
	11, // 10: walletnode.RegisterCascade.Session:output_type -> walletnode.SessionReply
	12, // 11: walletnode.RegisterCascade.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	13, // 12: walletnode.RegisterCascade.ConnectTo:output_type -> walletnode.ConnectToReply
	14, // 13: walletnode.RegisterCascade.MeshNodes:output_type -> walletnode.MeshNodesReply
	15, // 14: walletnode.RegisterCascade.SendRegMetadata:output_type -> walletnode.SendRegMetadataReply
	2,  // 15: walletnode.RegisterCascade.UploadAsset:output_type -> walletnode.UploadAssetReply
	16, // 16: walletnode.RegisterCascade.SendSignedActionTicket:output_type -> walletnode.SendSignedActionTicketReply
	17, // 17: walletnode.RegisterCascade.SendActionAct:output_type -> walletnode.SendActionActReply
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_register_cascade_wn_proto_init() }
func file_register_cascade_wn_proto_init() {
	if File_register_cascade_wn_proto != nil {
		return
	}
	file_common_wn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_register_cascade_wn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendSignedCascadeTicketRequest); i {
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
		file_register_cascade_wn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadAssetRequest); i {
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
		file_register_cascade_wn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadAssetReply); i {
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
		file_register_cascade_wn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadAssetRequest_MetaData); i {
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
	file_register_cascade_wn_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*UploadAssetRequest_AssetPiece)(nil),
		(*UploadAssetRequest_MetaData_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_register_cascade_wn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_register_cascade_wn_proto_goTypes,
		DependencyIndexes: file_register_cascade_wn_proto_depIdxs,
		MessageInfos:      file_register_cascade_wn_proto_msgTypes,
	}.Build()
	File_register_cascade_wn_proto = out.File
	file_register_cascade_wn_proto_rawDesc = nil
	file_register_cascade_wn_proto_goTypes = nil
	file_register_cascade_wn_proto_depIdxs = nil
}
