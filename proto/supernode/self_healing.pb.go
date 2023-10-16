// Copyright (c) 2021-2021 The Pastel Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: self_healing.proto

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

type SelfHealingDataMessageType int32

const (
	SelfHealingData_MessageType_UNKNOWN                           SelfHealingDataMessageType = 0
	SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE     SelfHealingDataMessageType = 1
	SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE     SelfHealingDataMessageType = 2
	SelfHealingData_MessageType_SELF_HEALING_VERIFICATION_MESSAGE SelfHealingDataMessageType = 3
)

// Enum value maps for SelfHealingDataMessageType.
var (
	SelfHealingDataMessageType_name = map[int32]string{
		0: "MessageType_UNKNOWN",
		1: "MessageType_SELF_HEALING_ISSUANCE_MESSAGE",
		2: "MessageType_SELF_HEALING_RESPONSE_MESSAGE",
		3: "MessageType_SELF_HEALING_VERIFICATION_MESSAGE",
	}
	SelfHealingDataMessageType_value = map[string]int32{
		"MessageType_UNKNOWN":                           0,
		"MessageType_SELF_HEALING_ISSUANCE_MESSAGE":     1,
		"MessageType_SELF_HEALING_RESPONSE_MESSAGE":     2,
		"MessageType_SELF_HEALING_VERIFICATION_MESSAGE": 3,
	}
)

func (x SelfHealingDataMessageType) Enum() *SelfHealingDataMessageType {
	p := new(SelfHealingDataMessageType)
	*p = x
	return p
}

func (x SelfHealingDataMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SelfHealingDataMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_self_healing_proto_enumTypes[0].Descriptor()
}

func (SelfHealingDataMessageType) Type() protoreflect.EnumType {
	return &file_self_healing_proto_enumTypes[0]
}

func (x SelfHealingDataMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SelfHealingDataMessageType.Descriptor instead.
func (SelfHealingDataMessageType) EnumDescriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{0, 0}
}

type SelfHealingDataStatus int32

const (
	SelfHealingData_Status_UNKNOWN                   SelfHealingDataStatus = 0
	SelfHealingData_Status_PENDING                   SelfHealingDataStatus = 1
	SelfHealingData_Status_RESPONDED                 SelfHealingDataStatus = 2
	SelfHealingData_Status_SUCCEEDED                 SelfHealingDataStatus = 3
	SelfHealingData_Status_FAILED_TIMEOUT            SelfHealingDataStatus = 4
	SelfHealingData_Status_FAILED_INCORRECT_RESPONSE SelfHealingDataStatus = 5
)

// Enum value maps for SelfHealingDataStatus.
var (
	SelfHealingDataStatus_name = map[int32]string{
		0: "Status_UNKNOWN",
		1: "Status_PENDING",
		2: "Status_RESPONDED",
		3: "Status_SUCCEEDED",
		4: "Status_FAILED_TIMEOUT",
		5: "Status_FAILED_INCORRECT_RESPONSE",
	}
	SelfHealingDataStatus_value = map[string]int32{
		"Status_UNKNOWN":                   0,
		"Status_PENDING":                   1,
		"Status_RESPONDED":                 2,
		"Status_SUCCEEDED":                 3,
		"Status_FAILED_TIMEOUT":            4,
		"Status_FAILED_INCORRECT_RESPONSE": 5,
	}
)

func (x SelfHealingDataStatus) Enum() *SelfHealingDataStatus {
	p := new(SelfHealingDataStatus)
	*p = x
	return p
}

func (x SelfHealingDataStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SelfHealingDataStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_self_healing_proto_enumTypes[1].Descriptor()
}

func (SelfHealingDataStatus) Type() protoreflect.EnumType {
	return &file_self_healing_proto_enumTypes[1]
}

func (x SelfHealingDataStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SelfHealingDataStatus.Descriptor instead.
func (SelfHealingDataStatus) EnumDescriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{0, 1}
}

type SelfHealingData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId                   string                        `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	MessageType                 SelfHealingDataMessageType    `protobuf:"varint,2,opt,name=message_type,json=messageType,proto3,enum=supernode.SelfHealingDataMessageType" json:"message_type,omitempty"`
	ChallengeStatus             SelfHealingDataStatus         `protobuf:"varint,3,opt,name=challenge_status,json=challengeStatus,proto3,enum=supernode.SelfHealingDataStatus" json:"challenge_status,omitempty"`
	MerklerootWhenChallengeSent string                        `protobuf:"bytes,4,opt,name=merkleroot_when_challenge_sent,json=merklerootWhenChallengeSent,proto3" json:"merkleroot_when_challenge_sent,omitempty"`
	ChallengingMasternodeId     string                        `protobuf:"bytes,5,opt,name=challenging_masternode_id,json=challengingMasternodeId,proto3" json:"challenging_masternode_id,omitempty"`
	RespondingMasternodeId      string                        `protobuf:"bytes,6,opt,name=responding_masternode_id,json=respondingMasternodeId,proto3" json:"responding_masternode_id,omitempty"`
	ReconstructedFileHash       []byte                        `protobuf:"bytes,7,opt,name=reconstructed_file_hash,json=reconstructedFileHash,proto3" json:"reconstructed_file_hash,omitempty"`
	ChallengeFile               *SelfHealingDataChallengeFile `protobuf:"bytes,8,opt,name=challenge_file,json=challengeFile,proto3" json:"challenge_file,omitempty"`
	ChallengeId                 string                        `protobuf:"bytes,9,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	RegTicketId                 string                        `protobuf:"bytes,10,opt,name=reg_ticket_id,json=regTicketId,proto3" json:"reg_ticket_id,omitempty"`
	ActionTicketId              string                        `protobuf:"bytes,11,opt,name=action_ticket_id,json=actionTicketId,proto3" json:"action_ticket_id,omitempty"`
	IsSenseTicket               bool                          `protobuf:"varint,12,opt,name=is_sense_ticket,json=isSenseTicket,proto3" json:"is_sense_ticket,omitempty"`
	SenseFileIds                []string                      `protobuf:"bytes,13,rep,name=senseFileIds,proto3" json:"senseFileIds,omitempty"`
}

func (x *SelfHealingData) Reset() {
	*x = SelfHealingData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelfHealingData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelfHealingData) ProtoMessage() {}

func (x *SelfHealingData) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelfHealingData.ProtoReflect.Descriptor instead.
func (*SelfHealingData) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{0}
}

func (x *SelfHealingData) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *SelfHealingData) GetMessageType() SelfHealingDataMessageType {
	if x != nil {
		return x.MessageType
	}
	return SelfHealingData_MessageType_UNKNOWN
}

func (x *SelfHealingData) GetChallengeStatus() SelfHealingDataStatus {
	if x != nil {
		return x.ChallengeStatus
	}
	return SelfHealingData_Status_UNKNOWN
}

func (x *SelfHealingData) GetMerklerootWhenChallengeSent() string {
	if x != nil {
		return x.MerklerootWhenChallengeSent
	}
	return ""
}

func (x *SelfHealingData) GetChallengingMasternodeId() string {
	if x != nil {
		return x.ChallengingMasternodeId
	}
	return ""
}

func (x *SelfHealingData) GetRespondingMasternodeId() string {
	if x != nil {
		return x.RespondingMasternodeId
	}
	return ""
}

func (x *SelfHealingData) GetReconstructedFileHash() []byte {
	if x != nil {
		return x.ReconstructedFileHash
	}
	return nil
}

func (x *SelfHealingData) GetChallengeFile() *SelfHealingDataChallengeFile {
	if x != nil {
		return x.ChallengeFile
	}
	return nil
}

func (x *SelfHealingData) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *SelfHealingData) GetRegTicketId() string {
	if x != nil {
		return x.RegTicketId
	}
	return ""
}

func (x *SelfHealingData) GetActionTicketId() string {
	if x != nil {
		return x.ActionTicketId
	}
	return ""
}

func (x *SelfHealingData) GetIsSenseTicket() bool {
	if x != nil {
		return x.IsSenseTicket
	}
	return false
}

func (x *SelfHealingData) GetSenseFileIds() []string {
	if x != nil {
		return x.SenseFileIds
	}
	return nil
}

type ProcessSelfHealingChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessSelfHealingChallengeRequest) Reset() {
	*x = ProcessSelfHealingChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessSelfHealingChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessSelfHealingChallengeRequest) ProtoMessage() {}

func (x *ProcessSelfHealingChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessSelfHealingChallengeRequest.ProtoReflect.Descriptor instead.
func (*ProcessSelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessSelfHealingChallengeRequest) GetData() *SelfHealingData {
	if x != nil {
		return x.Data
	}
	return nil
}

type ProcessSelfHealingChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessSelfHealingChallengeReply) Reset() {
	*x = ProcessSelfHealingChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessSelfHealingChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessSelfHealingChallengeReply) ProtoMessage() {}

func (x *ProcessSelfHealingChallengeReply) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessSelfHealingChallengeReply.ProtoReflect.Descriptor instead.
func (*ProcessSelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessSelfHealingChallengeReply) GetData() *SelfHealingData {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifySelfHealingChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifySelfHealingChallengeRequest) Reset() {
	*x = VerifySelfHealingChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifySelfHealingChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifySelfHealingChallengeRequest) ProtoMessage() {}

func (x *VerifySelfHealingChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifySelfHealingChallengeRequest.ProtoReflect.Descriptor instead.
func (*VerifySelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{3}
}

func (x *VerifySelfHealingChallengeRequest) GetData() *SelfHealingData {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifySelfHealingChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifySelfHealingChallengeReply) Reset() {
	*x = VerifySelfHealingChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifySelfHealingChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifySelfHealingChallengeReply) ProtoMessage() {}

func (x *VerifySelfHealingChallengeReply) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifySelfHealingChallengeReply.ProtoReflect.Descriptor instead.
func (*VerifySelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{4}
}

func (x *VerifySelfHealingChallengeReply) GetData() *SelfHealingData {
	if x != nil {
		return x.Data
	}
	return nil
}

type SelfHealingDataChallengeFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileHashToChallenge string `protobuf:"bytes,1,opt,name=file_hash_to_challenge,json=fileHashToChallenge,proto3" json:"file_hash_to_challenge,omitempty"`
}

func (x *SelfHealingDataChallengeFile) Reset() {
	*x = SelfHealingDataChallengeFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelfHealingDataChallengeFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelfHealingDataChallengeFile) ProtoMessage() {}

func (x *SelfHealingDataChallengeFile) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelfHealingDataChallengeFile.ProtoReflect.Descriptor instead.
func (*SelfHealingDataChallengeFile) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SelfHealingDataChallengeFile) GetFileHashToChallenge() string {
	if x != nil {
		return x.FileHashToChallenge
	}
	return ""
}

var File_self_healing_proto protoreflect.FileDescriptor

var file_self_healing_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x65, 0x6c, 0x66, 0x5f, 0x68, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x1a,
	0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xea, 0x08, 0x0a, 0x0f, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x64, 0x12, 0x49, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e,
	0x67, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x4c,
	0x0a, 0x10, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0f, 0x63, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x43, 0x0a, 0x1e,
	0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x77, 0x68, 0x65, 0x6e, 0x5f,
	0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x1b, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x72, 0x6f, 0x6f, 0x74,
	0x57, 0x68, 0x65, 0x6e, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x6e,
	0x74, 0x12, 0x3a, 0x0a, 0x19, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x67,
	0x5f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x67, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x38, 0x0a,
	0x18, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x16, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x63, 0x6f, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x65, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x4f, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x44,
	0x61, 0x74, 0x61, 0x2e, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x72, 0x65, 0x67, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x67, 0x54,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x49,
	0x64, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x73, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x53, 0x65,
	0x6e, 0x73, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x65, 0x6e,
	0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x65, 0x6e, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x1a, 0x44, 0x0a,
	0x0d, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x33,
	0x0a, 0x16, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x5f, 0x63,
	0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13,
	0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x54, 0x6f, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x22, 0xb7, 0x01, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x2d, 0x0a, 0x29,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x53, 0x45, 0x4c, 0x46,
	0x5f, 0x48, 0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x49, 0x53, 0x53, 0x55, 0x41, 0x4e, 0x43,
	0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x01, 0x12, 0x2d, 0x0a, 0x29, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f,
	0x48, 0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45,
	0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x02, 0x12, 0x31, 0x0a, 0x2d, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x48,
	0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x03, 0x22, 0x9d, 0x01,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01,
	0x12, 0x14, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f,
	0x4e, 0x44, 0x45, 0x44, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f, 0x54, 0x49,
	0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x04, 0x12, 0x24, 0x0a, 0x20, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f, 0x49, 0x4e, 0x43, 0x4f, 0x52, 0x52, 0x45,
	0x43, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x05, 0x22, 0x54, 0x0a,
	0x22, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c,
	0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65,
	0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x52, 0x0a, 0x20, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65,
	0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x53, 0x0a, 0x21, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69,
	0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x51, 0x0a, 0x1f,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e,
	0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65,
	0x61, 0x6c, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32,
	0xc3, 0x02, 0x0a, 0x0b, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x12,
	0x41, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01,
	0x30, 0x01, 0x12, 0x79, 0x0a, 0x1b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c,
	0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x12, 0x2d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67,
	0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2b, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43,
	0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x76, 0x0a,
	0x1a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69,
	0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x2c, 0x2e, 0x73, 0x75,
	0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65,
	0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66,
	0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73,
	0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x3b, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f,
	0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_self_healing_proto_rawDescOnce sync.Once
	file_self_healing_proto_rawDescData = file_self_healing_proto_rawDesc
)

func file_self_healing_proto_rawDescGZIP() []byte {
	file_self_healing_proto_rawDescOnce.Do(func() {
		file_self_healing_proto_rawDescData = protoimpl.X.CompressGZIP(file_self_healing_proto_rawDescData)
	})
	return file_self_healing_proto_rawDescData
}

var file_self_healing_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_self_healing_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_self_healing_proto_goTypes = []interface{}{
	(SelfHealingDataMessageType)(0),            // 0: supernode.SelfHealingData.messageType
	(SelfHealingDataStatus)(0),                 // 1: supernode.SelfHealingData.status
	(*SelfHealingData)(nil),                    // 2: supernode.SelfHealingData
	(*ProcessSelfHealingChallengeRequest)(nil), // 3: supernode.ProcessSelfHealingChallengeRequest
	(*ProcessSelfHealingChallengeReply)(nil),   // 4: supernode.ProcessSelfHealingChallengeReply
	(*VerifySelfHealingChallengeRequest)(nil),  // 5: supernode.VerifySelfHealingChallengeRequest
	(*VerifySelfHealingChallengeReply)(nil),    // 6: supernode.VerifySelfHealingChallengeReply
	(*SelfHealingDataChallengeFile)(nil),       // 7: supernode.SelfHealingData.challengeFile
	(*SessionRequest)(nil),                     // 8: supernode.SessionRequest
	(*SessionReply)(nil),                       // 9: supernode.SessionReply
}
var file_self_healing_proto_depIdxs = []int32{
	0,  // 0: supernode.SelfHealingData.message_type:type_name -> supernode.SelfHealingData.messageType
	1,  // 1: supernode.SelfHealingData.challenge_status:type_name -> supernode.SelfHealingData.status
	7,  // 2: supernode.SelfHealingData.challenge_file:type_name -> supernode.SelfHealingData.challengeFile
	2,  // 3: supernode.ProcessSelfHealingChallengeRequest.data:type_name -> supernode.SelfHealingData
	2,  // 4: supernode.ProcessSelfHealingChallengeReply.data:type_name -> supernode.SelfHealingData
	2,  // 5: supernode.VerifySelfHealingChallengeRequest.data:type_name -> supernode.SelfHealingData
	2,  // 6: supernode.VerifySelfHealingChallengeReply.data:type_name -> supernode.SelfHealingData
	8,  // 7: supernode.SelfHealing.Session:input_type -> supernode.SessionRequest
	3,  // 8: supernode.SelfHealing.ProcessSelfHealingChallenge:input_type -> supernode.ProcessSelfHealingChallengeRequest
	5,  // 9: supernode.SelfHealing.VerifySelfHealingChallenge:input_type -> supernode.VerifySelfHealingChallengeRequest
	9,  // 10: supernode.SelfHealing.Session:output_type -> supernode.SessionReply
	4,  // 11: supernode.SelfHealing.ProcessSelfHealingChallenge:output_type -> supernode.ProcessSelfHealingChallengeReply
	6,  // 12: supernode.SelfHealing.VerifySelfHealingChallenge:output_type -> supernode.VerifySelfHealingChallengeReply
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_self_healing_proto_init() }
func file_self_healing_proto_init() {
	if File_self_healing_proto != nil {
		return
	}
	file_common_sn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_self_healing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelfHealingData); i {
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
		file_self_healing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessSelfHealingChallengeRequest); i {
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
		file_self_healing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessSelfHealingChallengeReply); i {
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
		file_self_healing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifySelfHealingChallengeRequest); i {
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
		file_self_healing_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifySelfHealingChallengeReply); i {
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
		file_self_healing_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelfHealingDataChallengeFile); i {
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
			RawDescriptor: file_self_healing_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_self_healing_proto_goTypes,
		DependencyIndexes: file_self_healing_proto_depIdxs,
		EnumInfos:         file_self_healing_proto_enumTypes,
		MessageInfos:      file_self_healing_proto_msgTypes,
	}.Build()
	File_self_healing_proto = out.File
	file_self_healing_proto_rawDesc = nil
	file_self_healing_proto_goTypes = nil
	file_self_healing_proto_depIdxs = nil
}
