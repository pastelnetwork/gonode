// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: supernode/process_userdata_supernode.proto

package supernode

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

type MDLSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *MDLSessionRequest) Reset() {
	*x = MDLSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MDLSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MDLSessionRequest) ProtoMessage() {}

func (x *MDLSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MDLSessionRequest.ProtoReflect.Descriptor instead.
func (*MDLSessionRequest) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{0}
}

func (x *MDLSessionRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type MDLSessionReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
}

func (x *MDLSessionReply) Reset() {
	*x = MDLSessionReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MDLSessionReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MDLSessionReply) ProtoMessage() {}

func (x *MDLSessionReply) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MDLSessionReply.ProtoReflect.Descriptor instead.
func (*MDLSessionReply) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{1}
}

func (x *MDLSessionReply) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

type SuperNodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UserdataHash represents UserdataProcessRequest's hash value, to make sure
	// UserdataProcessRequest's integrity
	UserdataHash string `protobuf:"bytes,1,opt,name=userdata_hash,json=userdataHash,proto3" json:"userdata_hash,omitempty"`
	// UserdataResultHash represents UserdataReply's hash value, to make sure
	// walletnode's UserdataReply integrity
	UserdataResultHash string `protobuf:"bytes,2,opt,name=userdata_result_hash,json=userdataResultHash,proto3" json:"userdata_result_hash,omitempty"`
	// SuperNodeSignature is the digital signature created by supernode for the
	// [userdata_hash+userdata_result_hash]
	HashSignature string `protobuf:"bytes,3,opt,name=hash_signature,json=hashSignature,proto3" json:"hash_signature,omitempty"`
	// Supernode's pastelID of this supernode generate this SuperNodeRequest
	SupernodePastelID string `protobuf:"bytes,4,opt,name=supernode_pastelID,json=supernodePastelID,proto3" json:"supernode_pastelID,omitempty"`
	// Supernode's nodeID that init this SuperNodeRequest
	NodeID string `protobuf:"bytes,5,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *SuperNodeRequest) Reset() {
	*x = SuperNodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuperNodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuperNodeRequest) ProtoMessage() {}

func (x *SuperNodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuperNodeRequest.ProtoReflect.Descriptor instead.
func (*SuperNodeRequest) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{2}
}

func (x *SuperNodeRequest) GetUserdataHash() string {
	if x != nil {
		return x.UserdataHash
	}
	return ""
}

func (x *SuperNodeRequest) GetUserdataResultHash() string {
	if x != nil {
		return x.UserdataResultHash
	}
	return ""
}

func (x *SuperNodeRequest) GetHashSignature() string {
	if x != nil {
		return x.HashSignature
	}
	return ""
}

func (x *SuperNodeRequest) GetSupernodePastelID() string {
	if x != nil {
		return x.SupernodePastelID
	}
	return ""
}

func (x *SuperNodeRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type SuperNodeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Result of the request is success or not
	ResponseCode int32 `protobuf:"varint,1,opt,name=response_code,json=responseCode,proto3" json:"response_code,omitempty"`
	// The detail of why result is success/fail, depend on response_code
	Detail string `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
}

func (x *SuperNodeReply) Reset() {
	*x = SuperNodeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SuperNodeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SuperNodeReply) ProtoMessage() {}

func (x *SuperNodeReply) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SuperNodeReply.ProtoReflect.Descriptor instead.
func (*SuperNodeReply) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{3}
}

func (x *SuperNodeReply) GetResponseCode() int32 {
	if x != nil {
		return x.ResponseCode
	}
	return 0
}

func (x *SuperNodeReply) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

type UserdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Real name of the user
	RealName string `protobuf:"bytes,1,opt,name=real_name,json=realName,proto3" json:"real_name,omitempty"`
	// Facebook link of the user
	FacebookLink string `protobuf:"bytes,2,opt,name=facebook_link,json=facebookLink,proto3" json:"facebook_link,omitempty"`
	// Twitter link of the user
	TwitterLink string `protobuf:"bytes,3,opt,name=twitter_link,json=twitterLink,proto3" json:"twitter_link,omitempty"`
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency string `protobuf:"bytes,4,opt,name=native_currency,json=nativeCurrency,proto3" json:"native_currency,omitempty"`
	// Location of the user
	Location string `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	// Primary language of the user
	PrimaryLanguage string `protobuf:"bytes,6,opt,name=primary_language,json=primaryLanguage,proto3" json:"primary_language,omitempty"`
	// The categories of user's work
	Categories string `protobuf:"bytes,7,opt,name=categories,proto3" json:"categories,omitempty"`
	// Biography of the user
	Biography string `protobuf:"bytes,8,opt,name=biography,proto3" json:"biography,omitempty"`
	// Avatar image of the user
	AvatarImage *UserdataRequest_UserImageUpload `protobuf:"bytes,9,opt,name=avatar_image,json=avatarImage,proto3" json:"avatar_image,omitempty"`
	// Cover photo of the user
	CoverPhoto *UserdataRequest_UserImageUpload `protobuf:"bytes,10,opt,name=cover_photo,json=coverPhoto,proto3" json:"cover_photo,omitempty"`
	// Artist's PastelID
	ArtistPastelID string `protobuf:"bytes,11,opt,name=artist_pastelID,json=artistPastelID,proto3" json:"artist_pastelID,omitempty"`
	// Epoch Timestamp of the request
	Timestamp int64 `protobuf:"varint,12,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Previous block hash in the chain
	PreviousBlockHash string `protobuf:"bytes,13,opt,name=previous_block_hash,json=previousBlockHash,proto3" json:"previous_block_hash,omitempty"`
	// UserdataHash represents UserdataProcessRequest's hash value, to make sure
	// UserdataProcessRequest's integrity
	UserdataHash string `protobuf:"bytes,14,opt,name=userdata_hash,json=userdataHash,proto3" json:"userdata_hash,omitempty"`
	// Signature of the userdata_hash
	Signature string `protobuf:"bytes,15,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *UserdataRequest) Reset() {
	*x = UserdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserdataRequest) ProtoMessage() {}

func (x *UserdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserdataRequest.ProtoReflect.Descriptor instead.
func (*UserdataRequest) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{4}
}

func (x *UserdataRequest) GetRealName() string {
	if x != nil {
		return x.RealName
	}
	return ""
}

func (x *UserdataRequest) GetFacebookLink() string {
	if x != nil {
		return x.FacebookLink
	}
	return ""
}

func (x *UserdataRequest) GetTwitterLink() string {
	if x != nil {
		return x.TwitterLink
	}
	return ""
}

func (x *UserdataRequest) GetNativeCurrency() string {
	if x != nil {
		return x.NativeCurrency
	}
	return ""
}

func (x *UserdataRequest) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *UserdataRequest) GetPrimaryLanguage() string {
	if x != nil {
		return x.PrimaryLanguage
	}
	return ""
}

func (x *UserdataRequest) GetCategories() string {
	if x != nil {
		return x.Categories
	}
	return ""
}

func (x *UserdataRequest) GetBiography() string {
	if x != nil {
		return x.Biography
	}
	return ""
}

func (x *UserdataRequest) GetAvatarImage() *UserdataRequest_UserImageUpload {
	if x != nil {
		return x.AvatarImage
	}
	return nil
}

func (x *UserdataRequest) GetCoverPhoto() *UserdataRequest_UserImageUpload {
	if x != nil {
		return x.CoverPhoto
	}
	return nil
}

func (x *UserdataRequest) GetArtistPastelID() string {
	if x != nil {
		return x.ArtistPastelID
	}
	return ""
}

func (x *UserdataRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *UserdataRequest) GetPreviousBlockHash() string {
	if x != nil {
		return x.PreviousBlockHash
	}
	return ""
}

func (x *UserdataRequest) GetUserdataHash() string {
	if x != nil {
		return x.UserdataHash
	}
	return ""
}

func (x *UserdataRequest) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

type UserdataRequest_UserImageUpload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content  []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *UserdataRequest_UserImageUpload) Reset() {
	*x = UserdataRequest_UserImageUpload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supernode_process_userdata_supernode_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserdataRequest_UserImageUpload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserdataRequest_UserImageUpload) ProtoMessage() {}

func (x *UserdataRequest_UserImageUpload) ProtoReflect() protoreflect.Message {
	mi := &file_supernode_process_userdata_supernode_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserdataRequest_UserImageUpload.ProtoReflect.Descriptor instead.
func (*UserdataRequest_UserImageUpload) Descriptor() ([]byte, []int) {
	return file_supernode_process_userdata_supernode_proto_rawDescGZIP(), []int{4, 0}
}

func (x *UserdataRequest_UserImageUpload) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *UserdataRequest_UserImageUpload) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

var File_supernode_process_userdata_supernode_proto protoreflect.FileDescriptor

var file_supernode_process_userdata_supernode_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x75,
	0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x2b, 0x0a, 0x11, 0x4d, 0x44, 0x4c, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f,
	0x64, 0x65, 0x49, 0x44, 0x22, 0x29, 0x0a, 0x0f, 0x4d, 0x44, 0x4c, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x22,
	0xd7, 0x01, 0x0a, 0x10, 0x53, 0x75, 0x70, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x73, 0x65,
	0x72, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68, 0x12, 0x30, 0x0a, 0x14, 0x75, 0x73, 0x65,
	0x72, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x75, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x25, 0x0a, 0x0e, 0x68,
	0x61, 0x73, 0x68, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x68, 0x61, 0x73, 0x68, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x2d, 0x0a, 0x12, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11,
	0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49,
	0x44, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x4d, 0x0a, 0x0e, 0x53, 0x75, 0x70,
	0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x22, 0xc3, 0x05, 0x0a, 0x0f, 0x55, 0x73, 0x65,
	0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x72, 0x65, 0x61, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x72, 0x65, 0x61, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x61, 0x63,
	0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x21,
	0x0a, 0x0c, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x4c, 0x69, 0x6e,
	0x6b, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72,
	0x79, 0x5f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x12,
	0x4d, 0x0a, 0x0c, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x0b, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x4b,
	0x0a, 0x0b, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52,
	0x0a, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x27, 0x0a, 0x0f, 0x61,
	0x72, 0x74, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x50, 0x61, 0x73, 0x74,
	0x65, 0x6c, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x11, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x64,
	0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x1a, 0x47, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0xfa,
	0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x47, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e,
	0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x44, 0x4c, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x73, 0x75,
	0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4d, 0x44, 0x4c, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x4f, 0x0a, 0x15, 0x53,
	0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x54, 0x6f, 0x50, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x12, 0x1b, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x53, 0x75, 0x70, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x75,
	0x70, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4d, 0x0a, 0x14,
	0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x54, 0x6f, 0x4c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x75, 0x70,
	0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x31, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_supernode_process_userdata_supernode_proto_rawDescOnce sync.Once
	file_supernode_process_userdata_supernode_proto_rawDescData = file_supernode_process_userdata_supernode_proto_rawDesc
)

func file_supernode_process_userdata_supernode_proto_rawDescGZIP() []byte {
	file_supernode_process_userdata_supernode_proto_rawDescOnce.Do(func() {
		file_supernode_process_userdata_supernode_proto_rawDescData = protoimpl.X.CompressGZIP(file_supernode_process_userdata_supernode_proto_rawDescData)
	})
	return file_supernode_process_userdata_supernode_proto_rawDescData
}

var file_supernode_process_userdata_supernode_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_supernode_process_userdata_supernode_proto_goTypes = []interface{}{
	(*MDLSessionRequest)(nil),               // 0: supernode.MDLSessionRequest
	(*MDLSessionReply)(nil),                 // 1: supernode.MDLSessionReply
	(*SuperNodeRequest)(nil),                // 2: supernode.SuperNodeRequest
	(*SuperNodeReply)(nil),                  // 3: supernode.SuperNodeReply
	(*UserdataRequest)(nil),                 // 4: supernode.UserdataRequest
	(*UserdataRequest_UserImageUpload)(nil), // 5: supernode.UserdataRequest.UserImageUpload
}
var file_supernode_process_userdata_supernode_proto_depIdxs = []int32{
	5, // 0: supernode.UserdataRequest.avatar_image:type_name -> supernode.UserdataRequest.UserImageUpload
	5, // 1: supernode.UserdataRequest.cover_photo:type_name -> supernode.UserdataRequest.UserImageUpload
	0, // 2: supernode.ProcessUserdata.Session:input_type -> supernode.MDLSessionRequest
	2, // 3: supernode.ProcessUserdata.SendUserdataToPrimary:input_type -> supernode.SuperNodeRequest
	4, // 4: supernode.ProcessUserdata.SendUserdataToLeader:input_type -> supernode.UserdataRequest
	1, // 5: supernode.ProcessUserdata.Session:output_type -> supernode.MDLSessionReply
	3, // 6: supernode.ProcessUserdata.SendUserdataToPrimary:output_type -> supernode.SuperNodeReply
	3, // 7: supernode.ProcessUserdata.SendUserdataToLeader:output_type -> supernode.SuperNodeReply
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_supernode_process_userdata_supernode_proto_init() }
func file_supernode_process_userdata_supernode_proto_init() {
	if File_supernode_process_userdata_supernode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_supernode_process_userdata_supernode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MDLSessionRequest); i {
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
		file_supernode_process_userdata_supernode_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MDLSessionReply); i {
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
		file_supernode_process_userdata_supernode_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SuperNodeRequest); i {
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
		file_supernode_process_userdata_supernode_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SuperNodeReply); i {
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
		file_supernode_process_userdata_supernode_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserdataRequest); i {
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
		file_supernode_process_userdata_supernode_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserdataRequest_UserImageUpload); i {
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
			RawDescriptor: file_supernode_process_userdata_supernode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_supernode_process_userdata_supernode_proto_goTypes,
		DependencyIndexes: file_supernode_process_userdata_supernode_proto_depIdxs,
		MessageInfos:      file_supernode_process_userdata_supernode_proto_msgTypes,
	}.Build()
	File_supernode_process_userdata_supernode_proto = out.File
	file_supernode_process_userdata_supernode_proto_rawDesc = nil
	file_supernode_process_userdata_supernode_proto_goTypes = nil
	file_supernode_process_userdata_supernode_proto_depIdxs = nil
}
