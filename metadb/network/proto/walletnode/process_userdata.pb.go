// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.5.1
// source: process_userdata.proto

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

type SessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsPrimary bool `protobuf:"varint,1,opt,name=is_primary,json=isPrimary,proto3" json:"is_primary,omitempty"`
}

func (x *SessionRequest) Reset() {
	*x = SessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionRequest) ProtoMessage() {}

func (x *SessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionRequest.ProtoReflect.Descriptor instead.
func (*SessionRequest) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{0}
}

func (x *SessionRequest) GetIsPrimary() bool {
	if x != nil {
		return x.IsPrimary
	}
	return false
}

type SessionReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
}

func (x *SessionReply) Reset() {
	*x = SessionReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReply) ProtoMessage() {}

func (x *SessionReply) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReply.ProtoReflect.Descriptor instead.
func (*SessionReply) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{1}
}

func (x *SessionReply) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

type AcceptedNodesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AcceptedNodesRequest) Reset() {
	*x = AcceptedNodesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesRequest) ProtoMessage() {}

func (x *AcceptedNodesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesRequest.ProtoReflect.Descriptor instead.
func (*AcceptedNodesRequest) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{2}
}

type AcceptedNodesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peers []*AcceptedNodesReply_Peer `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (x *AcceptedNodesReply) Reset() {
	*x = AcceptedNodesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesReply) ProtoMessage() {}

func (x *AcceptedNodesReply) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesReply.ProtoReflect.Descriptor instead.
func (*AcceptedNodesReply) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{3}
}

func (x *AcceptedNodesReply) GetPeers() []*AcceptedNodesReply_Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

type ConnectToRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessID string `protobuf:"bytes,1,opt,name=sessID,proto3" json:"sessID,omitempty"`
	NodeID string `protobuf:"bytes,2,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *ConnectToRequest) Reset() {
	*x = ConnectToRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectToRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectToRequest) ProtoMessage() {}

func (x *ConnectToRequest) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectToRequest.ProtoReflect.Descriptor instead.
func (*ConnectToRequest) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{4}
}

func (x *ConnectToRequest) GetSessID() string {
	if x != nil {
		return x.SessID
	}
	return ""
}

func (x *ConnectToRequest) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type ConnectToReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectToReply) Reset() {
	*x = ConnectToReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectToReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectToReply) ProtoMessage() {}

func (x *ConnectToReply) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectToReply.ProtoReflect.Descriptor instead.
func (*ConnectToReply) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{5}
}

type UserdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Real name of the user
	Realname string `protobuf:"bytes,1,opt,name=realname,proto3" json:"realname,omitempty"`
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
	Timestamp int32 `protobuf:"varint,12,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Previous block hash in the chain
	PreviousBlockHash string `protobuf:"bytes,13,opt,name=previous_block_hash,json=previousBlockHash,proto3" json:"previous_block_hash,omitempty"`
	// UserdataHash represents UserdataProcessRequest's hash value, to make sure UserdataProcessRequest's integrity
	UserdataHash string `protobuf:"bytes,14,opt,name=userdata_hash,json=userdataHash,proto3" json:"userdata_hash,omitempty"`
	// Signature of the userdata_hash
	Signature string `protobuf:"bytes,15,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *UserdataRequest) Reset() {
	*x = UserdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserdataRequest) ProtoMessage() {}

func (x *UserdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[6]
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
	return file_process_userdata_proto_rawDescGZIP(), []int{6}
}

func (x *UserdataRequest) GetRealname() string {
	if x != nil {
		return x.Realname
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

func (x *UserdataRequest) GetTimestamp() int32 {
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

type UserdataReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Result of the request is success or not
	ResponseCode int32 `protobuf:"varint,1,opt,name=response_code,json=responseCode,proto3" json:"response_code,omitempty"`
	// The detail of why result is success/fail, depend on response_code
	Detail string `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
	// Error detail on realname
	Realname string `protobuf:"bytes,3,opt,name=realname,proto3" json:"realname,omitempty"`
	// Error detail on facebook_link
	FacebookLink string `protobuf:"bytes,4,opt,name=facebook_link,json=facebookLink,proto3" json:"facebook_link,omitempty"`
	// Error detail on twitter_link
	TwitterLink string `protobuf:"bytes,5,opt,name=twitter_link,json=twitterLink,proto3" json:"twitter_link,omitempty"`
	// Error detail on native_currency
	NativeCurrency string `protobuf:"bytes,6,opt,name=native_currency,json=nativeCurrency,proto3" json:"native_currency,omitempty"`
	// Error detail on location
	Location string `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
	// Error detail on primary_language
	PrimaryLanguage string `protobuf:"bytes,8,opt,name=primary_language,json=primaryLanguage,proto3" json:"primary_language,omitempty"`
	// Error detail on categories
	Categories string `protobuf:"bytes,9,opt,name=categories,proto3" json:"categories,omitempty"`
	// Error detail on biography
	Biography string `protobuf:"bytes,10,opt,name=biography,proto3" json:"biography,omitempty"`
	// Error detail on avatar
	AvatarImage string `protobuf:"bytes,11,opt,name=avatar_image,json=avatarImage,proto3" json:"avatar_image,omitempty"`
	// Error detail on cover photo
	CoverPhoto string `protobuf:"bytes,12,opt,name=cover_photo,json=coverPhoto,proto3" json:"cover_photo,omitempty"`
}

func (x *UserdataReply) Reset() {
	*x = UserdataReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserdataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserdataReply) ProtoMessage() {}

func (x *UserdataReply) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserdataReply.ProtoReflect.Descriptor instead.
func (*UserdataReply) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{7}
}

func (x *UserdataReply) GetResponseCode() int32 {
	if x != nil {
		return x.ResponseCode
	}
	return 0
}

func (x *UserdataReply) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *UserdataReply) GetRealname() string {
	if x != nil {
		return x.Realname
	}
	return ""
}

func (x *UserdataReply) GetFacebookLink() string {
	if x != nil {
		return x.FacebookLink
	}
	return ""
}

func (x *UserdataReply) GetTwitterLink() string {
	if x != nil {
		return x.TwitterLink
	}
	return ""
}

func (x *UserdataReply) GetNativeCurrency() string {
	if x != nil {
		return x.NativeCurrency
	}
	return ""
}

func (x *UserdataReply) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *UserdataReply) GetPrimaryLanguage() string {
	if x != nil {
		return x.PrimaryLanguage
	}
	return ""
}

func (x *UserdataReply) GetCategories() string {
	if x != nil {
		return x.Categories
	}
	return ""
}

func (x *UserdataReply) GetBiography() string {
	if x != nil {
		return x.Biography
	}
	return ""
}

func (x *UserdataReply) GetAvatarImage() string {
	if x != nil {
		return x.AvatarImage
	}
	return ""
}

func (x *UserdataReply) GetCoverPhoto() string {
	if x != nil {
		return x.CoverPhoto
	}
	return ""
}

type AcceptedNodesReply_Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *AcceptedNodesReply_Peer) Reset() {
	*x = AcceptedNodesReply_Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_userdata_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedNodesReply_Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedNodesReply_Peer) ProtoMessage() {}

func (x *AcceptedNodesReply_Peer) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedNodesReply_Peer.ProtoReflect.Descriptor instead.
func (*AcceptedNodesReply_Peer) Descriptor() ([]byte, []int) {
	return file_process_userdata_proto_rawDescGZIP(), []int{3, 0}
}

func (x *AcceptedNodesReply_Peer) GetNodeID() string {
	if x != nil {
		return x.NodeID
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
		mi := &file_process_userdata_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserdataRequest_UserImageUpload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserdataRequest_UserImageUpload) ProtoMessage() {}

func (x *UserdataRequest_UserImageUpload) ProtoReflect() protoreflect.Message {
	mi := &file_process_userdata_proto_msgTypes[9]
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
	return file_process_userdata_proto_rawDescGZIP(), []int{6, 0}
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

var File_process_userdata_proto protoreflect.FileDescriptor

var file_process_userdata_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x22, 0x2f, 0x0a, 0x0e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50, 0x72,
	0x69, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x26, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x73, 0x73, 0x49, 0x44, 0x22, 0x16, 0x0a,
	0x14, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6f, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65,
	0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x39, 0x0a, 0x05, 0x70,
	0x65, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64,
	0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52,
	0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x1a, 0x1e, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x42, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x73, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x73, 0x73,
	0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x10, 0x0a, 0x0e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0xc4, 0x05, 0x0a,
	0x0f, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x4c, 0x69, 0x6e,
	0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x5f, 0x6c, 0x69, 0x6e,
	0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72,
	0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x5f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70,
	0x68, 0x79, 0x12, 0x4e, 0x0a, 0x0c, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x0b, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74,
	0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x0a, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x50, 0x68, 0x6f, 0x74, 0x6f,
	0x12, 0x27, 0x0a, 0x0f, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x73, 0x74, 0x65,
	0x6c, 0x49, 0x44, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x72, 0x74, 0x69, 0x73,
	0x74, 0x50, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x72, 0x65, 0x76, 0x69,
	0x6f, 0x75, 0x73, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x64,
	0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x75, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1a, 0x47, 0x0a, 0x0f, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0xa2, 0x03, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x4c,
	0x69, 0x6e, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x69, 0x74, 0x74, 0x65, 0x72, 0x5f, 0x6c,
	0x69, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x77, 0x69, 0x74, 0x74,
	0x65, 0x72, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65,
	0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x70,
	0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x5f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x69, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61,
	0x70, 0x68, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x76, 0x61, 0x74,
	0x61, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x32, 0xba, 0x02, 0x0a, 0x0f, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x12, 0x43, 0x0a, 0x07,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x30,
	0x01, 0x12, 0x51, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64,
	0x65, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54,
	0x6f, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x48, 0x0a, 0x0c, 0x53,
	0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x28, 0x01, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x62, 0x2f,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_process_userdata_proto_rawDescOnce sync.Once
	file_process_userdata_proto_rawDescData = file_process_userdata_proto_rawDesc
)

func file_process_userdata_proto_rawDescGZIP() []byte {
	file_process_userdata_proto_rawDescOnce.Do(func() {
		file_process_userdata_proto_rawDescData = protoimpl.X.CompressGZIP(file_process_userdata_proto_rawDescData)
	})
	return file_process_userdata_proto_rawDescData
}

var file_process_userdata_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_process_userdata_proto_goTypes = []interface{}{
	(*SessionRequest)(nil),                  // 0: walletnode.SessionRequest
	(*SessionReply)(nil),                    // 1: walletnode.SessionReply
	(*AcceptedNodesRequest)(nil),            // 2: walletnode.AcceptedNodesRequest
	(*AcceptedNodesReply)(nil),              // 3: walletnode.AcceptedNodesReply
	(*ConnectToRequest)(nil),                // 4: walletnode.ConnectToRequest
	(*ConnectToReply)(nil),                  // 5: walletnode.ConnectToReply
	(*UserdataRequest)(nil),                 // 6: walletnode.UserdataRequest
	(*UserdataReply)(nil),                   // 7: walletnode.UserdataReply
	(*AcceptedNodesReply_Peer)(nil),         // 8: walletnode.AcceptedNodesReply.Peer
	(*UserdataRequest_UserImageUpload)(nil), // 9: walletnode.UserdataRequest.UserImageUpload
}
var file_process_userdata_proto_depIdxs = []int32{
	8, // 0: walletnode.AcceptedNodesReply.peers:type_name -> walletnode.AcceptedNodesReply.Peer
	9, // 1: walletnode.UserdataRequest.avatar_image:type_name -> walletnode.UserdataRequest.UserImageUpload
	9, // 2: walletnode.UserdataRequest.cover_photo:type_name -> walletnode.UserdataRequest.UserImageUpload
	0, // 3: walletnode.ProcessUserdata.Session:input_type -> walletnode.SessionRequest
	2, // 4: walletnode.ProcessUserdata.AcceptedNodes:input_type -> walletnode.AcceptedNodesRequest
	4, // 5: walletnode.ProcessUserdata.ConnectTo:input_type -> walletnode.ConnectToRequest
	6, // 6: walletnode.ProcessUserdata.SendUserdata:input_type -> walletnode.UserdataRequest
	1, // 7: walletnode.ProcessUserdata.Session:output_type -> walletnode.SessionReply
	3, // 8: walletnode.ProcessUserdata.AcceptedNodes:output_type -> walletnode.AcceptedNodesReply
	5, // 9: walletnode.ProcessUserdata.ConnectTo:output_type -> walletnode.ConnectToReply
	7, // 10: walletnode.ProcessUserdata.SendUserdata:output_type -> walletnode.UserdataReply
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_process_userdata_proto_init() }
func file_process_userdata_proto_init() {
	if File_process_userdata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_process_userdata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionRequest); i {
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
		file_process_userdata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionReply); i {
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
		file_process_userdata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesRequest); i {
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
		file_process_userdata_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesReply); i {
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
		file_process_userdata_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectToRequest); i {
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
		file_process_userdata_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectToReply); i {
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
		file_process_userdata_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_process_userdata_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserdataReply); i {
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
		file_process_userdata_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptedNodesReply_Peer); i {
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
		file_process_userdata_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_process_userdata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_process_userdata_proto_goTypes,
		DependencyIndexes: file_process_userdata_proto_depIdxs,
		MessageInfos:      file_process_userdata_proto_msgTypes,
	}.Build()
	File_process_userdata_proto = out.File
	file_process_userdata_proto_rawDesc = nil
	file_process_userdata_proto_goTypes = nil
	file_process_userdata_proto_depIdxs = nil
}
