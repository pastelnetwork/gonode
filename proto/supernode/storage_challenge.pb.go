// Copyright (c) 2021-2021 The Pastel Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: storage_challenge.proto

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

type StorageChallengeDataMessageType int32

const (
	StorageChallengeData_MessageType_UNKNOWN                                StorageChallengeDataMessageType = 0
	StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE     StorageChallengeDataMessageType = 1
	StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE     StorageChallengeDataMessageType = 2
	StorageChallengeData_MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE StorageChallengeDataMessageType = 3
)

// Enum value maps for StorageChallengeDataMessageType.
var (
	StorageChallengeDataMessageType_name = map[int32]string{
		0: "MessageType_UNKNOWN",
		1: "MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE",
		2: "MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE",
		3: "MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE",
	}
	StorageChallengeDataMessageType_value = map[string]int32{
		"MessageType_UNKNOWN":                                0,
		"MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE":     1,
		"MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE":     2,
		"MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE": 3,
	}
)

func (x StorageChallengeDataMessageType) Enum() *StorageChallengeDataMessageType {
	p := new(StorageChallengeDataMessageType)
	*p = x
	return p
}

func (x StorageChallengeDataMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StorageChallengeDataMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_storage_challenge_proto_enumTypes[0].Descriptor()
}

func (StorageChallengeDataMessageType) Type() protoreflect.EnumType {
	return &file_storage_challenge_proto_enumTypes[0]
}

func (x StorageChallengeDataMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StorageChallengeDataMessageType.Descriptor instead.
func (StorageChallengeDataMessageType) EnumDescriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{0, 0}
}

type StorageChallengeDataStatus int32

const (
	StorageChallengeData_Status_UNKNOWN                   StorageChallengeDataStatus = 0
	StorageChallengeData_Status_PENDING                   StorageChallengeDataStatus = 1
	StorageChallengeData_Status_RESPONDED                 StorageChallengeDataStatus = 2
	StorageChallengeData_Status_SUCCEEDED                 StorageChallengeDataStatus = 3
	StorageChallengeData_Status_FAILED_TIMEOUT            StorageChallengeDataStatus = 4
	StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE StorageChallengeDataStatus = 5
)

// Enum value maps for StorageChallengeDataStatus.
var (
	StorageChallengeDataStatus_name = map[int32]string{
		0: "Status_UNKNOWN",
		1: "Status_PENDING",
		2: "Status_RESPONDED",
		3: "Status_SUCCEEDED",
		4: "Status_FAILED_TIMEOUT",
		5: "Status_FAILED_INCORRECT_RESPONSE",
	}
	StorageChallengeDataStatus_value = map[string]int32{
		"Status_UNKNOWN":                   0,
		"Status_PENDING":                   1,
		"Status_RESPONDED":                 2,
		"Status_SUCCEEDED":                 3,
		"Status_FAILED_TIMEOUT":            4,
		"Status_FAILED_INCORRECT_RESPONSE": 5,
	}
)

func (x StorageChallengeDataStatus) Enum() *StorageChallengeDataStatus {
	p := new(StorageChallengeDataStatus)
	*p = x
	return p
}

func (x StorageChallengeDataStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StorageChallengeDataStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_storage_challenge_proto_enumTypes[1].Descriptor()
}

func (StorageChallengeDataStatus) Type() protoreflect.EnumType {
	return &file_storage_challenge_proto_enumTypes[1]
}

func (x StorageChallengeDataStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StorageChallengeDataStatus.Descriptor instead.
func (StorageChallengeDataStatus) EnumDescriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{0, 1}
}

type StorageChallengeData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId                    string                             `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	MessageType                  StorageChallengeDataMessageType    `protobuf:"varint,2,opt,name=message_type,json=messageType,proto3,enum=supernode.StorageChallengeDataMessageType" json:"message_type,omitempty"`
	ChallengeStatus              StorageChallengeDataStatus         `protobuf:"varint,3,opt,name=challenge_status,json=challengeStatus,proto3,enum=supernode.StorageChallengeDataStatus" json:"challenge_status,omitempty"`
	BlockNumChallengeSent        int32                              `protobuf:"varint,4,opt,name=block_num_challenge_sent,json=blockNumChallengeSent,proto3" json:"block_num_challenge_sent,omitempty"`
	BlockNumChallengeRespondedTo int32                              `protobuf:"varint,5,opt,name=block_num_challenge_responded_to,json=blockNumChallengeRespondedTo,proto3" json:"block_num_challenge_responded_to,omitempty"`
	BlockNumChallengeVerified    int32                              `protobuf:"varint,6,opt,name=block_num_challenge_verified,json=blockNumChallengeVerified,proto3" json:"block_num_challenge_verified,omitempty"`
	MerklerootWhenChallengeSent  string                             `protobuf:"bytes,7,opt,name=merkleroot_when_challenge_sent,json=merklerootWhenChallengeSent,proto3" json:"merkleroot_when_challenge_sent,omitempty"`
	ChallengingMasternodeId      string                             `protobuf:"bytes,8,opt,name=challenging_masternode_id,json=challengingMasternodeId,proto3" json:"challenging_masternode_id,omitempty"`
	RespondingMasternodeId       string                             `protobuf:"bytes,9,opt,name=responding_masternode_id,json=respondingMasternodeId,proto3" json:"responding_masternode_id,omitempty"`
	ChallengeFile                *StorageChallengeDataChallengeFile `protobuf:"bytes,10,opt,name=challenge_file,json=challengeFile,proto3" json:"challenge_file,omitempty"`
	ChallengeSliceCorrectHash    string                             `protobuf:"bytes,11,opt,name=challenge_slice_correct_hash,json=challengeSliceCorrectHash,proto3" json:"challenge_slice_correct_hash,omitempty"`
	ChallengeResponseHash        string                             `protobuf:"bytes,12,opt,name=challenge_response_hash,json=challengeResponseHash,proto3" json:"challenge_response_hash,omitempty"`
	ChallengeId                  string                             `protobuf:"bytes,13,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
}

func (x *StorageChallengeData) Reset() {
	*x = StorageChallengeData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageChallengeData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageChallengeData) ProtoMessage() {}

func (x *StorageChallengeData) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageChallengeData.ProtoReflect.Descriptor instead.
func (*StorageChallengeData) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{0}
}

func (x *StorageChallengeData) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *StorageChallengeData) GetMessageType() StorageChallengeDataMessageType {
	if x != nil {
		return x.MessageType
	}
	return StorageChallengeData_MessageType_UNKNOWN
}

func (x *StorageChallengeData) GetChallengeStatus() StorageChallengeDataStatus {
	if x != nil {
		return x.ChallengeStatus
	}
	return StorageChallengeData_Status_UNKNOWN
}

func (x *StorageChallengeData) GetBlockNumChallengeSent() int32 {
	if x != nil {
		return x.BlockNumChallengeSent
	}
	return 0
}

func (x *StorageChallengeData) GetBlockNumChallengeRespondedTo() int32 {
	if x != nil {
		return x.BlockNumChallengeRespondedTo
	}
	return 0
}

func (x *StorageChallengeData) GetBlockNumChallengeVerified() int32 {
	if x != nil {
		return x.BlockNumChallengeVerified
	}
	return 0
}

func (x *StorageChallengeData) GetMerklerootWhenChallengeSent() string {
	if x != nil {
		return x.MerklerootWhenChallengeSent
	}
	return ""
}

func (x *StorageChallengeData) GetChallengingMasternodeId() string {
	if x != nil {
		return x.ChallengingMasternodeId
	}
	return ""
}

func (x *StorageChallengeData) GetRespondingMasternodeId() string {
	if x != nil {
		return x.RespondingMasternodeId
	}
	return ""
}

func (x *StorageChallengeData) GetChallengeFile() *StorageChallengeDataChallengeFile {
	if x != nil {
		return x.ChallengeFile
	}
	return nil
}

func (x *StorageChallengeData) GetChallengeSliceCorrectHash() string {
	if x != nil {
		return x.ChallengeSliceCorrectHash
	}
	return ""
}

func (x *StorageChallengeData) GetChallengeResponseHash() string {
	if x != nil {
		return x.ChallengeResponseHash
	}
	return ""
}

func (x *StorageChallengeData) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

type ProcessStorageChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *StorageChallengeData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessStorageChallengeRequest) Reset() {
	*x = ProcessStorageChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessStorageChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessStorageChallengeRequest) ProtoMessage() {}

func (x *ProcessStorageChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessStorageChallengeRequest.ProtoReflect.Descriptor instead.
func (*ProcessStorageChallengeRequest) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessStorageChallengeRequest) GetData() *StorageChallengeData {
	if x != nil {
		return x.Data
	}
	return nil
}

type ProcessStorageChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *StorageChallengeData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessStorageChallengeReply) Reset() {
	*x = ProcessStorageChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessStorageChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessStorageChallengeReply) ProtoMessage() {}

func (x *ProcessStorageChallengeReply) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessStorageChallengeReply.ProtoReflect.Descriptor instead.
func (*ProcessStorageChallengeReply) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessStorageChallengeReply) GetData() *StorageChallengeData {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifyStorageChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *StorageChallengeData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifyStorageChallengeRequest) Reset() {
	*x = VerifyStorageChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyStorageChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyStorageChallengeRequest) ProtoMessage() {}

func (x *VerifyStorageChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyStorageChallengeRequest.ProtoReflect.Descriptor instead.
func (*VerifyStorageChallengeRequest) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{3}
}

func (x *VerifyStorageChallengeRequest) GetData() *StorageChallengeData {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifyStorageChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *StorageChallengeData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifyStorageChallengeReply) Reset() {
	*x = VerifyStorageChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyStorageChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyStorageChallengeReply) ProtoMessage() {}

func (x *VerifyStorageChallengeReply) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyStorageChallengeReply.ProtoReflect.Descriptor instead.
func (*VerifyStorageChallengeReply) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{4}
}

func (x *VerifyStorageChallengeReply) GetData() *StorageChallengeData {
	if x != nil {
		return x.Data
	}
	return nil
}

type StorageChallengeDataChallengeFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileHashToChallenge      string `protobuf:"bytes,1,opt,name=file_hash_to_challenge,json=fileHashToChallenge,proto3" json:"file_hash_to_challenge,omitempty"`
	ChallengeSliceStartIndex int64  `protobuf:"varint,2,opt,name=challenge_slice_start_index,json=challengeSliceStartIndex,proto3" json:"challenge_slice_start_index,omitempty"`
	ChallengeSliceEndIndex   int64  `protobuf:"varint,3,opt,name=challenge_slice_end_index,json=challengeSliceEndIndex,proto3" json:"challenge_slice_end_index,omitempty"`
}

func (x *StorageChallengeDataChallengeFile) Reset() {
	*x = StorageChallengeDataChallengeFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_challenge_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageChallengeDataChallengeFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageChallengeDataChallengeFile) ProtoMessage() {}

func (x *StorageChallengeDataChallengeFile) ProtoReflect() protoreflect.Message {
	mi := &file_storage_challenge_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageChallengeDataChallengeFile.ProtoReflect.Descriptor instead.
func (*StorageChallengeDataChallengeFile) Descriptor() ([]byte, []int) {
	return file_storage_challenge_proto_rawDescGZIP(), []int{0, 0}
}

func (x *StorageChallengeDataChallengeFile) GetFileHashToChallenge() string {
	if x != nil {
		return x.FileHashToChallenge
	}
	return ""
}

func (x *StorageChallengeDataChallengeFile) GetChallengeSliceStartIndex() int64 {
	if x != nil {
		return x.ChallengeSliceStartIndex
	}
	return 0
}

func (x *StorageChallengeDataChallengeFile) GetChallengeSliceEndIndex() int64 {
	if x != nil {
		return x.ChallengeSliceEndIndex
	}
	return 0
}

var File_storage_challenge_proto protoreflect.FileDescriptor

var file_storage_challenge_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x75, 0x70, 0x65, 0x72,
	0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf1, 0x0a, 0x0a, 0x14, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d,
	0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x4e, 0x0a,
	0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x51, 0x0a,
	0x10, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x0f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x37, 0x0a, 0x18, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x5f, 0x63, 0x68,
	0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x15, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x43, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x46, 0x0a, 0x20, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x74, 0x6f, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x1c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x64, 0x54,
	0x6f, 0x12, 0x3f, 0x0a, 0x1c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x5f, 0x63,
	0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x19, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75,
	0x6d, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x12, 0x43, 0x0a, 0x1e, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x72, 0x6f, 0x6f, 0x74,
	0x5f, 0x77, 0x68, 0x65, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f,
	0x73, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1b, 0x6d, 0x65, 0x72, 0x6b,
	0x6c, 0x65, 0x72, 0x6f, 0x6f, 0x74, 0x57, 0x68, 0x65, 0x6e, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x19, 0x63, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x63, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x49, 0x64, 0x12, 0x38, 0x0a, 0x18, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x5f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x54, 0x0a,
	0x0e, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x12, 0x3f, 0x0a, 0x1c, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x5f, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x19, 0x63, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x36, 0x0a, 0x17, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x49, 0x64, 0x1a,
	0xbe, 0x01, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x33, 0x0a, 0x16, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74,
	0x6f, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x13, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x54, 0x6f, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x3d, 0x0a, 0x1b, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x5f, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x18, 0x63, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x39, 0x0a, 0x19, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x5f, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x16, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x22, 0xc6, 0x01, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x17, 0x0a, 0x13, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x32, 0x0a, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x53, 0x54, 0x4f, 0x52, 0x41, 0x47, 0x45,
	0x5f, 0x43, 0x48, 0x41, 0x4c, 0x4c, 0x45, 0x4e, 0x47, 0x45, 0x5f, 0x49, 0x53, 0x53, 0x55, 0x41,
	0x4e, 0x43, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x01, 0x12, 0x32, 0x0a,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x53, 0x54, 0x4f,
	0x52, 0x41, 0x47, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4c, 0x4c, 0x45, 0x4e, 0x47, 0x45, 0x5f, 0x52,
	0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10,
	0x02, 0x12, 0x36, 0x0a, 0x32, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x5f, 0x53, 0x54, 0x4f, 0x52, 0x41, 0x47, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4c, 0x4c, 0x45, 0x4e,
	0x47, 0x45, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x03, 0x22, 0x9d, 0x01, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x44, 0x45, 0x44,
	0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x53, 0x55, 0x43,
	0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55,
	0x54, 0x10, 0x04, 0x12, 0x24, 0x0a, 0x20, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x5f, 0x49, 0x4e, 0x43, 0x4f, 0x52, 0x52, 0x45, 0x43, 0x54, 0x5f, 0x52,
	0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x05, 0x22, 0x55, 0x0a, 0x1e, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x53, 0x0a, 0x1c, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x33, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x54, 0x0a, 0x1d, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x52, 0x0a, 0x1b, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x33, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32,
	0xb0, 0x02, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x12, 0x41, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x19, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x28, 0x01, 0x30, 0x01, 0x12, 0x6d, 0x0a, 0x17, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x12, 0x29, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e,
	0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x6a, 0x0a, 0x16, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x12, 0x28, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67,
	0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x3b, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_challenge_proto_rawDescOnce sync.Once
	file_storage_challenge_proto_rawDescData = file_storage_challenge_proto_rawDesc
)

func file_storage_challenge_proto_rawDescGZIP() []byte {
	file_storage_challenge_proto_rawDescOnce.Do(func() {
		file_storage_challenge_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_challenge_proto_rawDescData)
	})
	return file_storage_challenge_proto_rawDescData
}

var file_storage_challenge_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_storage_challenge_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_storage_challenge_proto_goTypes = []interface{}{
	(StorageChallengeDataMessageType)(0),      // 0: supernode.StorageChallengeData.messageType
	(StorageChallengeDataStatus)(0),           // 1: supernode.StorageChallengeData.status
	(*StorageChallengeData)(nil),              // 2: supernode.StorageChallengeData
	(*ProcessStorageChallengeRequest)(nil),    // 3: supernode.ProcessStorageChallengeRequest
	(*ProcessStorageChallengeReply)(nil),      // 4: supernode.ProcessStorageChallengeReply
	(*VerifyStorageChallengeRequest)(nil),     // 5: supernode.VerifyStorageChallengeRequest
	(*VerifyStorageChallengeReply)(nil),       // 6: supernode.VerifyStorageChallengeReply
	(*StorageChallengeDataChallengeFile)(nil), // 7: supernode.StorageChallengeData.challengeFile
	(*SessionRequest)(nil),                    // 8: supernode.SessionRequest
	(*SessionReply)(nil),                      // 9: supernode.SessionReply
}
var file_storage_challenge_proto_depIdxs = []int32{
	0,  // 0: supernode.StorageChallengeData.message_type:type_name -> supernode.StorageChallengeData.messageType
	1,  // 1: supernode.StorageChallengeData.challenge_status:type_name -> supernode.StorageChallengeData.status
	7,  // 2: supernode.StorageChallengeData.challenge_file:type_name -> supernode.StorageChallengeData.challengeFile
	2,  // 3: supernode.ProcessStorageChallengeRequest.data:type_name -> supernode.StorageChallengeData
	2,  // 4: supernode.ProcessStorageChallengeReply.data:type_name -> supernode.StorageChallengeData
	2,  // 5: supernode.VerifyStorageChallengeRequest.data:type_name -> supernode.StorageChallengeData
	2,  // 6: supernode.VerifyStorageChallengeReply.data:type_name -> supernode.StorageChallengeData
	8,  // 7: supernode.StorageChallenge.Session:input_type -> supernode.SessionRequest
	3,  // 8: supernode.StorageChallenge.ProcessStorageChallenge:input_type -> supernode.ProcessStorageChallengeRequest
	5,  // 9: supernode.StorageChallenge.VerifyStorageChallenge:input_type -> supernode.VerifyStorageChallengeRequest
	9,  // 10: supernode.StorageChallenge.Session:output_type -> supernode.SessionReply
	4,  // 11: supernode.StorageChallenge.ProcessStorageChallenge:output_type -> supernode.ProcessStorageChallengeReply
	6,  // 12: supernode.StorageChallenge.VerifyStorageChallenge:output_type -> supernode.VerifyStorageChallengeReply
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_storage_challenge_proto_init() }
func file_storage_challenge_proto_init() {
	if File_storage_challenge_proto != nil {
		return
	}
	file_common_sn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_storage_challenge_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageChallengeData); i {
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
		file_storage_challenge_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessStorageChallengeRequest); i {
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
		file_storage_challenge_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessStorageChallengeReply); i {
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
		file_storage_challenge_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyStorageChallengeRequest); i {
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
		file_storage_challenge_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyStorageChallengeReply); i {
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
		file_storage_challenge_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageChallengeDataChallengeFile); i {
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
			RawDescriptor: file_storage_challenge_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storage_challenge_proto_goTypes,
		DependencyIndexes: file_storage_challenge_proto_depIdxs,
		EnumInfos:         file_storage_challenge_proto_enumTypes,
		MessageInfos:      file_storage_challenge_proto_msgTypes,
	}.Build()
	File_storage_challenge_proto = out.File
	file_storage_challenge_proto_rawDesc = nil
	file_storage_challenge_proto_goTypes = nil
	file_storage_challenge_proto_depIdxs = nil
}
