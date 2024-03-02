// Copyright (c) 2021-2021 The Pastel Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
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

type SelfHealingMessageMessageType int32

const (
	SelfHealingMessage_MessageType_UNKNOWN                           SelfHealingMessageMessageType = 0
	SelfHealingMessage_MessageType_SELF_HEALING_CHALLENGE_MESSAGE    SelfHealingMessageMessageType = 1
	SelfHealingMessage_MessageType_SELF_HEALING_RESPONSE_MESSAGE     SelfHealingMessageMessageType = 2
	SelfHealingMessage_MessageType_SELF_HEALING_VERIFICATION_MESSAGE SelfHealingMessageMessageType = 3
)

// Enum value maps for SelfHealingMessageMessageType.
var (
	SelfHealingMessageMessageType_name = map[int32]string{
		0: "MessageType_UNKNOWN",
		1: "MessageType_SELF_HEALING_CHALLENGE_MESSAGE",
		2: "MessageType_SELF_HEALING_RESPONSE_MESSAGE",
		3: "MessageType_SELF_HEALING_VERIFICATION_MESSAGE",
	}
	SelfHealingMessageMessageType_value = map[string]int32{
		"MessageType_UNKNOWN":                           0,
		"MessageType_SELF_HEALING_CHALLENGE_MESSAGE":    1,
		"MessageType_SELF_HEALING_RESPONSE_MESSAGE":     2,
		"MessageType_SELF_HEALING_VERIFICATION_MESSAGE": 3,
	}
)

func (x SelfHealingMessageMessageType) Enum() *SelfHealingMessageMessageType {
	p := new(SelfHealingMessageMessageType)
	*p = x
	return p
}

func (x SelfHealingMessageMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SelfHealingMessageMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_self_healing_proto_enumTypes[0].Descriptor()
}

func (SelfHealingMessageMessageType) Type() protoreflect.EnumType {
	return &file_self_healing_proto_enumTypes[0]
}

func (x SelfHealingMessageMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SelfHealingMessageMessageType.Descriptor instead.
func (SelfHealingMessageMessageType) EnumDescriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{6, 0}
}

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId string `protobuf:"bytes,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId    string `protobuf:"bytes,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	ReceiverId  string `protobuf:"bytes,2,opt,name=receiver_id,json=receiverId,proto3" json:"receiver_id,omitempty"`
	IsOnline    bool   `protobuf:"varint,3,opt,name=is_online,json=isOnline,proto3" json:"is_online,omitempty"`
	RespondedAt string `protobuf:"bytes,4,opt,name=responded_at,json=respondedAt,proto3" json:"responded_at,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *PingResponse) GetReceiverId() string {
	if x != nil {
		return x.ReceiverId
	}
	return ""
}

func (x *PingResponse) GetIsOnline() bool {
	if x != nil {
		return x.IsOnline
	}
	return false
}

func (x *PingResponse) GetRespondedAt() string {
	if x != nil {
		return x.RespondedAt
	}
	return ""
}

type ProcessSelfHealingChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessSelfHealingChallengeRequest) Reset() {
	*x = ProcessSelfHealingChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessSelfHealingChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessSelfHealingChallengeRequest) ProtoMessage() {}

func (x *ProcessSelfHealingChallengeRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ProcessSelfHealingChallengeRequest.ProtoReflect.Descriptor instead.
func (*ProcessSelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessSelfHealingChallengeRequest) GetData() *SelfHealingMessage {
	if x != nil {
		return x.Data
	}
	return nil
}

type ProcessSelfHealingChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProcessSelfHealingChallengeReply) Reset() {
	*x = ProcessSelfHealingChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessSelfHealingChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessSelfHealingChallengeReply) ProtoMessage() {}

func (x *ProcessSelfHealingChallengeReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ProcessSelfHealingChallengeReply.ProtoReflect.Descriptor instead.
func (*ProcessSelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessSelfHealingChallengeReply) GetData() *SelfHealingMessage {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifySelfHealingChallengeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifySelfHealingChallengeRequest) Reset() {
	*x = VerifySelfHealingChallengeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifySelfHealingChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifySelfHealingChallengeRequest) ProtoMessage() {}

func (x *VerifySelfHealingChallengeRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use VerifySelfHealingChallengeRequest.ProtoReflect.Descriptor instead.
func (*VerifySelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{4}
}

func (x *VerifySelfHealingChallengeRequest) GetData() *SelfHealingMessage {
	if x != nil {
		return x.Data
	}
	return nil
}

type VerifySelfHealingChallengeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *VerifySelfHealingChallengeReply) Reset() {
	*x = VerifySelfHealingChallengeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifySelfHealingChallengeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifySelfHealingChallengeReply) ProtoMessage() {}

func (x *VerifySelfHealingChallengeReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use VerifySelfHealingChallengeReply.ProtoReflect.Descriptor instead.
func (*VerifySelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{5}
}

func (x *VerifySelfHealingChallengeReply) GetData() *SelfHealingMessage {
	if x != nil {
		return x.Data
	}
	return nil
}

type SelfHealingMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageType     SelfHealingMessageMessageType `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=supernode.SelfHealingMessageMessageType" json:"message_type,omitempty"`
	TriggerId       string                        `protobuf:"bytes,2,opt,name=trigger_id,json=triggerId,proto3" json:"trigger_id,omitempty"`
	Data            []byte                        `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	SenderId        string                        `protobuf:"bytes,4,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	SenderSignature []byte                        `protobuf:"bytes,5,opt,name=sender_signature,json=senderSignature,proto3" json:"sender_signature,omitempty"`
}

func (x *SelfHealingMessage) Reset() {
	*x = SelfHealingMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelfHealingMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelfHealingMessage) ProtoMessage() {}

func (x *SelfHealingMessage) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelfHealingMessage.ProtoReflect.Descriptor instead.
func (*SelfHealingMessage) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{6}
}

func (x *SelfHealingMessage) GetMessageType() SelfHealingMessageMessageType {
	if x != nil {
		return x.MessageType
	}
	return SelfHealingMessage_MessageType_UNKNOWN
}

func (x *SelfHealingMessage) GetTriggerId() string {
	if x != nil {
		return x.TriggerId
	}
	return ""
}

func (x *SelfHealingMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SelfHealingMessage) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *SelfHealingMessage) GetSenderSignature() []byte {
	if x != nil {
		return x.SenderSignature
	}
	return nil
}

type BroadcastSelfHealingMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data            []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	SenderId        string `protobuf:"bytes,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	SenderSignature []byte `protobuf:"bytes,3,opt,name=sender_signature,json=senderSignature,proto3" json:"sender_signature,omitempty"`
	Type            int64  `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *BroadcastSelfHealingMetricsRequest) Reset() {
	*x = BroadcastSelfHealingMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastSelfHealingMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastSelfHealingMetricsRequest) ProtoMessage() {}

func (x *BroadcastSelfHealingMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastSelfHealingMetricsRequest.ProtoReflect.Descriptor instead.
func (*BroadcastSelfHealingMetricsRequest) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{7}
}

func (x *BroadcastSelfHealingMetricsRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *BroadcastSelfHealingMetricsRequest) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *BroadcastSelfHealingMetricsRequest) GetSenderSignature() []byte {
	if x != nil {
		return x.SenderSignature
	}
	return nil
}

func (x *BroadcastSelfHealingMetricsRequest) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

type BroadcastSelfHealingMetricsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastSelfHealingMetricsReply) Reset() {
	*x = BroadcastSelfHealingMetricsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_self_healing_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastSelfHealingMetricsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastSelfHealingMetricsReply) ProtoMessage() {}

func (x *BroadcastSelfHealingMetricsReply) ProtoReflect() protoreflect.Message {
	mi := &file_self_healing_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastSelfHealingMetricsReply.ProtoReflect.Descriptor instead.
func (*BroadcastSelfHealingMetricsReply) Descriptor() ([]byte, []int) {
	return file_self_healing_proto_rawDescGZIP(), []int{8}
}

var File_self_healing_proto protoreflect.FileDescriptor

var file_self_healing_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x65, 0x6c, 0x66, 0x5f, 0x68, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x1a,
	0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x73, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2a, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x8c, 0x01, 0x0a,
	0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69,
	0x73, 0x5f, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x69, 0x73, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x22, 0x57, 0x0a, 0x22, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e,
	0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66,
	0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x55, 0x0a, 0x20, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53,
	0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65,
	0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x56, 0x0a, 0x21, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67,
	0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48,
	0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x54, 0x0a, 0x1f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c,
	0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x31, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x98, 0x03, 0x0a, 0x12, 0x53, 0x65,
	0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x4c, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x29,
	0x0a, 0x10, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xb8, 0x01, 0x0a, 0x0b, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x2e, 0x0a, 0x2a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x48, 0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x43,
	0x48, 0x41, 0x4c, 0x4c, 0x45, 0x4e, 0x47, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45,
	0x10, 0x01, 0x12, 0x2d, 0x0a, 0x29, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x48, 0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x52,
	0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10,
	0x02, 0x12, 0x31, 0x0a, 0x2d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x48, 0x45, 0x41, 0x4c, 0x49, 0x4e, 0x47, 0x5f, 0x56, 0x45,
	0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41,
	0x47, 0x45, 0x10, 0x03, 0x22, 0x94, 0x01, 0x0a, 0x22, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x22, 0x0a, 0x20, 0x42,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c,
	0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32,
	0xf7, 0x03, 0x0a, 0x0b, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x12,
	0x41, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01,
	0x30, 0x01, 0x12, 0x37, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x79, 0x0a, 0x1b, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e,
	0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x2d, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65,
	0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x65, 0x6c,
	0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x76, 0x0a, 0x1a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x12, 0x2c, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69,
	0x6e, 0x67, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67,
	0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x79,
	0x0a, 0x1b, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x48,
	0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x2d, 0x2e,
	0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x73,
	0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x6c, 0x66, 0x48, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x6c, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67, 0x6f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x3b, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_self_healing_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_self_healing_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_self_healing_proto_goTypes = []interface{}{
	(SelfHealingMessageMessageType)(0),         // 0: supernode.SelfHealingMessage.messageType
	(*PingRequest)(nil),                        // 1: supernode.PingRequest
	(*PingResponse)(nil),                       // 2: supernode.PingResponse
	(*ProcessSelfHealingChallengeRequest)(nil), // 3: supernode.ProcessSelfHealingChallengeRequest
	(*ProcessSelfHealingChallengeReply)(nil),   // 4: supernode.ProcessSelfHealingChallengeReply
	(*VerifySelfHealingChallengeRequest)(nil),  // 5: supernode.VerifySelfHealingChallengeRequest
	(*VerifySelfHealingChallengeReply)(nil),    // 6: supernode.VerifySelfHealingChallengeReply
	(*SelfHealingMessage)(nil),                 // 7: supernode.SelfHealingMessage
	(*BroadcastSelfHealingMetricsRequest)(nil), // 8: supernode.BroadcastSelfHealingMetricsRequest
	(*BroadcastSelfHealingMetricsReply)(nil),   // 9: supernode.BroadcastSelfHealingMetricsReply
	(*SessionRequest)(nil),                     // 10: supernode.SessionRequest
	(*SessionReply)(nil),                       // 11: supernode.SessionReply
}
var file_self_healing_proto_depIdxs = []int32{
	7,  // 0: supernode.ProcessSelfHealingChallengeRequest.data:type_name -> supernode.SelfHealingMessage
	7,  // 1: supernode.ProcessSelfHealingChallengeReply.data:type_name -> supernode.SelfHealingMessage
	7,  // 2: supernode.VerifySelfHealingChallengeRequest.data:type_name -> supernode.SelfHealingMessage
	7,  // 3: supernode.VerifySelfHealingChallengeReply.data:type_name -> supernode.SelfHealingMessage
	0,  // 4: supernode.SelfHealingMessage.message_type:type_name -> supernode.SelfHealingMessage.messageType
	10, // 5: supernode.SelfHealing.Session:input_type -> supernode.SessionRequest
	1,  // 6: supernode.SelfHealing.Ping:input_type -> supernode.PingRequest
	3,  // 7: supernode.SelfHealing.ProcessSelfHealingChallenge:input_type -> supernode.ProcessSelfHealingChallengeRequest
	5,  // 8: supernode.SelfHealing.VerifySelfHealingChallenge:input_type -> supernode.VerifySelfHealingChallengeRequest
	8,  // 9: supernode.SelfHealing.BroadcastSelfHealingMetrics:input_type -> supernode.BroadcastSelfHealingMetricsRequest
	11, // 10: supernode.SelfHealing.Session:output_type -> supernode.SessionReply
	2,  // 11: supernode.SelfHealing.Ping:output_type -> supernode.PingResponse
	4,  // 12: supernode.SelfHealing.ProcessSelfHealingChallenge:output_type -> supernode.ProcessSelfHealingChallengeReply
	6,  // 13: supernode.SelfHealing.VerifySelfHealingChallenge:output_type -> supernode.VerifySelfHealingChallengeReply
	9,  // 14: supernode.SelfHealing.BroadcastSelfHealingMetrics:output_type -> supernode.BroadcastSelfHealingMetricsReply
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_self_healing_proto_init() }
func file_self_healing_proto_init() {
	if File_self_healing_proto != nil {
		return
	}
	file_common_sn_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_self_healing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
			switch v := v.(*PingResponse); i {
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
		file_self_healing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_self_healing_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_self_healing_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_self_healing_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelfHealingMessage); i {
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
		file_self_healing_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastSelfHealingMetricsRequest); i {
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
		file_self_healing_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastSelfHealingMetricsReply); i {
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
			NumEnums:      1,
			NumMessages:   9,
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
