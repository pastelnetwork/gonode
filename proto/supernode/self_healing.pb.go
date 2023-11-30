// Code generated by protoc-gen-go. DO NOT EDIT.
// source: self_healing.proto

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

type SelfHealingDataMessageType int32

const (
	SelfHealingData_MessageType_UNKNOWN                           SelfHealingDataMessageType = 0
	SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE     SelfHealingDataMessageType = 1
	SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE     SelfHealingDataMessageType = 2
	SelfHealingData_MessageType_SELF_HEALING_VERIFICATION_MESSAGE SelfHealingDataMessageType = 3
)

var SelfHealingDataMessageType_name = map[int32]string{
	0: "MessageType_UNKNOWN",
	1: "MessageType_SELF_HEALING_ISSUANCE_MESSAGE",
	2: "MessageType_SELF_HEALING_RESPONSE_MESSAGE",
	3: "MessageType_SELF_HEALING_VERIFICATION_MESSAGE",
}

var SelfHealingDataMessageType_value = map[string]int32{
	"MessageType_UNKNOWN":                           0,
	"MessageType_SELF_HEALING_ISSUANCE_MESSAGE":     1,
	"MessageType_SELF_HEALING_RESPONSE_MESSAGE":     2,
	"MessageType_SELF_HEALING_VERIFICATION_MESSAGE": 3,
}

func (x SelfHealingDataMessageType) String() string {
	return proto.EnumName(SelfHealingDataMessageType_name, int32(x))
}

func (SelfHealingDataMessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{2, 0}
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

var SelfHealingDataStatus_name = map[int32]string{
	0: "Status_UNKNOWN",
	1: "Status_PENDING",
	2: "Status_RESPONDED",
	3: "Status_SUCCEEDED",
	4: "Status_FAILED_TIMEOUT",
	5: "Status_FAILED_INCORRECT_RESPONSE",
}

var SelfHealingDataStatus_value = map[string]int32{
	"Status_UNKNOWN":                   0,
	"Status_PENDING":                   1,
	"Status_RESPONDED":                 2,
	"Status_SUCCEEDED":                 3,
	"Status_FAILED_TIMEOUT":            4,
	"Status_FAILED_INCORRECT_RESPONSE": 5,
}

func (x SelfHealingDataStatus) String() string {
	return proto.EnumName(SelfHealingDataStatus_name, int32(x))
}

func (SelfHealingDataStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{2, 1}
}

type SelfHealingMessageMessageType int32

const (
	SelfHealingMessage_MessageType_UNKNOWN                        SelfHealingMessageMessageType = 0
	SelfHealingMessage_MessageType_SELF_HEALING_CHALLENGE_MESSAGE SelfHealingMessageMessageType = 1
	SelfHealingMessage_MessageType_SELF_HEALING_RESPONSE_MESSAGE  SelfHealingMessageMessageType = 2
)

var SelfHealingMessageMessageType_name = map[int32]string{
	0: "MessageType_UNKNOWN",
	1: "MessageType_SELF_HEALING_CHALLENGE_MESSAGE",
	2: "MessageType_SELF_HEALING_RESPONSE_MESSAGE",
}

var SelfHealingMessageMessageType_value = map[string]int32{
	"MessageType_UNKNOWN":                        0,
	"MessageType_SELF_HEALING_CHALLENGE_MESSAGE": 1,
	"MessageType_SELF_HEALING_RESPONSE_MESSAGE":  2,
}

func (x SelfHealingMessageMessageType) String() string {
	return proto.EnumName(SelfHealingMessageMessageType_name, int32(x))
}

func (SelfHealingMessageMessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{7, 0}
}

type PingRequest struct {
	SenderId             string   `protobuf:"bytes,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{0}
}

func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (m *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(m, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

type PingResponse struct {
	SenderId             string   `protobuf:"bytes,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	ReceiverId           string   `protobuf:"bytes,2,opt,name=receiver_id,json=receiverId,proto3" json:"receiver_id,omitempty"`
	IsOnline             bool     `protobuf:"varint,3,opt,name=is_online,json=isOnline,proto3" json:"is_online,omitempty"`
	RespondedAt          string   `protobuf:"bytes,4,opt,name=responded_at,json=respondedAt,proto3" json:"responded_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{1}
}

func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (m *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(m, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

func (m *PingResponse) GetReceiverId() string {
	if m != nil {
		return m.ReceiverId
	}
	return ""
}

func (m *PingResponse) GetIsOnline() bool {
	if m != nil {
		return m.IsOnline
	}
	return false
}

func (m *PingResponse) GetRespondedAt() string {
	if m != nil {
		return m.RespondedAt
	}
	return ""
}

type SelfHealingData struct {
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
	XXX_NoUnkeyedLiteral        struct{}                      `json:"-"`
	XXX_unrecognized            []byte                        `json:"-"`
	XXX_sizecache               int32                         `json:"-"`
}

func (m *SelfHealingData) Reset()         { *m = SelfHealingData{} }
func (m *SelfHealingData) String() string { return proto.CompactTextString(m) }
func (*SelfHealingData) ProtoMessage()    {}
func (*SelfHealingData) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{2}
}

func (m *SelfHealingData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfHealingData.Unmarshal(m, b)
}
func (m *SelfHealingData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfHealingData.Marshal(b, m, deterministic)
}
func (m *SelfHealingData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfHealingData.Merge(m, src)
}
func (m *SelfHealingData) XXX_Size() int {
	return xxx_messageInfo_SelfHealingData.Size(m)
}
func (m *SelfHealingData) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfHealingData.DiscardUnknown(m)
}

var xxx_messageInfo_SelfHealingData proto.InternalMessageInfo

func (m *SelfHealingData) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *SelfHealingData) GetMessageType() SelfHealingDataMessageType {
	if m != nil {
		return m.MessageType
	}
	return SelfHealingData_MessageType_UNKNOWN
}

func (m *SelfHealingData) GetChallengeStatus() SelfHealingDataStatus {
	if m != nil {
		return m.ChallengeStatus
	}
	return SelfHealingData_Status_UNKNOWN
}

func (m *SelfHealingData) GetMerklerootWhenChallengeSent() string {
	if m != nil {
		return m.MerklerootWhenChallengeSent
	}
	return ""
}

func (m *SelfHealingData) GetChallengingMasternodeId() string {
	if m != nil {
		return m.ChallengingMasternodeId
	}
	return ""
}

func (m *SelfHealingData) GetRespondingMasternodeId() string {
	if m != nil {
		return m.RespondingMasternodeId
	}
	return ""
}

func (m *SelfHealingData) GetReconstructedFileHash() []byte {
	if m != nil {
		return m.ReconstructedFileHash
	}
	return nil
}

func (m *SelfHealingData) GetChallengeFile() *SelfHealingDataChallengeFile {
	if m != nil {
		return m.ChallengeFile
	}
	return nil
}

func (m *SelfHealingData) GetChallengeId() string {
	if m != nil {
		return m.ChallengeId
	}
	return ""
}

func (m *SelfHealingData) GetRegTicketId() string {
	if m != nil {
		return m.RegTicketId
	}
	return ""
}

func (m *SelfHealingData) GetActionTicketId() string {
	if m != nil {
		return m.ActionTicketId
	}
	return ""
}

func (m *SelfHealingData) GetIsSenseTicket() bool {
	if m != nil {
		return m.IsSenseTicket
	}
	return false
}

func (m *SelfHealingData) GetSenseFileIds() []string {
	if m != nil {
		return m.SenseFileIds
	}
	return nil
}

type SelfHealingDataChallengeFile struct {
	FileHashToChallenge  string   `protobuf:"bytes,1,opt,name=file_hash_to_challenge,json=fileHashToChallenge,proto3" json:"file_hash_to_challenge,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfHealingDataChallengeFile) Reset()         { *m = SelfHealingDataChallengeFile{} }
func (m *SelfHealingDataChallengeFile) String() string { return proto.CompactTextString(m) }
func (*SelfHealingDataChallengeFile) ProtoMessage()    {}
func (*SelfHealingDataChallengeFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{2, 0}
}

func (m *SelfHealingDataChallengeFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfHealingDataChallengeFile.Unmarshal(m, b)
}
func (m *SelfHealingDataChallengeFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfHealingDataChallengeFile.Marshal(b, m, deterministic)
}
func (m *SelfHealingDataChallengeFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfHealingDataChallengeFile.Merge(m, src)
}
func (m *SelfHealingDataChallengeFile) XXX_Size() int {
	return xxx_messageInfo_SelfHealingDataChallengeFile.Size(m)
}
func (m *SelfHealingDataChallengeFile) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfHealingDataChallengeFile.DiscardUnknown(m)
}

var xxx_messageInfo_SelfHealingDataChallengeFile proto.InternalMessageInfo

func (m *SelfHealingDataChallengeFile) GetFileHashToChallenge() string {
	if m != nil {
		return m.FileHashToChallenge
	}
	return ""
}

type ProcessSelfHealingChallengeRequest struct {
	Data                 *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ProcessSelfHealingChallengeRequest) Reset()         { *m = ProcessSelfHealingChallengeRequest{} }
func (m *ProcessSelfHealingChallengeRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessSelfHealingChallengeRequest) ProtoMessage()    {}
func (*ProcessSelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{3}
}

func (m *ProcessSelfHealingChallengeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessSelfHealingChallengeRequest.Unmarshal(m, b)
}
func (m *ProcessSelfHealingChallengeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessSelfHealingChallengeRequest.Marshal(b, m, deterministic)
}
func (m *ProcessSelfHealingChallengeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessSelfHealingChallengeRequest.Merge(m, src)
}
func (m *ProcessSelfHealingChallengeRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessSelfHealingChallengeRequest.Size(m)
}
func (m *ProcessSelfHealingChallengeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessSelfHealingChallengeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessSelfHealingChallengeRequest proto.InternalMessageInfo

func (m *ProcessSelfHealingChallengeRequest) GetData() *SelfHealingMessage {
	if m != nil {
		return m.Data
	}
	return nil
}

type ProcessSelfHealingChallengeReply struct {
	Data                 *SelfHealingMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ProcessSelfHealingChallengeReply) Reset()         { *m = ProcessSelfHealingChallengeReply{} }
func (m *ProcessSelfHealingChallengeReply) String() string { return proto.CompactTextString(m) }
func (*ProcessSelfHealingChallengeReply) ProtoMessage()    {}
func (*ProcessSelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{4}
}

func (m *ProcessSelfHealingChallengeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessSelfHealingChallengeReply.Unmarshal(m, b)
}
func (m *ProcessSelfHealingChallengeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessSelfHealingChallengeReply.Marshal(b, m, deterministic)
}
func (m *ProcessSelfHealingChallengeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessSelfHealingChallengeReply.Merge(m, src)
}
func (m *ProcessSelfHealingChallengeReply) XXX_Size() int {
	return xxx_messageInfo_ProcessSelfHealingChallengeReply.Size(m)
}
func (m *ProcessSelfHealingChallengeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessSelfHealingChallengeReply.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessSelfHealingChallengeReply proto.InternalMessageInfo

func (m *ProcessSelfHealingChallengeReply) GetData() *SelfHealingMessage {
	if m != nil {
		return m.Data
	}
	return nil
}

type VerifySelfHealingChallengeRequest struct {
	Data                 *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *VerifySelfHealingChallengeRequest) Reset()         { *m = VerifySelfHealingChallengeRequest{} }
func (m *VerifySelfHealingChallengeRequest) String() string { return proto.CompactTextString(m) }
func (*VerifySelfHealingChallengeRequest) ProtoMessage()    {}
func (*VerifySelfHealingChallengeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{5}
}

func (m *VerifySelfHealingChallengeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifySelfHealingChallengeRequest.Unmarshal(m, b)
}
func (m *VerifySelfHealingChallengeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifySelfHealingChallengeRequest.Marshal(b, m, deterministic)
}
func (m *VerifySelfHealingChallengeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifySelfHealingChallengeRequest.Merge(m, src)
}
func (m *VerifySelfHealingChallengeRequest) XXX_Size() int {
	return xxx_messageInfo_VerifySelfHealingChallengeRequest.Size(m)
}
func (m *VerifySelfHealingChallengeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifySelfHealingChallengeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VerifySelfHealingChallengeRequest proto.InternalMessageInfo

func (m *VerifySelfHealingChallengeRequest) GetData() *SelfHealingData {
	if m != nil {
		return m.Data
	}
	return nil
}

type VerifySelfHealingChallengeReply struct {
	Data                 *SelfHealingData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *VerifySelfHealingChallengeReply) Reset()         { *m = VerifySelfHealingChallengeReply{} }
func (m *VerifySelfHealingChallengeReply) String() string { return proto.CompactTextString(m) }
func (*VerifySelfHealingChallengeReply) ProtoMessage()    {}
func (*VerifySelfHealingChallengeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{6}
}

func (m *VerifySelfHealingChallengeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifySelfHealingChallengeReply.Unmarshal(m, b)
}
func (m *VerifySelfHealingChallengeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifySelfHealingChallengeReply.Marshal(b, m, deterministic)
}
func (m *VerifySelfHealingChallengeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifySelfHealingChallengeReply.Merge(m, src)
}
func (m *VerifySelfHealingChallengeReply) XXX_Size() int {
	return xxx_messageInfo_VerifySelfHealingChallengeReply.Size(m)
}
func (m *VerifySelfHealingChallengeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifySelfHealingChallengeReply.DiscardUnknown(m)
}

var xxx_messageInfo_VerifySelfHealingChallengeReply proto.InternalMessageInfo

func (m *VerifySelfHealingChallengeReply) GetData() *SelfHealingData {
	if m != nil {
		return m.Data
	}
	return nil
}

type SelfHealingMessage struct {
	MessageType          SelfHealingMessageMessageType `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=supernode.SelfHealingMessageMessageType" json:"message_type,omitempty"`
	ChallengeId          string                        `protobuf:"bytes,2,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	Data                 []byte                        `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	SenderId             string                        `protobuf:"bytes,4,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	SenderSignature      []byte                        `protobuf:"bytes,5,opt,name=sender_signature,json=senderSignature,proto3" json:"sender_signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SelfHealingMessage) Reset()         { *m = SelfHealingMessage{} }
func (m *SelfHealingMessage) String() string { return proto.CompactTextString(m) }
func (*SelfHealingMessage) ProtoMessage()    {}
func (*SelfHealingMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae6b4541a7de850f, []int{7}
}

func (m *SelfHealingMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfHealingMessage.Unmarshal(m, b)
}
func (m *SelfHealingMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfHealingMessage.Marshal(b, m, deterministic)
}
func (m *SelfHealingMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfHealingMessage.Merge(m, src)
}
func (m *SelfHealingMessage) XXX_Size() int {
	return xxx_messageInfo_SelfHealingMessage.Size(m)
}
func (m *SelfHealingMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfHealingMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SelfHealingMessage proto.InternalMessageInfo

func (m *SelfHealingMessage) GetMessageType() SelfHealingMessageMessageType {
	if m != nil {
		return m.MessageType
	}
	return SelfHealingMessage_MessageType_UNKNOWN
}

func (m *SelfHealingMessage) GetChallengeId() string {
	if m != nil {
		return m.ChallengeId
	}
	return ""
}

func (m *SelfHealingMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SelfHealingMessage) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

func (m *SelfHealingMessage) GetSenderSignature() []byte {
	if m != nil {
		return m.SenderSignature
	}
	return nil
}

func init() {
	proto.RegisterEnum("supernode.SelfHealingDataMessageType", SelfHealingDataMessageType_name, SelfHealingDataMessageType_value)
	proto.RegisterEnum("supernode.SelfHealingDataStatus", SelfHealingDataStatus_name, SelfHealingDataStatus_value)
	proto.RegisterEnum("supernode.SelfHealingMessageMessageType", SelfHealingMessageMessageType_name, SelfHealingMessageMessageType_value)
	proto.RegisterType((*PingRequest)(nil), "supernode.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "supernode.PingResponse")
	proto.RegisterType((*SelfHealingData)(nil), "supernode.SelfHealingData")
	proto.RegisterType((*SelfHealingDataChallengeFile)(nil), "supernode.SelfHealingData.challengeFile")
	proto.RegisterType((*ProcessSelfHealingChallengeRequest)(nil), "supernode.ProcessSelfHealingChallengeRequest")
	proto.RegisterType((*ProcessSelfHealingChallengeReply)(nil), "supernode.ProcessSelfHealingChallengeReply")
	proto.RegisterType((*VerifySelfHealingChallengeRequest)(nil), "supernode.VerifySelfHealingChallengeRequest")
	proto.RegisterType((*VerifySelfHealingChallengeReply)(nil), "supernode.VerifySelfHealingChallengeReply")
	proto.RegisterType((*SelfHealingMessage)(nil), "supernode.SelfHealingMessage")
}

func init() { proto.RegisterFile("self_healing.proto", fileDescriptor_ae6b4541a7de850f) }

var fileDescriptor_ae6b4541a7de850f = []byte{
	// 978 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x56, 0xdd, 0x72, 0xda, 0x46,
	0x14, 0xae, 0xb0, 0xe3, 0x98, 0x03, 0x06, 0xcd, 0xa6, 0xb1, 0x31, 0x9e, 0x34, 0x58, 0xd3, 0xc9,
	0x60, 0xb7, 0xc6, 0x0d, 0x99, 0xe9, 0x5f, 0xae, 0x28, 0xc8, 0xb6, 0xa6, 0x58, 0xb8, 0x12, 0xc4,
	0x33, 0xbd, 0xd1, 0x28, 0x68, 0x11, 0x3b, 0x16, 0x2b, 0xaa, 0x5d, 0x92, 0xe1, 0x01, 0x7a, 0xd7,
	0x57, 0xe8, 0x73, 0xf4, 0x99, 0xda, 0x67, 0xe8, 0x45, 0x47, 0x2b, 0x21, 0x24, 0x1c, 0x88, 0x9b,
	0x3b, 0xf4, 0x9d, 0xef, 0x7c, 0x9c, 0x3d, 0xfa, 0xce, 0xd1, 0x02, 0x62, 0xd8, 0x1b, 0x59, 0x63,
	0x6c, 0x7b, 0x84, 0xba, 0x8d, 0x69, 0xe0, 0x73, 0x1f, 0xe5, 0xd9, 0x6c, 0x8a, 0x03, 0xea, 0x3b,
	0xb8, 0x5a, 0x1e, 0xfa, 0x93, 0x89, 0x4f, 0x2d, 0x46, 0xa3, 0x98, 0x72, 0x0a, 0x85, 0x1b, 0x42,
	0x5d, 0x03, 0xff, 0x36, 0xc3, 0x8c, 0xa3, 0x23, 0xc8, 0x33, 0x4c, 0x1d, 0x1c, 0x58, 0xc4, 0xa9,
	0x48, 0x35, 0xa9, 0x9e, 0x37, 0x76, 0x23, 0x40, 0x73, 0x94, 0x3f, 0x24, 0x28, 0x46, 0x64, 0x36,
	0xf5, 0x29, 0xc3, 0x1b, 0xd9, 0xe8, 0x39, 0x14, 0x02, 0x3c, 0xc4, 0xe4, 0x5d, 0x14, 0xce, 0x89,
	0x30, 0x2c, 0x20, 0xcd, 0x09, 0xb3, 0x09, 0xb3, 0x7c, 0xea, 0x11, 0x8a, 0x2b, 0x5b, 0x35, 0xa9,
	0xbe, 0x6b, 0xec, 0x12, 0xd6, 0x13, 0xcf, 0xe8, 0x18, 0x8a, 0x81, 0xf8, 0x1b, 0x07, 0x3b, 0x96,
	0xcd, 0x2b, 0xdb, 0x22, 0xbd, 0x90, 0x60, 0x2d, 0xae, 0xfc, 0xb3, 0x0b, 0x65, 0x13, 0x7b, 0xa3,
	0xab, 0xe8, 0xb0, 0x1d, 0x9b, 0xdb, 0xe8, 0x19, 0xc0, 0x04, 0x33, 0x66, 0xbb, 0x78, 0x59, 0x52,
	0x3e, 0x46, 0x34, 0x07, 0x69, 0x50, 0x5c, 0x84, 0xf9, 0x7c, 0x8a, 0x45, 0x51, 0xa5, 0xe6, 0x8b,
	0x46, 0xd2, 0xa0, 0xc6, 0x8a, 0x60, 0x23, 0xa6, 0xf7, 0xe7, 0x53, 0x6c, 0x14, 0x52, 0x0f, 0xa8,
	0x0b, 0xf2, 0x70, 0x6c, 0x7b, 0x1e, 0xa6, 0x2e, 0xb6, 0x18, 0xb7, 0xf9, 0x8c, 0x89, 0x43, 0x94,
	0x9a, 0xc7, 0x1b, 0xe4, 0x22, 0xa2, 0x51, 0x4e, 0x52, 0x4d, 0x01, 0xa0, 0x36, 0x7c, 0x31, 0xc1,
	0xc1, 0x9d, 0x87, 0x03, 0xdf, 0xe7, 0xd6, 0xfb, 0x31, 0xa6, 0x56, 0x4a, 0x1d, 0xd3, 0x45, 0x03,
	0x8e, 0x96, 0xac, 0xdb, 0x31, 0xa6, 0xed, 0x44, 0x06, 0x53, 0x8e, 0x7e, 0x84, 0xc3, 0x45, 0x12,
	0xa1, 0xae, 0x35, 0xb1, 0x19, 0x8f, 0xca, 0x08, 0x7b, 0xf1, 0x48, 0xe4, 0x1f, 0xa4, 0x08, 0xd7,
	0x49, 0x5c, 0x73, 0xd0, 0xf7, 0x50, 0x89, 0x7b, 0x7b, 0x3f, 0x75, 0x47, 0xa4, 0xee, 0x2f, 0xe3,
	0x99, 0xcc, 0x6f, 0xe1, 0x20, 0xc0, 0x43, 0x9f, 0x32, 0x1e, 0xcc, 0x86, 0x1c, 0x3b, 0xd6, 0x88,
	0x78, 0xd8, 0x1a, 0xdb, 0x6c, 0x5c, 0x79, 0x5c, 0x93, 0xea, 0x45, 0xe3, 0x69, 0x26, 0x7c, 0x41,
	0x3c, 0x7c, 0x65, 0xb3, 0x31, 0xea, 0x41, 0x69, 0x79, 0xc4, 0x30, 0xa7, 0xb2, 0x5b, 0x93, 0xea,
	0x85, 0x66, 0x7d, 0x43, 0xfb, 0x92, 0x84, 0x50, 0xc5, 0xd8, 0xcb, 0x3c, 0x86, 0x96, 0x59, 0x0a,
	0x12, 0xa7, 0x92, 0x8f, 0x2c, 0x93, 0x60, 0x9a, 0x83, 0x14, 0xd8, 0x0b, 0xb0, 0x6b, 0x71, 0x32,
	0xbc, 0xc3, 0x3c, 0xe4, 0xc0, 0xc2, 0x56, 0x6e, 0x5f, 0x60, 0x9a, 0x83, 0xea, 0x20, 0xdb, 0x43,
	0x4e, 0x7c, 0x9a, 0xa2, 0x15, 0x04, 0xad, 0x14, 0xe1, 0x09, 0xf3, 0x05, 0x94, 0x09, 0x0b, 0xdf,
	0x0e, 0xc3, 0x31, 0xb7, 0x52, 0x14, 0x36, 0xde, 0x23, 0xcc, 0x0c, 0xd1, 0x88, 0x89, 0x14, 0x28,
	0x0a, 0x52, 0x58, 0xa5, 0xe6, 0xb0, 0xca, 0x5e, 0x6d, 0xab, 0x9e, 0x37, 0x32, 0x58, 0xb5, 0x03,
	0x2b, 0xa7, 0x79, 0x05, 0xfb, 0x49, 0x23, 0x2d, 0xee, 0x2f, 0xed, 0x10, 0xbb, 0xfa, 0xc9, 0x28,
	0x6e, 0x64, 0xdf, 0x4f, 0x5c, 0xa0, 0xfc, 0x25, 0x41, 0xc6, 0xa4, 0x07, 0xf0, 0xe4, 0x7a, 0xf9,
	0x68, 0x0d, 0xf4, 0x9f, 0xf5, 0xde, 0xad, 0x2e, 0x7f, 0x86, 0xce, 0xe0, 0x24, 0x1d, 0x30, 0xd5,
	0xee, 0x85, 0x75, 0xa5, 0xb6, 0xba, 0x9a, 0x7e, 0x69, 0x69, 0xa6, 0x39, 0x68, 0xe9, 0x6d, 0xd5,
	0xba, 0x56, 0x4d, 0xb3, 0x75, 0xa9, 0xca, 0xd2, 0x46, 0xba, 0xa1, 0x9a, 0x37, 0x3d, 0xdd, 0x5c,
	0xd2, 0x73, 0xe8, 0x25, 0x9c, 0xad, 0xa5, 0xbf, 0x51, 0x0d, 0xed, 0x42, 0x6b, 0xb7, 0xfa, 0x5a,
	0x4f, 0x4f, 0x52, 0xb6, 0x94, 0x3f, 0x25, 0xd8, 0x89, 0x86, 0x03, 0x21, 0x28, 0x45, 0x53, 0x91,
	0xaa, 0x77, 0x89, 0xdd, 0xa8, 0x7a, 0x47, 0xd3, 0x2f, 0x65, 0x09, 0x7d, 0x0e, 0x72, 0x8c, 0x45,
	0x25, 0x74, 0xd4, 0x8e, 0x9c, 0x4b, 0xa1, 0xe6, 0xa0, 0xdd, 0x56, 0xd5, 0x10, 0xdd, 0x42, 0x87,
	0xf0, 0x34, 0x46, 0x2f, 0x5a, 0x5a, 0x57, 0xed, 0x58, 0x7d, 0xed, 0x5a, 0xed, 0x0d, 0xfa, 0xf2,
	0x36, 0xfa, 0x12, 0x6a, 0xd9, 0x90, 0xa6, 0xb7, 0x7b, 0x86, 0xa1, 0xb6, 0xfb, 0xc9, 0xd1, 0xe4,
	0x47, 0xca, 0x2d, 0x28, 0x37, 0x81, 0x3f, 0xc4, 0x8c, 0xa5, 0x3c, 0x99, 0x34, 0x7e, 0xb1, 0x3e,
	0x5f, 0xc2, 0xb6, 0x63, 0x73, 0x5b, 0xbc, 0xa2, 0x42, 0xf3, 0xd9, 0x87, 0x9d, 0x1c, 0xb7, 0xc6,
	0x10, 0x54, 0x65, 0x00, 0xb5, 0x8d, 0xc2, 0x53, 0x6f, 0xfe, 0x29, 0xb2, 0x26, 0x1c, 0xbf, 0xc1,
	0x01, 0x19, 0xcd, 0x37, 0x95, 0xdb, 0xc8, 0xe8, 0x56, 0xd7, 0x0f, 0x5e, 0x2c, 0xfa, 0x0b, 0x3c,
	0xdf, 0x24, 0x1a, 0x96, 0xfa, 0x7f, 0x25, 0xff, 0xce, 0x01, 0xba, 0x7f, 0x08, 0xd4, 0x5d, 0x59,
	0xd4, 0x92, 0xd8, 0xac, 0x27, 0x1b, 0x4f, 0xbe, 0x7e, 0x57, 0xaf, 0x6e, 0x86, 0xdc, 0xfd, 0xcd,
	0x80, 0xe2, 0xba, 0xb7, 0xc4, 0xca, 0x12, 0xbf, 0xb3, 0x9f, 0xb7, 0xed, 0x95, 0xcf, 0xdb, 0x09,
	0xc8, 0x71, 0x90, 0x11, 0x97, 0xda, 0x7c, 0x16, 0x60, 0xb1, 0x63, 0x8b, 0x46, 0x39, 0xc2, 0xcd,
	0x05, 0xac, 0xfc, 0xfe, 0xd0, 0xa9, 0x6c, 0xc0, 0xe9, 0xda, 0xb9, 0x69, 0x5f, 0xb5, 0xba, 0x5d,
	0x55, 0xbf, 0xfc, 0xf4, 0xb1, 0x6c, 0xfe, 0x9b, 0x83, 0x42, 0xaa, 0x6d, 0xa8, 0x05, 0x8f, 0x4d,
	0xcc, 0x18, 0xf1, 0x29, 0x3a, 0xcc, 0x74, 0x56, 0x60, 0xb1, 0x49, 0xaa, 0x07, 0x1f, 0x0a, 0x4d,
	0xbd, 0x79, 0x5d, 0xfa, 0x46, 0x42, 0xdf, 0xc1, 0x76, 0x78, 0x23, 0x40, 0xfb, 0x29, 0x52, 0xea,
	0x3e, 0x91, 0x49, 0xce, 0x5c, 0x1d, 0xe6, 0x70, 0xb4, 0xc1, 0xf6, 0xe8, 0x2c, 0x9d, 0xf7, 0xd1,
	0xb9, 0xab, 0x7e, 0xf5, 0x50, 0x7a, 0x68, 0xd1, 0x77, 0x50, 0x5d, 0xef, 0x62, 0xf4, 0x75, 0x4a,
	0xea, 0xa3, 0x13, 0x54, 0x3d, 0x7d, 0x20, 0x7b, 0xea, 0xcd, 0x7f, 0x7a, 0xfd, 0xeb, 0x0f, 0x2e,
	0xe1, 0xe3, 0xd9, 0xdb, 0xc6, 0xd0, 0x9f, 0x9c, 0x4f, 0xc3, 0x6f, 0xa8, 0x47, 0x31, 0x7f, 0xef,
	0x07, 0x77, 0xe7, 0xae, 0x1f, 0x4a, 0x9c, 0x8b, 0x3b, 0xd9, 0x79, 0x22, 0xf9, 0x3a, 0xf9, 0xf5,
	0x76, 0x47, 0x84, 0x5e, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0xba, 0x69, 0xdf, 0x8b, 0xe0, 0x09,
	0x00, 0x00,
}
