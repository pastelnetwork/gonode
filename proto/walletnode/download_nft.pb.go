// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: download_nft.proto

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

type DownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid      string `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Timestamp string `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature string `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	Ttxid     string `protobuf:"bytes,4,opt,name=ttxid,proto3" json:"ttxid,omitempty"`
	Ttype     string `protobuf:"bytes,5,opt,name=ttype,proto3" json:"ttype,omitempty"`
	SendHash  bool   `protobuf:"varint,6,opt,name=send_hash,json=sendHash,proto3" json:"send_hash,omitempty"`
}

func (x *DownloadRequest) Reset() {
	*x = DownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadRequest) ProtoMessage() {}

func (x *DownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadRequest.ProtoReflect.Descriptor instead.
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{0}
}

func (x *DownloadRequest) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

func (x *DownloadRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *DownloadRequest) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *DownloadRequest) GetTtxid() string {
	if x != nil {
		return x.Ttxid
	}
	return ""
}

func (x *DownloadRequest) GetTtype() string {
	if x != nil {
		return x.Ttype
	}
	return ""
}

func (x *DownloadRequest) GetSendHash() bool {
	if x != nil {
		return x.SendHash
	}
	return false
}

type DownloadReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File []byte `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *DownloadReply) Reset() {
	*x = DownloadReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadReply) ProtoMessage() {}

func (x *DownloadReply) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadReply.ProtoReflect.Descriptor instead.
func (*DownloadReply) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{1}
}

func (x *DownloadReply) GetFile() []byte {
	if x != nil {
		return x.File
	}
	return nil
}

type DownloadThumbnailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid     string `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Numnails int32  `protobuf:"varint,2,opt,name=numnails,proto3" json:"numnails,omitempty"`
}

func (x *DownloadThumbnailRequest) Reset() {
	*x = DownloadThumbnailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadThumbnailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadThumbnailRequest) ProtoMessage() {}

func (x *DownloadThumbnailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadThumbnailRequest.ProtoReflect.Descriptor instead.
func (*DownloadThumbnailRequest) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{2}
}

func (x *DownloadThumbnailRequest) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

func (x *DownloadThumbnailRequest) GetNumnails() int32 {
	if x != nil {
		return x.Numnails
	}
	return 0
}

type DownloadThumbnailReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Thumbnailone []byte `protobuf:"bytes,1,opt,name=thumbnailone,proto3" json:"thumbnailone,omitempty"`
	Thumbnailtwo []byte `protobuf:"bytes,2,opt,name=thumbnailtwo,proto3" json:"thumbnailtwo,omitempty"`
}

func (x *DownloadThumbnailReply) Reset() {
	*x = DownloadThumbnailReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadThumbnailReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadThumbnailReply) ProtoMessage() {}

func (x *DownloadThumbnailReply) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadThumbnailReply.ProtoReflect.Descriptor instead.
func (*DownloadThumbnailReply) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadThumbnailReply) GetThumbnailone() []byte {
	if x != nil {
		return x.Thumbnailone
	}
	return nil
}

func (x *DownloadThumbnailReply) GetThumbnailtwo() []byte {
	if x != nil {
		return x.Thumbnailtwo
	}
	return nil
}

type DownloadDDAndFingerprintsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid string `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
}

func (x *DownloadDDAndFingerprintsRequest) Reset() {
	*x = DownloadDDAndFingerprintsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadDDAndFingerprintsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadDDAndFingerprintsRequest) ProtoMessage() {}

func (x *DownloadDDAndFingerprintsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadDDAndFingerprintsRequest.ProtoReflect.Descriptor instead.
func (*DownloadDDAndFingerprintsRequest) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{4}
}

func (x *DownloadDDAndFingerprintsRequest) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

type DownloadDDAndFingerprintsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File []byte `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *DownloadDDAndFingerprintsReply) Reset() {
	*x = DownloadDDAndFingerprintsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_download_nft_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadDDAndFingerprintsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadDDAndFingerprintsReply) ProtoMessage() {}

func (x *DownloadDDAndFingerprintsReply) ProtoReflect() protoreflect.Message {
	mi := &file_download_nft_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadDDAndFingerprintsReply.ProtoReflect.Descriptor instead.
func (*DownloadDDAndFingerprintsReply) Descriptor() ([]byte, []int) {
	return file_download_nft_proto_rawDescGZIP(), []int{5}
}

func (x *DownloadDDAndFingerprintsReply) GetFile() []byte {
	if x != nil {
		return x.File
	}
	return nil
}

var File_download_nft_proto protoreflect.FileDescriptor

var file_download_nft_proto_rawDesc = []byte{
	0x0a, 0x12, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x66, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x77, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xaa, 0x01, 0x0a, 0x0f, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x74, 0x78, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x74, 0x78, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x48, 0x61, 0x73, 0x68, 0x22, 0x23,
	0x0a, 0x0d, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x22, 0x4a, 0x0a, 0x18, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x54,
	0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x78, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x6e, 0x61, 0x69, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x75, 0x6d, 0x6e, 0x61, 0x69, 0x6c, 0x73, 0x22,
	0x60, 0x0a, 0x16, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x68, 0x75, 0x6d, 0x62,
	0x6e, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x68, 0x75,
	0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x6f, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0c, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x6f, 0x6e, 0x65, 0x12, 0x22, 0x0a,
	0x0c, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x74, 0x77, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0c, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x74, 0x77,
	0x6f, 0x22, 0x36, 0x0a, 0x20, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x44, 0x41,
	0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x22, 0x34, 0x0a, 0x1e, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x44, 0x41, 0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x32,
	0xf0, 0x02, 0x0a, 0x0b, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x66, 0x74, 0x12,
	0x44, 0x0a, 0x08, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1b, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x30, 0x01, 0x12, 0x5d, 0x0a, 0x11, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x54, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x12, 0x24, 0x2e, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x54, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x22, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x75, 0x0a, 0x19, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x44, 0x44, 0x41, 0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74,
	0x73, 0x12, 0x2c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x44, 0x41, 0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x44, 0x41, 0x6e, 0x64, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x4d, 0x4e, 0x73, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67,
	0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_download_nft_proto_rawDescOnce sync.Once
	file_download_nft_proto_rawDescData = file_download_nft_proto_rawDesc
)

func file_download_nft_proto_rawDescGZIP() []byte {
	file_download_nft_proto_rawDescOnce.Do(func() {
		file_download_nft_proto_rawDescData = protoimpl.X.CompressGZIP(file_download_nft_proto_rawDescData)
	})
	return file_download_nft_proto_rawDescData
}

var file_download_nft_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_download_nft_proto_goTypes = []interface{}{
	(*DownloadRequest)(nil),                  // 0: walletnode.DownloadRequest
	(*DownloadReply)(nil),                    // 1: walletnode.DownloadReply
	(*DownloadThumbnailRequest)(nil),         // 2: walletnode.DownloadThumbnailRequest
	(*DownloadThumbnailReply)(nil),           // 3: walletnode.DownloadThumbnailReply
	(*DownloadDDAndFingerprintsRequest)(nil), // 4: walletnode.DownloadDDAndFingerprintsRequest
	(*DownloadDDAndFingerprintsReply)(nil),   // 5: walletnode.DownloadDDAndFingerprintsReply
	(*GetTopMNsRequest)(nil),                 // 6: walletnode.GetTopMNsRequest
	(*GetTopMNsReply)(nil),                   // 7: walletnode.GetTopMNsReply
}
var file_download_nft_proto_depIdxs = []int32{
	0, // 0: walletnode.DownloadNft.Download:input_type -> walletnode.DownloadRequest
	2, // 1: walletnode.DownloadNft.DownloadThumbnail:input_type -> walletnode.DownloadThumbnailRequest
	4, // 2: walletnode.DownloadNft.DownloadDDAndFingerprints:input_type -> walletnode.DownloadDDAndFingerprintsRequest
	6, // 3: walletnode.DownloadNft.GetTopMNs:input_type -> walletnode.GetTopMNsRequest
	1, // 4: walletnode.DownloadNft.Download:output_type -> walletnode.DownloadReply
	3, // 5: walletnode.DownloadNft.DownloadThumbnail:output_type -> walletnode.DownloadThumbnailReply
	5, // 6: walletnode.DownloadNft.DownloadDDAndFingerprints:output_type -> walletnode.DownloadDDAndFingerprintsReply
	7, // 7: walletnode.DownloadNft.GetTopMNs:output_type -> walletnode.GetTopMNsReply
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_download_nft_proto_init() }
func file_download_nft_proto_init() {
	if File_download_nft_proto != nil {
		return
	}
	file_common_wn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_download_nft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadRequest); i {
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
		file_download_nft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadReply); i {
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
		file_download_nft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadThumbnailRequest); i {
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
		file_download_nft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadThumbnailReply); i {
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
		file_download_nft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadDDAndFingerprintsRequest); i {
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
		file_download_nft_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadDDAndFingerprintsReply); i {
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
			RawDescriptor: file_download_nft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_download_nft_proto_goTypes,
		DependencyIndexes: file_download_nft_proto_depIdxs,
		MessageInfos:      file_download_nft_proto_msgTypes,
	}.Build()
	File_download_nft_proto = out.File
	file_download_nft_proto_rawDesc = nil
	file_download_nft_proto_goTypes = nil
	file_download_nft_proto_depIdxs = nil
}
